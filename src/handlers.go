package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

//WhitelistRequest is the request format to validate an IP's country and if it belongs in the passed whitelist
type WhitelistRequest struct {
	WhitelistedCountries []string `json:"whitelisted_countries"`
}

//ResponseStruct is the return response for application handlers
type ResponseStruct struct {
	Response string `json:"response"`
}

//checkWhitelistHandler decodes the request and calls the CheckWhitelist function to validate
//if the passed ip is a whitelisted country
func checkWhitelistHandler(w http.ResponseWriter, r *http.Request) {
	var req WhitelistRequest
	w.Header().Add("Content-Type", "application/json")
	if !strings.EqualFold(r.Method, "Get") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := fmt.Errorf("invalid request type")
		Log(log.ErrorLevel, err.Error(), flag.Lookup("test.v") == nil)
		response := ResponseStruct{Response: err.Error()}
		jsoniter.NewEncoder(w).Encode(response)
		return
	}
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article we
	// wish to delete
	ip := vars["ip"]

	err := jsoniter.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Log(log.ErrorLevel, err.Error(), flag.Lookup("test.v") == nil)
		w.WriteHeader(http.StatusBadRequest)
		response := ResponseStruct{Response: err.Error()}
		jsoniter.NewEncoder(w).Encode(response)
		return
	}

	if ip == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := fmt.Errorf("empty ip value")
		Log(log.ErrorLevel, err.Error(), flag.Lookup("test.v") == nil)
		response := ResponseStruct{Response: err.Error()}
		jsoniter.NewEncoder(w).Encode(response)
		return
	}
	found, err := CheckWhitelist(ip, req.WhitelistedCountries)
	if err != nil {
		Log(log.ErrorLevel, err.Error(), flag.Lookup("test.v") == nil)
		w.WriteHeader(http.StatusInternalServerError)
		response := ResponseStruct{Response: err.Error()}
		jsoniter.NewEncoder(w).Encode(response)
		return
	}
	response := ResponseStruct{}
	switch found {
	case true:
		w.WriteHeader(http.StatusOK)
		response.Response = "whitelisted"
	case false:
		w.WriteHeader(http.StatusOK)
		response.Response = "not whitelisted"
	}
	jsoniter.NewEncoder(w).Encode(response)
	return
}

//getStatusHandler returns the status/heartbeat of the application. in future implementations, we can
//update this to throw different status' based on if the server is updating, there are issues reading
//data, etc.
func getStatusHandler(w http.ResponseWriter, r *http.Request) {
	statusReturn := map[string]string{}
	w.Header().Add("Content-Type", "application/json")

	statusReturn["status"] = "200 - OK"

	jsoniter.NewEncoder(w).Encode(statusReturn)
	return
}
