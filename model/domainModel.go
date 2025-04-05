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

func CreateDomain(d Domain, id int) (bool, error) {
	log.Println("INFO: Creating domain " + d.DomainName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to create domain " + d.DomainName)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to create domain " + d.DomainName)
			t.Rollback()
		}
	}()

	q, err := t.Prepare("INSERT INTO Domains (DomainName, CreatorId) VALUES (?, ?)")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	_, err = q.Exec(d.DomainName, id)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: Domain " + d.DomainName + " created successfully")
	return true, nil
}

func DeleteDomain(domain string) (bool, error) {
	log.Println("INFO: Deleting domain " + domain)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Failed to begin transaction")
		return false, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed to delete domain " + domain)
			t.Rollback()
		}
		if err != nil {
			log.Println("ERROR: Failed to delete domain " + domain)
			t.Rollback()
		}
	}()

	q, err := DB.Prepare("DELETE FROM Domains WHERE DomainName IS ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return false, err
	}

	_, err = q.Exec(domain)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return false, err
	}

	err = t.Commit()
	if err != nil {
		log.Println("ERROR: Failed to commit transaction")
		return false, err
	}

	log.Println("INFO: Domain " + domain + " deleted successfully")
	return true, nil
}

func GetDomainById(id int) (Domain, error) {
	idStr := strconv.Itoa(id)
	log.Println("INFO: Getting domain by id " + idStr)
	rec, err := DB.Prepare("SELECT * FROM Domains WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return Domain{}, err
	}
	defer rec.Close()

	domain := Domain{}
	r, err := rec.Query(id)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return Domain{}, err
	}
	defer r.Close()

	err = r.Scan(
		&domain.Id,
		&domain.DomainName,
		&domain.CreatorId,
		&domain.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No rows found")
			return Domain{}, nil
		}
		log.Println("ERROR: Failed to scan rows")
		return Domain{}, err
	}

	log.Println("INFO: Domain " + domain.DomainName + " found")
	return domain, nil
}

func GetDomainByDomainName(domainname string) (Domain, error) {
	log.Println("INFO: Getting domain by name " + domainname)
	rec, err := DB.Prepare("SELECT * FROM Domains WHERE DomainName = ?")
	if err != nil {
		log.Println("ERROR: Failed to prepare statement")
		return Domain{}, err
	}
	defer rec.Close()

	domain := Domain{}

	r, err := rec.Query(domainname)
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
		return Domain{}, err
	}
	defer r.Close()

	err = r.Scan(
		&domain.Id,
		&domain.DomainName,
		&domain.CreatorId,
		&domain.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No rows found")
			return Domain{}, nil
		}
		log.Println("ERROR: Failed to scan rows")
		return Domain{}, err
	}

	log.Println("INFO: Domain " + domain.DomainName + " found")
	return domain, nil
}

func GetDomains() ([]Domain, error) {
	log.Println("INFO: Getting all domains")
	rows, err := DB.Query("SELECT * FROM Domains")
	if err != nil {
		log.Println("ERROR: Failed to execute statement")
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
			if err == sql.ErrNoRows {
				log.Println("ERROR: No rows found")
				return nil, nil
			}
			log.Println("ERROR: Failed to scan rows")
			return nil, err
		}
		domains = append(domains, domain)
	}

	log.Println("INFO: Found " + strconv.Itoa(len(domains)) + " domains")
	return domains, nil
}
