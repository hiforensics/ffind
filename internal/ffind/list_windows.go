//go:build windows

package ffind

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/mxk/go-vss"

	"github.com/cuhsat/futils/pkg/sys"
)

const (
	symlink = "ffind"
)

func ShadowCopy(drv string) (dir string, err error) {
	dir = filepath.Join(os.TempDir(), symlink)

	if _, err = os.Stat(dir); !os.IsNotExist(err) {
		os.Remove(dir)
	}

	err = vss.CreateLink(dir, drv)

	return
}

func Header() []string {
	return []string{
		"Filename",
		"Path",
		"Size (bytes)",
		"Modified",
		"Accessed",
		"Created",
	}
}

func Record(f string) (l []string) {
	l = append(l, filepath.Base(f))
	l = append(l, f)

	fi, err := os.Stat(f)

	if err != nil {
		sys.Error(err)
		return
	}

	fd := fi.Sys().(*syscall.Win32FileAttributeData)

	m := time.Unix(0, fd.LastWriteTime.Nanoseconds())
	a := time.Unix(0, fd.LastAccessTime.Nanoseconds())
	c := time.Unix(0, fd.CreationTime.Nanoseconds())

	l = append(l, fmt.Sprintf("%d", fi.Size()))
	l = append(l, m.String(), a.String(), c.String())

	return
}
