package smcfixcli

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const SUFFIX = "-smcfix"

func CleanFolder(path string, outDir string, overwrite bool) {
	files := find(path, ".smc")
	fmt.Printf("Found %d .smc files in %s\n", len(files), path)
	for _, element := range files {
		CleanFile(element, outDir, overwrite)
	}
}

func CleanFile(path string, outDir string, overwrite bool) {
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
