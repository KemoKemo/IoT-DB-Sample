package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Person is information of a person
type Person struct {
	Name  string
	Phone string
}

var (
	dbaddr = flag.String("dbaddr", "localhost", "an address of the mongodb")
)

const (
	exitCodeOK int = iota
	exitCodeFailed
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	session, err := mgo.Dial(*dbaddr)
	if err != nil {
		log.Println(err)
		return exitCodeFailed
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Println(err)
		return exitCodeFailed
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Println(err)
		return exitCodeFailed
	}

	fmt.Println("Phone:", result.Phone)
	return exitCodeOK
}
