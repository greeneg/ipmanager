package model

import "database/sql"

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
