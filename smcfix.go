package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	// "io/ioutil"
	"io/fs"
	"path/filepath"
)

const HEADER_CHECK_SIZE = 200
const HEADER_SIZE = 512
const SUFFIX = "-smcfix"

type options struct {
	dir       string
	file      string
	out       string
	overwrite bool
	help      bool
}

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	var opt options

	flag.StringVar(&opt.dir, "dir", workingDir, "Directory to scan for SMC files.")
	flag.StringVar(&opt.file, "file", "", "Single SMC file to check and clean.")
	flag.StringVar(&opt.out, "out", "", "Specify alternate output directory.")
	flag.BoolVar(&opt.overwrite, "overwrite", false, "Overwrite or create new e.g. \"[filename]-smcfix.smc\" (default false)")
	flag.BoolVar(&opt.help, "help", false, "Show this help.")
	flag.Parse()

	// Usage demo
	opt.validateFlags()
	opt.validateFile()
	opt.updateOut()

	if opt.file != "" {
		cleanFile(opt.file, opt.out, opt.overwrite)
	} else {
		cleanFolder(opt.dir, opt.out, opt.overwrite)
	}
}

func (opt *options) validateFile() {
	ext := filepath.Ext(opt.file)
	if ext != ".smc" {
		fmt.Println("Error: Not a .smc file.")
		os.Exit(0)
	}
}

func (opt *options) updateOut() {
	if isFlagPassed("out") {
		return
	}
	if isFlagPassed("file") {
		opt.out = filepath.Dir(opt.file)
	} else {
		opt.out = opt.dir
	}
}

func (opt *options) validateFlags() {
	if opt.help {
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

func cleanFolder(path string, outDir string, overwrite bool) {
	files := find(path, ".smc")
	fmt.Printf("Found %d .smc files in %s\n", len(files), path)
	for _, element := range files {
		cleanFile(element, outDir, overwrite)
	}
}

func cleanFile(path string, outDir string, overwrite bool) {
	fin, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fin.Close()
	filename := filepath.Base(path)

	fmt.Print(filename)
	headerSize, err := headerSize(fin)
	if err != nil {
		fmt.Print(" 💀 Failed.\n")
		return
	}
	if headerSize > 0 { // has smc header
		// create output file name
		outfile := createOutFilepath(path, outDir)
		outfile = outfile[:len(outfile)-4]
		if !overwrite {
			outfile = outfile + SUFFIX
		}
		outfile = outfile + ".smc"
		// check if it already exists.
		if fileExists(outfile) {
			fmt.Printf(" ♻️ File Exists: %s \n", outfile)
			return
		}

		func() {
			fout, err := os.Create(outfile)
			if err != nil {
				fmt.Print("\nError: ")
				log.Fatal(err)
			}
			defer fout.Close()

			// Offset is the number of bytes you want to exclude
			_, err = fin.Seek(headerSize, io.SeekStart)
			if err != nil {
				fmt.Print("\nError: ")
				log.Fatal(err)
			}

			_, err = io.Copy(fout, fin)
			if err != nil {
				fmt.Printf("\nError: Couldn't write %s\n", outfile)
				log.Fatal(err)
			}

			if overwrite {
				moveFile(path, outfile)
				outfile = path
			}
			fmt.Printf(" 🧼 " + outfile + " ✅\n")
		}()
	} else {
		fmt.Printf(" ✅\n")
	}
}

func moveFile(src string, dst string) {
	if err := os.Remove(src); err != nil {
		fmt.Printf("\nError: Couldn't overwrite %s with %s\n", src, dst)
		log.Fatal(err)
	}
	if err := os.Rename(dst, src); err != nil {
		fmt.Printf("\nError: Couldn't overwrite %s with %s\n", src, dst)
		log.Fatal(err)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return false
}

func createOutFilepath(infile string, outdir string) string {
	filename := filepath.Base(infile)
	// does outdir exist
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		log.Fatal(err)
	}
	return filepath.Join(outdir, filename)
}

func headerSize(file *os.File) (int64, error) {
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size() % 1024, nil
}

func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext && !strings.Contains(s, SUFFIX) {
			a = append(a, s)
		}
		return nil
	})
	return a
}
