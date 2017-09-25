package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	s, err := getSettings("collector.yml")
	if err != nil {
		log.Printf("error: %v", err)
		return exitCodeFailed
	}

	for {
		err = getSensorAndPostToDB(s.Sensor, s.DB)
		if err != nil {
			log.Printf("error: %v", err)
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

func getSensorAndPostToDB(sensor, db string) error {
	resp, err := http.Get(sensor)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respDB, err := http.Post(db, "application/json", resp.Body)
	if err != nil {
		return err
	}
	defer respDB.Body.Close()

	return nil
}
