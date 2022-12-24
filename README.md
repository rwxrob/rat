# Functional PEG packrat parser in Go

![WIP](https://img.shields.io/badge/status-wip-red.svg)
[![GoDoc](https://godoc.org/github.com/rwxrob/rat?status.svg)](https://godoc.org/github.com/rwxrob/rat)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rwxrob/rat)](https://goreportcard.com/report/github.com/rwxrob/rat)
[![Coverage](https://gocover.io/_badge/github.com/rwxrob/rat)](https://gocover.io/github.com/rwxrob/rat)

Inspired by Bryan Ford's PEG packrat parser paper (say that 10 times fast) and Bruce Hill's great overview of how to create one from scratch. Quint and I created [one for PEGN](https://github.com/rwxrob/pegn-go) without realizing it at the time (but we got carried away, the type switching on `structs` and pseudo-code where a bit much). This is a happy medium between maintainability, simplicity, and performance that most developers can get started with right away.

One particularly amazing thing about Go is the slice abstraction allows a reference to the underlying array to be embedded into *every* `rat.Result` without any cost beyond a single pointer reference. Go also has the built-in concept of `[]rune` slices that automatically distinguishes between unicode code points of different byte size.

* Packrat Parsing and Parsing Expression Grammars  
  <https://bford.info/packrat/>
* Packrat Parsing from Scratch -- Naming Things  
  <https://blog.bruce-hill.com/packrat-parsing-from-scratch>
