package syncer

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/piyushkumar/codeforces-calendar-cron/internal/codeforces"
)

// GroupedContest represents one or more parallel contests occurring at the same time.
type GroupedContest struct {
	PrimaryID        int64
	ContestIDs       []int64
	Title            string
	ContestType      string
	StartTimeSeconds int64
	DurationSeconds  int64
	Contests         []codeforces.Contest
}

var divRegex = regexp.MustCompile(`(?i)\(Div\.?\s*(\d+)\)`)

// GroupContests groups contests that have the exact same start time and duration.
func GroupContests(contests []codeforces.Contest) []GroupedContest {
	type key struct {
		start    int64
		duration int64
	}

	groupsMap := make(map[key][]codeforces.Contest)
	var orderedKeys []key

	for _, c := range contests {
		k := key{start: c.StartTimeSeconds, duration: c.DurationSeconds}
		if _, exists := groupsMap[k]; !exists {
			orderedKeys = append(orderedKeys, k)
		}
		groupsMap[k] = append(groupsMap[k], c)
	}

	var result []GroupedContest
	for _, k := range orderedKeys {
		group := groupsMap[k]
		if len(group) == 1 {
			c := group[0]
			result = append(result, GroupedContest{
				PrimaryID:        c.ID,
				ContestIDs:       []int64{c.ID},
				Title:            c.Name,
				ContestType:      c.Type,
				StartTimeSeconds: c.StartTimeSeconds,
				DurationSeconds:  c.DurationSeconds,
				Contests:         group,
			})
			continue
		}

		// Sort by contest ID
		sort.Slice(group, func(i, j int) bool {
			return group[i].ID < group[j].ID
		})

		var ids []int64
		for _, c := range group {
			ids = append(ids, c.ID)
		}

		title := mergeTitles(group)

		result = append(result, GroupedContest{
			PrimaryID:        group[0].ID,
			ContestIDs:       ids,
			Title:            title,
			ContestType:      group[0].Type,
			StartTimeSeconds: k.start,
			DurationSeconds:  k.duration,
			Contests:         group,
		})
	}

	return result
}

func mergeTitles(contests []codeforces.Contest) string {
	var divs []string
	baseName := ""

	allMatchDiv := true
	for _, c := range contests {
		matches := divRegex.FindStringSubmatch(c.Name)
		if len(matches) > 1 {
			divs = append(divs, fmt.Sprintf("Div. %s", matches[1]))
			if baseName == "" {
				baseName = strings.TrimSpace(divRegex.ReplaceAllString(c.Name, ""))
			}
		} else {
			allMatchDiv = false
		}
	}

	if allMatchDiv && baseName != "" && len(divs) > 1 {
		return fmt.Sprintf("%s (%s)", baseName, strings.Join(divs, " + "))
	}

	// Fallback: join distinct titles with " / "
	var titles []string
	for _, c := range contests {
		titles = append(titles, c.Name)
	}
	return strings.Join(titles, " / ")
}
