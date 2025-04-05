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
	"encoding/json"
	"log"
	"strconv"
)

func CreateHost(h Host, id int) (bool, error) {
	log.Println("INFO: Creating host " + h.HostName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to create host " + h.HostName)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to create host " + h.HostName)
			t.Rollback()
		}
	}()

	q, err := t.Prepare("INSERT INTO Hosts (HostName, MacAddresses, CreatorId) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	strJsonMacAddressesSlice, err := json.Marshal(h.MacAddresses)
	if err != nil {
		log.Println("ERROR: Failed to marshal mac addresses")
		return false, err
	}

	_, err = q.Exec(h.HostName, strJsonMacAddressesSlice, id)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: Host " + h.HostName + " created successfully")
	return true, nil
}

func DeleteHostname(hostname string) (bool, error) {
	log.Println("INFO: Deleting host " + hostname)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to delete host " + hostname)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to delete host " + hostname)
			t.Rollback()
		}
	}()

	q, err := t.Prepare("DELETE FROM Hosts WHERE HostName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	_, err = q.Exec(hostname)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: Host " + hostname + " deleted successfully")
	return true, nil
}

func UpdateMacAddresses(hostname string, data []string) (bool, error) {
	log.Println("INFO: Updating MAC addresses for host " + hostname)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to update MAC addresses for host " + hostname)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to update MAC addresses for host " + hostname)
			t.Rollback()
		}
	}()

	q, err := t.Prepare("UPDATE Hosts SET MacAddresses = ? WHERE HostName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	macAddressSlice, err := json.Marshal(data)
	if err != nil {
		log.Println("ERROR: Failed to marshal mac addresses")
		return false, err
	}

	_, err = q.Exec(macAddressSlice, hostname)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: MAC addresses for host " + hostname + " updated successfully")
	return true, nil
}

func GetHostById(id int) (Host, error) {
	idStr := strconv.Itoa(id)
	log.Println("INFO: Getting host by ID: " + idStr)
	rec, err := DB.Prepare("SELECT * FROM Hosts WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return Host{}, err
	}
	defer rec.Close()

	strHost := StringHost{}

	r, err := rec.Query(id)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return Host{}, err
	}
	defer r.Close()

	err = r.Scan(
		&strHost.Id,
		&strHost.HostName,
		&strHost.MacAddresses,
		&strHost.CreatorId,
		&strHost.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No host found with ID: " + idStr)
			return Host{}, nil
		}
		log.Println("ERROR: Failed to scan rows")
		return Host{}, err
	}
	// process the MAC addresses
	unmarshalledMacAddresses := make([]string, 0)
	_ = json.Unmarshal([]byte(strHost.MacAddresses), &unmarshalledMacAddresses)
	host := Host{}
	host.Id = strHost.Id
	host.HostName = strHost.HostName
	host.MacAddresses = unmarshalledMacAddresses
	host.CreatorId = strHost.CreatorId
	host.CreationDate = strHost.CreationDate

	log.Println("INFO: Host " + host.HostName + " found")
	return host, nil
}

func GetHostByHostName(hostname string) (Host, error) {
	log.Println("INFO: Getting host by name: " + hostname)
	rec, err := DB.Prepare("SELECT * FROM Hosts WHERE HostName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return Host{}, err
	}
	defer rec.Close()

	strHost := StringHost{}
	r, err := rec.Query(hostname)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return Host{}, err
	}
	defer r.Close()

	err = r.Scan(
		&strHost.Id,
		&strHost.HostName,
		&strHost.MacAddresses,
		&strHost.CreatorId,
		&strHost.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No host found with name: " + hostname)
			return Host{}, nil
		}
		log.Println("ERROR: Failed to scan rows")
		return Host{}, err
	}

	// process the MAC addresses
	unmarshalledMacAddresses := make([]string, 0)
	_ = json.Unmarshal([]byte(strHost.MacAddresses), &unmarshalledMacAddresses)
	host := Host{}
	host.Id = strHost.Id
	host.HostName = strHost.HostName
	host.MacAddresses = unmarshalledMacAddresses
	host.CreatorId = strHost.CreatorId
	host.CreationDate = strHost.CreationDate

	log.Println("INFO: Host " + host.HostName + " found")
	return host, nil
}

func GetHosts() ([]Host, error) {
	log.Println("INFO: Getting all hosts")
	rows, err := DB.Query("SELECT * FROM Hosts")
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return nil, err
	}
	defer rows.Close()

	hosts := make([]Host, 0)
	for rows.Next() {
		strHost := StringHost{}
		err = rows.Scan(
			&strHost.Id,
			&strHost.HostName,
			&strHost.MacAddresses,
			&strHost.CreatorId,
			&strHost.CreationDate,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("ERROR: No hosts found")
				return nil, nil
			}
			log.Println("ERROR: Failed to scan rows")
			return nil, err
		}

		// process json of MacAddresses to remove unneeded escapes
		unmarshalledMacAddresses := make([]string, 0)
		_ = json.Unmarshal([]byte(strHost.MacAddresses), &unmarshalledMacAddresses)
		for i := 0; i < len(unmarshalledMacAddresses); i++ {
			log.Println("DEBUG: Unmarshalled MAC Addresses: " + unmarshalledMacAddresses[i])
		}

		host := Host{}
		host.Id = strHost.Id
		host.HostName = strHost.HostName
		host.MacAddresses = unmarshalledMacAddresses
		host.CreatorId = strHost.CreatorId
		host.CreationDate = strHost.CreationDate

		hosts = append(hosts, host)
	}

	log.Println("INFO: Found " + strconv.Itoa(len(hosts)) + " hosts")
	return hosts, nil
}
