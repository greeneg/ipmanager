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

type Address struct {
	Id           int    `json:"Id"`
	Address      string `json:"Address"`
	HostNameId   int    `json:"HostNameId"`
	DomainId     int    `json:"DomainId"`
	SubnetId     int    `json:"SubnetId"`
	CreatorId    int    `json:"CreatorId"`
	CreationDate string `json:"CreationDate"`
}

type Domain struct {
	Id           int    `json:"Id"`
	DomainName   string `json:"DomainName"`
	CreatorId    int    `json:"CreatorId"`
	CreationDate string `json:"CreationDate"`
}

type Host struct {
	Id           int      `json:"Id"`
	HostName     string   `json:"HostName"`
	MacAddresses []string `json:"MacAddresses"`
	CreatorId    int      `json:"CreatorId"`
	CreationDate string   `json:"CreationDate"`
}

type StringHost struct {
	Id           int    `json:"Id"`
	HostName     string `json:"HostName"`
	MacAddresses string `json:"MacAddresses"`
	CreatorId    int    `json:"CreatorId"`
	CreationDate string `json:"CreationDate"`
}

type Subnet struct {
	Id             int    `json:"Id"`
	NetworkName    string `json:"NetworkName"`
	NetworkPrefix  string `json:"NetworkPrefix"`
	BitMask        int    `json:"BitMask"`
	GatewayAddress string `json:"GatewayAddress"`
	DomainId       int    `json:"DomainId"`
	CreatorId      int    `json:"CreatorId"`
	CreationDate   string `json:"CreationDate"`
}

type User struct {
	Id              int    `json:"Id"`
	UserName        string `json:"UserName"`
	Status          string `json:"Status"`
	PasswordHash    string `json:"PasswordHash"`
	CreationDate    string `json:"CreationDate"`
	LastChangedDate string `json:"LastChangedDate"`
}

type SubnetUpdate struct {
	NetworkPrefix  string `json:"NetworkPrefix"`
	BitMask        int    `json:"BitMask"`
	GatewayAddress string `json:"GatewayAddress"`
	DomainName     string `json:"DomainName"`
}

type ProposedUser struct {
	Id           int    `json:"Id"`
	UserName     string `json:"UserName"`
	Status       string `json:"Status"`
	Password     string `json:"Password"`
	CreationDate string `json:"CreationDate"`
}

type UserStatus struct {
	Status string `json:"status"`
}

type FailureMsg struct {
	Error string `json:"error"`
}

type SuccessMsg struct {
	Message string `json:"message"`
}

type UserStatusMsg struct {
	Message    string `json:"message"`
	UserStatus string `json:"userStatus"`
}

type DomainList struct {
	Data []Domain `json:"data"`
}

type HostList struct {
	Data []Host `json:"data"`
}

type MacAddressList struct {
	Data []string `json:"data"`
}

type PasswordChange struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type Subnets struct {
	Data []Subnet `json:"data"`
}

type UsersList struct {
	Data []User `json:"data"`
}
