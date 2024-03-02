package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/greeneg/ipmanager/helpers"
	"github.com/greeneg/ipmanager/model"
)

func (i *IpManager) GetHosts(c *gin.Context) {
	hosts, err := model.GetHosts()
	helpers.CheckError(err)

	if hosts == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": hosts})
	}
}

func (i *IpManager) GetHostByHostName(c *gin.Context) {
	host := c.Param("hostname")
	ent, err := model.GetHostByHostName(host)
	helpers.CheckError(err)

	if ent.HostName == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with host name " + host})
	} else {
		c.IndentedJSON(http.StatusOK, ent)
	}
}

func (i *IpManager) GetHostById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("hostid"))
	ent, err := model.GetHostById(id)
	helpers.CheckError(err)

	if ent.HostName == "" {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with host id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, ent)
	}
}
