package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getJobAction(c *gin.Context) {

	id := c.Param("id")

	data, err := GetJobLogs(id)
	if err != nil {
		panic(err)
	}

	if data == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, string(data))
}
