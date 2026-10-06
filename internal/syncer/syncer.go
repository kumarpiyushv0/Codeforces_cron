package syncer

import (
	"context"
	"fmt"
	"log"

	"github.com/piyushkumar/codeforces-calendar-cron/internal/calendar"
	"github.com/piyushkumar/codeforces-calendar-cron/internal/codeforces"
)

// SyncResult holds summary metrics from a sync run.
type SyncResult struct {
	TotalUpcomingContests int      `json:"total_upcoming_contests"`
	CreatedCount          int      `json:"created_count"`
	UpdatedCount          int      `json:"updated_count"`
	UnchangedCount        int      `json:"unchanged_count"`
	Errors                []string `json:"errors,omitempty"`
}

// Syncer coordinates syncing between Codeforces and Google Calendar.
type Syncer struct {
	cfClient        *codeforces.Client
	calClient       *calendar.Client
	reminderMinutes int64
}

// NewSyncer creates a new Syncer instance.
func NewSyncer(cfClient *codeforces.Client, calClient *calendar.Client, reminderMinutes int64) *Syncer {
	if reminderMinutes <= 0 {
		reminderMinutes = 30
	}
	return &Syncer{
		cfClient:        cfClient,
		calClient:       calClient,
		reminderMinutes: reminderMinutes,
	}
}

// Sync performs the sync operation.
func (s *Syncer) Sync(ctx context.Context) (*SyncResult, error) {
	log.Println("[Syncer] Fetching upcoming contests from Codeforces API...")
	contests, err := s.cfClient.FetchUpcomingContests(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch codeforces contests: %w", err)
	}

	result := &SyncResult{
		TotalUpcomingContests: len(contests),
	}

	grouped := GroupContests(contests)
	log.Printf("[Syncer] Found %d upcoming contest(s) grouped into %d event(s). Syncing with Google Calendar...",
		len(contests), len(grouped))

	for _, g := range grouped {
		var secondaryIDs []int64
		if len(g.ContestIDs) > 1 {
			secondaryIDs = g.ContestIDs[1:]
		}

		action, err := s.calClient.UpsertGroupedEvent(
			ctx,
			g.PrimaryID,
			secondaryIDs,
			g.Title,
			g.ContestType,
			g.StartTimeSeconds,
			g.DurationSeconds,
			g.Contests,
			s.reminderMinutes,
		)
		if err != nil {
			errMsg := fmt.Sprintf("failed to sync contest group %s (ID %d): %v", g.Title, g.PrimaryID, err)
			log.Printf("[Syncer Error] %s", errMsg)
			result.Errors = append(result.Errors, errMsg)
			continue
		}

		switch action {
		case "created":
			log.Printf("[Syncer] [+] Created event for: %s", g.Title)
			result.CreatedCount++
		case "updated":
			log.Printf("[Syncer] [~] Updated event for: %s", g.Title)
			result.UpdatedCount++
		case "unchanged":
			result.UnchangedCount++
		}
	}

	log.Printf("[Syncer] Completed sync run: %d created, %d updated, %d unchanged, %d error(s)",
		result.CreatedCount, result.UpdatedCount, result.UnchangedCount, len(result.Errors))

	return result, nil
}
