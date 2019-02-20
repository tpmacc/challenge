package challenge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
	"io/ioutil"
)

var tableData fileTableData

//type TableData interface {
//	position(string) (Position, error)
//	calculateTotals([]string) (Position, error)
//}

type LeagueTable struct {

	Competition Competition
	Season Season
	Standings []LeagueStandings
}

//func (standings Standings) Total() (LeagueStandings, error) {
//
//	for _, leagueStandings := range standings.Standings {
//		if leagueStandings.Type == "TOTAL" {
//			return leagueStandings, nil
//		}
//	}
//
//	return LeagueStandings{}, errors.New("No total standings found")
//}

//func (standings Standings) position(teamName string) (Position, error) {
//
//	LeagueStandings, err := standings.Total()
//	if err != nil {
//		return Position{}, nil
//	}
//
//	for _, position := range LeagueStandings.Table {
//		if position.Team.Name == teamName {
//			return Position{}, nil
//		}
//	}
//
//	return Position{}, fmt.Errorf("No team named %s in league", teamName)
//
//}

type Competition struct {
	ID int
	Area Area
	Name string
	Code string
	Plan string
	LastUpdated string
}

type Season struct {
	ID int
	StartDate string
	EndDate string
	CurrentMatchday int
	Winner string
}

type LeagueStandings struct {
	Stage string
	Type string
	Group interface{}
	Table []Position
}

type Position struct {
	Position int
	Team Team
	PlayedGames int
	Won int
	Draw int
	Lost int
	Points int
	GoalsFor int
	GoalsAgainst int
	GoalDifference int
}

type Team struct {
	ID int
	Name string
	CrestUrl string
}

type Area struct {
	ID int
	Name string
}



type Table struct {
	LeagueCaption string `json:"leagueCaption"`
	Matchday int	`json:"matchDay"`
	Standing []Standing	`json:standing`
}

type Standing struct {
	Position int `json:position`
	TeamName string `json:teamName`
	PlayedGames int `json:playedGames`
	Points int `json:points`
	Goals int `json:goals`
	GoalsAgainst int `json:goalsAgainst`
	GoalDifference int `json:goalDifference`
	Wins int `json:wins`
	Draws int `json:draws`
	Losses int `json:losses`
}

type Challenger struct {
	Name string
	Teams []Position
	Totals Position
}

type Pick struct {
	Name string
	Teams []string
}

type fileTableData struct{
	Table
}

func (data *fileTableData) read(fileName string) error {

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Error reading test data file: data/current.json: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(content)).Decode(&data); err != nil {
		return fmt.Errorf("Error unmarshalling file %s to json: %v", fileName, err)
	}

	return nil
}

func (data fileTableData) standing(target string) (Standing, error) {

	for _, standing := range data.Standing {
		if standing.TeamName == target {
			return standing, nil
		}
	}

	return Standing{}, fmt.Errorf("No team named %s in league", target)
}

func ReadData() TemplateData {

	ds := FileDataSource{dataFile:"data/v2standings.json"}

	leagueTable, err := ds.FetchLeagueTable()
	if err != nil {
		return TemplateData{err: errors.Annotate(err, "Error building challenge standings")}
	}

	return TemplateData{data:leagueTable}
}

func BuildChallengeStandings() TemplateData {

	//rawData := ReadData().data
	//
	//tableData, ok := rawData.(TableData)
	//if !ok {
	//	//fmt.Printf("Unexpected data type.  Expected TableData but got %v", rawData)
	//	return TemplateData{err:errors.Errorf("Unexpected data type.  Expected TableData but got %v", rawData)}
	//}

	ds := FileDataSource{dataFile:"data/v2standings.json"}

	leagueTable, err := ds.FetchLeagueTable()
	if err != nil {
		return TemplateData{err: errors.Annotate(err, "Error building challenge standings")}
	}

	challengers, err := buildChallengeStandings(leagueTable)
	if err != nil {
		//fmt.Printf("Error building challenge standings: %s", err.Error())
		return TemplateData{err: errors.Annotate(err, "Error building challenge standings")}
	}

	return TemplateData{data:challengers}
}

 func buildChallengeStandings(data LeagueTable) ([]Challenger, error) {

	var challengers []Challenger

	for _, pick := range picks() {

		var challenger Challenger

		challenger.Name = pick.Name
		for _, team := range pick.Teams {
			standing, err := statistics(data, team)
			if err != nil {
				return nil, err
			}
			challenger.Teams = append(challenger.Teams, standing)
		}

		totals, err := calculateTotals(data, pick.Teams)
		if err != nil {
			return nil, err
		}

		challenger.Totals = totals
		challengers = append(challengers, challenger)
	}

	return challengers, nil
}

func picks() []Pick {
	return []Pick {
		Pick{Name:"Gleannrua",
			Teams:[]string {
				"Manchester City FC",
				"Chelsea FC",
				"West Ham United FC",
				"Burnley FC",
				"Fulham FC",
			},
		},
		Pick{Name:"Farmingdale",
			Teams:[]string {
				"Liverpool FC",
				"Tottenham Hotspur FC",
				"Everton FC",
				"AFC Bournemouth",
				"Wolverhampton Wanderers FC",
			},
		},
	}
}
