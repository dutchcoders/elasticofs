package fs

import (
	"path"
	"path/filepath"

	"github.com/dutchcoders/elasticofs/json"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type Nodes struct {
	Dir
}

func (d *Nodes) Views() map[string]string {
	return map[string]string{
		"stats": "stats",
	}
}

func (d *Nodes) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if p, ok := d.Views()[name]; ok {
		return &View{
			mfs:  d.mfs,
			Path: path.Join(d.Path, p),
			/*template*/
		}, nil
	}

	return &Node{
		Dir{
			mfs:  d.mfs,
			Path: path.Join(d.Path, name),
			/*template*/
		},
	}, nil

	return nil, fuse.ENOENT
}

func (dir *Nodes) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var entries = []fuse.Dirent{}
	for k, _ := range dir.Views() {
		entries = append(entries, fuse.Dirent{
			Inode: 0, Name: k, Type: fuse.DT_File,
		})
	}

	// retrieve nodes
	path := filepath.Join(dir.Path, "stats")

	req, err := dir.mfs.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := dir.mfs.client.Do(req, &resp); err != nil {
		return nil, err
	}

	if v, ok := resp["nodes"]; !ok {
	} else if nodes, ok := v.(json.M); !ok {
	} else {
		for k, _ := range nodes {
			node := &Node{
				Dir{
					mfs:  dir.mfs,
					Path: k,
				},
			}

			entries = append(entries, node.Dirent())
		}
	}

	return entries, nil
}
