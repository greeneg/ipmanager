package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/greeneg/ipmanager/controllers"
	_ "github.com/greeneg/ipmanager/docs"
	"github.com/greeneg/ipmanager/globals"
	"github.com/greeneg/ipmanager/helpers"
	"github.com/greeneg/ipmanager/middleware"
	"github.com/greeneg/ipmanager/model"
	"github.com/greeneg/ipmanager/routes"
)

// @title IpManager API
// @version 0.0.2
// @description A simple API for managing networks

// @contact.name Gary Greene
// @contact.url https://github.com/greeneg/ipmanager

// @securityDefinitions.basic BasicAuth

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api/v1

// @schemas http
func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// lets get our working directory
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	helpers.CheckError(err)

	// config path is derived from app working directory
	configDir := filepath.Join(appdir, "config")

	// now that we have our appdir and configDir, lets read in our app config
	// and marshall it to the Config struct
	config := globals.Config{}
	jsonContent, err := os.ReadFile(filepath.Join(configDir, "config.json"))
	helpers.CheckError(err)
	err = json.Unmarshal(jsonContent, &config)
	helpers.CheckError(err)

	// create an app object that contains our routes and the configuration
	IpManager := new(controllers.IpManager)
	IpManager.AppPath = appdir
	IpManager.ConfigPath = configDir
	IpManager.ConfStruct = config

	err = model.ConnectDatabase(IpManager.ConfStruct.DbPath)
	helpers.CheckError(err)

	// some defaults for using session support
	r.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))

	// API
	public := r.Group("/api/v1")
	routes.PublicRoutes(public, IpManager)

	private := r.Group("/api/v1")
	private.Use(middleware.AuthCheck)
	routes.PrivateRoutes(private, IpManager)

	// swagger doc
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	tcpPort := strconv.Itoa(IpManager.ConfStruct.TcpPort)
	r.Run(":" + tcpPort)
}
