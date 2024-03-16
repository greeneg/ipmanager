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
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/ipmanager/helpers"
	"github.com/greeneg/ipmanager/model"
)

func (i *IpManager) GetAddresses(c *gin.Context) {
	addrs, err := model.GetAddresses()
	helpers.CheckError(err)

	if addrs == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": addrs})
	}
}

func (i *IpManager) GetAddressByHostName(c *gin.Context) {
	hostName := c.Param("hostname")
	ent, err := model.GetAddressByHostName(hostName)
	helpers.CheckError(err)

	if ent.Address == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found for " + hostName})
	} else {
		c.IndentedJSON(http.StatusOK, ent)
	}
}

func (i *IpManager) GetAddressByHostNameId(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("hostid"))
	ent, err := model.GetAddressByHostNameId(id)
	helpers.CheckError(err)

	if ent.Address == "" {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found for id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, ent)
	}
}

func (i *IpManager) GetAddressById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ent, err := model.GetAddressById(id)
	helpers.CheckError(err)

	if ent.Address == "" {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found for id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, ent)
	}
}

func (i *IpManager) GetAddressByIpAddress(c *gin.Context) {
	ip := c.Param("ip")
	ent, err := model.GetAddressByIpAddress(ip)
	helpers.CheckError(err)

	if ent.Address == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found for IP address " + ip})
	} else {
		c.IndentedJSON(http.StatusOK, ent)
	}
}

func (i *IpManager) GetAddressesByDomainId(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("domainid"))
	ent, err := model.GetAddressesByDomainId(id)
	helpers.CheckError(err)

	if ent == nil {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with domain id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": ent})
	}
}

func (i *IpManager) GetAddressesByDomainName(c *gin.Context) {
	domainname := c.Param("domainname")
	ent, err := model.GetAddressesByDomainName(domainname)
	helpers.CheckError(err)

	if ent == nil {
		c.IndentedJSON(http.StatusOK, gin.H{"error": "no records found with domain name " + domainname})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": ent})
	}
}

func (i *IpManager) GetAddressesBySubnetId(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("subnetid"))
	ent, err := model.GetAddressesBySubnetId(id)
	helpers.CheckError(err)

	if ent == nil {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with subnet id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": ent})
	}
}

func (i *IpManager) GetAddressesBySubnetName(c *gin.Context) {
	subnetname := c.Param("subnetname")
	ent, err := model.GetAddressesBySubnetName(subnetname)
	helpers.CheckError(err)

	if ent == nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with subnet name " + subnetname})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": ent})
	}
}
