package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oschwald/maxminddb-golang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//config paths for application
const configPath = "./"
const configFile = "config"

//Port that server is listening on
var Port string

//LogPath is the directory for logs
var LogPath string

func main() {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		Log(log.FatalLevel, err.Error(), flag.Lookup("test.v") == nil)
	}

	//read in config values
	Port = viper.GetString("port")
	databasePath := viper.GetString("database path")
	LogPath = viper.GetString("LogPath")
	err = setupDB(databasePath)
	if err != nil {
		Log(log.FatalLevel, err.Error(), flag.Lookup("test.v") == nil)
	}

	router := setupRouter()
	srv := &http.Server{
		Addr:    ":" + Port,
		Handler: router,
	}

	fmt.Printf("------- project is now listening on %v --------- \n", Port)
	log.Fatal(srv.ListenAndServe())
}

//setupRouter is a basic router function that sets up the application handlers
func setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/checkWhitelist/{ip}", checkWhitelistHandler)
	router.HandleFunc("/", getStatusHandler)
	return router
}

//setupDB reads the database file from a specified mmdb path.
//for future releases, we can run a curl job via a cron job or a background
//process that can download and place files into the stage folder. from that point
// we would want to downtime and move, rewrite database files (with timestamps) to the
//rollback folder, and move the stage mmdb file into the data folder and re-run the setupDB
//function. the curl we would want to run against:
//https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=YOUR_LICENSE_KEY&suffix=tar.gz
func setupDB(databasePath string) error {
	db, err := maxminddb.Open(databasePath)
	if err != nil {
		return err
	}
	CountryDatabase = db
	return nil
}
