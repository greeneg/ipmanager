package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/ipmanager/model"
)

type Config struct {
	TcpPort string `json:"tcpPort"`
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

func (i *IpManager) getDomains(c *gin.Context) {
	domains, err := model.GetDomains()
	checkError(err)

	if domains == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": domains})
	}
}

func (i *IpManager) getHosts(c *gin.Context) {
	hosts, err := model.GetHosts()
	checkError(err)

	if hosts == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": hosts})
	}
}

func (i *IpManager) getSubnets(c *gin.Context) {
	snets, err := model.GetSubnets()
	checkError(err)

	if snets == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": snets})
	}
}

func (i *IpManager) getUsers(c *gin.Context) {
	users, err := model.GetUsers()
	checkError(err)

	if users == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": users})
	}
}

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

	r.Run(IpManager.confStruct.TcpPort)
}
