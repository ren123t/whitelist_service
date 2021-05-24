package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	maxminddb "github.com/oschwald/maxminddb-golang"
	log "github.com/sirupsen/logrus"
)

//CountryDatabase Persistant database for country data from maxmind mmdb file
var CountryDatabase *maxminddb.Reader

//Country is the model for our country data, currently we only care about the english name value from
//mmdb's country dataset
type Country struct {
	Name    string `json:"name"`
	IsoCode string `json:"iso_code"`
}

//CheckWhitelist pulls ip country information through the GetCountryData call and validates if it
//is found in the passed whitelisted country string slice. if found, it will return true
func CheckWhitelist(ipString string, whitelistedCountry []string) (bool, error) {
	whitelisted := false
	country, err := GetCountryData(ipString)
	if err != nil {
		Log(log.ErrorLevel, err.Error(), flag.Lookup("test.v") == nil)
		return whitelisted, err
	}
	for _, v := range whitelistedCountry {
		if strings.ToUpper(v) == strings.ToUpper(country.Name) {
			whitelisted = true
		}
	}
	return whitelisted, nil
}

//GetCountryData parses the IP string value and returns a populated Country struct for use from the
//mmdb file
func GetCountryData(ipString string) (Country, error) {
	var country Country
	var record map[string]interface{}
	ip := net.ParseIP(ipString)
	err := CountryDatabase.Lookup(ip, &record)
	if err != nil {
		Log(log.ErrorLevel, err.Error(), flag.Lookup("test.v") == nil)
		return country, err
	}

	//this solution is assuming that the maxmind records will always have this uniform data structure on returning,
	// the way this is implemented, it allows for safe unboxing in the off-chance that there are missing map values
	countryData, ok := record["country"].(map[string]interface{})
	if !ok {
		err := fmt.Errorf("failed to find country value")
		Log(log.ErrorLevel, err.Error(), flag.Lookup("test.v") == nil)
		return country, err
	}
	countryNames, ok := countryData["names"].(map[string]interface{})
	if !ok {
		err := fmt.Errorf("failed to find country names value")
		Log(log.ErrorLevel, err.Error(), flag.Lookup("test.v") == nil)
		return country, err
	}
	name, ok := countryNames["en"].(string)
	if !ok {
		err := fmt.Errorf("failed to find country name english value")
		Log(log.ErrorLevel, err.Error(), flag.Lookup("test.v") == nil)
		return country, err
	}

	country.Name = name

	//version 1.0.0 calls for a list of regular names, so this data is supplementary; however
	//we may want to look toward this in the future since it seems to be a more uniform datatype,
	//allowing universal support for non-english users
	isoCode, ok := countryData["iso_code"].(string)
	if !ok {
		country.IsoCode = "UNKNOWN"
	} else {
		country.IsoCode = isoCode
	}

	return country, nil
}

//Log is a basic logging function. look to implement rolling logs in the future releases
func Log(level log.Level, msg string, runLog bool) error {
	if runLog {
		log.SetLevel(level)
		logpath := fmt.Sprintf("%v%v.log", LogPath, level)
		f, logErr := os.OpenFile(logpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if logErr != nil {
			fmt.Println(logErr)
			return logErr
		}
		defer f.Close()
		log.SetOutput(f)
		switch level {
		case log.PanicLevel:
			log.Panicf(msg)
		case log.FatalLevel:
			log.Fatalf(msg)
		case log.ErrorLevel:
			log.Errorf(msg)
		case log.WarnLevel:
			log.Warnf(msg)
		case log.InfoLevel:
			log.Infof(msg)
		case log.DebugLevel:
			log.Debugf(msg)
		}
	}
	return nil
}
