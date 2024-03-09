package model

import "database/sql"

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
