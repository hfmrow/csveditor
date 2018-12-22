// main.go

package main

import (
	"fmt"
	"os"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"

	. "github.com/hfmrow/csveditor/genLib"
)

// Main
func main() {
	// Get command line argument if exist.
	if len(os.Args) > 1 {
		argFilename = os.Args[1]
	}

	// Create TempDir
	tempDir = TempMake(ReplaceSpace(appName)) + PathSep() // Init Tempdir
	defer TempRemove(tempDir)

	// Start UI ...
	ChkRow = NewCheckedRow() // Initialise checkedRows
	ui.Main(setupUI)         // Launching ui
}

// Setup GUI
func setupUI() {
	mainwin = ui.NewWindow(fmt.Sprintln(appName, appVers, appCreat, copyRight), 800, 600, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		returnFlag := false
		if CsvProfileList.Initialised { // Check for existing profile
			if CsvProfileList.Modified { // Check for modified profile ...
				mainwin.Disable()
				DialogBoxAsk3(mainwin, "File modified.", "Yes", "No", "Cancel",
					func() { // Yes button pushed, save profile and CSV before exiting
						doCheckAndSaveOnExit(&CsvProfileList)
						ui.Quit()
						returnFlag = true
					}, func() { // No button pushed, simply quit
						ui.Quit()
						returnFlag = true
					}, func() { // Cancel button pushed, don't exit and return to main window
						mainwin.Enable()
						returnFlag = false
					}, "Save file ?", TruncateString(CsvProfileList.FileName, "...", 80, 2))
			} else { // Profile not modified, exit
				ui.Quit()
				returnFlag = true
			}
		} else { // Profile not initialised, exit
			ui.Quit()
			returnFlag = true
		}
		return returnFlag
	})

	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})
	// Make Tabs
	mainTab = ui.NewTab()
	mainwin.SetChild(mainTab)
	mainwin.SetMargined(true)

	mainTab.InsertAt("Table", 0, tabTable(argFilename))
	mainTab.InsertAt("Options", 1, tabOptions())
	mainTab.InsertAt("Sort", 2, tabSort())
	mainTab.SetMargined(0, true)

	mainwin.Show()
}
