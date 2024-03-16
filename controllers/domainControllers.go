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

// CreateDomain Add a domain
//
//	@Summary		Create a new domain
//	@Description	Create a new domain
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Param			domain	body	model.Domain	true	"Domain data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/domain [post]
func (i *IpManager) CreateDomain(c *gin.Context) {
	var json model.Domain
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

	s, err := model.CreateDomain(json, userObject.Id)
	if s {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Domain has been added to system"})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// DeleteDomain Add a domain
//
//	@Summary		Delete a domain
//	@Description	Delete a domain
//	@Tags			domain
//	@Produce		json
//	@Param			domainname	path	string	true	"Domain name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/domain/{domainname} [delete]
func (i *IpManager) DeleteDomain(c *gin.Context) {
	domain := c.Param("domainname")
	status, err := model.DeleteDomain(domain)
	if err != nil {
		log.Println("ERROR: Cannot delete domain: " + string(err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove domain! " + string(err.Error())})
		return
	}

	if status {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "domain " + domain + " has been removed from system"})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove domain!"})
	}
}

// GetDomains Retrieve a list of domain
//
//	@Summary		Retrieve a list of domain
//	@Description	Retrieve a list of domain
//	@Tags			domain
//	@Produce		json
//	@Success		200	{object}	model.DomainList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/domains [get]
func (i *IpManager) GetDomains(c *gin.Context) {
	domains, err := model.GetDomains()
	helpers.CheckError(err)

	if domains == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": domains})
	}
}

// GetDomainById Retrieve a domain by Id
//
//	@Summary		Retrieve a domain by Id
//	@Description	Retrieve a domain by Id
//	@Tags			domain
//	@Produce		json
//	@Param			domainid	path	string	true	"Domain Id"
//	@Success		200	{object}	model.Domain
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/domain/id/{domainid} [get]
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

// GetDomainByName Retrieve a domain by DomainName
//
//	@Summary		Retrieve a domain by DomainName
//	@Description	Retrieve a domain by DomainName
//	@Tags			domain
//	@Produce		json
//	@Param			domainname	path	string	true	"Domain name"
//	@Success		200	{object}	model.Domain
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/domain/name/{domainname} [get]
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
