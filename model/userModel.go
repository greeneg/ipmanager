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
	q, err := DB.Prepare("SELECT PasswordHash FROM Users WHERE UserName = ?")
	if err != nil {
		return "", err
	}

	passwordHash := ""
	err = q.QueryRow(username).Scan(
		&passwordHash,
	)
	if err != nil {
		return "", err
	}

	return passwordHash, nil
}

func storeNewPassword(hashedPassword string, username string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	// now we need to create a new transaction to SET the password hash into the DB
	q, err := DB.Prepare("UPDATE Users SET PasswordHash = ?, LastChangedDate = ? WHERE UserName = ?")
	if err != nil {
		return false, err
	}

	// get time stamp
	tStamp := time.Now().Format("2006-01-02 15:04:05") // force into SQL DateTime format

	_, err = q.Exec(hashedPassword, tStamp, username)
	if err != nil {
		return false, err
	}

	t.Commit()

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
	rec, err := DB.Prepare("SELECT * FROM Users WHERE Id = ?")
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = rec.QueryRow(id).Scan(
		&user.Id,
		&user.UserName,
		&user.Status,
		&user.PasswordHash,
		&user.CreationDate,
		&user.LastChangedDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}

	return user, nil
}

func GetUserByUserName(username string) (User, error) {
	rec, err := DB.Prepare("SELECT * FROM Users WHERE UserName = ?")
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = rec.QueryRow(username).Scan(
		&user.Id,
		&user.UserName,
		&user.Status,
		&user.PasswordHash,
		&user.CreationDate,
		&user.LastChangedDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}

	return user, nil
}

func CreateUser(p ProposedUser) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := t.Prepare("INSERT INTO Users (UserName, PasswordHash) VALUES (?, ?)")
	if err != nil {
		return false, err
	}

	// take password and hash it
	hash := sha512.Sum512([]byte(p.Password))
	passwdHash := hex.EncodeToString(hash[:])

	_, err = q.Exec(p.UserName, passwdHash)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}

func DeleteUser(username string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM Users WHERE UserName IS ?")
	if err != nil {
		return false, err
	}

	_, err = q.Exec(username)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}

func GetUsers() ([]User, error) {
	rows, err := DB.Query("SELECT * FROM Users")
	if err != nil {
		return nil, err
	}

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
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserStatus(username string) (string, error) {
	t, err := DB.Begin()
	if err != nil {
		return "", err
	}

	q, err := DB.Prepare("SELECT Status FROM Users WHERE UserName IS ?")
	if err != nil {
		return "", err
	}

	status := ""
	err = q.QueryRow(username).Scan(
		&status,
	)
	if err != nil {
		return "", err
	}

	t.Commit()

	log.Println("INFO: User '" + username + "' status: " + status)

	return status, nil
}

func SetUserStatus(username string, j UserStatus) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := DB.Prepare("UPDATE Users SET Status = ? WHERE UserName = ?")
	if err != nil {
		return false, err
	}
	// ensure the UserStatus.Status value is either 'enabled' or 'locked'
	log.Println("INFO: user to set status of: " + username)
	log.Println("INFO: requested state to set user to: " + j.Status)
	if j.Status != "enabled" && j.Status != "locked" {
		return false, &InvalidStatusValue{Err: errors.New("invalid value: " + j.Status)}
	}

	result, err := q.Exec(j.Status, username)
	if err != nil {
		return false, err
	}
	numberOfRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	log.Println("INFO: SQL result: Rows: " + strconv.Itoa(int(numberOfRows)))

	t.Commit()

	return true, nil
}
