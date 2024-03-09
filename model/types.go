package model

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
	Id           int    `json:"Id"`
	HostName     string `json:"HostName"`
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
	Id           int    `json:"Id"`
	UserName     string `json:"UserName"`
	PasswordHash string `json:"PasswordHash"`
	CreationDate string `json:"CreationDate"`
}

type ProposedUser struct {
	Id           int    `json:"Id"`
	UserName     string `json:"UserName"`
	Password     string `json:"Password"`
	CreationDate string `json:"CreationDate"`
}
