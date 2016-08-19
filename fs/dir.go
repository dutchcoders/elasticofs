package fs

import (
	"os"
	"path"

	"golang.org/x/net/context"

	"bazil.org/fuse"
)

type Dir struct {
	Path  string
	Inode uint64
	mfs   *ElasticoFS
}

func (dir *Dir) Dirent() fuse.Dirent {
	return fuse.Dirent{
		Inode: dir.Inode, Name: path.Base(dir.Path), Type: fuse.DT_Dir,
	}
}

func (dir *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	*a = fuse.Attr{
		Mode: os.ModeDir | dir.mfs.config.mode | 0110,
	}
	return nil
}
