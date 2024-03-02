package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

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
	rec, err := DB.Prepare("SELECT * FROM Domains WHERE Id = ?")
	if err != nil {
		return Domain{}, err
	}

	domain := Domain{}
	err = rec.QueryRow(id).Scan(
		&domain.Id,
		&domain.DomainName,
		&domain.CreatorId,
		&domain.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Domain{}, nil
		}
		return Domain{}, err
	}

	return domain, nil
}

func GetDomainByDomainName(domainname string) (Domain, error) {
	rec, err := DB.Prepare("SELECT * FROM Domains WHERE DomainName = ?")
	if err != nil {
		return Domain{}, err
	}

	domain := Domain{}
	err = rec.QueryRow(domainname).Scan(
		&domain.Id,
		&domain.DomainName,
		&domain.CreatorId,
		&domain.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Domain{}, nil
		}
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
	rec, err := DB.Prepare("SELECT * FROM Hosts WHERE Id = ?")
	if err != nil {
		return Host{}, err
	}

	host := Host{}
	err = rec.QueryRow(id).Scan(
		&host.Id,
		&host.HostName,
		&host.CreatorId,
		&host.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Host{}, nil
		}
		return Host{}, err
	}

	return host, nil
}

func GetHostByHostName(hostname string) (Host, error) {
	rec, err := DB.Prepare("SELECT * FROM Hosts WHERE HostName = ?")
	if err != nil {
		return Host{}, err
	}

	host := Host{}
	err = rec.QueryRow(hostname).Scan(
		&host.Id,
		&host.HostName,
		&host.CreatorId,
		&host.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Host{}, nil
		}
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
	rec, err := DB.Prepare("SELECT * FROM Subnets WHERE Id = ?")
	if err != nil {
		return Subnet{}, err
	}

	subnet := Subnet{}
	err = rec.QueryRow(id).Scan(
		&subnet.Id,
		&subnet.NetworkName,
		&subnet.NetworkPrefix,
		&subnet.BitMask,
		&subnet.GatewayAddress,
		&subnet.DomainId,
		&subnet.CreatorId,
		&subnet.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Subnet{}, err
		}
		return Subnet{}, err
	}

	return subnet, nil
}

func GetSubnetByNetworkName(snetname string) (Subnet, error) {
	rec, err := DB.Prepare("SELECT * FROM Subnets WHERE NetworkName = ?")
	if err != nil {
		return Subnet{}, err
	}

	subnet := Subnet{}
	err = rec.QueryRow(snetname).Scan(
		&subnet.Id,
		&subnet.NetworkName,
		&subnet.NetworkPrefix,
		&subnet.BitMask,
		&subnet.GatewayAddress,
		&subnet.DomainId,
		&subnet.CreatorId,
		&subnet.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Subnet{}, err
		}
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
			&snet.NetworkPrefix,
			&snet.BitMask,
			&snet.GatewayAddress,
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
			&snet.NetworkPrefix,
			&snet.BitMask,
			&snet.GatewayAddress,
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
			&snet.NetworkPrefix,
			&snet.BitMask,
			&snet.GatewayAddress,
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
	rec, err := DB.Prepare("SELECT * FROM Users WHERE Id = ?")
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = rec.QueryRow(id).Scan(
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
