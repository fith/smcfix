package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	var opt Options

	flag.StringVar(&opt.Dir, "dir", workingDir, "Directory to scan for SMC files.")
	flag.StringVar(&opt.File, "file", "", "Single SMC file to check and clean.")
	flag.StringVar(&opt.Out, "out", "", "Specify alternate output directory.")
	flag.BoolVar(&opt.Overwrite, "overwrite", false, "Overwrite or create new e.g. \"[filename]-smcfix.smc\" (default false)")
	flag.BoolVar(&opt.Help, "help", false, "Show this help.")
	flag.Parse()

	if opt.Count() == 0 {
		gui()
	} else {
		cli(opt)
	}
}

func cli(opt Options) {
	// Usage demo
	opt.ValidateFlags()
	opt.ValidateFile()
	opt.UpdateOut()

	var s Cli

	if opt.File != "" {
		s.CleanFile(opt.File, opt.Out, opt.Overwrite)
	} else {
		s.CleanFolder(opt.Dir, opt.Out, opt.Overwrite)
	}
	fmt.Println("Results\n")
	fmt.Printf("Checked: \t%d\n", s.Results.Total)
	fmt.Printf("Updated: \t%d\n", s.Results.Updated)
	fmt.Printf("Failed: \t%d\n", s.Results.Failed)
}

func gui() {
	var gui Gui
	gui.Start()
}
