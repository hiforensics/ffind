package ffind

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/cuhsat/futils/pkg/hash"
	"github.com/cuhsat/futils/pkg/sys"
	"github.com/cuhsat/futils/pkg/zip"

	"github.com/cuhsat/ffind/internal/windows"
)

const (
	fallback = "ffind"
	zipExt   = ".zip"
	csvExt   = ".csv"
	rLimit   = 1024
)

const (
	Version = "0.5.2"
)

type ffind struct {
	wg sync.WaitGroup

	root string
	zip  string
	csv  string
	hsh  string
	rp   bool
	sc   bool
	so   bool
	uo   bool
}

type fstep func(in <-chan string, out chan<- string)

func Find(root, zip, csv, hsh string, rp, sc, so, uo bool) (files []string) {
	ff := &ffind{
		root: root,
		zip:  zip,
		csv:  csv,
		hsh:  hsh,
		rp:   rp,
		sc:   sc,
		so:   so,
		uo:   uo,
	}

	// Switch to live mode
	if len(ff.root)+len(ff.zip)+len(ff.csv)+len(ff.hsh) == 0 {
		if runtime.GOOS == "windows" {
			ff.live()
		}
	}

	var ch [4]chan string
	var ci = 0

	for i := range ch {
		ch[i] = make(chan string, rLimit)
	}

	add := func(fn fstep) {
		ff.wg.Add(1)

		go fn(ch[ci], ch[ci+1])

		ci++
	}

	ff.wg.Add(1)

	go ff.enum(ch[ci])

	if len(ff.zip) > 0 {
		add(ff.comp)
	}

	if len(ff.csv) > 0 {
		add(ff.list)
	}

	if len(ff.hsh) > 0 {
		add(ff.hash)
	}

	for f := range ch[ci] {
		if len(ff.hsh) == 0 {
			f = ff.path(f)
		}

		if sys.Progress != nil {
			sys.Progress(f)
		}

		files = append(files, f)
	}

	ff.wg.Wait()

	return
}

func (ff *ffind) enum(out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	if ff.sc {
		sc, err := ShadowCopy(windows.SystemDrive())

		if err != nil {
			sys.Fatal(err)
		}

		ff.root = sc
	}

	if len(ff.root) > 0 {
		fi, err := os.Stat(ff.root)

		if err != nil {
			sys.Fatal(err)
		}

		if !fi.IsDir() {
			sys.Fatal("not a directory")
		}
	}

	if !ff.uo {
		windows.EnumSystem(ff.root, out)
	}

	if !ff.so {
		windows.EnumUsers(ff.root, out)
	}
}

func (ff *ffind) comp(in <-chan string, out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	z, err := zip.NewZip(ff.zip, fmt.Sprintf("ffind %s", Version))

	if err != nil {
		sys.Fatal(err)
	}

	defer z.Close()

	for artifact := range in {
		err := z.Write(artifact, ff.path(artifact))

		if err != nil {
			sys.Error(err)
		}

		out <- artifact
	}
}

func (ff *ffind) list(in <-chan string, out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	f, err := os.Create(ff.csv)

	if err != nil {
		sys.Fatal(err)
	}

	defer f.Close()

	w := csv.NewWriter(f)

	err = w.Write(Header())

	if err != nil {
		sys.Fatal(err)
	}

	defer w.Flush()

	for artifact := range in {
		err = w.Write(Record(artifact))

		if err != nil {
			sys.Error(err)
		}

		out <- artifact
	}
}

func (ff *ffind) hash(in <-chan string, out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	for artifact := range in {
		s, err := hash.Sum(artifact, ff.hsh)

		if err != nil {
			sys.Error(err)
			continue
		}

		out <- fmt.Sprintf("%x%s%s", s, hash.Sep, ff.path(artifact))
	}
}

func (ff *ffind) live() {
	host, err := os.Hostname()

	if err != nil {
		sys.Error(err)

		host = fallback
	}

	ff.zip = host + zipExt
	ff.csv = host + csvExt

	ff.rp = true
	ff.sc = true
}

func (ff *ffind) path(f string) string {
	if !ff.rp {
		return f
	}

	r := len(ff.root)

	if r > 0 {
		r++
	}

	return f[r:]
}
