package calendar

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/piyushkumar/codeforces-calendar-cron/internal/codeforces"
	googlecalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// Client wraps Google Calendar API operations.
type Client struct {
	service    *googlecalendar.Service
	calendarID string
}

// Config holds Google Calendar client configuration.
type Config struct {
	// CredentialsJSON can be a raw JSON string or base64 encoded JSON string
	CredentialsJSON string
	// CredentialsFile is a local file path to the service account JSON file
	CredentialsFile string
	// CalendarID is the target Google Calendar ID (default is "primary")
	CalendarID string
}

// NewClient initializes a Google Calendar service using Service Account credentials.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.CalendarID == "" {
		cfg.CalendarID = "primary"
	}

	var opts []option.ClientOption

	if cfg.CredentialsJSON != "" {
		rawJSON := strings.TrimSpace(cfg.CredentialsJSON)
		// Check if it's base64 encoded
		if !strings.HasPrefix(rawJSON, "{") {
			decoded, err := base64.StdEncoding.DecodeString(rawJSON)
			if err == nil && strings.HasPrefix(strings.TrimSpace(string(decoded)), "{") {
				rawJSON = string(decoded)
			}
		}
		opts = append(opts, option.WithCredentialsJSON([]byte(rawJSON)))
	} else if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	} else {
		return nil, fmt.Errorf("no google credentials provided (set GOOGLE_CREDENTIALS_JSON or GOOGLE_APPLICATION_CREDENTIALS)")
	}

	opts = append(opts, option.WithScopes(googlecalendar.CalendarEventsScope))

	srv, err := googlecalendar.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create google calendar service: %w", err)
	}

	return &Client{
		service:    srv,
		calendarID: cfg.CalendarID,
	}, nil
}

// FormatEventID generates a deterministic, Google Calendar compliant event ID.
// Google Calendar event ID requirement: 5-1024 characters, characters from [a-v0-9].
func FormatEventID(contestID int64) string {
	return fmt.Sprintf("cfcontest%s", strconv.FormatInt(contestID, 10))
}

// UpsertGroupedEvent creates or updates a calendar event for a grouped Codeforces contest.
// It also cleans up any individual events that are now part of this merged group.
func (c *Client) UpsertGroupedEvent(
	ctx context.Context,
	primaryID int64,
	secondaryIDs []int64,
	title string,
	contestType string,
	startTimeSeconds int64,
	durationSeconds int64,
	contests []codeforces.Contest,
	reminderMinutes int64,
) (string, error) {
	eventID := FormatEventID(primaryID)
	startTime := time.Unix(startTimeSeconds, 0).UTC()
	endTime := time.Unix(startTimeSeconds+durationSeconds, 0).UTC()

	var descBuilder strings.Builder
	descBuilder.WriteString(fmt.Sprintf("%s\n\n", title))
	descBuilder.WriteString(fmt.Sprintf("Type: %s\nDuration: %d minutes\n\nContest Links:\n", contestType, durationSeconds/60))

	primaryURL := fmt.Sprintf("https://codeforces.com/contest/%d", primaryID)
	for _, cnt := range contests {
		descBuilder.WriteString(fmt.Sprintf("- %s: https://codeforces.com/contest/%d\n", cnt.Name, cnt.ID))
	}

	description := descBuilder.String()

	desiredEvent := &googlecalendar.Event{
		Id:          eventID,
		Summary:     title,
		Description: description,
		Location:    primaryURL,
		Start: &googlecalendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &googlecalendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		Status: "confirmed",
	}

	if reminderMinutes > 0 {
		desiredEvent.Reminders = &googlecalendar.EventReminders{
			UseDefault: false,
			Overrides: []*googlecalendar.EventReminder{
				{
					Method:  "popup",
					Minutes: reminderMinutes,
				},
				{
					Method:  "popup",
					Minutes: 10,
				},
			},
			ForceSendFields: []string{"UseDefault"},
		}
	} else {
		desiredEvent.Reminders = &googlecalendar.EventReminders{
			UseDefault:      true,
			ForceSendFields: []string{"UseDefault"},
		}
	}

	action := "unchanged"

	// Check if primary event already exists
	existingEvent, err := c.service.Events.Get(c.calendarID, eventID).Context(ctx).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
			// Event does not exist, insert it
			_, insertErr := c.service.Events.Insert(c.calendarID, desiredEvent).Context(ctx).Do()
			if insertErr != nil {
				return "", fmt.Errorf("failed to insert event %s: %w", eventID, insertErr)
			}
			action = "created"
		} else {
			return "", fmt.Errorf("failed to check existing event %s: %w", eventID, err)
		}
	} else {
		// Check if update is needed
		needsUpdate := false
		if existingEvent.Summary != desiredEvent.Summary {
			existingEvent.Summary = desiredEvent.Summary
			needsUpdate = true
		}
		if existingEvent.Description != desiredEvent.Description {
			existingEvent.Description = desiredEvent.Description
			needsUpdate = true
		}
		if existingEvent.Start == nil || existingEvent.Start.DateTime != desiredEvent.Start.DateTime {
			existingEvent.Start = desiredEvent.Start
			needsUpdate = true
		}
		if existingEvent.End == nil || existingEvent.End.DateTime != desiredEvent.End.DateTime {
			existingEvent.End = desiredEvent.End
			needsUpdate = true
		}
		if existingEvent.Status != "confirmed" {
			existingEvent.Status = "confirmed"
			needsUpdate = true
		}
		existingEvent.Reminders = desiredEvent.Reminders

		if needsUpdate {
			_, updateErr := c.service.Events.Update(c.calendarID, eventID, existingEvent).Context(ctx).Do()
			if updateErr != nil {
				return "", fmt.Errorf("failed to update event %s: %w", eventID, updateErr)
			}
			action = "updated"
		}
	}

	// Clean up any secondary IDs that were previously created as standalone events
	for _, secID := range secondaryIDs {
		secEventID := FormatEventID(secID)
		_ = c.DeleteEventIfExists(ctx, secEventID)
	}

	return action, nil
}

// DeleteEventIfExists removes an event by ID if it is currently in the calendar.
func (c *Client) DeleteEventIfExists(ctx context.Context, eventID string) error {
	err := c.service.Events.Delete(c.calendarID, eventID).Context(ctx).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
			return nil
		}
		return err
	}
	log.Printf("[Sync] [-] Removed redundant event %s from calendar.", eventID)
	return nil
}
