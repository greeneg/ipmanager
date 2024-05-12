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
	g.GET("/address/:id", i.GetAddressById)                                 // get an addresses details by id
	g.GET("/address/host/id/:hostid", i.GetAddressByHostNameId)             // get address by host's id
	g.GET("/address/host/name/:hostname", i.GetAddressByHostName)           // get address by the hosts name
	g.GET("/address/ip/:ip", i.GetAddressByIpAddress)                       // get the address details by the IP address
	g.GET("/addresses", i.GetAddresses)                                     // get all addresses
	g.GET("/addresses/domain/id/:domainid", i.GetAddressesByDomainId)       // get all addresses by domain id
	g.GET("/addresses/domain/name/:domainname", i.GetAddressesByDomainName) // get all addresses by domain name
	g.GET("/addresses/subnet/id/:subnetid", i.GetAddressesBySubnetId)       // get all addresses from the subnet id
	g.GET("/addresses/subnet/name/:subnetname", i.GetAddressesBySubnetName) // get all addresses by the subnet name
	g.GET("/addresses/subnet/name/:subnetname/unassigned")                  // get all unassigned addresses
	// domain related routes
	g.GET("/domain/id/:domainid", i.GetDomainById)             // get the domain by id
	g.GET("/domain/name/:domainname", i.GetDomainByDomainName) // get the domain by its domain name
	g.GET("/domains", i.GetDomains)                            // get all domains
	// host related routes
	g.GET("/host/id/:hostid", i.GetHostById)           // get a host's details by its host id
	g.GET("/host/name/:hostname", i.GetHostByHostName) // get a host's details by its host name
	g.GET("/hosts", i.GetHosts)                        // get all hosts
	// subnet related routes
	g.GET("/subnet/id/:subnetid", i.GetSubnetById)                      // get a subnet by its id
	g.GET("/subnet/name/:subnetname", i.GetSubnetByNetworkName)         // get a subnet by its name
	g.GET("/subnets", i.GetSubnets)                                     // get all subnets
	g.GET("/subnets/domain/id/:domainid", i.GetSubnetsByDomainId)       // get all subnets by domain id
	g.GET("/subnets/domain/name/:domainname", i.GetSubnetsByDomainName) // get all subnets by domain name
	// user related routes
	g.GET("/user/id/:id", i.GetUserById)           // get a user by id
	g.GET("/user/name/:name", i.GetUserByUserName) // get a user by name
	g.GET("/users", i.GetUsers)                    // get all users
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
	g.PATCH("/subnet/:networkname", i.ModifySubnet)  // update a subnet's network information
	g.DELETE("/subnet/:networkname", i.DeleteSubnet) // trash a subnet
	// user related routes
	g.POST("/user", i.CreateUser)                   // create new user
	g.PATCH("/user/:name", i.ChangeAccountPassword) // update a user password
	g.PATCH("/user/:name/status", i.SetUserStatus)  // lock a user
	g.GET("/user/:name/status", i.GetUserStatus)    // get whether a user is locked or not
	g.DELETE("/user/:name", i.DeleteUser)           // trash a user
}
