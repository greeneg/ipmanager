package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/ipmanager/helpers"
	"github.com/greeneg/ipmanager/model"
)

func (i *IpManager) CreateUser(c *gin.Context) {
	var json model.ProposedUser
	if err := c.ShouldBindJSON(&json); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s, err := model.CreateUser(json)
	if s {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "User has been added to system"})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (i *IpManager) GetUsers(c *gin.Context) {
	users, err := model.GetUsers()
	helpers.CheckError(err)

	safeUsers := make([]SafeUser, 0)
	for _, user := range users {
		safeUser := SafeUser{}
		safeUser.Id = user.Id
		safeUser.UserName = user.UserName
		safeUser.CreationDate = user.CreationDate

		safeUsers = append(safeUsers, safeUser)
	}

	if users == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": safeUsers})
	}
}

func (i *IpManager) GetUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ent, err := model.GetUserById(id)
	helpers.CheckError(err)

	// don't return the password hash
	safeUser := new(SafeUser)
	safeUser.Id = ent.Id
	safeUser.UserName = ent.UserName
	safeUser.CreationDate = ent.CreationDate

	if ent.UserName == "" {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with user id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, safeUser)
	}
}

func (i *IpManager) GetUserByUserName(c *gin.Context) {
	username := c.Param("name")
	ent, err := model.GetUserByUserName(username)
	helpers.CheckError(err)

	// don't return the password hash
	safeUser := new(SafeUser)
	safeUser.Id = ent.Id
	safeUser.UserName = ent.UserName
	safeUser.CreationDate = ent.CreationDate

	if ent.UserName == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found with user name " + username})
	} else {
		c.IndentedJSON(http.StatusOK, safeUser)
	}
}
