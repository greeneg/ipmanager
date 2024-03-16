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

import "database/sql"

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
