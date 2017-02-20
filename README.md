# TAME 
[![Build Status](https://travis-ci.org/mono83/tame.svg)](https://travis-ci.org/mono83/tame)
[![Go Report Card](https://goreportcard.com/badge/github.com/mono83/tame)](https://goreportcard.com/report/github.com/mono83/tame)
[![GoDoc](https://godoc.org/github.com/mono83/tame?status.svg)](https://godoc.org/github.com/mono83/tame)

Simple HTTP wrapper to retrieve arbitrary web pages. Usage:

```go

// Initialize User object
u := user.New()

// p contains html page with some getter methods
p, err := u.Get("https://google.com")
```