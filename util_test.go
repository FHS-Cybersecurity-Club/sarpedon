package main

import "testing"

func TestApplyConfiguredTeamAliases(t *testing.T) {
	previousConfig := sarpConfig
	t.Cleanup(func() { sarpConfig = previousConfig })
	sarpConfig.Team = []teamData{{ID: "team-1", Alias: "New Name"}}

	scores := []scoreEntry{
		{Team: teamData{ID: "team-1", Alias: "Old Name", Email: "saved@example.org"}, Points: 10},
		{Team: teamData{ID: "team-1", Alias: "Older Name", Email: "saved@example.org"}, Points: 5},
		{Team: teamData{ID: "removed-team", Alias: "Historical Name"}, Points: 3},
	}

	applyConfiguredTeamAliases(scores)

	if scores[0].Team.Alias != "New Name" || scores[1].Team.Alias != "New Name" {
		t.Fatalf("configured alias was not applied to all scores: %+v", scores)
	}
	if scores[0].Team.ID != "team-1" || scores[0].Team.Email != "saved@example.org" {
		t.Fatalf("score identity or email changed: %+v", scores[0].Team)
	}
	if scores[2].Team.Alias != "Historical Name" {
		t.Fatalf("alias for removed team changed: %+v", scores[2].Team)
	}

	teams, err := parseScoresIntoTeams(scores)
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 2 || teams[0].Alias != "New Name" || teams[0].Score != 15 {
		t.Fatalf("leaderboard did not show the configured alias: %+v", teams)
	}
}
