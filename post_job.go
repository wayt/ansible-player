package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func postJobAction(c *gin.Context) {

	f := new(CreateJobForm)
	if err := c.Bind(f); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := f.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Get assotiated playbook
	job, err := GetJob(f.Name)
	if err != nil {
		panic(err)
	}

	if job == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	job.Run()

	status := http.StatusOK
	if job.Error != nil {
		status = http.StatusInternalServerError
	}

	c.JSON(status, job)
}
