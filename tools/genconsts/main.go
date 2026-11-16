package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"strings"
)

var events = strings.Fields("NONE DOWN UP MOTION SCROLL KEY HOTKEY CONTROL MONITORS QUIT CLOSE FOCUS DROP TEXT PREEDIT")
var backends = strings.Fields("AUTO WAYLAND GNOME PORTAL X11 WIN32 MACOS TEST OFFSCREEN")
var codes = map[string]string{"OK": "OK", "ENOTSUP": "CodeNotSupported", "EDENIED": "CodeDenied", "ENODISPLAY": "CodeNoDisplay", "EGONE": "CodeGone", "EAGAIN": "CodeAgain", "ENOMEM": "CodeNoMem", "EARG": "CodeArg"}
var mods = map[string]string{"SHIFT": "Shift", "CTRL": "Ctrl", "ALT": "Alt", "SUPER": "Super"}
var upper = map[string]string{"UI": "UI", "FPS": "FPS", "ID": "ID", "RGBA": "RGBA", "BGRA": "BGRA", "ARGB": "ARGB", "OK": "OK", "EW": "EW", "NS": "NS", "NWSE": "NWSE", "NESW": "NESW", "PGUP": "PgUp", "PGDN": "PgDn", "IBEAM": "IBeam", "MACOS": "MacOS", "LTR": "LTR", "RTL": "RTL", "MENUITEM": "MenuItem"}
var fkey = regexp.MustCompile(`^[A-Z]\d+$`)

func camel(s string) string {
	var b strings.Builder
	for _, p := range strings.Split(s, "_") {
		if u, ok := upper[p]; ok {
			b.WriteString(u)
		} else if fkey.MatchString(p) {
			b.WriteString(p)
		} else if p != "" {
			b.WriteString(p[:1] + strings.ToLower(p[1:]))
		}
	}
	return b.String()
}

func has(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func goName(n string) string {
	s := strings.TrimPrefix(n, "GOOME_")
	switch {
	case has(events, s):
		return "Event" + camel(s)
	case has(backends, s):
		return "Backend" + camel(s)
	case codes[s] != "":
		return codes[s]
	case mods[s] != "":
		return "Mod" + mods[s]
	case strings.HasPrefix(s, "K_"):
		return "Key" + camel(s[2:])
	case strings.HasPrefix(s, "S_"):
		return "Style" + camel(s[2:])
	}
	return camel(s)
}

func main() {
	src, err := os.ReadFile("native/goome.h")
	if err != nil {
		fmt.Fprintln(os.Stderr, "run it from the repo root:", err)
		os.Exit(1)
	}
	head := strings.SplitN(string(src), "#ifdef GOOME_IMPLEMENTATION", 2)[0]
	comments := regexp.MustCompile(`//[^\n]*`)
	names := regexp.MustCompile(`(?m)^\s*(GOOME_[A-Z0-9_]+)`)
	var all []string
	seen := map[string]bool{}
	for _, body := range regexp.MustCompile(`(?s)enum\s*\{(.*?)\}`).FindAllStringSubmatch(head, -1) {
		b := strings.ReplaceAll(comments.ReplaceAllString(body[1], ""), ",", "\n")
		for _, m := range names.FindAllStringSubmatch(b, -1) {
			if !seen[m[1]] {
				seen[m[1]] = true
				all = append(all, m[1])
			}
		}
	}
	var out bytes.Buffer
	out.WriteString("package goome\n\n/*\n#include \"native/goome.h\"\n*/\nimport \"C\"\n\nconst (\n")
	used := map[string]string{}
	for _, n := range all {
		g := goName(n)
		if other, ok := used[g]; ok {
			fmt.Fprintf(os.Stderr, "%s and %s both turn into %s\n", other, n, g)
			os.Exit(1)
		}
		used[g] = n
		fmt.Fprintf(&out, "\t%s = C.%s\n", g, n)
	}
	out.WriteString(")\n\nvar cConsts = map[string]int{\n")
	for _, n := range all {
		fmt.Fprintf(&out, "\t%q: %s,\n", n, goName(n))
	}
	out.WriteString("}\n")
	code, err := format.Source(out.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile("consts.go", code, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("consts.go has %d constants\n", len(all))
}
