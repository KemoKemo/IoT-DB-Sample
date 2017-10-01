package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	exitCodeOK int = iota
	exitCodeFailed
)

type sensorData struct {
	Date       time.Time `json:"date"`
	SensorList []sensor  `json:"sensor-list"`
}

type sensor struct {
	Number       int     `json:"number"`
	Name         string  `json:"name"`
	TemperatureC float64 `json:"temp_c"`
}

const (
	dbname = "tempdata"
)

var (
	dbaddr = flag.String("dbaddr", "localhost", "an address of the mongodb")
)

func init() {
	flag.Parse()
}

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	s, err := mgo.Dial(*dbaddr)
	if err != nil {
		log.Println("Failed to dial to db", err)
		return exitCodeFailed
	}
	defer s.Close()

	list := []sensorData{}
	err = s.DB(dbname).C("sensors").Find(
		bson.M{"date": bson.M{"$lte": time.Now()}},
	).Sort("-$natural").Limit(10).All(&list)
	if err != nil {
		log.Println("Failed to find", err)
		return exitCodeFailed
	}

	fmt.Println("sensors:")
	for _, item := range list {
		fmt.Printf(" %v\n", item)
	}

	return exitCodeOK
}
