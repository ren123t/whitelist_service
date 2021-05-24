package main

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type ModelSuite struct {
	suite.Suite
}

//RUN ME - suite for all tests
func TestModelSuite(t *testing.T) {
	modelSuite := new(ModelSuite)
	suite.Run(t, modelSuite)
}

func (suite *ModelSuite) SetupSuite() {
	LogPath = "./logs/"
	Log(log.InfoLevel, "=============== Running Model Suite ======================", true)
	setupDB("./test-data/test-data.mmdb")
}

func (suite *ModelSuite) TearDownSuite() {
	fmt.Println("========== Model Testsuite completed ===========")
	Log(log.InfoLevel, "=============== Model Testsuite completed ======================", true)
	CountryDatabase.Close()
}

func (suite *ModelSuite) TestInvalidCheckWhitelist() {
	Log(log.InfoLevel, fmt.Sprintf("====== Running TestInvalidCheckWhitelist ==========="), true)
	validCountryList := []string{"china", "United States", "coASTa RiCa"}
	validIP := "1.207.235.255"
	testCases := []struct {
		casename             string
		ip                   string
		whitelistedCountries []string
	}{
		{"Empty IP", "", validCountryList},
		{"Invalid IP", "Invalid ip", validCountryList},
	}

	for _, testcase := range testCases {
		_, err := CheckWhitelist(testcase.ip, testcase.whitelistedCountries)

		if !suite.Error(err, "was expecting an error, returned ok") {
			Log(log.InfoLevel, fmt.Sprintf("was expecting an error, returned ok on case %v", testcase.casename), true)
		}
	}

	CountryDatabase.Close()

	_, err := CheckWhitelist(validIP, validCountryList)

	if !suite.Error(err, "was expecting an error, returned ok") {
		Log(log.InfoLevel, "was expecting an error, returned ok on closed dataset", true)
	}

	//reload database data
	setupDB("./test-data/test-data.mmdb")

	fmt.Println("TestInvalidCheckWhitelist Completed")
}

func (suite *ModelSuite) TestInvalidGetCountryData() {
	Log(log.InfoLevel, fmt.Sprintf("====== Running TestInvalidGetCountryData ==========="), true)
	validIP := "1.207.235.255"
	testCases := []struct {
		casename string
		ip       string
	}{
		{"Empty IP", ""},
		{"Invalid IP", "Invalid ip"},
	}

	for _, testcase := range testCases {
		_, err := GetCountryData(testcase.ip)

		if !suite.Error(err, "was expecting an error, returned ok") {
			Log(log.InfoLevel, fmt.Sprintf("was expecting an error, returned ok on case %v", testcase.casename), true)
		}
	}

	CountryDatabase.Close()

	_, err := GetCountryData(validIP)

	if !suite.Error(err, "was expecting an error, returned ok on closed dataset") {
		Log(log.InfoLevel, "was expecting an error, returned ok on closed dataset", true)
	}
	setupDB("./test-data/test-data.mmdb")

	fmt.Println("============ TestInvalidCheckWhitelist Completed ==================")
}

func (suite *ModelSuite) TestGetCountryData() {
	Log(log.InfoLevel, fmt.Sprintf("====== Running TestGetCountryData ==========="), true)
	//test constants from static test-data.mmdb file
	constant := Country{Name: "China", IsoCode: "CN"}
	constantIP := "1.207.235.255"

	//validate correct return results from mmdb file
	result, err := GetCountryData(constantIP)
	if !suite.NoError(err, "was expecting no error, returned %v", err) {
		Log(log.InfoLevel, fmt.Sprintf("was expecting no error, returned %v", err), true)
	}
	if !suite.NotNil(result, "was expecting an item, returned this value: %v", result) {
		Log(log.InfoLevel, fmt.Sprintf("was expecting an item, returned this value: %v", result), true)
	}
	if !suite.Equal(constant.Name, result.Name) {
		Log(log.InfoLevel, fmt.Sprintf("was expecting an %v, returned this value: %v", constant.Name, result.Name), true)
	}
	if !suite.Equal(constant.IsoCode, result.IsoCode) {
		Log(log.InfoLevel, fmt.Sprintf("was expecting an %v, returned this value: %v", constant.IsoCode, result.IsoCode), true)
	}
	fmt.Println("============ TestGetCountryData Completed ==================")
}

func (suite *ModelSuite) TestCheckWhitelist() {
	Log(log.InfoLevel, fmt.Sprintf("====== Running TestCheckWhitelist ==========="), true)
	//test constants from static test-data.mmdb file
	constantWhitelistTrue := []string{"china", "united states"}
	constantWhitelistFalse := []string{"united states"}
	constantIP := "1.207.235.255"

	//check whitelisted slice expected outcome
	result, err := CheckWhitelist(constantIP, constantWhitelistTrue)
	if !suite.NoError(err, "was expecting no error, returned %v", err) {
		Log(log.InfoLevel, fmt.Sprintf("was expecting no error, returned %v", err), true)
	}
	if !suite.NotNil(result, "was expecting an item, returned this value: %v", result) {
		Log(log.InfoLevel, fmt.Sprintf("was expecting an item, returned this value: %v", result), true)
	}
	if !suite.True(result, "expected true value, got false") {
		Log(log.InfoLevel, "expected false value, got false", true)
	}

	//Check non-whitelisted slice expected outcome
	result, err = CheckWhitelist(constantIP, constantWhitelistFalse)
	if !suite.NoError(err, "was expecting no error, returned %v", err) {
		Log(log.InfoLevel, fmt.Sprintf("was expecting no error, returned %v", err), true)
	}
	if !suite.NotNil(result, "was expecting an item, returned this value: %v", result) {
		Log(log.InfoLevel, fmt.Sprintf("was expecting an item, returned this value: %v", result), true)
	}
	if !suite.False(result, "expected false value, got true") {
		Log(log.InfoLevel, "expected false value, got true", true)
	}

	fmt.Println("================ TestGetCountryData Completed =================")
}

func (suite *ModelSuite) TestLog() {

	err := Log(log.DebugLevel, "testing", true)
	if !suite.NoError(err, "was expecting no error, returned %v", err) {
		fmt.Printf("was expecting no error, returned %v\n", err)
	}

	fmt.Println("=========== TestLog Completed ==================")
}
