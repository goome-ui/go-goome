package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/goome-ui/go-goome"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: grab [monitor] [x y w h] out.ppm")
		os.Exit(2)
	}
	var nums []int
	for _, a := range args[:len(args)-1] {
		if len(nums) == 5 {
			break
		}
		n, _ := strconv.Atoi(a)
		nums = append(nums, n)
	}
	g, err := goome.Open(goome.Options{Title: "goome grab"})
	if err != nil {
		fmt.Fprintln(os.Stderr, "cant open:", err)
		os.Exit(1)
	}
	defer g.Close()
	mon := 0
	if len(nums) == 1 || len(nums) == 5 {
		mon = nums[0]
	}
	var r *goome.Rect
	if len(nums) >= 4 {
		n := nums[len(nums)-4:]
		r = &goome.Rect{X: n[0], Y: n[1], W: n[2], H: n[3]}
	}
	img, err := g.CaptureWait(mon, r)
	if err == nil {
		err = img.SavePPM(args[len(args)-1])
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "grab failed: %v\n", err)
		os.Exit(1)
	}
	a := img.Area()
	extra := ""
	if img.HasSelf() {
		extra += ", has our overlay"
	}
	if img.HasCursor() {
		extra += ", has the cursor"
	}
	fmt.Printf("%dx%d px of %d,%d %dx%d on monitor %d%s\n", img.W(), img.H(), a.X, a.Y, a.W, a.H, mon, extra)
}
