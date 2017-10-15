package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const (
	exitCodeOK int = iota
	exitCodeFailed
)

// Settings describes the settings to collect information
// Memo: The member parameters must be public to unmarshal
type Settings struct {
	Duration int    `yaml:"duration"`
	Sensor   string `yaml:"sensor"`
	DB       string `yaml:"db"`
}

// DataSet is the data set of the sensors.
type DataSet struct {
	Date       time.Time `json:"date"`
	SensorList []Sensor  `json:"sensor_list"`
}

// Sensor is the data of a sensor.
type Sensor struct {
	Number       int     `json:"number"`
	Name         string  `json:"name"`
	TemperatureC float64 `json:"temp_c"`
}

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	s, err := getSettings("collector.yml")
	if err != nil {
		log.Printf("Failed to get settings: %v", err)
		return exitCodeFailed
	}
	log.Printf("Sensor's URI:%s, DB's URI:%s", s.Sensor, s.DB)

	for {
		err = getSensorAndPostToDB(s.Sensor, s.DB)
		if err != nil {
			log.Printf("Faild to get and post: %v", err)
			break
		}
		time.Sleep(time.Duration(s.Duration) * time.Minute)
	}

	return exitCodeOK
}

func getSettings(path string) (*Settings, error) {
	s := &Settings{}
	f, err := os.Open(path)
	if err != nil {
		return s, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return s, err
	}

	err = yaml.Unmarshal(b, s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func getSensorAndPostToDB(sensorURI, dbURI string) error {
	resp, err := http.Get(sensorURI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var data DataSet
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	data.Date = time.Now()

	bj, err := json.Marshal(data)
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(bj))
	respDB, err := http.Post(dbURI, "application/json", reader)
	if err != nil {
		return err
	}
	defer respDB.Body.Close()

	return nil
}
