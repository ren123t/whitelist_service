
# Whitelist Service API

## This is a service that allows you to verify if an IP is within a group of countries. Here are the instructions:

### *note* this is for local run only. this will not support DNS the way its currently set up

* to run, run ./whitelist_service within the src folder (in linux)

* verify the port in the configuration file is where you want to be running the service on

* call localhost:PORT/checkWhitelisted/IP with a json body of {whitelisted_countries: []string} to verify if an ip is whitelisted or not