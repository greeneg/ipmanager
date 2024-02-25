package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/ipmanager/model"
)

type Config struct {
	TcpPort int    `json:"tcpPort"`
	DbPath  string `json:"dbPath"`
}

type IpManager struct {
	appPath    string
	configPath string
	confStruct Config
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (i *IpManager) getAddresses(c *gin.Context) {
	addrs, err := model.GetAddresses()
	checkError(err)

	if addrs == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": addrs})
	}
}

/*
func (i *IpManager) getAddressesByDomainId(c *gin.Context)
func (i *IpManager) getAddressesByDomainName(c *gin.Context)
func (i *IpManager) getAddressesBySubnetId(c *gin.Context)
func (i *IpManager) getAddressesBySubnetName(c *gin.Context)
func (i *IpManager) getAddressByHostName(c *gin.Context)
func (i *IpManager) getAddressByHostNameId(c *gin.Context)
func (i *IpManager) getAddressById(c *gin.Context)
func (i *IpManager) getAddressByIpAddress(c *gin.Context)
*/

func (i *IpManager) getDomains(c *gin.Context) {
	domains, err := model.GetDomains()
	checkError(err)

	if domains == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": domains})
	}
}

/*
func (i *IpManager) getDomainById(c *gin.Context)
func (i *IpManager) getDomainByDomainName(c *gin.Context)
*/

func (i *IpManager) getHosts(c *gin.Context) {
	hosts, err := model.GetHosts()
	checkError(err)

	if hosts == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": hosts})
	}
}

// func (i *IpManager) getHostByHostName(c *gin.Context)
// func (i *IpManager) getHostById(c *gin.Context)

func (i *IpManager) getSubnets(c *gin.Context) {
	snets, err := model.GetSubnets()
	checkError(err)

	if snets == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": snets})
	}
}

/*
func (i *IpManager) GetSubnetById(c *gin.Context)
func (i *IpManager) GetSubnetByNetworkName(c *gin.Context)
func (i *IpManager) GetSubnestByDomainId(c *gin.Context)
func (i *IpManager) GetSubnestByDomainName(c *gin.Context)
*/

func (i *IpManager) getUsers(c *gin.Context) {
	users, err := model.GetUsers()
	checkError(err)

	if users == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": users})
	}
}

//func (i *IpManager) GetUserById(c *gin.Context)
//func (i *IpManager) GetUserByUserName(c *gin.Context)

func main() {
	r := gin.Default()

	// lets get our working directory
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkError(err)

	// config path is derived from app working directory
	configDir := filepath.Join(appdir, "config")

	// now that we have our appdir and configDir, lets read in our app config
	// and marshall it to the Config struct
	config := Config{}
	jsonContent, err := os.ReadFile(filepath.Join(configDir, "config.json"))
	checkError(err)
	err = json.Unmarshal(jsonContent, &config)
	checkError(err)

	// create an app object that contains our routes and the configuration
	IpManager := new(IpManager)
	IpManager.appPath = appdir
	IpManager.configPath = configDir
	IpManager.confStruct = config

	err = model.ConnectDatabase(IpManager.confStruct.DbPath)
	checkError(err)

	// API
	router := r.Group("/api/v1")
	{
		router.GET("/addresses", IpManager.getAddresses)
		router.GET("/domains", IpManager.getDomains)
		router.GET("/hosts", IpManager.getHosts)
		router.GET("/subnets", IpManager.getSubnets)
		router.GET("/users", IpManager.getUsers)
	}

	tcpPort := strconv.Itoa(IpManager.confStruct.TcpPort)
	r.Run(":" + tcpPort)
}
