package ffind

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/cuhsat/futils/pkg/sys"
	"github.com/cuhsat/futils/pkg/zip"
)

var (
	tmp, _  = os.MkdirTemp(os.TempDir(), "ffind")
	archive = filepath.Join(tmp, "ffind.zip")
	sysroot = filepath.Join(tmp, "sysroot")
)

func TestMain(m *testing.M) {
	sys.Progress = nil

	os.Exit(m.Run())
}

func TestFind(t *testing.T) {
	cases := []struct {
		name, file string
	}{
		{
			name: "Test find on a Windows system",
			file: filepath.Join("..", "testdata", "windows.zip"),
		},
	}

	for _, tt := range cases {
		if err := zip.Unzip(tt.file, sysroot); err != nil {
			t.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			a, err := zip.Index(tt.file)

			if err != nil {
				t.Fatal(err)
			}

			c := Find(sysroot, archive, "", "", true, false, false, false)

			b, err := zip.Index(archive)

			if err != nil {
				t.Fatal(err)
			}

			slices.Sort(c)

			compare(t, a, b)
			compare(t, a, c)
		})

		if _, err := os.Stat(sysroot); !os.IsNotExist(err) {
			os.RemoveAll(sysroot)
		}
	}
}

func BenchmarkFind(b *testing.B) {
	b.Run("Benchmark find", func(b *testing.B) {
		file := filepath.Join("..", "testdata", "windows.zip")

		if err := zip.Unzip(file, sysroot); err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			Find(sysroot, "", "", "", true, false, false, false)
		}

		b.StopTimer()

		if _, err := os.Stat(sysroot); !os.IsNotExist(err) {
			os.RemoveAll(sysroot)
		}
	})
}

func compare(t *testing.T, a, b []string) {
	if !reflect.DeepEqual(a, b) {
		for _, f := range a {
			t.Log(">>>", f)
		}

		for _, f := range b {
			t.Log("<<<", f)
		}

		t.Fatal("Lists are not equal")
	}
}
