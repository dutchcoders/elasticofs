package fs

import (
	"bytes"

	"bazil.org/fuse"
	"golang.org/x/net/context"
)

type View struct {
	data []byte
	Path string
	mfs  *ElasticoFS
}

func (s *View) Attr(ctx context.Context, a *fuse.Attr) error {
	s.load()
	*a = fuse.Attr{
		Size: uint64(len(s.data)),
		Mode: s.mfs.config.mode,
	}

	return nil
}

func (s *View) ReadAll(ctx context.Context) ([]byte, error) {
	s.load()
	return s.data, nil
}

func (s *View) load() error {
	// todo(nl5887): once.Do()
	req, err := s.mfs.client.NewRequest("GET", s.Path, nil)
	if err != nil {
		return err
	}

	var resp bytes.Buffer
	if err := s.mfs.client.Do(req, &resp); err != nil {
		return err
	}

	s.data = resp.Bytes()
	return nil
}
