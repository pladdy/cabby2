package sqlite

import (
	"database/sql"

	// import sqlite dependency
	_ "github.com/mattn/go-sqlite3"

	cabby "github.com/pladdy/cabby2"
)

// DiscoveryService implements a SQLite version of the DiscoveryService interface
type DiscoveryService struct {
	DB *sql.DB
}

// Discovery will read from the data store and return the resource
func (s DiscoveryService) Discovery() (cabby.Discovery, error) {
	resource, action := "Discovery", "read"
	start := cabby.LogServiceStart(resource, action)
	result, err := s.discovery()
	cabby.LogServiceEnd(resource, action, start)
	return result, err
}

func (s DiscoveryService) discovery() (cabby.Discovery, error) {
	sql := `select td.title, td.description, td.contact, td.default_url,
						 case
							 when tar.api_root_path is null then 'No API Roots defined' else tar.api_root_path
						 end api_root_path
					 from
						 taxii_discovery td
						 left join taxii_api_root tar
							 on td.id = tar.discovery_id`

	d := cabby.Discovery{}
	var apiRoots []string

	rows, err := s.DB.Query(sql)
	if err != nil {
		return d, err
	}

	for rows.Next() {
		var apiRoot string
		if err := rows.Scan(&d.Title, &d.Description, &d.Contact, &d.Default, &apiRoot); err != nil {
			return d, err
		}
		if apiRoot != "No API Roots defined" {
			apiRoots = append(apiRoots, apiRoot)
		}
	}

	err = rows.Err()
	d.APIRoots = apiRoots
	return d, err
}
