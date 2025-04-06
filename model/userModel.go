package model

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
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"strconv"
	"time"
)

func getStoredPasswordHash(username string) (string, error) {
	log.Println("INFO: Getting stored password hash for user: " + username)
	q, err := DB.Prepare("SELECT PasswordHash FROM Users WHERE UserName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return "", err
	}
	defer q.Close()

	passwordHash := ""
	r, err := q.Query(username)
	if err != nil {
		log.Println("ERROR: Failed to query for password hash")
		return "", err
	}
	defer r.Close()
	err = r.Scan(
		&passwordHash,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No password hash found for user: " + username)
			return "", nil
		}
		log.Println("ERROR: Failed to scan for password hash")
		return "", err
	}

	log.Println("INFO: Retrieved password hash for user: " + username)
	return passwordHash, nil
}

func storeNewPassword(hashedPassword string, username string) (bool, error) {
	log.Println("INFO: Storing new password hash for user: " + username)
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to store new password hash for user: " + username)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to store new password hash for user: " + username)
			t.Rollback()
		}
	}()

	// now we need to create a new transaction to SET the password hash into the DB
	q, err := DB.Prepare("UPDATE Users SET PasswordHash = ?, LastChangedDate = ? WHERE UserName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	// get time stamp
	tStamp := time.Now().Format("2006-01-02 15:04:05") // force into SQL DateTime format

	_, err = q.Exec(hashedPassword, tStamp, username)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: Stored new password hash for user: " + username)
	return true, nil
}

func ChangeAccountPassword(username string, oldPassword string, newPassword string) (bool, error) {
	hashedOldPassword := sha512.Sum512([]byte(oldPassword))
	encodedHashedOldPassword := hex.EncodeToString(hashedOldPassword[:])

	storedHash, err := getStoredPasswordHash(username)
	if err != nil {
		return false, err
	}
	log.Println("INFO: Retrieved stored hash")

	// now get password hash from the db
	if storedHash != encodedHashedOldPassword {
		p := new(PasswordHashMismatch)
		return false, p
	}

	// matches, so hash new password
	hashedNewPassword := sha512.Sum512([]byte(newPassword))
	encodedHashedNewPassword := hex.EncodeToString(hashedNewPassword[:])
	_, err = storeNewPassword(encodedHashedNewPassword, username)
	if err != nil {
		return false, err
	}
	log.Println("INFO: Stored updated hash")

	return true, nil
}

func GetUserById(id int) (User, error) {
	log.Println("INFO: Getting user by ID: " + strconv.Itoa(id))
	idStr := strconv.Itoa(id)
	rec, err := DB.Prepare("SELECT * FROM Users WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return User{}, err
	}
	defer rec.Close()

	user := User{}
	r, err := rec.Query(id)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return User{}, err
	}
	defer r.Close()
	err = r.Scan(
		&user.Id,
		&user.UserName,
		&user.Status,
		&user.PasswordHash,
		&user.CreationDate,
		&user.LastChangedDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No user found with ID: " + idStr)
			return User{}, nil
		}
		log.Println("ERROR: Failed to scan rows")
		return User{}, err
	}

	log.Println("INFO: User " + user.UserName + " found")
	return user, nil
}

func GetUserByUserName(username string) (User, error) {
	log.Println("INFO: Getting user by name: " + username)
	rec, err := DB.Prepare("SELECT * FROM Users WHERE UserName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return User{}, err
	}
	defer rec.Close()

	user := User{}
	r, err := rec.Query(username)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return User{}, err
	}
	defer r.Close()

	err = r.Scan(
		&user.Id,
		&user.UserName,
		&user.Status,
		&user.PasswordHash,
		&user.CreationDate,
		&user.LastChangedDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No user found with name: " + username)
			return User{}, nil
		}
		log.Println("ERROR: Failed to scan rows")
		return User{}, err
	}

	log.Println("INFO: User " + user.UserName + " found")
	return user, nil
}

func CreateUser(p ProposedUser) (bool, error) {
	log.Println("INFO: Creating user: " + p.UserName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to create user: " + p.UserName)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to create user: " + p.UserName)
			t.Rollback()
		}
	}()

	q, err := t.Prepare("INSERT INTO Users (UserName, PasswordHash) VALUES (?, ?)")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	// take password and hash it
	hash := sha512.Sum512([]byte(p.Password))
	passwdHash := hex.EncodeToString(hash[:])

	_, err = q.Exec(p.UserName, passwdHash)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: User " + p.UserName + " created successfully")
	return true, nil
}

func DeleteUser(username string) (bool, error) {
	log.Println("INFO: Deleting user: " + username)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to delete user: " + username)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to delete user: " + username)
			t.Rollback()
		}
	}()

	q, err := DB.Prepare("DELETE FROM Users WHERE UserName IS ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	_, err = q.Exec(username)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: User " + username + " deleted successfully")
	return true, nil
}

func GetUsers() ([]User, error) {
	log.Println("INFO: Getting all users")
	rows, err := DB.Query("SELECT * FROM Users")
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		err = rows.Scan(
			&user.Id,
			&user.UserName,
			&user.Status,
			&user.PasswordHash,
			&user.CreationDate,
			&user.LastChangedDate,
		)
		if err != nil {
			log.Println("ERROR: Failed to scan rows")
			return nil, err
		}
		users = append(users, user)
	}

	log.Println("INFO: Found " + strconv.Itoa(len(users)) + " users")
	return users, nil
}

func GetUserStatus(username string) (string, error) {
	log.Println("INFO: Getting user status for: " + username)
	q, err := DB.Prepare("SELECT Status FROM Users WHERE UserName IS ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return "", err
	}
	defer q.Close()

	status := ""
	r, err := q.Query(username)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return "", err
	}
	defer r.Close()

	err = r.Scan(
		&status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No user found with name: " + username)
			return "", nil
		}
		log.Println("ERROR: Failed to scan rows")
		return "", err
	}

	log.Println("INFO: User '" + username + "' status: " + status)
	return status, nil
}

func SetUserStatus(username string, j UserStatus) (bool, error) {
	log.Println("INFO: Setting user status for: " + username)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to set user status for: " + username)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to set user status for: " + username)
			t.Rollback()
		}
	}()

	q, err := DB.Prepare("UPDATE Users SET Status = ? WHERE UserName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}
	// ensure the UserStatus.Status value is either 'enabled' or 'locked'
	log.Println("INFO: user to set status of: " + username)
	log.Println("INFO: requested state to set user to: " + j.Status)
	if j.Status != "enabled" && j.Status != "locked" {
		log.Println("ERROR: Invalid status value: " + j.Status)
		return false, &InvalidStatusValue{Err: errors.New("invalid value: " + j.Status)}
	}

	result, err := q.Exec(j.Status, username)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}
	numberOfRows, err := result.RowsAffected()
	if err != nil {
		log.Println("ERROR: Failed to get number of rows affected")
		return false, err
	}
	log.Println("INFO: SQL result: Rows: " + strconv.Itoa(int(numberOfRows)))

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: User " + username + " status set to: " + j.Status)
	return true, nil
}
