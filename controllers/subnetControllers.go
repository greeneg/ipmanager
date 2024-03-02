package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/ipmanager/helpers"
	"github.com/greeneg/ipmanager/model"
)

func (i *IpManager) GetSubnets(c *gin.Context) {
	snets, err := model.GetSubnets()
	helpers.CheckError(err)

	if snets == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": snets})
	}
}

func (i *IpManager) GetSubnetById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("subnetid"))
	ent, err := model.GetSubnetById(id)
	helpers.CheckError(err)

	if ent.NetworkName == "" {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with subnet id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, ent)
	}
}

func (i *IpManager) GetSubnetByNetworkName(c *gin.Context) {
	netname := c.Param("subnetname")
	ent, err := model.GetSubnetByNetworkName(netname)
	helpers.CheckError(err)

	if ent.NetworkName == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with subnet name " + netname})
	} else {
		c.IndentedJSON(http.StatusOK, ent)
	}
}

func (i *IpManager) GetSubnetsByDomainId(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("domainid"))
	ent, err := model.GetSubnestByDomainId(id)
	helpers.CheckError(err)

	if ent == nil {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with domain id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": ent})
	}
}

func (i *IpManager) GetSubnetsByDomainName(c *gin.Context) {
	domainname := c.Param("domainname")
	ent, err := model.GetSubnestByDomainName(domainname)
	helpers.CheckError(err)

	if ent == nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with domain name " + domainname})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": ent})
	}
}
