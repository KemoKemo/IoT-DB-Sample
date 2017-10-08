package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"gopkg.in/mgo.v2"
)

// Exit code
const (
	exitCodeOK int = iota
	exitCodeFailed
)

const (
	dbName     = "tempdata"
	columnName = "sensors"
)

var (
	dbaddr   = flag.String("dbaddr", "localhost", "an address of the mongodb")
	addr     = flag.String("addr", ":9000", "an address of this api service")
	varsLock sync.RWMutex
	vars     map[*gin.Context]map[string]interface{}
)

func init() {
	vars = map[*gin.Context]map[string]interface{}{}
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

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/sensors", withVars(withData(s, epGetSensors)))
	r.POST("/sensors", withVars(withData(s, epPostSensors)))
	err = r.Run(*addr)
	if err != nil {
		log.Println("Failed to run a DB API service", err)
		return exitCodeFailed
	}

	return exitCodeOK
}
