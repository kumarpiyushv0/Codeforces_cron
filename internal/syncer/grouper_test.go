package syncer

import (
	"testing"

	"github.com/piyushkumar/codeforces-calendar-cron/internal/codeforces"
)

func TestGroupContests(t *testing.T) {
	contests := []codeforces.Contest{
		{
			ID:               2273,
			Name:             "Codeforces Round (Div. 1)",
			Type:             "CF",
			Phase:            "BEFORE",
			DurationSeconds:  9000,
			StartTimeSeconds: 1791743700,
		},
		{
			ID:               2274,
			Name:             "Codeforces Round (Div. 2)",
			Type:             "CF",
			Phase:            "BEFORE",
			DurationSeconds:  9000,
			StartTimeSeconds: 1791743700,
		},
		{
			ID:               2275,
			Name:             "Codeforces Round 1125 (Div. 3)",
			Type:             "ICPC",
			Phase:            "BEFORE",
			DurationSeconds:  9000,
			StartTimeSeconds: 1791383700,
		},
	}

	grouped := GroupContests(contests)
	if len(grouped) != 2 {
		t.Fatalf("expected 2 grouped contests, got %d", len(grouped))
	}

	if grouped[0].Title != "Codeforces Round (Div. 1 + Div. 2)" {
		t.Errorf("expected merged title 'Codeforces Round (Div. 1 + Div. 2)', got %s", grouped[0].Title)
	}

	if len(grouped[0].ContestIDs) != 2 {
		t.Errorf("expected 2 contest IDs in group 0, got %d", len(grouped[0].ContestIDs))
	}
}
