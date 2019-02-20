package challenge

import (
	"fmt"
)

func statistics(leagueTable LeagueTable, team string) (Position, error) {

	for _, next := range leagueTable.Standings[0].Table {
			if next.Team.Name == team {
				return next, nil
			}
		}

	return Position{}, fmt.Errorf("No team named %s in league", team)
}


func calculateTotals(leagueTable LeagueTable, teams []string) (Position, error) {

	var totals Position

	for _, team := range teams {
		stats, err := statistics(leagueTable, team)
		if err != nil {
			return Position{}, nil
		}

		totals.PlayedGames += stats.PlayedGames
		totals.Won += stats.Won
		totals.Draw += stats.Draw
		totals.Lost += stats.Lost
		totals.GoalsFor += stats.GoalsFor
		totals.GoalsAgainst += stats.GoalsAgainst
		totals.GoalDifference += stats.GoalDifference
		totals.Points += stats.Points
	}

	return totals, nil
}

