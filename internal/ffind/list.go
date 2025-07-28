//go:build !windows

package ffind

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hiforensics/utils/pkg/sys"
)

func ShadowCopy(drive string) (dir string, err error) {
	return "", errors.ErrUnsupported
}

func Header() []string {
	return []string{
		"Filename",
		"Path",
		"Size (bytes)",
		"Modified",
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

	l = append(l, fmt.Sprintf("%d", fi.Size()))
	l = append(l, fi.ModTime().String())

	return
}
