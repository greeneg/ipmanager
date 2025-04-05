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
	"database/sql"
	"log"
	"strconv"
)

func GetAddressById(id int) (Address, error) {
	log.Println("INFO: Getting address by id: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM AssignedAddresses WHERE id = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement for GetAddressById")
		return Address{}, err
	}
	defer rec.Close()

	addr := Address{}
	r, err := rec.Query(id)
	if err != nil {
		log.Println("ERROR: Failed to query address by id")
		return Address{}, err
	}
	defer r.Close()

	err = r.Scan(
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

	log.Println("INFO: Address found by id: " + strconv.Itoa(addr.Id))
	return addr, nil
}

func GetHostIdByHostname(hostname string) (int, error) {
	var hostNameId int = 0
	log.Println("INFO: Getting host id by hostname: " + hostname)
	rec, err := DB.Prepare("SELECT Id FROM Hosts WHERE HostName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement for GetHostIdByHostname")
		return hostNameId, err
	}
	defer rec.Close()

	r, err := rec.Query(hostname)
	if err != nil {
		log.Println("ERROR: Failed to query host id by hostname")
		return hostNameId, err
	}
	defer r.Close()

	err = r.Scan(
		&hostNameId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No hostname found")
			return hostNameId, nil
		}
		log.Println("ERROR: Failed to scan hostname id")
		return hostNameId, err
	}

	log.Println("INFO: Host id found by hostname: " + strconv.Itoa(hostNameId))
	return hostNameId, nil
}

func GetAddressByHostName(hostname string) (Address, error) {
	log.Println("INFO: Getting address by hostname: " + hostname)
	hostNameId, err := GetHostIdByHostname(hostname)
	if err != nil {
		log.Println("ERROR: Failed to get host id by hostname")
		return Address{}, err
	}
	if hostNameId == 0 {
		log.Println("ERROR: No hostname found")
		return Address{}, nil
	}

	rec, err := DB.Prepare("SELECT * FROM AssignedAddresses WHERE HostNameId = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement for GetAddressByHostName")
		return Address{}, err
	}
	defer rec.Close()

	addr := Address{}

	r, err := rec.Query(hostNameId)
	if err != nil {
		log.Println("ERROR: Failed to query address by hostname")
		return Address{}, err
	}
	defer r.Close()

	err = r.Scan(
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
			log.Println("ERROR: No address found for hostname")
			return Address{}, nil
		}
		log.Println("ERROR: Failed to scan address by hostname")
		return Address{}, err
	}

	log.Println("INFO: Address found by hostname: " + strconv.Itoa(addr.Id))
	return addr, nil
}

func GetAddressByHostNameId(id int) (Address, error) {
	log.Println("INFO: Getting address by hostname id: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM AssignedAddresses WHERE HostNameId = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement for GetAddressByHostNameId")
		return Address{}, err
	}
	defer rec.Close()

	addr := Address{}

	r, err := rec.Query(id)
	if err != nil {
		log.Println("ERROR: Failed to query address by hostname id")
		return Address{}, err
	}
	defer r.Close()

	err = r.Scan(
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
			log.Println("ERROR: No address found for hostname id")
			return Address{}, nil
		}
		log.Println("ERROR: Failed to scan address by hostname id")
		return Address{}, err
	}

	log.Println("INFO: Address found by hostname id: " + strconv.Itoa(addr.Id))
	return addr, nil
}

func GetAddressByIpAddress(ip string) (Address, error) {
	log.Println("INFO: Getting address by ip address: " + ip)
	rec, err := DB.Prepare("SELECT * FROM AssignedAddresses WHERE Address = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement for GetAddressByIpAddress")
		return Address{}, err
	}
	defer rec.Close()

	addr := Address{}

	r, err := rec.Query(ip)
	if err != nil {
		log.Println("ERROR: Failed to query address by ip address")
		return Address{}, err
	}
	defer r.Close()

	err = r.Scan(
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
			log.Println("ERROR: No address found for ip address")
			return Address{}, nil
		}
		log.Println("ERROR: Failed to scan address by ip address")
		return Address{}, err
	}

	log.Println("INFO: Address found by ip address: " + strconv.Itoa(addr.Id))
	return addr, nil
}

func GetAddressesByDomainId(id int) ([]Address, error) {
	log.Println("INFO: Getting addresses by domain id: " + strconv.Itoa(id))
	rows, err := DB.Query("SELECT * FROM AssignedAddresses WHERE DomainId = ?", id)
	if err != nil {
		log.Println("ERROR: Failed to query addresses by domain id")
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
			log.Println("ERROR: Failed to scan address by domain id")
			return nil, err
		}
		addresses = append(addresses, address)
	}

	log.Println("INFO: Addresses found by domain id: " + strconv.Itoa(id))
	return addresses, nil
}

func GetDomainIdByDomainName(domainname string) (int, error) {
	log.Println("INFO: Getting domain id by domain name: " + domainname)
	rec, err := DB.Prepare("SELECT Id FROM Domains WHERE DomainName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement for GetDomainIdByDomainName")
		return 0, err
	}
	defer rec.Close()

	var id int
	err = rec.QueryRow(domainname).Scan(
		&id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No domain found")
			return 0, nil
		}
		log.Println("ERROR: Failed to scan domain id")
		return 0, err
	}

	log.Println("INFO: Domain id found by domain name: " + strconv.Itoa(id))
	return id, nil
}

func GetAddressesByDomainName(domainname string) ([]Address, error) {
	log.Println("INFO: Getting addresses by domain name: " + domainname)
	id, err := GetDomainIdByDomainName(domainname)
	if err != nil {
		log.Println("ERROR: Failed to get domain id by domain name")
		return nil, err
	}
	if id == 0 {
		log.Println("ERROR: No domain found")
		return nil, nil
	}

	rows, err := DB.Query("SELECT * FROM AssignedAddresses WHERE DomainId = ?", id)
	if err != nil {
		log.Println("ERROR: Failed to query addresses by domain id")
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
			log.Println("ERROR: Failed to scan address by domain id")
			return nil, err
		}
		addresses = append(addresses, address)
	}

	log.Println("INFO: Addresses found by domain name: " + domainname)
	return addresses, nil
}

func GetAddressesBySubnetId(id int) ([]Address, error) {
	log.Println("INFO: Getting addresses by subnet id: " + strconv.Itoa(id))
	rows, err := DB.Query("SELECT * FROM AssignedAddresses WHERE SubnetId = ?", id)
	if err != nil {
		log.Println("ERROR: Failed to query addresses by subnet id")
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
			log.Println("ERROR: Failed to scan address by subnet id")
			return nil, err
		}
		addresses = append(addresses, address)
	}

	log.Println("INFO: Addresses found by subnet id: " + strconv.Itoa(id))
	return addresses, nil
}

func GetSubnetIdBySubnetName(snetname string) (int, error) {
	log.Println("INFO: Getting subnet id by subnet name: " + snetname)
	rec, err := DB.Prepare("SELECT Id FROM Subnets WHERE NetworkName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement for GetSubnetIdBySubnetName")
		return 0, err
	}
	defer rec.Close()

	var id int
	err = rec.QueryRow(snetname).Scan(
		&id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No subnet found")
			return 0, nil
		}
		log.Println("ERROR: Failed to scan subnet id")
		return 0, err
	}

	log.Println("INFO: Subnet id found by subnet name: " + strconv.Itoa(id))
	return id, nil
}

func GetAddressesBySubnetName(snetname string) ([]Address, error) {
	log.Println("INFO: Getting addresses by subnet name: " + snetname)
	id, err := GetSubnetIdBySubnetName(snetname)
	if err != nil {
		log.Println("ERROR: Failed to get subnet id by subnet name")
		return nil, err
	}
	if id == 0 {
		log.Println("ERROR: No subnet found")
		return nil, nil
	}

	rows, err := DB.Query("SELECT * FROM AssignedAddresses WHERE SubnetId = ?", id)
	if err != nil {
		log.Println("ERROR: Failed to query addresses by subnet id")
		return nil, err
	}

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
			log.Println("ERROR: Failed to scan address by subnet id")
			return nil, err
		}
		addresses = append(addresses, address)
	}

	log.Println("INFO: Addresses found by subnet name: " + snetname)
	return addresses, nil
}

func GetAddresses() ([]Address, error) {
	log.Println("INFO: Getting all addresses")
	rows, err := DB.Query("SELECT * FROM AssignedAddresses")
	if err != nil {
		log.Println("ERROR: Failed to query all addresses")
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
			log.Println("ERROR: Failed to scan address")
			return nil, err
		}
		addresses = append(addresses, address)
	}

	log.Println("INFO: All addresses found")
	return addresses, nil
}
