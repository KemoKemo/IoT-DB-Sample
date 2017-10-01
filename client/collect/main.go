package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	exitCodeOK int = iota
	exitCodeFailed
)

var (
	sensorBase = flag.String("sensor", "http://localhost:5000", "an address of the sensor api")
	dbBase     = flag.String("db", "http://localhost:9000", "an address of the db api")
	sensorURI  string
	dbURI      string
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

func init() {
	flag.Parse()
	sensorURI = strings.Join([]string{*sensorBase, "api", "v1", "sensors"}, "/")
	dbURI = strings.Join([]string{*dbBase, "sensors"}, "/")
}

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	resp, err := http.Get(sensorURI)
	if err != nil {
		log.Println("Failed to get sensor data", err)
		return exitCodeFailed
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body", err)
		return exitCodeFailed
	}
	fmt.Printf("Status: %v, Body: %v\n", resp.Status, string(b))

	var data sensorData
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Println("Failed to unmarshal", err)
		return exitCodeFailed
	}
	data.Date = time.Now()
	fmt.Printf("data: %v\n", data)

	bj, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal", err)
		return exitCodeFailed
	}
	reader := strings.NewReader(string(bj))

	respDB, err := http.Post(dbURI, "application/json", reader)
	if err != nil {
		log.Println("Failed to post", err)
		return exitCodeFailed
	}
	defer respDB.Body.Close()

	b2, err := ioutil.ReadAll(respDB.Body)
	if err != nil {
		log.Println("Failed to read response body", err)
		return exitCodeFailed
	}
	fmt.Printf("Status: %v, Body: %v\n", respDB.Status, string(b2))

	return exitCodeOK
}
