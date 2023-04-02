package internal

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

const (
	cacheFilePerm fs.FileMode = 0o600
)

type Cache struct {
	Accounts map[string]account
	Config   string
}

func GetCaches(prefix string) ([]string, error) {
	stateHome, err := GetStateHome()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(stateHome)
	if err != nil {
		return nil, err
	}

	var result []string

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := filepath.Ext(entry.Name())

			if ext == ".json" {
				config := strings.TrimSuffix(entry.Name(), ext)

				if strings.HasPrefix(config, prefix) {
					result = append(result, config)
				}
			}
		}
	}

	return result, nil
}

func NewCache(config string) (*Cache, error) {
	result := &Cache{
		Accounts: make(map[string]account),
		Config:   config,
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

func (c *Cache) clean(clock clock, expiry time.Duration) {
	now := clock.now()

	for _, acct := range c.Accounts {
		if !now.Before(acct.Timestamp.Add(expiry)) {
			delete(c.Accounts, acct.Object)
		}
	}
}

func (c *Cache) Len() int {
	return len(c.Accounts)
}

func (c *Cache) Remove() error {
	file, err := c.filePath()
	if err != nil {
		return err
	}

	return os.Remove(file)
}

func (c *Cache) SortedAccounts(prefix string, exclusions []string) []string {
	result := make([]string, 0, c.Len())

	for name := range c.Accounts {
		lowerName := strings.ToLower(name)

		if (prefix == "" || strings.HasPrefix(lowerName, prefix)) &&
			!ContainsFunc(exclusions, func(s string) bool {
				return strings.ToLower(s) == lowerName
			}) {
			result = append(result, name)
		}
	}

	sort.Strings(result)

	return result
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

	return filepath.Join(home, c.fileName()), nil
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

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, cacheFilePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(bytes)

	return err
}

func (c *Cache) merge(accts []account) {
	for _, acct := range accts {
		if acct.ok() {
			c.Accounts[acct.Object] = acct
		}
	}
}
