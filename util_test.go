package goome

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func fnv(px []uint32, w, h, stride int) uint32 {
	v := uint32(2166136261)
	for y := 0; y < h; y++ {
		for _, p := range px[y*stride : y*stride+w] {
			for k := 0; k < 4; k++ {
				v = (v ^ (p >> (k * 8) & 255)) * 16777619
			}
		}
	}
	return v
}

func data(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func cc() string {
	if c := os.Getenv("CC"); c != "" {
		return c
	}
	out, err := exec.Command("go", "env", "CC").Output()
	if err == nil {
		if c := strings.Fields(string(out)); len(c) > 0 {
			return c[0]
		}
	}
	return "cc"
}

func runC(t *testing.T, src string) string {
	t.Helper()
	dir := t.TempDir()
	exe := filepath.Join(dir, "prog")
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	out, err := exec.Command(cc(), "-std=c99", "-Inative", "-o", exe, src).CombinedOutput()
	if err != nil {
		t.Fatalf("cant build %s: %v\n%s", src, err, out)
	}
	res, err := exec.Command(exe).Output()
	if err != nil {
		t.Fatalf("%s failed: %v", src, err)
	}
	return strings.ReplaceAll(string(res), "\r\n", "\n")
}
