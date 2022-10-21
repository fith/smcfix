package main

import (
	"flag"
	"log"
	"os"
	"smcfix/smcfixmod"
	"smcfix/smcfixopt"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	var opt smcfixopt.Smcfixopt

	flag.StringVar(&opt.Dir, "dir", workingDir, "Directory to scan for SMC files.")
	flag.StringVar(&opt.File, "file", "", "Single SMC file to check and clean.")
	flag.StringVar(&opt.Out, "out", "", "Specify alternate output directory.")
	flag.BoolVar(&opt.Overwrite, "overwrite", false, "Overwrite or create new e.g. \"[filename]-smcfix.smc\" (default false)")
	flag.BoolVar(&opt.Help, "help", false, "Show this help.")
	flag.Parse()

	cli(opt)
}

func cli(opt smcfixopt.Smcfixopt) {
	// Usage demo
	opt.ValidateFlags()
	opt.ValidateFile()
	opt.UpdateOut()

	if opt.File != "" {
		smcfixmod.CleanFile(opt.File, opt.Out, opt.Overwrite)
	} else {
		smcfixmod.CleanFolder(opt.Dir, opt.Out, opt.Overwrite)
	}
}
