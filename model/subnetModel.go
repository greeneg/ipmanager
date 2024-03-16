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
