package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

func TestHandlersSuite(t *testing.T) {
	handlerSuite := new(HandlerSuite)
	suite.Run(t, handlerSuite)
}

type HandlerSuite struct {
	suite.Suite
	Request WhitelistRequest
}

func (suite *HandlerSuite) SetupSuite() {
	LogPath = "./logs/"
	Log(log.InfoLevel, "=============== Running Handlers Suite ======================", true)
	setupDB("./test-data/test-data.mmdb")
	Port = "8080"
}

func (suite *HandlerSuite) TearDownSuite() {
	Log(log.InfoLevel, fmt.Sprintf("========== Handlers Testsuite completed ==========="), true)
	fmt.Println("========== Handlers Testsuite completed ===========")
	CountryDatabase.Close()
}

func (suite *HandlerSuite) TestCheckWhitelistHandler() {
	Log(log.InfoLevel, fmt.Sprintf("====== Running TestCheckWhitelistHandler ==========="), true)
	var wl interface{}
	tt := []struct {
		testName string
		request  WhitelistRequest
		expected string
	}{
		// --------------------TODO--------------------
		// Add Test Cases to get challenge and accept response to assure
		// that the writer is writing the correct response each time.
		{"Empty Values But Valid format", suite.Request, "empty ip value"},
		{"Not Found", suite.Request, "not whitelisted"},
		{"Not Get", suite.Request, "invalid request type"},
		{"Found", suite.Request, "whitelisted"},
		{"Invalid Json", suite.Request, "readObjectStart: expect { or n, but found \", error found in #1 byte of ...|\"INVALID#!#|..., bigger context ...|\"INVALID#!#!\"|..."},
		{"Closed database", suite.Request, "cannot call Lookup on a closed database"},
	}

	// Edit variables so they are realated to tc
	for _, tc := range tt {
		err := CountryDatabase.Lookup(net.ParseIP("1.207.235.255"), nil)
		if err.Error() == "cannot call Lookup on a closed database" {
			setupDB("./test-data/test-data.mmdb")
		}
		httpMethod := http.MethodGet
		wl = tc.request
		ip := ""
		switch tc.testName {
		case "Not Get":
			httpMethod = http.MethodPost
			ip = "1.207.235.255"
		case "Empty Values But Valid format":
			request := wl.(WhitelistRequest)
			request.WhitelistedCountries = []string{}
			wl = request
		case "Invalid Json":
			ip = "1.207.235.255"
			wl = "INVALID#!#!"
		case "Not Found":
			request := wl.(WhitelistRequest)
			ip = "1.207.235.255"
			request.WhitelistedCountries = []string{"United States", "Brazil"}
			wl = request
		case "Found":
			request := wl.(WhitelistRequest)
			ip = "1.207.235.255"
			request.WhitelistedCountries = []string{"China", "Brazil"}
			wl = request
		case "Closed database":
			request := wl.(WhitelistRequest)
			ip = "1.207.235.255"
			request.WhitelistedCountries = []string{"China", "Brazil"}
			wl = request
			CountryDatabase.Close()
		}

		toSend, err := jsoniter.Marshal(wl)
		if err != nil {
			fmt.Println("Marshalling of request failed")
		}
		req, err := http.NewRequest(httpMethod, fmt.Sprintf("localhost:%v/checkWhitelist", Port), bytes.NewBuffer(toSend))
		if err != nil {
			fmt.Printf("%v error in %v request to %v\n", err, "GET", fmt.Sprintf("localhost:%v/checkWhitelist/%v", Port, ip))
		}

		// Edit the header to display correct information
		req.Header.Add("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{
			"ip": ip,
		})
		// Create a test recorder to act as http writer
		rec := httptest.NewRecorder()

		// Call Name Handler with rec as writer and req as request
		checkWhitelistHandler(rec, req)

		var resp ResponseStruct

		jsoniter.NewDecoder(rec.Body).Decode(&resp)

		if !suite.NotEmpty(resp.Response) {
			Log(log.InfoLevel, fmt.Sprintf("Received a nil response "), true)
		}
		if !suite.Equal(tc.expected, resp.Response,
			fmt.Sprintf("Received a response other than %v, received %v instead from %v", tc.expected, resp.Response, tc.testName)) {
			Log(log.InfoLevel, fmt.Sprintf("Received a response other than %v, received %v instead from %v", tc.expected, resp.Response, tc.testName), true)
		}
	}

	fmt.Println("============== TestCheckWhitelistHandler Completed ================")
}

func (suite *HandlerSuite) TestStatusHandler() {
	Log(log.InfoLevel, fmt.Sprintf("====== Running TesStatusHandler ==========="), true)
	req, err := http.NewRequest(http.MethodGet, "localhost:"+Port+"/", nil)
	if err != nil {
		fmt.Printf("%v error in %v request to %v\n", err, "GET", "localhost:"+Port+"/")
	}

	// Edit the header to display correct information
	req.Header.Add("Content-Type", "application/json")

	// Create a test recorder to act as http writer
	rec := httptest.NewRecorder()

	// Call Name Handler with rec as writer and req as request
	getStatusHandler(rec, req)

	var resp map[string]interface{}
	jsoniter.NewDecoder(rec.Body).Decode(&resp)

	if !suite.NotEmpty(resp) {
		Log(log.InfoLevel, fmt.Sprintf("Received a nil response "), true)
	}
	if !suite.Equal(map[string]interface{}{"status": "200 - OK"}, resp,
		fmt.Sprintf("Received a response other than %v, received %v instead", map[string]interface{}{"status": "200 - OK"}, resp)) {
		Log(log.InfoLevel, fmt.Sprintf("Received a response other than %v, received %v instead", map[string]interface{}{"status": "200 - OK"}, resp), true)
	}

	fmt.Println("========== TestStatusHandler Completed ===================")
}
