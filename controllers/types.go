package controllers

import "github.com/greeneg/ipmanager/globals"

type IpManager struct {
	AppPath    string
	ConfigPath string
	ConfStruct globals.Config
}

type SafeUser struct {
	Id           int
	UserName     string
	CreationDate string
}
