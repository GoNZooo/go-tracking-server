package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	trackingserver "tracking-server/tracking-server"
)

func main() {
	programName := os.Args[0]
	if len(os.Args) < 2 {
		fmt.Println("Usage:", programName, "<port>")
		os.Exit(1)
	}
	arguments := os.Args
	port, err := strconv.ParseInt(arguments[1], 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Starting server on port", port)

	databaseOptions := trackingserver.DatabaseOptions{
		Host:     requireEnvironmentValue("DATABASE_HOST"),
		Port:     requireEnvironmentInteger("DATABASE_PORT"),
		Database: requireEnvironmentValue("DATABASE_DATABASE"),
		User:     requireEnvironmentValue("DATABASE_USER"),
		Password: requireEnvironmentValue("DATABASE_PASSWORD"),
	}

	server := trackingserver.NewServer()
	server.Serve(int(port), databaseOptions)
}

func requireEnvironmentValue(key string) string {
	value, hasValue := os.LookupEnv(key)
	if !hasValue {
		log.Panicln("Required environment value", key, "not present")
	}

	return value
}

func requireEnvironmentInteger(key string) int {
	integerString := requireEnvironmentValue(key)

	value, err := strconv.ParseInt(integerString, 10, 32)
	if err != nil {
		log.Panicln("Required environment value", key, "is not parsable as integer")
	}

	return int(value)
}
