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
