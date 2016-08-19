package fs

import (
	"path/filepath"

	"github.com/dutchcoders/elasticofs/json"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type Index struct {
	Dir
}

func (d *Index) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if p, ok := d.Views()[name]; ok {
		return &View{
			mfs:  d.mfs,
			Path: filepath.Join(d.Path, p),
			/*template*/
		}, nil
	} else {
		return &IndexType{
			Dir{
				mfs:  d.mfs,
				Path: filepath.Join(d.Path, name),
			},
			/*template*/
		}, nil
	}

	return nil, fuse.ENOENT
}

func (d *Index) Views() map[string]string {
	return map[string]string{
		"_mapping":     "_mapping",
		"_stats":       "_stats",
		"_field_stats": "_field_stats",
	}
}

func (index *Index) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var entries = []fuse.Dirent{}
	for k, _ := range index.Views() {
		entries = append(entries, fuse.Dirent{
			Inode: 0, Name: k, Type: fuse.DT_File,
		})
	}

	path := "/"
	path = filepath.Join(path, index.Path)
	path = filepath.Join(path, "_mapping")

	req, err := index.mfs.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp json.M
	if err := index.mfs.client.Do(req, &resp); err != nil {
		return nil, err
	}

	if v, ok := resp[index.Path]; !ok {
	} else if indices, ok := v.(json.M); !ok {
	} else if m, ok := indices["mappings"]; !ok {
	} else if mappings, ok := m.(json.M); !ok {
	} else {
		for k, _ := range mappings {
			index := &IndexType{
				Dir{
					mfs:  index.mfs,
					Path: filepath.Join(index.Path, k),
				},
			}

			entries = append(entries, index.Dirent())
		}
	}

	return entries, nil
}
