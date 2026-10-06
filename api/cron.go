package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/piyushkumar/codeforces-calendar-cron/internal/calendar"
	"github.com/piyushkumar/codeforces-calendar-cron/internal/codeforces"
	"github.com/piyushkumar/codeforces-calendar-cron/internal/syncer"
)

// Handler is the Vercel serverless function entrypoint.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Optional security check: if CRON_SECRET is configured, enforce Bearer authorization
	expectedSecret := os.Getenv("CRON_SECRET")
	if expectedSecret != "" {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") || strings.TrimPrefix(authHeader, "Bearer ") != expectedSecret {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
	}

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Google credentials not configured. Please set GOOGLE_CREDENTIALS_JSON.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Minute)
	defer cancel()

	calClient, err := calendar.NewClient(ctx, calendar.Config{
		CredentialsJSON: credentialsJSON,
		CredentialsFile: credentialsFile,
		CalendarID:      calendarID,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	cfClient := codeforces.NewClient(nil)
	syncService := syncer.NewSyncer(cfClient, calClient, reminderMinutes)

	result, err := syncService.Sync(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
