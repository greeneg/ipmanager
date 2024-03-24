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

// CreateHost Register a host into the system
//
//	@Summary		Register host
//	@Description	Add a new host
//	@Tags			host
//	@Accept			json
//	@Produce		json
//	@Param			host	body	model.Host	true	"Host Data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/host [post]
func (i *IpManager) CreateHost(c *gin.Context) {
	var json model.Host
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

	s, err := model.CreateHost(json, userObject.Id)
	if s {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Host has been added to system"})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// DeleteHostname Remove a host from the system
//
//	@Summary		Delete a host
//	@Description	Delete a host
//	@Tags			host
//	@Accept			json
//	@Produce		json
//	@Param			hostname	path	string	true	"Hostname"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/host/{hostname} [delete]
func (i *IpManager) DeleteHostname(c *gin.Context) {
	hostname := c.Param("hostname")
	status, err := model.DeleteDomain(hostname)
	if err != nil {
		log.Println("ERROR: Cannot delete host: " + string(err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove host! " + string(err.Error())})
		return
	}

	if status {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "host " + hostname + " has been removed from system"})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove host!"})
	}
}

// UpdateMacAddresses Update a host's MAC address list
//
//	@Summary		Update MAC address list
//	@Description	Update MAC address list
//	@Tags			host
//	@Accept			json
//	@Produce		json
//	@Param			hostname	path	string	true	"Hostname"
//	@Param			updateMacAddresses	body	model.MacAddressList	true	"MAC address list"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/host/{hostname} [patch]
func (i *IpManager) UpdateMacAddresses(c *gin.Context) {
	hostname := c.Param("hostname")
	var json model.MacAddressList
	if err := c.ShouldBindJSON(&json); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := model.UpdateMacAddresses(hostname, json.Data)
	if err != nil {
		log.Println("ERROR: Cannot update host's MAC address list: " + string(err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to update host's MAC address list: " + string(err.Error())})
		return
	}

	if status {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "host " + hostname + "'s MAC address list has been modified"})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to update host's MAC address list"})
	}
}

// GetHosts Retrieve list of all hosts
//
//	@Summary		Retrieve list of all hosts
//	@Description	Retrieve list of all hosts
//	@Tags			host
//	@Produce		json
//	@Success		200	{object}	model.HostList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/hosts [get]
func (i *IpManager) GetHosts(c *gin.Context) {
	hosts, err := model.GetHosts()
	helpers.CheckError(err)

	if hosts == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": hosts})
	}
}

// GetHostByHostName Retrieve a host by its hostname
//
//	@Summary		Retrieve a host by its hostname
//	@Description	Retrieve a host by its hostname
//	@Tags			host
//	@Produce		json
//	@Param			hostname	path	string	true	"hostname"
//	@Success		200	{object}	model.Host
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/host/name/{hostname} [get]
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

// GetHostById Retrieve a host by its Id
//
//	@Summary		Retrieve a host by its Id
//	@Description	Retrieve a host by its Id
//	@Tags			host
//	@Produce		json
//	@Param			hostid	path	string	true	"host Id"
//	@Success		200	{object}	model.Host
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/host/id/{hostid} [get]
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
