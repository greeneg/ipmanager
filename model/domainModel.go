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

func CreateDomain(d Domain, id int) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := t.Prepare("INSERT INTO Domains (DomainName, CreatorId) VALUES (?, ?)")
	if err != nil {
		return false, err
	}

	_, err = q.Exec(d.DomainName, id)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}

func DeleteDomain(domain string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM Domains WHERE DomainName IS ?")
	if err != nil {
		return false, err
	}

	_, err = q.Exec(domain)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}

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
