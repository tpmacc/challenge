package challenge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/icedream/go-footballdata"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestParseStandings(t *testing.T) {

	data := new(LeagueTable)

	content, err := ioutil.ReadFile("data/v2standings.json")
	assert.NoError(t, err, "Error reading test data file: data/current.json: %v", err)

	err = json.NewDecoder(bytes.NewReader(content)).Decode(&data)
	assert.NoError(t, err, "Error unmarshalling content to json")

	assert.Equal(t, "Premier League", data.Competition.Name)
	assert.Equal(t, "Manchester City FC", data.Standings[0].Table[0].Team.Name)
}

func TestBuildChallengeStandings(t *testing.T) {
	// Create client (optionally with auth token)
	client := new(footballdata.Client).
		WithToken(viper.GetString("authToken"))

	// Get list of seasons...
	competitions, err := client.Competitions().Do()
	if err != nil {
		panic(err)
	}

	// ...and print them
	for _, competition := range competitions {
		fmt.Println(competition.Id, competition.Caption)
	}
}
