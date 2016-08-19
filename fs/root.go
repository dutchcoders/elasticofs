package fs

import (
	"path/filepath"

	"github.com/dutchcoders/elasticofs/json"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type Root struct {
	Dir
}

func (d *Root) Views() map[string]string {
	return map[string]string{
		"_stats":       "_stats",
		"_field_stats": "_field_stats",
	}
}

func (d *Root) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if p, ok := d.Views()[name]; ok {
		return &View{
			mfs:  d.mfs,
			Path: filepath.Join(d.Path, p),
			/*template*/
		}, nil
	}

	switch {
	case name == "_nodes":
		return &Nodes{
			Dir: Dir{
				mfs:  d.mfs,
				Path: "_nodes",
			},
		}, nil
	case name == "_cat":
		return &Cat{
			Dir: Dir{
				mfs:  d.mfs,
				Path: "_cat",
			},
		}, nil
	case name == "_cluster":
		return &Cluster{
			Dir: Dir{
				mfs:  d.mfs,
				Path: "_cluster",
			},
		}, nil
	default:
		return &Index{
			Dir: Dir{
				mfs:  d.mfs,
				Path: name,
			},
		}, nil
	}
	return nil, fuse.ENOENT
}

func (dir *Root) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var entries = []fuse.Dirent{}
	for k, _ := range dir.Views() {
		entries = append(entries, fuse.Dirent{
			Inode: 0, Name: k, Type: fuse.DT_File,
		})
	}

	entries = append(entries, fuse.Dirent{
		Inode: 0, Name: "_cluster", Type: fuse.DT_Dir,
	})

	entries = append(entries, fuse.Dirent{
		Inode: 0, Name: "_cat", Type: fuse.DT_Dir,
	})

	entries = append(entries, fuse.Dirent{
		Inode: 0, Name: "_nodes", Type: fuse.DT_Dir,
	})

	// retrieve indexes
	path := filepath.Join(dir.Path, "_stats")

	req, err := dir.mfs.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := dir.mfs.client.Do(req, &resp); err != nil {
		return nil, err
	}

	if v, ok := resp["indices"]; !ok {
	} else if indices, ok := v.(json.M); !ok {
	} else {
		for k, _ := range indices {
			index := &Index{
				Dir{
					mfs:  dir.mfs,
					Path: k,
				},
			}

			entries = append(entries, index.Dirent())
		}
	}

	return entries, nil
}
