package main

import (
	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
)

func withData(s *mgo.Session, f func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		session := s.Copy()
		defer session.Close()
		setVars(c, "db", session.DB(dbName))
		f(c)
	}
}

func withVars(f func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		openVars(c)
		defer closeVars(c)
		f(c)
	}
}

func openVars(c *gin.Context) {
	varsLock.Lock()
	vars[c] = map[string]interface{}{}
	varsLock.Unlock()
}

func closeVars(c *gin.Context) {
	varsLock.Lock()
	delete(vars, c)
	varsLock.Unlock()
}

func setVars(c *gin.Context, key string, value interface{}) {
	varsLock.Lock()
	vars[c][key] = value
	varsLock.Unlock()
}

func getVars(c *gin.Context, key string) interface{} {
	varsLock.RLock()
	value := vars[c][key]
	varsLock.RUnlock()
	return value
}
