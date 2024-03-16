package model

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
)

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
		&user.PasswordHash,
		&user.CreationDate,
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

	defer q.Close()

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
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserStatus(username string) (string, error) {
	status := ""
	return status, nil
}
