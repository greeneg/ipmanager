package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

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
	Id                         int    `json:"Id"`
	NetworkName                string `json:"NetworkName"`
	AddressStart               string `json:"AddressStart"`
	AddressEnd                 string `json:"AddressEnd"`
	BitMask                    int    `json:"BitMask"`
	GatewayAddress             string `json:"GatewayAddress"`
	NumberOfAvailableAddresses int    `json:"NumberOfAvailableAddresses"`
	DomainId                   int    `json:"DomainId"`
	CreatorId                  int    `json:"CreatorId"`
	CreationDate               string `json:"CreationDate"`
}

type User struct {
	Id           int    `json:"Id"`
	UserName     string `json:"UserName"`
	PasswordHash string `json:"PasswordHash"`
	CreationDate string `json:"CreationDate"`
}

func ConnectDatabase(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func GetAddressById(id int) (Address, error) {
	rec, err := DB.Prepare("SELECT * FROM AssignedAddresses WHERE id = ?")
	if err != nil {
		return Address{}, err
	}

	addr := Address{}
	err = rec.QueryRow(id).Scan(
		&addr.Id,
		&addr.Address,
		&addr.HostNameId,
		&addr.DomainId,
		&addr.SubnetId,
		&addr.CreatorId,
		&addr.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Address{}, nil
		}
		return Address{}, err
	}

	return addr, nil
}

func GetAddressByHostName(hostname string) (Address, error) {
	rec, err := DB.Prepare("SELECT Id FROM Hosts WHERE HostName = ?")
	if err != nil {
		return Address{}, err
	}
	var hostNameId int
	err = rec.QueryRow(hostname).Scan(
		&hostNameId,
	)
	if err != nil {
		return Address{}, err
	}

	rec, err = DB.Prepare("SELECT * FROM AssignedAddresses WHERE HostNameId = ?")
	if err != nil {
		return Address{}, err
	}

	addr := Address{}
	qErr := rec.QueryRow(hostNameId).Scan(
		&addr.Id,
		&addr.Address,
		&addr.HostNameId,
		&addr.DomainId,
		&addr.SubnetId,
		&addr.CreatorId,
		&addr.CreationDate,
	)
	if qErr != nil {
		if qErr == sql.ErrNoRows {
			return Address{}, nil
		}
		return Address{}, qErr
	}

	return addr, nil
}

func GetAddressByHostNameId(id int) (Address, error) {
	rec, err := DB.Prepare("SELECT * FROM AssignedAddresses WHERE HostNameId = ?")
	if err != nil {
		return Address{}, err
	}

	addr := Address{}
	qErr := rec.QueryRow(id).Scan(
		&addr.Id,
		&addr.Address,
		&addr.HostNameId,
		&addr.DomainId,
		&addr.SubnetId,
		&addr.CreatorId,
		&addr.CreationDate,
	)
	if qErr != nil {
		if qErr == sql.ErrNoRows {
			return Address{}, nil
		}
		return Address{}, qErr
	}

	return addr, nil
}

func GetAddressByIpAddress(ip string) (Address, error) {
	rec, err := DB.Prepare("SELECT * FROM AssignedAddresses WHERE Address = ?")
	if err != nil {
		return Address{}, err
	}

	addr := Address{}
	qErr := rec.QueryRow(ip).Scan(
		&addr.Id,
		&addr.Address,
		&addr.HostNameId,
		&addr.DomainId,
		&addr.SubnetId,
		&addr.CreatorId,
		&addr.CreationDate,
	)
	if qErr != nil {
		if qErr == sql.ErrNoRows {
			return Address{}, nil
		}
		return Address{}, qErr
	}

	return addr, nil
}

func GetAddressesByDomainId(id int) ([]Address, error) {
	rows, err := DB.Query("SELECT * FROM AssignedAddresses WHERE DomainId = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := make([]Address, 0)
	for rows.Next() {
		address := Address{}
		err = rows.Scan(
			&address.Id,
			&address.Address,
			&address.HostNameId,
			&address.DomainId,
			&address.SubnetId,
			&address.CreatorId,
			&address.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func GetAddressesByDomainName(domainname string) ([]Address, error) {
	rec, err := DB.Prepare("SELECT Id FROM Domains WHERE DomainName = ?")
	if err != nil {
		return nil, err
	}
	var id int
	err = rec.QueryRow(domainname).Scan(
		&id,
	)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query("SELECT * FROM AssignedAddresses WHERE DomainId = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := make([]Address, 0)
	for rows.Next() {
		address := Address{}
		err = rows.Scan(
			&address.Id,
			&address.Address,
			&address.HostNameId,
			&address.DomainId,
			&address.SubnetId,
			&address.CreatorId,
			&address.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func GetAddressesBySubnetId(id int) ([]Address, error) {
	rows, err := DB.Query("SELECT * FROM AssignedAddresses WHERE SubnetId = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := make([]Address, 0)
	for rows.Next() {
		address := Address{}
		err = rows.Scan(
			&address.Id,
			&address.Address,
			&address.HostNameId,
			&address.DomainId,
			&address.SubnetId,
			&address.CreatorId,
			&address.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func GetAddressesBySubnetName(snetname string) ([]Address, error) {
	rec, err := DB.Prepare("SELECT Id FROM Subnets WHERE NetworkName = ?")
	if err != nil {
		return nil, err
	}
	var id int
	err = rec.QueryRow(snetname).Scan(
		&id,
	)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query("SELECT * FROM AssignedAddresses WHERE SubnetId = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := make([]Address, 0)
	for rows.Next() {
		address := Address{}
		err = rows.Scan(
			&address.Id,
			&address.Address,
			&address.HostNameId,
			&address.DomainId,
			&address.SubnetId,
			&address.CreatorId,
			&address.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func GetAddresses() ([]Address, error) {
	rows, err := DB.Query("SELECT * FROM AssignedAddresses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := make([]Address, 0)
	for rows.Next() {
		address := Address{}
		err = rows.Scan(
			&address.Id,
			&address.Address,
			&address.SubnetId,
			&address.HostNameId,
			&address.DomainId,
			&address.CreatorId,
			&address.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func GetDomainById(id int) (Domain, error) {
	rows, err := DB.Query("SELECT * FROM Domains WHERE Id = ?", id)
	if err != nil {
		return Domain{}, err
	}
	defer rows.Close()

	domain := Domain{}
	err = rows.Scan(
		&domain.Id,
		&domain.DomainName,
		&domain.CreatorId,
		&domain.CreationDate,
	)
	if err != nil {
		return Domain{}, err
	}

	return domain, nil
}

func GetDomainByDomainName(domainname string) (Domain, error) {
	rows, err := DB.Query("SELECT * FROM Domains WHERE DomainName = ?", domainname)
	if err != nil {
		return Domain{}, err
	}
	defer rows.Close()

	domain := Domain{}
	err = rows.Scan(
		&domain.Id,
		&domain.DomainName,
		&domain.CreatorId,
		&domain.CreationDate,
	)
	if err != nil {
		return Domain{}, err
	}

	return domain, nil
}

func GetDomains() ([]Domain, error) {
	rows, err := DB.Query("SELECT * FROM Domains")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	domains := make([]Domain, 0)
	for rows.Next() {
		domain := Domain{}
		err = rows.Scan(
			&domain.Id,
			&domain.DomainName,
			&domain.CreatorId,
			&domain.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	return domains, nil
}

func GetHostById(id int) (Host, error) {
	rows, err := DB.Query("SELECT * FROM Hosts WHERE Id = ?", id)
	if err != nil {
		return Host{}, err
	}
	defer rows.Close()

	host := Host{}
	err = rows.Scan(
		&host.Id,
		&host.HostName,
		&host.CreatorId,
		&host.CreationDate,
	)
	if err != nil {
		return Host{}, err
	}

	return host, nil
}

func GetHostByHostName(hostname string) (Host, error) {
	rows, err := DB.Query("SELECT * FROM Hosts WHERE HostName = ?", hostname)
	if err != nil {
		return Host{}, err
	}
	defer rows.Close()

	host := Host{}
	err = rows.Scan(
		&host.Id,
		&host.HostName,
		&host.CreatorId,
		&host.CreationDate,
	)
	if err != nil {
		return Host{}, err
	}

	return host, nil
}

func GetHosts() ([]Host, error) {
	rows, err := DB.Query("SELECT * FROM Hosts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hosts := make([]Host, 0)
	for rows.Next() {
		host := Host{}
		err = rows.Scan(
			&host.Id,
			&host.HostName,
			&host.CreatorId,
			&host.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func GetSubnetById(id int) (Subnet, error) {
	rows, err := DB.Query("SELECT * FROM Subnets WHERE Id = ?", id)
	if err != nil {
		return Subnet{}, err
	}
	defer rows.Close()

	subnet := Subnet{}
	err = rows.Scan(
		&subnet.Id,
		&subnet.NetworkName,
		&subnet.AddressStart,
		&subnet.AddressEnd,
		&subnet.BitMask,
		&subnet.GatewayAddress,
		&subnet.NumberOfAvailableAddresses,
		&subnet.DomainId,
		&subnet.CreatorId,
		&subnet.CreationDate,
	)
	if err != nil {
		return Subnet{}, err
	}

	return subnet, nil
}

func GetSubnetByNetworkName(snetname string) (Subnet, error) {
	rows, err := DB.Query("SELECT * FROM Subnets WHERE NetworkName = ?", snetname)
	if err != nil {
		return Subnet{}, err
	}
	defer rows.Close()

	subnet := Subnet{}
	err = rows.Scan(
		&subnet.Id,
		&subnet.NetworkName,
		&subnet.AddressStart,
		&subnet.AddressEnd,
		&subnet.BitMask,
		&subnet.GatewayAddress,
		&subnet.NumberOfAvailableAddresses,
		&subnet.DomainId,
		&subnet.CreatorId,
		&subnet.CreationDate,
	)
	if err != nil {
		return Subnet{}, err
	}

	return subnet, nil
}

func GetSubnestByDomainId(id int) ([]Subnet, error) {
	rows, err := DB.Query("SELECT * FROM Subnets WHERE DomainId = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subnets := make([]Subnet, 0)
	for rows.Next() {
		snet := Subnet{}
		err = rows.Scan(
			&snet.Id,
			&snet.NetworkName,
			&snet.AddressStart,
			&snet.AddressEnd,
			&snet.BitMask,
			&snet.GatewayAddress,
			&snet.NumberOfAvailableAddresses,
			&snet.DomainId,
			&snet.CreatorId,
			&snet.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, snet)
	}

	return subnets, nil
}

func GetSubnestByDomainName(domainname string) ([]Subnet, error) {
	rec, err := DB.Prepare("SELECT Id FROM Domains WHERE DomainName = ?")
	if err != nil {
		return nil, err
	}
	var id int
	err = rec.QueryRow(domainname).Scan(
		&id,
	)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query("SELECT * FROM Subnets WHERE DomainId = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subnets := make([]Subnet, 0)
	for rows.Next() {
		snet := Subnet{}
		err = rows.Scan(
			&snet.Id,
			&snet.NetworkName,
			&snet.AddressStart,
			&snet.AddressEnd,
			&snet.BitMask,
			&snet.GatewayAddress,
			&snet.NumberOfAvailableAddresses,
			&snet.DomainId,
			&snet.CreatorId,
			&snet.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, snet)
	}

	return subnets, nil
}

func GetSubnets() ([]Subnet, error) {
	rows, err := DB.Query("SELECT * FROM Subnets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subnets := make([]Subnet, 0)
	for rows.Next() {
		snet := Subnet{}
		err = rows.Scan(
			&snet.Id,
			&snet.NetworkName,
			&snet.AddressStart,
			&snet.AddressEnd,
			&snet.BitMask,
			&snet.GatewayAddress,
			&snet.NumberOfAvailableAddresses,
			&snet.DomainId,
			&snet.CreatorId,
			&snet.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, snet)
	}

	return subnets, nil
}

func GetUserById(id int) (User, error) {
	rows, err := DB.Query("SELECT * FROM Users WHERE Id = ?", id)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	user := User{}
	err = rows.Scan(
		&user.Id,
		&user.UserName,
		&user.PasswordHash,
		&user.CreationDate,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserByUserName(username string) (User, error) {
	rows, err := DB.Query("SELECT * FROM Users WHERE UserName = ?", username)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	user := User{}
	err = rows.Scan(
		&user.Id,
		&user.UserName,
		&user.PasswordHash,
		&user.CreationDate,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
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
