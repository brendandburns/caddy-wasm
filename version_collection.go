package caddy_wasm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v55/github"
)

type VersionCollection interface {
	// Get Versions should always refresh from source
	GetVersions() ([]string, error)
	// GetWebAssembly should work, even if GetVersions() has never been called.
	GetWebAssembly(version string) ([]byte, error)
}

type fileSystemVersionCollection struct {
	fs   FileSystemInterface
	dir  string
	glob string
}

func ForFilesystemGlob(dir, glob string) VersionCollection {
	return &fileSystemVersionCollection{osfs{}, dir, glob}
}

func (f *fileSystemVersionCollection) GetVersions() ([]string, error) {
	files, err := f.fs.Glob(filepath.Join(f.dir, f.glob))
	if err != nil {
		return nil, err
	}

	for i, file := range files {
		files[i] = filepath.Base(file)
	}
	return files, nil
}

func (f *fileSystemVersionCollection) GetWebAssembly(version string) ([]byte, error) {
	return f.fs.ReadFile(filepath.Join(f.dir, version))
}

type urlVersionCollection struct {
	url string
}

func ForURL(url string) VersionCollection {
	return &urlVersionCollection{url}
}

func (u *urlVersionCollection) GetVersions() ([]string, error) {
	url, err := url.Parse(u.url)
	if err != nil {
		return nil, err
	}
	ix := strings.LastIndex(url.Path, "/")
	return []string{url.Path[ix+1:]}, nil
}

func (u *urlVersionCollection) GetWebAssembly(version string) ([]byte, error) {
	res, err := http.Get(u.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

type githubReleaseCollection struct {
	owner    string
	repo     string
	fileName string
	idMap    map[string]int64
	client   *github.Client
}

func ForGithubRepository(owner, repo, fileName string) (VersionCollection, error) {
	client := github.NewClient(nil)
	return &githubReleaseCollection{
		owner:    owner,
		repo:     repo,
		fileName: fileName,
		idMap:    make(map[string]int64),
		client:   client}, nil
}

func (g *githubReleaseCollection) GetVersions() ([]string, error) {
	r, _, e := g.client.Repositories.ListReleases(context.Background(), g.owner, g.repo, nil)
	if e != nil {
		return nil, e
	}
	v := []string{}
	for _, release := range r {
		for _, asset := range release.Assets {
			if *asset.Name == g.fileName {
				v = append(v, *release.Name)
				g.idMap[*release.Name] = *asset.ID
				break
			}
		}
	}
	return v, nil
}

func (g *githubReleaseCollection) GetWebAssembly(version string) ([]byte, error) {
	id, found := g.idMap[version]
	if !found {
		if _, e := g.GetVersions(); e != nil {
			return nil, e
		}
		id, found = g.idMap[version]
		if !found {
			return nil, fmt.Errorf("couldn't find version: %v", version)
		}
	}
	data, _, e := g.client.Repositories.DownloadReleaseAsset(context.Background(), g.owner, g.repo, id, g.client.Client())
	if e != nil {
		return nil, e
	}
	defer data.Close()
	return io.ReadAll(data)
}
