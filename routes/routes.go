package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/greeneg/ipmanager/controllers"
)

func PublicRoutes(g *gin.RouterGroup, i *controllers.IpManager) {
	// address related routes
	g.GET("/address/:id", i.GetAddressById)
	g.GET("/address/filter/host/id/:hostid", i.GetAddressByHostNameId)
	g.GET("/address/filter/host/name/:hostname", i.GetAddressByHostName)
	g.GET("/address/filter/ip/:ip", i.GetAddressByIpAddress)
	g.GET("/addresses", i.GetAddresses)
	g.GET("/addresses/filter/domain/id/:domainid", i.GetAddressesByDomainId)
	g.GET("/addresses/filter/domain/name/:domainname", i.GetAddressesByDomainName)
	g.GET("/addresses/filter/subnet/id/:subnetid", i.GetAddressesBySubnetId)
	g.GET("/addresses/filter/subnet/name/:subnetname", i.GetAddressesBySubnetName)
	g.GET("/addresses/filter/subnet/name/:subnetname/unassigned")
	// domain related routes
	g.GET("/domain/filter/id/:domainid", i.GetDomainById)
	g.GET("/domain/filter/name/:domainname", i.GetDomainByDomainName)
	g.GET("/domains", i.GetDomains)
	// host related routes
	g.GET("/host/filter/id/:hostid", i.GetHostById)
	g.GET("/host/filter/name/:hostname", i.GetHostByHostName)
	g.GET("/hosts", i.GetHosts)
	// subnet related routes
	g.GET("/subnet/filter/id/:subnetid", i.GetSubnetById)
	g.GET("/subnet/filter/name/:subnetname", i.GetSubnetByNetworkName)
	g.GET("/subnets", i.GetSubnets)
	g.GET("/subnets/filter/domain/id/:domainid", i.GetSubnetsByDomainId)
	g.GET("/subnets/filter/domain/name/:domainname", i.GetSubnetsByDomainName)
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
	g.POST("/host", i.CreateHost) // create a host
	g.PATCH("/host/:hostname")    // replace a host's MAC addresses
	g.DELETE("/host/:hostname")   // trash a host
	// subnet related routes
	g.POST("/subnet")                // create new subnet
	g.PATCH("/subnet/:networkname")  // update a subnet's network information
	g.DELETE("/subnet/:networkname") // trash a subnet
	// user related routes
	g.POST("/user", i.CreateUser)         // create new user
	g.PATCH("/user/:name")                // update a user password
	g.PATCH("/user/:name/status")         // lock a user
	g.GET("/user/:name/status")           // get whether a user is locked or not
	g.DELETE("/user/:name", i.DeleteUser) // trash a user
}
