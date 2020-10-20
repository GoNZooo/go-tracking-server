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

	databaseHost := requireEnvironmentValue("DATABASE_HOST")
	databasePort := requireEnvironmentInteger("DATABASE_PORT")
	databaseDatabase := requireEnvironmentValue("DATABASE_DATABASE")
	databaseUser := requireEnvironmentValue("DATABASE_USER")
	databasePassword := requireEnvironmentValue("DATABASE_PASSWORD")

	trackingserver.Serve(
		int(port),
		databaseHost,
		databasePort,
		databaseDatabase,
		databaseUser,
		databasePassword,
	)
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
