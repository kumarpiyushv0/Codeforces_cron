package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/piyushkumar/codeforces-calendar-cron/internal/calendar"
	"github.com/piyushkumar/codeforces-calendar-cron/internal/codeforces"
	"github.com/piyushkumar/codeforces-calendar-cron/internal/syncer"
)

func main() {
	log.Println("[Codeforces Cron] Starting sync execution...")

	credentialsJSON := os.Getenv("GOOGLE_CREDENTIALS_JSON")
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	calendarID := os.Getenv("CALENDAR_ID")
	if calendarID == "" {
		calendarID = "primary"
	}

	reminderMinutesStr := os.Getenv("REMINDER_MINUTES")
	reminderMinutes := int64(30)
	if reminderMinutesStr != "" {
		if val, err := strconv.ParseInt(reminderMinutesStr, 10, 64); err == nil && val > 0 {
			reminderMinutes = val
		}
	}

	if credentialsJSON == "" && credentialsFile == "" {
		log.Fatalf("Fatal: neither GOOGLE_CREDENTIALS_JSON nor GOOGLE_APPLICATION_CREDENTIALS is set.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	calClient, err := calendar.NewClient(ctx, calendar.Config{
		CredentialsJSON: credentialsJSON,
		CredentialsFile: credentialsFile,
		CalendarID:      calendarID,
	})
	if err != nil {
		log.Fatalf("Fatal: failed to initialize Google Calendar client: %v", err)
	}

	cfClient := codeforces.NewClient(nil)
	syncService := syncer.NewSyncer(cfClient, calClient, reminderMinutes)

	result, err := syncService.Sync(ctx)
	if err != nil {
		log.Fatalf("Fatal: sync failed: %v", err)
	}

	fmt.Printf("\n=== Sync Summary ===\n")
	fmt.Printf("Total upcoming contests: %d\n", result.TotalUpcomingContests)
	fmt.Printf("Events created:          %d\n", result.CreatedCount)
	fmt.Printf("Events updated:          %d\n", result.UpdatedCount)
	fmt.Printf("Events unchanged:        %d\n", result.UnchangedCount)
	fmt.Printf("Errors encountered:      %d\n", len(result.Errors))

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}
