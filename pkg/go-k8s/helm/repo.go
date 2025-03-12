package helm

import (
	"context"
	"fmt"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// RepoAddAndUpdate adds repo with given name and url and updates charts for all helm repos
func (c *Client) RepoAddAndUpdate(name, url string) error {
	err := c.RepoAdd(name, url)
	if err != nil {
		return err
	}
	return c.RepoUpdate()
}

// RepoAdd adds repo with given name and url
func (c *Client) RepoAdd(name, url string) error {
	repoFile := c.Settings.RepositoryConfig

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), 0750)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(repoFile), err)
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer func(fileLock *flock.Flock) {
			err := fileLock.Unlock()
			if err != nil {
				log.Fatal(err)
			}
		}(fileLock)
	}
	if err != nil {
		return fmt.Errorf("failed to lock file %s: %w", fileLock, err)
	}

	b, err := os.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		log.Fatal(err)
	}

	if f.Has(name) {
		// repo already exists
		return nil
	}

	clt := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&clt, getter.All(c.Settings))
	if err != nil {
		return fmt.Errorf("failed to create chart repository: %w", err)
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
		return fmt.Errorf("failed to download index file: %w", err)
	}

	f.Update(&clt)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		return fmt.Errorf("failed to write repository file: %w", err)
	}

	return nil
}

// RepoUpdate updates charts for all helm repos
func (c *Client) RepoUpdate() error {
	repoFile := c.Settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		return errors.New("no repositories found. You must add one before updating")
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(c.Settings))
		if err != nil {
			return fmt.Errorf("failed to create chart repository: %w", err)
		}
		repos = append(repos, r)
	}

	if c.Settings.Debug {
		fmt.Printf("Hang tight while we grab the latest from your chart repositories...\n")
	}
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if c.Settings.Debug {
				if _, err := re.DownloadIndexFile(); err != nil {
					fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
				} else {
					fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
				}
			}
		}(re)
	}
	wg.Wait()
	if c.Settings.Debug {
		fmt.Printf("Update Complete. ⎈ Happy Helming!⎈\n")
	}

	return nil
}
