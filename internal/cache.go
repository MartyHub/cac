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

type cache struct {
	config string
	accts  []account
}

func newCache(config string) (*cache, error) {
	result := &cache{
		config: config,
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

	if err = json.Unmarshal(bytes, &result.accts); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *cache) exists() bool {
	return c.config != ""
}

func (c *cache) fileName() string {
	return c.config + ".json"
}

func (c *cache) filePath() (string, error) {
	home, err := GetStateHome()

	if err != nil {
		return "", err
	}

	return path.Join(home, c.fileName()), nil
}

func (c *cache) checkPermissions(file string) error {
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

func (c *cache) save() error {
	if !c.exists() {
		return nil
	}

	bytes, err := json.Marshal(c.accts)

	if err != nil {
		return err
	}

	filePath, err := c.filePath()

	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, cacheFilePerm)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(bytes)

	return err
}

func (c *cache) get(object string) (account, bool) {
	for _, acct := range c.accts {
		if acct.Object == object {
			return acct, true
		}
	}

	return account{
		Object: object,
	}, false
}

func (c *cache) set(acct account) {
	for i := range c.accts {
		if c.accts[i].Object == acct.Object {
			c.accts[i] = acct

			return
		}
	}

	c.accts = append(c.accts, acct)
}

func (c *cache) merge(accts []account) {
	for _, acct := range accts {
		if acct.ok() {
			c.set(acct)
		}
	}
}
