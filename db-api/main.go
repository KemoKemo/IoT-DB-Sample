package main

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"gopkg.in/mgo.v2"
)

type sensorData struct {
	Date       time.Time `json:"date" binding:"required"`
	SensorList []sensor  `json:"sensor-list" binding:"required"`
}

type sensor struct {
	Number       int     `json:"number" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	TemperatureC float64 `json:"temp_c" binding:"required"`
}

// Exit code
const (
	exitCodeOK int = iota
	exitCodeFailed
)

const (
	dbname = "tempdata"
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
	r.GET("/sensors", withVars(withData(s, epGetSensors)))
	r.POST("/sensors", withVars(withData(s, epPostSensors)))
	err = r.Run(*addr)
	if err != nil {
		log.Println("Failed to run a web service", err)
		return exitCodeFailed
	}

	return exitCodeOK
}
