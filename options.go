package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Options struct {
	Dir       string
	File      string
	Out       string
	Overwrite bool
	Help      bool
}

func (opt *Options) ValidateFile() {
	ext := filepath.Ext(opt.File)
	if ext != ".smc" {
		fmt.Println("Error: Not a .smc file.")
		os.Exit(0)
	}
}

func (opt *Options) UpdateOut() {
	if isFlagPassed("out") {
		return
	}
	if isFlagPassed("file") {
		opt.Out = filepath.Dir(opt.File)
	} else {
		opt.Out = opt.Dir
	}
}

func (opt *Options) ValidateFlags() {
	if opt.Help {
		flag.Usage()
		os.Exit(0)
	}
	if !isFlagPassed("dir") && !isFlagPassed("file") {
		fmt.Println("Error: No target to clean. Specify a directory or file.")
		flag.Usage()
		os.Exit(0)
	}
	if isFlagPassed("dir") && isFlagPassed("file") {
		fmt.Println("Error: Specify a dir OR a file. Can't do both.")
		flag.Usage()
		os.Exit(0)
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func (opt *Options) Count() int {
	return flag.NFlag()
}
