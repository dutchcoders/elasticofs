/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package fs

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"

	"github.com/dutchcoders/elasticofs/client"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
)

type ElasticoFS struct {
	config *Config
	client *client.Client
}

func New(options ...func(*Config)) (*ElasticoFS, error) {

	// set defaults
	cfg := &Config{
		gid:  0,
		uid:  0,
		mode: os.FileMode(0440),
	}

	if u, err := user.Current(); err != nil {
		if v, err := strconv.Atoi(u.Gid); err == nil {
			cfg.gid = uint32(v)
		}

		if v, err := strconv.Atoi(u.Uid); err == nil {
			cfg.uid = uint32(v)
		}
	}

	for _, optionFn := range options {
		optionFn(cfg)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	fs := &ElasticoFS{
		config: cfg,
	}
	return fs, nil
}

func (mfs *ElasticoFS) mount() (*fuse.Conn, error) {
	return fuse.Mount(
		mfs.config.mountpoint,
		fuse.FSName("ElasticoFS"),
		fuse.Subtype("ElasticoFS"),
		fuse.LocalVolume(),
		fuse.VolumeName(mfs.config.bucket), // bucket?
	)
}

func (mfs *ElasticoFS) Serve() error {
	if mfs.config.debug {
		fuse.Debug = func(msg interface{}) {
			//	fmt.Printf("%#v\n", msg)
		}
	}

	fmt.Println("Initializing elastico client...")

	if client, err := client.New(mfs.config.target); err != nil {
		return err
	} else {
		mfs.client = client
	}

	fmt.Println("Mounting target....")
	// mount the drive
	c, err := mfs.mount()
	if err != nil {
		return err
	}

	defer c.Close()

	fmt.Println("Mounted... Have fun!")

	// serve the filesystem
	if err := fs.Serve(c, mfs); err != nil {
		return err
	}

	// todo(nl5887): implement this
	fmt.Println("HOW TO QUIT?")

	// todo(nl5887): move trap signals to Main, this is not supposed to be in Serve
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGUSR1)

loop:
	for {
		// check if the mount process has an error to report
		select {
		case <-c.Ready:
			if err := c.MountError; err != nil {
				log.Fatal(err)
			}
		case s := <-signalCh:
			if s == syscall.SIGUSR1 {
				fmt.Println("PRINT STATS")
				continue
			}

			break loop
		}
	}

	return nil
}

func (mfs *ElasticoFS) Root() (fs.Node, error) {
	return &Root{
		Dir: Dir{
			mfs:  mfs,
			Path: "/",
		},
		/*template*/
	}, nil
}
