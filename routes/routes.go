package routes

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
	"github.com/gin-gonic/gin"

	"github.com/greeneg/ipmanager/controllers"
)

func PublicRoutes(g *gin.RouterGroup, i *controllers.IpManager) {
	// address related routes
	g.GET("/address/:id", i.GetAddressById)
	g.GET("/address/host/id/:hostid", i.GetAddressByHostNameId)
	g.GET("/address/host/name/:hostname", i.GetAddressByHostName)
	g.GET("/address/ip/:ip", i.GetAddressByIpAddress)
	g.GET("/addresses", i.GetAddresses)
	g.GET("/addresses/domain/id/:domainid", i.GetAddressesByDomainId)
	g.GET("/addresses/domain/name/:domainname", i.GetAddressesByDomainName)
	g.GET("/addresses/subnet/id/:subnetid", i.GetAddressesBySubnetId)
	g.GET("/addresses/subnet/name/:subnetname", i.GetAddressesBySubnetName)
	g.GET("/addresses/subnet/name/:subnetname/unassigned")
	// domain related routes
	g.GET("/domain/id/:domainid", i.GetDomainById)
	g.GET("/domain/name/:domainname", i.GetDomainByDomainName)
	g.GET("/domains", i.GetDomains)
	// host related routes
	g.GET("/host/id/:hostid", i.GetHostById)
	g.GET("/host/name/:hostname", i.GetHostByHostName)
	g.GET("/hosts", i.GetHosts)
	// subnet related routes
	g.GET("/subnet/id/:subnetid", i.GetSubnetById)
	g.GET("/subnet/name/:subnetname", i.GetSubnetByNetworkName)
	g.GET("/subnets", i.GetSubnets)
	g.GET("/subnets/domain/id/:domainid", i.GetSubnetsByDomainId)
	g.GET("/subnets/domain/name/:domainname", i.GetSubnetsByDomainName)
	// user related routes
	g.GET("/user/id/:id", i.GetUserById)
	g.GET("/user/name/:name", i.GetUserByUserName)
	g.GET("/users", i.GetUsers)
	// service related routes
	g.OPTIONS("/")   // API options
	g.GET("/health") // service health
}

func PrivateRoutes(g *gin.RouterGroup, i *controllers.IpManager) {
	// address assignment related routes
	g.POST("/address")            // assign an address to a host
	g.PATCH("/address/:address")  // update an address' assignment
	g.DELETE("/address/:address") // trash an address assignment
	// domain related routes
	g.POST("/domain", i.CreateDomain)               // create a domain
	g.DELETE("/domain/:domainname", i.DeleteDomain) // trash a domain
	// host related routes
	g.POST("/host", i.CreateHost)                    // create a host
	g.PATCH("/host/:hostname", i.UpdateMacAddresses) // replace a host's MAC addresses
	g.DELETE("/host/:hostname", i.DeleteHostname)    // trash a host
	// subnet related routes
	g.POST("/subnet", i.CreateSubnet)                // create new subnet
	g.PATCH("/subnet/:networkname")                  // update a subnet's network information
	g.DELETE("/subnet/:networkname", i.DeleteSubnet) // trash a subnet
	// user related routes
	g.POST("/user", i.CreateUser)                   // create new user
	g.PATCH("/user/:name", i.ChangeAccountPassword) // update a user password
	g.PATCH("/user/:name/status", i.SetUserStatus)  // lock a user
	g.GET("/user/:name/status", i.GetUserStatus)    // get whether a user is locked or not
	g.DELETE("/user/:name", i.DeleteUser)           // trash a user
}
