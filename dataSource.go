package challenge

import (
	"bytes"
	"encoding/json"
	"github.com/icedream/go-footballdata"
	"github.com/juju/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"fmt"
)

type DataSource interface {

	FetchLeagueTable() (LeagueTable, error)
}

type FileDataSource struct {
	dataFile string
}

func (ds FileDataSource) FetchLeagueTable() (LeagueTable, error) {

	leagueTable := LeagueTable{}

	content, err := ioutil.ReadFile(ds.dataFile)
	if err != nil {
		return leagueTable, fmt.Errorf("Error reading data file: %s: %v", ds.dataFile, err)
	}

	//pretty, err := json.MarshalIndent(string(content), "", "    ")
	//fmt.Printf("Here's the payload: %s", pretty)

	if err := json.NewDecoder(bytes.NewReader(content)).Decode(&leagueTable); err != nil {
		return leagueTable, errors.Annotate(err, fmt.Sprintf("Error decoding %s", ds.dataFile))
	}

	return leagueTable, nil
}

type ApiDataSource struct {
	client *footballdata.Client
}

func (ds ApiDataSource) FetchLeagueTable() (footballdata.LeagueTable, error) {

	if ds.client == nil {
		ds.client = new(footballdata.Client).
			WithToken(viper.GetString("authcode"))
	}

	return ds.client.LeagueTableOfCompetition(1234).Do()
}


//var ShowStandings = boot.EndpointHandler {
//
//	MakeEndpoint: func (app boot.Application) endpoint.Endpoint {
//		return func(ctx context.Context, request interface{}) (interface{}, error) {
//
//			var table fileTableData
//
//			err := table.read("data/current.json")
//
//			return TemplateData{table, err}, nil
//		}
//	},
//
//	RequestDecoder: func (_ context.Context, request *http.Request) (interface{}, error) {
//		return nil, nil
//	},
//
//	ResponseEncoder: func (_ context.Context, w http.ResponseWriter, response interface{}) error {
//
//		return RenderTemplate(w, "tables/leagueTable", response.(TemplateData))
//
//	},
//
//	//Middleware: []boot.MiddlewareBinder {boot.EntitlementsMiddleware, boot.OIDCAuthenticator},
//
//}


//func FetchTableData() (TableData, error) {
//
//	httpClient := &http.Client{Timeout:10 * time.Second}
//
//	_, err := httpClient.Get("http://api.football-data.org/v2/competitions/2021/standings")
//	if err != nil {
//		return nil, err
//	}
//
//	return nil, nil
//	//return marshallTableData(response)
//}

//func marshallTableData(response *http.Response) (*Standings, error) {
//
//	result := new(Standings)
//
//	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
//		return nil, err
//	}
//
//	return result, nil
//}
