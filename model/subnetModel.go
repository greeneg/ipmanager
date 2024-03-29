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
	"fmt"
	"log"
	"strconv"

	"github.com/seancfoley/ipaddress-go/ipaddr"
)

func createDynamicNetworkTable(networkName string, networkPrefix string, bitmask int) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	createStatement := `CREATE TABLE IF NOT EXISTS ` + networkName + `(
		Id	            INTEGER PRIMARY KEY AUTOINCREMENT
						UNIQUE	NOT NULL,
		IpAddress		STRING	UNIQUE	NOT NULL,
		AssignmentState	BOOL	NOT NULL	DEFAULT (0)
	)`

	q, err := t.Prepare(createStatement)
	if err != nil {
		return false, err
	}
	_, err = q.Exec()
	if err != nil {
		return false, err
	}

	t.Commit()

	log.Println("INFO: Populating addresses for dynamic table '" + networkName + "'")
	_, err = populateAddresses(networkName, networkPrefix, bitmask)
	if err != nil {
		// TODO: revert adding table if error occurs when populating it
		return false, err
	}

	return true, nil
}

func populateAddresses(networkName string, networkPrefix string, bitmask int) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	subnet := ipaddr.NewIPAddressString(networkPrefix + "/" + strconv.Itoa(bitmask)).GetAddress().WithoutPrefixLen()
	netAddr := subnet.GetNetIP()
	bcastAddr := subnet.GetUpper()
	iterator := subnet.Iterator()
	for next := iterator.Next(); next != nil; next = iterator.Next() {
		q, err := t.Prepare("INSERT INTO " + networkName + " (IpAddress) VALUES (?)")
		if err != nil {
			return false, err
		}
		address := fmt.Sprintf("%s", next)
		if address == netAddr.String() {
			continue
		}
		if address == bcastAddr.String() {
			continue
		}
		_, err = q.Exec(address)
		if err != nil {
			return false, err
		}
	}

	t.Commit()

	return true, nil
}

func CreateSubnet(s Subnet, id int) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := t.Prepare("INSERT INTO Subnets (NetworkName, NetworkPrefix, BitMask, GatewayAddress, DomainId, CreatorId) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return false, err
	}
	_, err = q.Exec(s.NetworkName, s.NetworkPrefix, s.BitMask, s.GatewayAddress, s.DomainId, id)
	if err != nil {
		return false, err
	}

	t.Commit()

	// now create the dynamic table for the ip range
	log.Printf("INFO: Creating dynamic table for '" + s.NetworkName + "' addresses")
	_, err = createDynamicNetworkTable(s.NetworkName, s.NetworkPrefix, s.BitMask)
	if err != nil {
		// TODO: revert transaction if subnet dynamic table cannot be built
		return false, err
	}

	return true, nil
}

func dropDynamicNetworkTable(subnetName string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := DB.Prepare("DROP TABLE IF EXISTS " + subnetName)
	if err != nil {
		return false, err
	}

	_, err = q.Exec()
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}

func DeleteSubnet(subnetName string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM Subnets WHERE NetworkName IS ?")
	if err != nil {
		return false, err
	}

	_, err = q.Exec(subnetName)
	if err != nil {
		return false, err
	}

	t.Commit()

	// now drop the network's table
	log.Println("INFO: Dropping dynamic table '" + subnetName + "'")
	_, err = dropDynamicNetworkTable(subnetName)
	if err != nil {
		// TODO: revert transaction some how if drop fails
		return false, err
	}

	return true, nil
}

func checkAddressTableInUse(subnetName string) error {
	t, err := DB.Begin()
	if err != nil {
		return err
	}

	q, err := t.Prepare("SELECT COUNT(*) FROM " + subnetName + " WHERE AssignmentState = 1")
	if err != nil {
		return err
	}
	var val string
	err = q.QueryRow().Scan(&val)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return err
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if num != 0 {
		a := new(AddressTableInUse)
		return a
	}

	t.Commit()

	return nil
}

func ModifySubnet(subnetName string, json SubnetUpdate) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	// check if the address table for the network has any entries that are in use. if so, error out with
	// network in use error
	err = checkAddressTableInUse(subnetName)
	if err != nil { // first deal with any errors
		return false, err
	}

	q, err := t.Prepare("UPDATE Subnets SET NetworkPrefix =?, BitMask = ?, GatewayAddress = ?, DomainId =? WHERE NetworkName = ?")
	if err != nil {
		return false, err
	}

	// get our needed values from the JSON sent in
	s := new(SubnetUpdate)
	s.NetworkPrefix = json.NetworkPrefix
	s.BitMask = json.BitMask
	s.GatewayAddress = json.GatewayAddress
	s.DomainName = json.DomainName

	// get the DomainId from the DomainName
	d, err := GetDomainByDomainName(s.DomainName)
	if err != nil {
		return false, err
	}

	_, err = q.Exec(s.NetworkPrefix, s.BitMask, s.GatewayAddress, d.Id, subnetName)
	if err != nil {
		return false, err
	}

	// Now that we've updated the prefix, we need to drop the address table and recreate it with the new
	// settings
	// first, drop old table
	_, err = dropDynamicNetworkTable(subnetName)
	if err != nil {
		return false, err
	}
	// now create the new dynamic table
	_, err = createDynamicNetworkTable(subnetName, s.NetworkPrefix, s.BitMask)
	if err != nil {
		// TODO: revert transaction (somehow) if subnet dynamic table cannot be built
		return false, err
	}

	t.Commit()

	return true, nil
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
