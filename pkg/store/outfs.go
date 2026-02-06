package store

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

type NamespaceInfo struct {
	Name         string
	ItemCount    int
	Updated      time.Time
	UpdatedHuman string
}

func EnsureOut(outDir string) error {
	return os.MkdirAll(outDir, 0o755)
}

func ListNamespaces(outDir string) ([]NamespaceInfo, error) {
	ents, err := os.ReadDir(outDir)
	if err != nil {
		return nil, err
	}

	var res []NamespaceInfo
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		ns := e.Name()
		nsDir := filepath.Join(outDir, ns)

		itemCount := 0
		latest := time.Time{}

		// Walk only one level deep is usually enough; but full walk is fine too
		filepath.WalkDir(nsDir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			itemCount++
			if info, err := d.Info(); err == nil {
				if info.ModTime().After(latest) {
					latest = info.ModTime()
				}
			}
			return nil
		})

		human := ""
		if !latest.IsZero() {
			human = latest.Format("2006-01-02 15:04")
		}

		res = append(res, NamespaceInfo{
			Name:         ns,
			ItemCount:    itemCount,
			Updated:      latest,
			UpdatedHuman: human,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res, nil
}
