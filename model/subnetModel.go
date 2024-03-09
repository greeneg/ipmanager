package model

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
