<h1 align="center">🛠 Contributing to go-goome 🛠</h1>

<p align="center">
Thanks for wanting to help! Small PRs that fix one thing are the easiest for me to check and merge, so keep them small and theyll get in way faster.<br>
This repo is only the Go side. If the bug is in goome itself (drawing, a backend, the ui kit) it goes to <a href="https://github.com/goome-ui/c-goome">c-goome</a> instead, docs fixes go to <a href="https://github.com/goome-ui/docs">goome-ui/docs</a> (thats <a href="https://goome-ui.github.io/docs/">goome-ui.github.io/docs</a>) and new example projects to <a href="https://github.com/goome-ui/examples">goome-ui/examples</a>.
</p>

<details>
<summary align="center"><b>📋 The short version 📋</b></summary>

<p align="center">
• One thing per PR. One fix or one feature, if it does two things it gets closed<br>
• Open your PR against the <b>dev</b> branch, not main<br>
• Never touch <code>native/</code>, goome.h in there is copied from c-goome<br>
• Dont touch the protected files, a bot closes the PR if you do<br>
• Test it and say how you tested it, and on which OS<br>
• CI has to pass<br>
• AI is fine, but say you used it, and you need to know what every line does
</p>
</details>

<details>
<summary align="center"><b>💬 Before you start 💬</b></summary>

<p align="center">
Small fixes can go straight to a PR.<br>
For a new feature or anything big, <a href="https://github.com/goome-ui/go-goome/issues">open an issue</a> first so we can talk about it before you spend a weekend on it.<br> <br>
Same if you want to rename something or change what a function takes or returns, people already use this so I want to talk that through first.
</p>
</details>

<details>
<summary align="center"><b>🌿 Branches 🌿</b></summary>

<p align="center">
Stuff goes dev, then beta, then main. I move it along, you only ever PR into <b>dev</b>.
</p>
</details>

<details>
<summary align="center"><b>🚫 What gets closed automatically 🚫</b></summary>

<p align="center">
A bot checks every PR and closes it with a comment if it hits any of these, fix it and open a new one.<br> <br>
It has to be against <b>dev</b>, main and beta only get stuff that already sat on dev for a while.<br> <br>
It cant touch <code>.github/</code>, since thats CI and these rules and a PR could just turn them off for itself. <code>native/</code> is off limits too, <code>goome.h</code> in there gets copied from c-goome by <code>tools/update_native.sh</code> and cgo compiles it on everyones computer. Same for <code>go.mod</code>, <code>go.sum</code>, <code>goome_c.go</code>, <code>helpers.h</code> and <code>tools/</code>, thats the cgo build and what runs when people build it. The <code>LICENSE</code> and the logo stay mine to change.<br> <br>
• <b>More than 3000 changed files</b> - way too big to check<br> <br>
Renaming or moving any of those files counts too. Think one of them is wrong? <a href="https://github.com/goome-ui/go-goome/issues">Open an issue</a> and I'll fix it.
</p>
</details>

<details>
<summary align="center"><b>🤖 Using AI 🤖</b></summary>

<p align="center">
Using AI to help write your PR is allowed, as long as you say so in the PR and you know what your code does. If you get asked about something you should be able to answer without going back to the AI.<br> <br>
And test it yourself. Untested AI dumps get closed
</p>
</details>

<details>
<summary align="center"><b>🧪 Testing 🧪</b></summary>

<p align="center">
Run the tests. They all run on goome's test backend, so nothing shows up on your screen:
</p>

```sh
go vet ./...
go test ./...
gofmt -l .
```

<p align="center">
You need a C compiler for cgo, on Windows thats mingw gcc on your PATH with <code>CGO_ENABLED=1</code>.<br> <br>
<code>layout_test.go</code> compiles <code>testdata/layout.c</code> and checks every struct the Go side mirrors (<code>Theme</code>, <code>UIConfig</code>, <code>Radii</code>, <code>Skin</code>, <code>Combo</code>) plus the cgo ones against what C says, and that every constant in goome.h has a Go name.<br> <br>
The drawing and ui tests check the same pixel hashes as the <code>tests/unit.c</code> in c-goome, and <code>TestFrameMatchesC</code> builds <code>testdata/frame.c</code> and checks an offscreen frame hashes the same from C and from Go. Those come from the C code, so if one moves something in the bindings is passing the wrong thing. Dont just update the number.<br> <br>
CI also runs the tests with <code>GOEXPERIMENT=cgocheck2</code>, which catches Go pointers getting kept by C. If you add a function that hands C a pointer it holds on to, put that memory in C (the ui uses its own arena for this).<br> <br>
In the PR write a few lines on what you ran and on which OS and Go version. If it changes something you can see add a before and after screenshot.<br> <br>
CI runs on linux (plus Go 1.23), windows and both macs when you open the PR. It has to be green before I merge.
</p>
</details>

<details>
<summary align="center"><b>✏️ Code style ✏️</b></summary>

<p align="center">
• Go 1.23 and up, nothing from a newer stdlib<br>
• No dependencies. Only the standard library and cgo<br>
• gofmt, and <code>go vet</code> has to be clean<br>
• Every function in goome.h stays reachable. New enum values come from <code>go run ./tools/genconsts</code>, new struct fields go in the Go struct and <code>testdata/layout.c</code><br>
• Errors come back as values (<code>*goome.Error</code>), dont panic for things that can just fail. Using something after <code>Free</code> or <code>Close</code> is the only panic<br>
• Anything that holds C memory gets a <code>Free()</code> and freeing twice is fine<br>
• Anything with a begin and an end in C gets closed by <code>ui.End()</code>, and calls the C end exactly when <a href="https://goome-ui.github.io/docs/ui-app/">the docs</a> say C wants it<br>
• Match the code around it. Short names are fine inside functions<br>
• No comments unless something really isnt obvious<br>
• Dont reformat whole files, it makes the diff impossible to read
</p>
</details>

<details>
<summary align="center"><b>💬 Discussions and first issues 💬</b></summary>
<p align="center">
Questions and stuff you made with goome go in <a href="https://github.com/goome-ui/c-goome/discussions">Discussions</a>.<br>
If you want something to work on, the <a href="https://github.com/goome-ui/go-goome/labels/good%20first%20issue">good first issue</a> ones are small. Comment on one before you start.
</p>
</details>

<p align="center">
Found a security hole? Dont open a public issue or PR for it, read <a href="https://github.com/goome-ui/.github/blob/main/SECURITY.md">SECURITY.md</a> and tell me privately first.<br> <br>
By sending a PR you agree your code goes out under the same <a href="LICENSE">GPL-3.0 license</a> as the rest of go-goome.
</p>
