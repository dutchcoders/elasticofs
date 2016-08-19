package fs

import (
	"path/filepath"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type Node struct {
	Dir
}

func (n *Node) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var entries = []fuse.Dirent{}
	for k, _ := range n.Views() {
		entries = append(entries, fuse.Dirent{
			Inode: 0, Name: k, Type: fuse.DT_File,
		})
	}

	return entries, nil
}

func (n *Node) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if p, ok := n.Views()[name]; ok {
		return &View{
			mfs:  n.mfs,
			Path: filepath.Join(n.Path, p),
			/*template*/
		}, nil
	}

	return nil, fuse.ENOENT
}

func (d *Node) Views() map[string]string {
	return map[string]string{
		"stats": "stats",
	}
}
