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
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/ipmanager/helpers"
	"github.com/greeneg/ipmanager/model"
)

// CreateUser Register a user for authentication and authorization
//
//	@Summary		Register user
//	@Description	Add a new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body	model.ProposedUser	true	"User Data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user [post]
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

// DeleteUser Remove a user for authentication and authorization
//
//	@Summary		Delete user
//	@Description	Delete a user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body	model.User	true	"User Data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/:name [post]
func (i *IpManager) DeleteUser(c *gin.Context) {
	username := c.Param("name")
	status, err := model.DeleteUser(username)
	if err != nil {
		log.Println("ERROR: Cannot delete user: " + string(err.Error()))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove user! " + string(err.Error())})
		return
	}

	if status {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "User " + username + " has been removed from system"})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove user!"})
	}
}

// GetUserStatus Retrieve the active status of a user. Can be either 'enabled' or 'locked'
//
//	@Summary		Retrieve a user's active status. Can be either 'enabled' or 'locked'
//	@Description	Retrieve a user's active status
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body	model.User.UserName	true	"User Data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.UserStatusMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/:name/status [get]
func (i *IpManager) GetUserStatus(c *gin.Context) {
	username := c.Param("name")
	status, err := model.GetUserStatus(username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to get the user " + username + " status: " + string(err.Error())})
		return
	}

	if status != "" {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "User status: " + status, "userStatus": status})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to retrieve user status"})
	}
}

func (i *IpManager) SetUserStatus(c *gin.Context) {
	username := c.Param("name")
	var json model.UserStatus
	if err := c.ShouldBindJSON(&json); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := model.SetUserStatus(username, json)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
		return
	}

	if status {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "User '" + username + "' has been " + json.Status})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
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
