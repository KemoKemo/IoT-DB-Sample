package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
)

func epGetSensors(c *gin.Context) {
	q, err := rParseQuery(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	db := getVars(c, "db").(*mgo.Database)
	result, err := getChartData(q, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func epPostSensors(c *gin.Context) {
	var data DataSet
	err := c.BindJSON(&data)
	if err != nil {
		log.Println("Failed to bind", err)
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"message": err.Error()})
		return
	}

	db := getVars(c, "db").(*mgo.Database)
	err = db.C(columnName).Insert(data)
	if err != nil {
		log.Println("Failed to insert data", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
