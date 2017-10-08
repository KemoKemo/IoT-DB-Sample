package main

import "time"

// ChartData describes the data of the chart.
type ChartData struct {
	DataSetList []DataSet `json:"data_set_list" binding:"required"`
}

// DataSet is the data set of the sensors.
type DataSet struct {
	Date       time.Time `json:"date" binding:"required"`
	SensorList []Sensor  `json:"sensor_list" binding:"required"`
}

// Sensor is the data of a sensor.
type Sensor struct {
	Number       int     `json:"number" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	TemperatureC float64 `json:"temp_c" binding:"required"`
}
