package internal

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

const driverName = "sqlite3"

type DBCache struct {
	db *sql.DB
}

func NewDBCache() (DBCache, error) {
	var result DBCache

	home, err := GetStateHome()
	if err != nil {
		return result, err
	}

	result.db, err = sql.Open(driverName, filepath.Join(home, "accounts.sqlite"))
	if err != nil {
		return result, err
	}

	return result, result.init()
}

func (c DBCache) Close() {
	_ = c.db.Close()
}

func (c DBCache) Configs(prefix string) ([]string, error) {
	rows, err := c.db.Query(
		"select distinct config from account where lower(config) like ? order by config",
		strings.ToLower(prefix)+"%",
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []string

	for rows.Next() {
		var config string

		err = rows.Scan(&config)
		if err != nil {
			return nil, err
		}

		result = append(result, config)
	}

	return result, nil
}

func (c DBCache) RemoveAll(config string) error {
	_, err := c.db.Exec("delete from account where config = ?", config)

	return err
}

func (c DBCache) SortedAccounts(config, prefix string, exclusions []string) ([]Account, error) {
	rows, err := c.db.Query("select name, value, created_at from account where config = ? order by name", config)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []Account

	for rows.Next() {
		var (
			acct      Account
			createdAt int64
		)

		err = rows.Scan(&acct.Object, &acct.Value, &createdAt)
		if err != nil {
			return nil, err
		}

		acct.Timestamp = time.Unix(createdAt, 0)

		lowerName := strings.ToLower(acct.Object)

		if (prefix == "" || strings.HasPrefix(lowerName, prefix)) &&
			!ContainsFunc(exclusions, func(s string) bool {
				return strings.ToLower(s) == lowerName
			}) {
			result = append(result, acct)
		}
	}

	return result, nil
}

func (c DBCache) init() error {
	_, err := c.db.Exec(`
		create table if not exists account (
		    config     text not null,
		    name       text not null,
		    value      text not null,
		    created_at int  not null,
		    primary key (config, name)
		) strict
	`)

	return err
}

func (c DBCache) clean(clock clock, expiry time.Duration) error {
	minCreatedDate := clock.now().Add(expiry * -1)

	_, err := c.db.Exec("delete from account where created_at < ?", minCreatedDate)

	return err
}

func (c DBCache) get(config, name string) (Account, error) {
	result := Account{
		Object:     name,
		StatusCode: http.StatusOK,
	}

	row := c.db.QueryRow("select value, created_at from account where config = ? and name = ?",
		config,
		name,
	)

	var createdAt int64

	err := row.Scan(&result.Value, &createdAt)

	result.Timestamp = time.Unix(createdAt, 0)

	return result, err
}

func (c DBCache) merge(config string, accounts []Account) error {
	for _, acct := range accounts {
		if _, err := c.db.Exec(
			"insert into account (config, name, value, created_at) values(?, ?, ?, ?) on conflict do nothing",
			config,
			acct.Object,
			acct.Value,
			acct.Timestamp.Unix(),
		); err != nil {
			return err
		}
	}

	return nil
}
