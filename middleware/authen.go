package middleware

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
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/greeneg/ipmanager/globals"
	"github.com/greeneg/ipmanager/helpers"
)

func processAuthorizationHeader(authHeader string) (string, string) {
	// split the header value at the space
	encodedString := strings.Split(authHeader, " ")

	// remove base64 encoding
	decodedString, _ := base64.StdEncoding.DecodeString(encodedString[1])

	// now lets return both the
	authValues := strings.Split(string(decodedString), ":")

	return authValues[0], authValues[1]
}

func AuthCheck(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		log.Println("INFO: No session found. Attempting to check for authentication headers")
		baHeader := c.GetHeader("Authorization")
		if baHeader == "" {
			log.Println("ERROR: No authentication header found. Aborting")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authorized!"})
			c.Abort()
			return
		}
		// otherwise, lets process that header
		username, password := processAuthorizationHeader(baHeader)
		authStatus := helpers.CheckUserPass(username, password)
		if authStatus {
			session.Set(globals.UserKey, username)
			if err := session.Save(); err != nil {
				c.IndentedJSON(http.StatusInternalServerError,
					gin.H{"error": "failed to save user session"})
			}
			log.Println("INFO: Authenticated")
		} else {
			log.Println("ERROR: Authentication failed. Aborting")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authorized!"})
			c.Abort()
			return
		}
	} else {
		userString := fmt.Sprintf("%v", user)
		log.Println("INFO: Session found: User: " + userString)
		log.Println("INFO: Authenticated")
	}
	c.Next()
}
