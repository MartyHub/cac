package internal

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"runtime"
)

const (
	cacheFilePerm fs.FileMode = 0o600
)

type Cache struct {
	Config   string
	Accounts []account
}

func NewCache(config string) (*Cache, error) {
	result := &Cache{
		Config: config,
	}

	if !result.exists() {
		return result, nil
	}

	file, err := result.filePath()

	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(file)

	if os.IsNotExist(err) {
		return result, nil
	}

	if err != nil {
		return nil, err
	}

	err = result.checkPermissions(file)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(bytes, &result.Accounts); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Cache) exists() bool {
	return c.Config != ""
}

func (c *Cache) fileName() string {
	return c.Config + ".json"
}

func (c *Cache) filePath() (string, error) {
	home, err := GetStateHome()

	if err != nil {
		return "", err
	}

	return path.Join(home, c.fileName()), nil
}

func (c *Cache) checkPermissions(file string) error {
	if runtime.GOOS != "windows" {
		stat, err := os.Stat(file)

		if err != nil {
			return err
		}

		permissions := stat.Mode().Perm()

		if permissions != cacheFilePerm {
			if err = os.Chmod(file, cacheFilePerm); err != nil {
				return fmt.Errorf(
					"incorrect permissions %v for file %s (must be %v): %w",
					permissions,
					file,
					cacheFilePerm,
					err,
				)
			}
		}
	}

	return nil
}

func (c *Cache) save() error {
	if !c.exists() {
		return nil
	}

	bytes, err := json.MarshalIndent(c.Accounts, "", "  ")

	if err != nil {
		return err
	}

	filePath, err := c.filePath()

	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, cacheFilePerm)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(bytes)

	return err
}

func (c *Cache) get(object string) (account, bool) {
	for _, acct := range c.Accounts {
		if acct.Object == object {
			return acct, true
		}
	}

	return account{
		Object: object,
	}, false
}

func (c *Cache) set(acct account) {
	for i := range c.Accounts {
		if c.Accounts[i].Object == acct.Object {
			c.Accounts[i] = acct

			return
		}
	}

	c.Accounts = append(c.Accounts, acct)
}

func (c *Cache) merge(accts []account) {
	for _, acct := range accts {
		if acct.ok() {
			c.set(acct)
		}
	}
}
