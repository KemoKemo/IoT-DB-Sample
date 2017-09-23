package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
)

func epGetSensors(c *gin.Context) {
	db := getVars(c, "db").(*mgo.Database)
	column := db.C("sensors")
	q := column.Find(nil)
	var result []*sensorData
	err := q.All(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messge": err.Error()})
	}
	c.JSON(http.StatusOK, &result)
}

func epPostSensors(c *gin.Context) {
	db := getVars(c, "db").(*mgo.Database)
	column := db.C("sensors")
	var data sensorData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"message": err.Error()})
		return
	}
	log.Printf("%v", data)
	err = column.Insert(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
