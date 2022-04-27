package main

import (
	"github.com/IAmRadek/importas"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(importas.Analyzer)
}
