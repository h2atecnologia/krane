package handler

import (
	"errors"

	"github.com/biensupernice/krane/internal/api/http"
	"github.com/biensupernice/krane/internal/deployment"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/gin-gonic/gin"
)

// CreateDeployment : using deployment spec
func CreateDeployment(c *gin.Context) {
	var t deployment.Template
	err := c.ShouldBindJSON(&t)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	// Compare with a zero value composite literal because all fields are comparable
	d := *deployment.FindDeployment(t.Name)
	if d != (deployment.Template{}) {
		http.BadRequest(c, "Deployment with that name already exists")
		return
	}

	deployment.SaveDeployment(&t)

	// Start new deployment thread
	go deployment.Start(t)

	http.Accepted(c)
}

// GetDeployment : get single deployment by name
func GetDeployment(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		errMsg := errors.New("Error getting deployment `name` from params")
		http.BadRequest(c, errMsg)
		return
	}

	// Get deployment by name
	d := deployment.FindDeployment(name)

	// compare an empty deployment against the one found in the store
	if *d == (deployment.Template{}) {
		logger.Debug("Unable to find a deployment by that name")
		http.Ok(c, nil)
		return
	}

	http.Ok(c, &d)
}

// GetDeployments : get all deployments
func GetDeployments(c *gin.Context) { http.Ok(c, deployment.FindAllDeployments()) }