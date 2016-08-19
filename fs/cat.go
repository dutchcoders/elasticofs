package fs

import (
	"path/filepath"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type Cat struct {
	Dir
}

func (d *Cat) Views() map[string]string {
	return map[string]string{
		"aliases":      "aliases",
		"count":        "count",
		"allocation":   "allocation",
		"health":       "health",
		"master":       "master",
		"nodes":        "nodes",
		"plugins":      "plugins",
		"repositories": "repositories",
		"tasks":        "tasks",
	}
}

func (c *Cat) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if p, ok := c.Views()[name]; ok {
		return &View{
			mfs:  c.mfs,
			Path: filepath.Join(c.Path, p),
			/*template*/
		}, nil
	}

	return nil, fuse.ENOENT
}

func (c *Cat) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var entries = []fuse.Dirent{}
	for k, _ := range c.Views() {
		entries = append(entries, fuse.Dirent{
			Inode: 0, Name: k, Type: fuse.DT_File,
		})
	}

	return entries, nil
}
