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
)

func CreateHost(h Host, id int) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := t.Prepare("INSERT INTO Hosts (HostName, MacAddresses, CreatorId) VALUES (?, ?, ?)")
	if err != nil {
		return false, err
	}

	strJsonMacAddressesSlice, err := json.Marshal(h.MacAddresses)
	if err != nil {
		return false, err
	}
	_, err = q.Exec(h.HostName, strJsonMacAddressesSlice, id)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}

func DeleteHostname(hostname string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := t.Prepare("DELETE FROM Hosts WHERE HostName = ?")
	if err != nil {
		return false, err
	}

	_, err = q.Exec(hostname)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}

func UpdateMacAddresses(hostname string, data []string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := t.Prepare("UPDATE Hosts SET MacAddresses = ? WHERE HostName = ?")
	if err != nil {
		return false, err
	}

	macAddressSlice, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	_, err = q.Exec(macAddressSlice, hostname)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}

func GetHostById(id int) (Host, error) {
	rec, err := DB.Prepare("SELECT * FROM Hosts WHERE Id = ?")
	if err != nil {
		return Host{}, err
	}

	strHost := StringHost{}
	err = rec.QueryRow(id).Scan(
		&strHost.Id,
		&strHost.HostName,
		&strHost.MacAddresses,
		&strHost.CreatorId,
		&strHost.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Host{}, nil
		}
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

	return host, nil
}

func GetHostByHostName(hostname string) (Host, error) {
	rec, err := DB.Prepare("SELECT * FROM Hosts WHERE HostName = ?")
	if err != nil {
		return Host{}, err
	}

	strHost := StringHost{}
	err = rec.QueryRow(hostname).Scan(
		&strHost.Id,
		&strHost.HostName,
		&strHost.MacAddresses,
		&strHost.CreatorId,
		&strHost.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Host{}, nil
		}
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

	return host, nil
}

func GetHosts() ([]Host, error) {
	rows, err := DB.Query("SELECT * FROM Hosts")
	if err != nil {
		return nil, err
	}

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

	return hosts, nil
}
