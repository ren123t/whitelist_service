package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

func TestMainSuite(t *testing.T) {

	mainSuite := new(MainSuite)
	suite.Run(t, mainSuite)
}

type MainSuite struct {
	suite.Suite
}

func (suite *MainSuite) SetupSuite() {
	LogPath = "./logs/"
	Log(log.InfoLevel, "=============== Running Main Test Suite ======================", true)
	Port = "8080"
}

func (suite *MainSuite) TearDownSuite() {
	Log(log.InfoLevel, "========== Main Testsuite completed ===========", true)
	fmt.Println("========== Main Testsuite completed ===========")
}

//TestSetupDB runs setupDB and uses invalid and valid paths. throws expected errors on invalid mmdb paths or files
func (suite *MainSuite) TestSetupDB() {
	Log(log.InfoLevel, fmt.Sprintf("====== Running TestSetupDB ==========="), true)
	tt := []struct {
		testName string
		dbPath   string
		expected string
	}{
		//Both valid and invalid test cases
		{"Valid MMDB", "./test-data/test-data.mmdb", ""},
		{"Invalid MMDB", "./test-data/badFile.mmdb", "invalid argument"},
		{"Invalid MMDB", "INVALIDPATH#!", "open INVALIDPATH#!: no such file or directory"},
	}
	for _, tc := range tt {

		err := setupDB(tc.dbPath)
		switch tc.testName {
		case "Valid MMDB":
			if !suite.NoError(err, "was expecting no Error, returned error") {
				Log(log.InfoLevel, fmt.Sprintf("was expecting no Error, returned error %v", err), true)
			}
		case "Invalid MMDB":
			if !suite.Error(err, "was expecting an error, returned ok") {
				Log(log.InfoLevel, "was expecting an error, returned ok", true)
			}
			if !suite.Equal(tc.expected, err.Error(), "was expecting %v, recieved %v", tc.expected, err) {
				Log(log.InfoLevel, fmt.Sprintf("was expecting %v, recieved %v", tc.expected, err), true)
			}
		}
	}

}

//TestCreateRouter Runs setupRouter and validates it returns information
func (suite *MainSuite) TestSetupRouter() {
	Log(log.InfoLevel, fmt.Sprintf("====== Running TestSetupRouter ==========="), true)
	var route *mux.Router

	route = setupRouter()

	srv := &http.Server{
		Addr:    ":" + Port,
		Handler: route,
	}
	defer srv.Close()
	go func(srv *http.Server) {
		srv.ListenAndServe()
	}(srv)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%v/", Port), nil)
	if err != nil {
		fmt.Printf("%v error in %v request to %v\n", err, "GET", "localhost:"+Port+"/")
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if !suite.NoError(err, "was expecting no error, returned %v", err) {
		Log(log.InfoLevel, fmt.Sprintf("was expecting no Error, returned error %v", err), true)
	}
	if !suite.NotEmpty(resp) {
		Log(log.InfoLevel, fmt.Sprintf("was expecting non-nil response"), true)
	}
}
