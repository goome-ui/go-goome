package goome

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestLayoutMatchesC(t *testing.T) {
	lines := strings.Split(strings.TrimSpace(runC(t, "testdata/layout.c")), "\n")
	sort.Strings(lines)
	want := strings.Join(lines, "\n")
	got := layoutReport()
	if got != want {
		g, w := strings.Split(got, "\n"), lines
		seen := map[string]bool{}
		for _, l := range g {
			seen[l] = true
		}
		for _, l := range w {
			if !seen[l] {
				t.Errorf("C has %q, go doesnt", l)
			}
		}
		have := map[string]bool{}
		for _, l := range w {
			have[l] = true
		}
		for _, l := range g {
			if !have[l] {
				t.Errorf("go has %q, C doesnt", l)
			}
		}
	}
}

func TestEveryConstant(t *testing.T) {
	src, err := os.ReadFile("native/goome.h")
	if err != nil {
		t.Fatal(err)
	}
	head := strings.SplitN(string(src), "#ifdef GOOME_IMPLEMENTATION", 2)[0]
	comments := regexp.MustCompile(`//[^\n]*`)
	names := regexp.MustCompile(`(?m)^\s*(GOOME_[A-Z0-9_]+)`)
	n := 0
	for _, body := range regexp.MustCompile(`(?s)enum\s*\{(.*?)\}`).FindAllStringSubmatch(head, -1) {
		b := strings.ReplaceAll(comments.ReplaceAllString(body[1], ""), ",", "\n")
		for _, m := range names.FindAllStringSubmatch(b, -1) {
			n++
			if _, ok := cConsts[m[1]]; !ok {
				t.Errorf("%s has no go constant", m[1])
			}
		}
	}
	if n < 200 || len(cConsts) < n {
		t.Errorf("found %d enum names, have %d constants", n, len(cConsts))
	}
	if KeyLeft != 0x90 || ModAlt != 0x400 || BackendTest != 7 || CodeArg != -7 || IconCount != 22 {
		t.Error("constants dont match goome.h")
	}
	if Version == "" {
		t.Error("no version")
	}
}
