package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const WINDOW_WIDTH = 800
const WINDOW_HEIGHT = 200
const DIALOG_HEIGHT = 600
const DIALOG_WIDTH = 800

type Gui struct {
	w      fyne.Window
	width  int
	height int
	opt    Options
	s      Cli

	filepathentry     *widget.Entry
	directorycheckbox *widget.Check
	overwritecheckbox *widget.Check
	okbutton          *widget.Button

	luri fyne.ListableURI
}

func (gui *Gui) Start() {
	a := app.New()
	gui.w = a.NewWindow("SMC Fix")
	content := gui.viewMain()
	gui.w.SetContent(content)
	gui.w.Resize(content.MinSize())
	gui.w.ShowAndRun()
}

func (gui *Gui) viewMain() fyne.CanvasObject {
	// |-------------------------------------------|
	// | [file path label]                         | row1
	// |-------------------------------------------|
	// | [file path entry] | [file picker button]  | row2
	// |-------------------------------------------|
	// | [dir check][overwrite check] | [ok button]| row3
	// |___________________________________________|

	// row1
	// [file path label]
	filepathlabel := widget.NewLabel("Choose the location of an .smc file…")
	row1 := container.NewMax(filepathlabel)

	// row2
	// [file path entry]
	gui.filepathentry = widget.NewEntry()
	gui.filepathentry.SetPlaceHolder("Folder or file path…")
	gui.filepathentry.OnChanged = func(s string) {
		if gui.filepathentry.Text != "" {
			gui.okbutton.Enable()
		} else {
			gui.okbutton.Disable()
		}
	}
	// [file picker button]
	filepickerbutton := widget.NewButtonWithIcon("Browse…", theme.FileIcon(), func() {
		dialogWindow := fyne.CurrentApp().NewWindow("Choose a .smc file…")
		dialogWindow.Resize(fyne.NewSize(DIALOG_WIDTH, DIALOG_HEIGHT))
		dialogWindow.SetFixedSize(true)
		file_Dialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if file != nil {
				filepath := file.URI().Path()
				gui.filepathentry.SetText(filepath)
			}
		}, dialogWindow)
		file_Dialog.Resize(fyne.NewSize(DIALOG_WIDTH, DIALOG_HEIGHT))
		// fiter to open .smc files only
		file_Dialog.SetFilter(storage.NewExtensionFileFilter([]string{".smc"}))
		file_Dialog.Show()
		if gui.filepathentry.Text != "" {
			uri := storage.NewFileURI(filepath.Dir(gui.filepathentry.Text))
			luri, err := storage.ListerForURI(uri)
			if err != nil {
				log.Fatal(err)
			}
			file_Dialog.SetLocation(luri)
		}
		file_Dialog.SetOnClosed(func() {
			dialogWindow.Hide()
		})
		dialogWindow.Show()
		// Show file selection dialog.
	})

	row2 := container.NewBorder(nil, nil, nil, filepickerbutton, gui.filepathentry)

	// row3
	// [overwrite checkbox]
	gui.directorycheckbox = widget.NewCheck("Update all in directory", nil)
	gui.overwritecheckbox = widget.NewCheck("Overwrite files", nil)
	gui.okbutton = widget.NewButton("Strip Headers", func() {
		path := gui.filepathentry.Text
		if isValidPath(path) {
			gui.s.Reset()
			if gui.directorycheckbox.Checked {
				// dir
				gui.s.CleanFolder(filepath.Dir(path), filepath.Dir(path), gui.overwritecheckbox.Checked)
			} else {
				// file
				gui.s.CleanFile(path, filepath.Dir(path), gui.overwritecheckbox.Checked)
			}
			resultsWindow := fyne.CurrentApp().NewWindow("Results")
			resultsWindow.Resize(fyne.NewSize(300, 300))
			resultsWindow.SetFixedSize(true)
			c := container.NewVBox(
				widget.NewLabel("Total files: "+strconv.Itoa(gui.s.Results.Total)),
				widget.NewLabel("Already good: "+strconv.Itoa(gui.s.Results.Done)),
				widget.NewLabel("Files updated: "+strconv.Itoa(gui.s.Results.Updated)),
				widget.NewLabel("Failed: "+strconv.Itoa(gui.s.Results.Failed)))
			resultsWindow.Resize(c.MinSize())
			resultsWindow.SetFixedSize(true)
			resultsWindow.SetContent(c)
			resultsWindow.Show()
		}
	})
	gui.okbutton.Disable()
	row3 := container.NewHBox(gui.directorycheckbox, gui.overwritecheckbox, layout.NewSpacer(), gui.okbutton)

	c := container.NewVBox(row1, row2, row3)

	pad := container.NewPadded(c)
	return pad
}

func isValidPath(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		// does not exist
		return false
	}
	if stat.IsDir() || filepath.Ext(path) == ".smc" {
		// not a dir or smc
		return true
	}
	return false
}
