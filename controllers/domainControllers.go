package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/ipmanager/helpers"
	"github.com/greeneg/ipmanager/model"
)

func (i *IpManager) GetDomains(c *gin.Context) {
	domains, err := model.GetDomains()
	helpers.CheckError(err)

	if domains == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": domains})
	}
}

func (i *IpManager) GetDomainById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("domainid"))
	ent, err := model.GetDomainById(id)
	helpers.CheckError(err)

	if ent.DomainName == "" {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with domain id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, ent)
	}
}

func (i *IpManager) GetDomainByDomainName(c *gin.Context) {
	domain := c.Param("domainname")
	ent, err := model.GetDomainByDomainName(domain)
	helpers.CheckError(err)

	if ent.DomainName == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with domain name " + domain})
	} else {
		c.IndentedJSON(http.StatusOK, ent)
	}
}
