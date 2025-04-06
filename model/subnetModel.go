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
	log.Println("INFO: Creating dynamic table for '" + networkName + "' addresses")
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to create dynamic table for '" + networkName + "' addresses")
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to create dynamic table for '" + networkName + "' addresses")
			t.Rollback()
		}
	}()

	createStatement := `CREATE TABLE IF NOT EXISTS ` + networkName + `(
		Id	            INTEGER PRIMARY KEY AUTOINCREMENT
						UNIQUE	NOT NULL,
		IpAddress		STRING	UNIQUE	NOT NULL,
		AssignmentState	BOOL	NOT NULL	DEFAULT (0)
	)`

	q, err := t.Prepare(createStatement)
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	_, err = q.Exec()
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: Populating addresses for dynamic table '" + networkName + "'")
	_, err = populateAddresses(networkName, networkPrefix, bitmask)
	if err != nil {
		log.Println("ERROR: Failed to populate addresses for dynamic table '" + networkName + "'")
		return false, err
	}

	log.Println("INFO: Dynamic table for '" + networkName + "' addresses created successfully")
	return true, nil
}

func populateAddresses(networkName string, networkPrefix string, bitmask int) (bool, error) {
	log.Println("INFO: Populating addresses for dynamic table '" + networkName + "'")
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to populate addresses for dynamic table '" + networkName + "'")
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to populate addresses for dynamic table '" + networkName + "'")
			t.Rollback()
		}
	}()

	subnet := ipaddr.NewIPAddressString(networkPrefix + "/" + strconv.Itoa(bitmask)).GetAddress().WithoutPrefixLen()
	netAddr := subnet.GetNetIP()
	bcastAddr := subnet.GetUpper()
	iterator := subnet.Iterator()
	for next := iterator.Next(); next != nil; next = iterator.Next() {
		q, err := t.Prepare("INSERT INTO " + networkName + " (IpAddress) VALUES (?)")
		if err != nil {
			log.Println("ERROR: Failed to prepare statement")
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
			log.Println("ERROR: Failed to execute statement")
			return false, err
		}
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: Addresses for dynamic table '" + networkName + "' populated successfully")
	return true, nil
}

func CreateSubnet(s Subnet, id int) (bool, error) {
	log.Println("INFO: Creating subnet " + s.NetworkName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to create subnet " + s.NetworkName)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to create subnet " + s.NetworkName)
			t.Rollback()
		}
	}()

	q, err := t.Prepare("INSERT INTO Subnets (NetworkName, NetworkPrefix, BitMask, GatewayAddress, DomainId, CreatorId) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	_, err = q.Exec(s.NetworkName, s.NetworkPrefix, s.BitMask, s.GatewayAddress, s.DomainId, id)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	// now create the dynamic table for the ip range
	log.Printf("INFO: Creating dynamic table for '" + s.NetworkName + "' addresses")
	_, err = createDynamicNetworkTable(s.NetworkName, s.NetworkPrefix, s.BitMask)
	if err != nil {
		log.Println("ERROR: Failed to create dynamic table for '" + s.NetworkName + "' addresses")
		return false, err
	}

	log.Println("INFO: Subnet " + s.NetworkName + " created successfully")
	return true, nil
}

func dropDynamicNetworkTable(subnetName string) (bool, error) {
	log.Println("INFO: Dropping dynamic table for '" + subnetName + "' addresses")
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to drop dynamic table for '" + subnetName + "' addresses")
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to drop dynamic table for '" + subnetName + "' addresses")
			t.Rollback()
		}
	}()

	q, err := DB.Prepare("DROP TABLE IF EXISTS " + subnetName)
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	_, err = q.Exec()
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: Dynamic table for '" + subnetName + "' addresses dropped successfully")
	return true, nil
}

func DeleteSubnet(subnetName string) (bool, error) {
	log.Println("INFO: Deleting subnet " + subnetName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to delete subnet " + subnetName)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to delete subnet " + subnetName)
			t.Rollback()
		}
	}()

	q, err := DB.Prepare("DELETE FROM Subnets WHERE NetworkName IS ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	_, err = q.Exec(subnetName)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	// now drop the network's table
	log.Println("INFO: Dropping dynamic table '" + subnetName + "'")
	_, err = dropDynamicNetworkTable(subnetName)
	if err != nil {
		log.Println("ERROR: Failed to drop dynamic table '" + subnetName + "'")
		return false, err
	}

	log.Println("INFO: Subnet " + subnetName + " deleted successfully")
	return true, nil
}

func checkAddressTableInUse(subnetName string) error {
	log.Println("INFO: Checking if address table '" + subnetName + "' is in use")
	q, err := DB.Prepare("SELECT COUNT(*) FROM " + subnetName + " WHERE AssignmentState = 1")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return err
	}
	defer q.Close()

	var val string
	r, err := q.Query()
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return err
	}
	defer r.Close()

	err = r.Scan(&val)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No rows found")
			return err
		}
		log.Println("ERROR: Failed to scan result")
		return err
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		log.Println("ERROR: Failed to convert string to int")
		return err
	}
	if num != 0 {
		a := new(AddressTableInUse)
		return a
	}

	log.Println("INFO: Address table '" + subnetName + "' is not in use")
	return nil
}

func ModifySubnet(subnetName string, json SubnetUpdate) (bool, error) {
	log.Println("INFO: Modifying subnet " + subnetName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to modify subnet " + subnetName)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to modify subnet " + subnetName)
			t.Rollback()
		}
	}()

	// check if the address table for the network has any entries that are in use. if so, error out with
	// network in use error
	err = checkAddressTableInUse(subnetName)
	if err != nil { // first deal with any errors
		if _, ok := err.(*AddressTableInUse); ok {
			log.Println("ERROR: Address table '" + subnetName + "' is in use")
			return false, err
		}
		log.Println("ERROR: Failed to check address table '" + subnetName + "'")
		return false, err
	}

	q, err := t.Prepare("UPDATE Subnets SET NetworkPrefix =?, BitMask = ?, GatewayAddress = ?, DomainId =? WHERE NetworkName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
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
		log.Println("ERROR: Failed to get DomainId from DomainName")
		return false, err
	}

	_, err = q.Exec(s.NetworkPrefix, s.BitMask, s.GatewayAddress, d.Id, subnetName)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	// Now that we've updated the prefix, we need to drop the address table and recreate it with the new
	// settings
	// first, drop old table
	_, err = dropDynamicNetworkTable(subnetName)
	if err != nil {
		log.Println("ERROR: Failed to drop dynamic table '" + subnetName + "'")
		return false, err
	}
	// now create the new dynamic table
	_, err = createDynamicNetworkTable(subnetName, s.NetworkPrefix, s.BitMask)
	if err != nil {
		log.Println("ERROR: Failed to create dynamic table for '" + subnetName + "' addresses")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: Subnet " + subnetName + " modified successfully")
	return true, nil
}

func GetSubnetById(id int) (Subnet, error) {
	log.Println("INFO: Getting subnet by id " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM Subnets WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return Subnet{}, err
	}
	defer rec.Close()

	subnet := Subnet{}
	r, err := rec.Query(id)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return Subnet{}, err
	}
	defer r.Close()

	err = r.Scan(
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
			log.Println("ERROR: No rows found")
			return Subnet{}, err
		}
		log.Println("ERROR: Failed to scan rows")
		return Subnet{}, err
	}

	log.Println("INFO: Subnet " + subnet.NetworkName + " found")
	return subnet, nil
}

func GetSubnetByNetworkName(snetname string) (Subnet, error) {
	log.Println("INFO: Getting subnet by name " + snetname)
	rec, err := DB.Prepare("SELECT * FROM Subnets WHERE NetworkName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return Subnet{}, err
	}
	defer rec.Close()

	subnet := Subnet{}

	r, err := rec.Query(snetname)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No rows found")
			return Subnet{}, err
		}
		log.Println("ERROR: Failed to execute statement")
		return Subnet{}, err
	}
	defer r.Close()

	err = r.Scan(
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
			log.Println("ERROR: No rows found")
			return Subnet{}, err
		}
		log.Println("ERROR: Failed to scan rows")
		return Subnet{}, err
	}

	log.Println("INFO: Subnet " + subnet.NetworkName + " found")
	return subnet, nil
}

func GetSubnestByDomainId(id int) ([]Subnet, error) {
	log.Println("INFO: Getting subnets by domain id " + strconv.Itoa(id))
	rows, err := DB.Query("SELECT * FROM Subnets WHERE DomainId = ?", id)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
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
			if err == sql.ErrNoRows {
				log.Println("ERROR: No rows found")
				return nil, err
			}
			log.Println("ERROR: Failed to scan rows")
			return nil, err
		}
		subnets = append(subnets, snet)
	}

	log.Println("INFO: Found " + strconv.Itoa(len(subnets)) + " subnets")
	return subnets, nil
}

func GetSubnestByDomainName(domainname string) ([]Subnet, error) {
	log.Println("INFO: Getting subnets by domain name " + domainname)
	id, err := GetDomainIdByDomainName(domainname)
	if err != nil {
		log.Println("ERROR: Failed to get domain id by domain name")
		return nil, err
	}
	if id == 0 {
		log.Println("ERROR: No domain found with name " + domainname)
		return nil, fmt.Errorf("no domain found with name %s", domainname)
	}

	rows, err := DB.Query("SELECT * FROM Subnets WHERE DomainId = ?", id)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
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
			if err == sql.ErrNoRows {
				log.Println("ERROR: No rows found")
				return nil, err
			}
			log.Println("ERROR: Failed to scan rows")
			return nil, err
		}
		subnets = append(subnets, snet)
	}

	log.Println("INFO: Found " + strconv.Itoa(len(subnets)) + " subnets")
	return subnets, nil
}

func GetSubnets() ([]Subnet, error) {
	log.Println("INFO: Getting all subnets")
	rows, err := DB.Query("SELECT * FROM Subnets")
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
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
			if err == sql.ErrNoRows {
				log.Println("ERROR: No rows found")
				return nil, err
			}
			log.Println("ERROR: Failed to scan rows")
			return nil, err
		}
		subnets = append(subnets, snet)
	}

	log.Println("INFO: Found " + strconv.Itoa(len(subnets)) + " subnets")
	return subnets, nil
}
