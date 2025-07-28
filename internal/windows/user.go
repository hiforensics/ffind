package windows

import (
	"fmt"
	"os"
	"path/filepath"
)

func UserDrive() string {
	if env := os.ExpandEnv("$HOMEDRIVE"); len(env) > 0 {
		return env
	} else {
		return "C:"
	}
}

func EnumUsers(sysroot string, out chan<- string) {
	if len(sysroot) == 0 {
		sysroot = UserDrive()
	}

	root := filepath.Join(sysroot, "Users")

	if _, err := os.Stat(root); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	r, err := os.Open(root)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	dir, err := r.Readdir(-1)

	r.Close()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, fi := range dir {
		if !fi.IsDir() {
			continue
		}

		for _, artifact := range [...]string{
			"/[Nn][Tt][Uu][Ss][Ee][Rr].[Dd][Aa][Tt]",
			"[Aa]pp[Dd]ata/[Ll]ocal/[Mm]icrosoft/[Ww]indows/[Uu]sr[Cc]lass.dat",
			"[Aa]pp[Dd]ata/[Rr]oaming/[Mm]icrosoft/[Ww]indows/[Rr]ecent/*/*.*[Dd]estinations-ms",
			"[Aa]pp[Dd]ata/[Ll]ocal/[Mm]icrosoft/[Ww]indows/[Ww]eb[Cc]ache/[Ww]eb[Cc]ache[Vv]??.dat",
			"[Aa]pp[Dd]ata/[Ll]ocal/*/[Uu]ser [Dd]ata/*/[Hh]istory",
			"[Aa]pp[Dd]ata/[Ll]ocal/*/*/[Uu]ser [Dd]ata/*/[Hh]istory",
			"[Aa]pp[Dd]ata/[Rr]oaming/*/*/[Pp]rofiles/*/places.sqlite",
			"[Aa]pp[Dd]ata/[Rr]oaming/*/*/*/[Hh]istory",
		} {
			files, err := filepath.Glob(filepath.Join(root, fi.Name(), artifact))

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			for _, file := range files {
				out <- filepath.ToSlash(file)
			}
		}
	}
}
