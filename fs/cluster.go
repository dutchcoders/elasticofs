package fs

import (
	"path/filepath"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type Cluster struct {
	Dir
}

func (d *Cluster) Views() map[string]string {
	return map[string]string{
		"health":        "health",
		"stats":         "stats",
		"settings":      "settings",
		"pending_tasks": "pending_tasks",
	}
}

func (c *Cluster) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if p, ok := c.Views()[name]; ok {
		return &View{
			mfs:  c.mfs,
			Path: filepath.Join(c.Path, p),
			/*template*/
		}, nil
	}

	return nil, fuse.ENOENT
}

func (cluster *Cluster) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var entries = []fuse.Dirent{}
	for k, _ := range cluster.Views() {
		entries = append(entries, fuse.Dirent{
			Inode: 0, Name: k, Type: fuse.DT_File,
		})
	}

	return entries, nil
}
