package fs

import (
	"path/filepath"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type IndexType struct {
	Dir
}

func (d *IndexType) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if p, ok := d.Views()[name]; ok {
		return &View{
			mfs:  d.mfs,
			Path: filepath.Join(d.Path, p),
			/*template*/
		}, nil
	}

	return nil, fuse.ENOENT
}

func (it *IndexType) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var entries = []fuse.Dirent{}
	for k, _ := range it.Views() {
		entries = append(entries, fuse.Dirent{
			Inode: 0, Name: k, Type: fuse.DT_File,
		})
	}

	return entries, nil
}

func (d *IndexType) Views() map[string]string {
	return map[string]string{
		"_mapping": "_mapping",
	}
}
