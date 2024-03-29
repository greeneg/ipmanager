package controllers

/*

  Copyright 2024, YggdrasilSoft, LLC.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/greeneg/ipmanager/helpers"
	"github.com/greeneg/ipmanager/model"
)

// CreateSubnet Register a subnet with the system
//
//	@Summary		Register subnet
//	@Description	Add a new subnet
//	@Tags			subnet
//	@Accept			json
//	@Produce		json
//	@Param			subnet	body	model.Subnet	true	"Subnet Data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/subnet [post]
func (i *IpManager) CreateSubnet(c *gin.Context) {
	var json model.Subnet
	if err := c.ShouldBindJSON(&json); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// need to get our current user context to get the CreatorId
	session := sessions.Default(c)
	user := session.Get("user")
	// if nil, we have an issue
	if user == nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
		return
	}

	// convert user interface to a string
	username := fmt.Sprintf("%v", user)
	// lets output our session user
	log.Println("INFO: Session user: " + username)
	// get our user id
	userObject, err := model.GetUserByUserName(username)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// what is our user Id
	log.Println("INFO: Session user's ID: " + strconv.Itoa(userObject.Id))

	s, err := model.CreateSubnet(json, userObject.Id)
	if s {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Subnet '" + json.NetworkName + "' has been added to system"})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// DeleteSubnet Remove a subnet
//
//	@Summary		Delete subnet
//	@Description	Delete a subnet
//	@Tags			subnet
//	@Accept			json
//	@Produce		json
//	@Param			networkname	path	string	true	"Network name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/subnet/{networkname} [delete]
func (i *IpManager) DeleteSubnet(c *gin.Context) {
	subnetName := c.Param("networkname")
	status, err := model.DeleteSubnet(subnetName)
	if err != nil {
		log.Println("ERROR: Cannot delete subnet: " + string(err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove subnet! " + string(err.Error())})
		return
	}

	if status {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Subnet '" + subnetName + "' has been removed from system"})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove user!"})
	}

}

// ModifySubnet Change a subnet's network information
//
//	@Summary		Change subnet network information
//	@Description	Change subnet network information
//	@Tags			subnet
//	@Accept			json
//	@Produce		json
//	@Param			networkname	path	string	true	"Network name"
//	@Param			subnetUpdate	body	model.SubnetUpdate	true	"Subnet data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/subnet/{networkname} [patch]
func (i *IpManager) ModifySubnet(c *gin.Context) {
	subnetName := c.Param("networkname")
	var json model.SubnetUpdate
	if err := c.ShouldBindJSON(&json); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := model.ModifySubnet(subnetName, json)
	if err != nil {
		log.Println("ERROR: Cannot modify subnet '" + subnetName + "'! " + string(err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to modify subnet '" + subnetName + "'! " + string(err.Error())})
		return
	}

	if status {
		c.IndentedJSON(http.StatusOK, gin.H{})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{})
	}
}

// GetSubnets Retrieve list of all subnets
//
//	@Summary		Retrieve list of all subnets
//	@Description	Retrieve list of all subnets
//	@Tags			subnet
//	@Produce		json
//	@Success		200	{object}	model.Subnets
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/subnets [get]
func (i *IpManager) GetSubnets(c *gin.Context) {
	snets, err := model.GetSubnets()
	helpers.CheckError(err)

	if snets == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": snets})
	}
}

// GetSubnetById Retrieve a subnet by its Id
//
//	@Summary		Retrieve a subnet by its Id
//	@Description	Retrieve a subnet by its Id
//	@Tags			subnet
//	@Produce		json
//	@Param			subnetid	path	string	true	"Subnet Id"
//	@Success		200	{object}	model.Subnet
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/subnet/id/{subnetname} [get]
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

// GetSubnetByNetworkName Retrieve a subnet by its network name
//
//	@Summary		Retrieve a subnet by its network name
//	@Description	Retrieve a subnet by its network name
//	@Tags			subnet
//	@Produce		json
//	@Param			subnetname	path	string	true	"Subnet name"
//	@Success		200	{object}	model.Subnet
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/subnet/name/{subnetname} [get]
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

// GetSubnetsByDomainId Retrieve the list of subnets assigned to a domain Id
//
//	@Summary		Retrieve a list of subnets assigned to a domain Id
//	@Description	Retrieve a list of subnets assigned to a domain Id
//	@Tags			subnet
//	@Produce		json
//	@Param			domainname	path	string	true	"Domain name"
//	@Success		200	{object}	model.Subnets
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/subnets/domain/id/{domainid} [get]
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

// GetSubnetsByDomainName Retrieve the list of subnets assigned to a domain name
//
//	@Summary		Retrieve a list of subnets assigned to a domain name
//	@Description	Retrieve a list of subnets assigned to a domain name
//	@Tags			subnet
//	@Produce		json
//	@Param			domainname	path	string	true	"Domain name"
//	@Success		200	{object}	model.Subnets
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/subnets/domain/name/{domainname} [get]
func (i *IpManager) GetSubnetsByDomainName(c *gin.Context) {
	domainname := c.Param("domainname")
	ent, err := model.GetSubnestByDomainName(domainname)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve records matching domain name " + domainname})
	}

	if ent == nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with domain name " + domainname})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": ent})
	}
}
