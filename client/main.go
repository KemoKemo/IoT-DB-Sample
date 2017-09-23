package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

const (
	exitCodeOK int = iota
	exitCodeFailed
)

var (
	sensor = flag.String("sensor", "localhost:5000", "an address of the sensor api")
	db     = flag.String("db", "localhost:9000", "an address of the db api")
)

func init() {
	flag.Parse()
}

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	resp, err := http.Get(*sensor)
	if err != nil {
		log.Println(err)
		return exitCodeFailed
	}
	defer resp.Body.Close()

	respDB, err := http.Post(*db, "application/json", resp.Body)
	if err != nil {
		log.Println(err)
		return exitCodeFailed
	}
	defer respDB.Body.Close()

	return exitCodeOK
}
