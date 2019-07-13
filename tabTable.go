// tabTable.go

package main

import (
	"fmt"

	"github.com/andlabs/ui"
	"github.com/hfmrow/csveditor/genLib"
)

func tabTable(filename string) ui.Control {
	// Read CSV table if the filename is provided ... if not, display blank tab
	if filename != "" {
		// Initialising storeProfiles
		storeProfiles = NewStoreProfiles()
		// Load, Store and set CSV
		addAndSetCsv(filename, "Main")
		// Refresh tabs
		refreshTable()
	}

	// Create TableSheet
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(false)
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(false)
	vbox.Append(hbox, false)

	// Open file button
	openFileButton := ui.NewButton("Open CSV file")
	hbox.Append(openFileButton, true)
	openFileButton.OnClicked(func(*ui.Button) {

		// Open fileOpen dialog
		filename := DialogBoxOpn(mainwin)
		if filename != "" {
			// Initialising storeProfiles
			storeProfiles = NewStoreProfiles()
			// Load, Store and set CSV
			addAndSetCsv(filename, "Main")
			ChkRow.Reset()         // Reset stored checked rows history
			sortFieldsFlag = false // Reset sort state
			// Refresh tabs
			refreshTable()
			refreshOption()
			refreshSort()
		}
	})
	// New file button
	newFileButton := ui.NewButton("New CSV file")
	hbox.Append(newFileButton, true)
	newFileButton.OnClicked(func(*ui.Button) {
		tableBox = vbox
		newTableWindow()
	})
	// Add/Edit button
	buttonAdd := ui.NewButton("Add/Edit row")
	hbox.Append(buttonAdd, true)
	buttonAdd.OnClicked(func(*ui.Button) {
		if CsvProfileList.Initialised { // If there is something to save, do it ...
			editRowWindow()
		}
	})
	// Delete button
	buttonDel := ui.NewButton("Delete row")
	hbox.Append(buttonDel, true)
	buttonDel.OnClicked(func(*ui.Button) {
		if CsvProfileList.Initialised { // If there is something to save, do it ...
			tableDeleteRow()
		}
	})

	// Search button
	buttonSearch := ui.NewButton("Search")
	hbox.Append(buttonSearch, true)
	buttonSearch.OnClicked(func(*ui.Button) {
		if CsvProfileList.Initialised { // If there is something to save, do it ...
			DialogBoxEntrySearch(mainwin, "Search in :  ", "Ok", "Cancel", func(entry string, cs, ww, rx bool) {
				filename := tempDir + genLib.GenFileName() + ".tmp" // Create temp file to store result ...
				foundDatas, err := genLib.SearchSl(entry, tableDatas, cs, false, false, rx, ww)
				if err == nil {
					var tmpDatas [][]string
					// Add column names to new table
					tmpDatas = append(tmpDatas, CsvProfileList.FieldNames)
					foundDatas = append(tmpDatas, foundDatas...)
					// Write it
					genLib.WriteCsv(filename, ",", foundDatas)
					// Load, Store and set CSV
					addAndSetCsv(filename, "Search")
					ChkRow.Reset() // Reset stored checked rows history
					// Refresh tabs
					refreshTable()
				} else {
					DialogBoxMsg(mainwin, "Information", fmt.Sprintf("%s", err))
				}
			}, func() {}, genLib.TruncateString(CsvProfileList.FileName, " ... ", 60, 2))
		}
	})

	// Save button
	saveButton := ui.NewButton("Save CSV file")
	hbox.Append(saveButton, true)
	saveButton.OnClicked(func(*ui.Button) {

		if CsvProfileList.Initialised { // If there is something to save, do it ...
			filename := DialogBoxSve(mainwin)
			if filename != "" {
				CsvProfileList.FileName = filename
				doCheckAndSaveOnExit(&CsvProfileList)
				mainwin.SetTitle(genLib.TruncateString(CsvProfileList.FileName, "...", 80, 1))
				refreshOption()
			}
		}
	})

	tableBox = makeTable(vbox)
	return tableBox
}

func refreshTable() {
	if tableBox != nil {
		tableBox.Delete(2)
		tableBox.Delete(1)
		tableBox = makeTable(tableBox)
	}
}

func makeTable(vbox *ui.Box) *ui.Box {
	mh = newModelHandler()
	model = ui.NewTableModel(mh)
	table := ui.NewTable(&ui.TableParams{
		Model:                         model,
		RowBackgroundColorModelColumn: -1,
	})

	vbox.Append(table, true)
	// Row zero for checkbox, one for row number
	table.AppendCheckboxTextColumn("", 0, ui.TableModelColumnAlwaysEditable, 1, ui.TableModelColumnNeverEditable, nil)
	// Append each text columns
	for idx, column := range CsvProfileList.FieldNames {
		if CsvProfileList.FieldDisplay[idx] == getValChar("y", 0).(string) { // "Yes" to display this column
			table.AppendTextColumn(column, idx+aditionalVisibleColumns, ui.TableModelColumnAlwaysEditable, nil)
		}
	}

	// Make Foot-Bar
	gridNavProf := ui.NewGrid()
	gridNavProf.SetPadded(true)
	vbox.Append(gridNavProf, false)

	onTheFlyEditCheckBox := ui.NewCheckbox("On the fly cell edition")
	onTheFlyEditCheckBox.OnToggled(func(*ui.Checkbox) { onTheFlyEdit = onTheFlyEditCheckBox.Checked() })
	titleLabel := ui.NewLabel(fmt.Sprint(appName, appVers, appCreat, copyRight))
	storedProfilesCbx := ui.NewCombobox()
	profFootBar = profFootBar[:0] // Reset tmp footBarStorage history
	for idx := 0; idx < len(storeProfiles); idx++ {
		if storeProfiles[idx].Initialised {
			profFootBar = append(profFootBar, storeProfiles[idx].Name)
			storedProfilesCbx.Append(storeProfiles[idx].Name)
			//			name := storeProfiles[idx].Name
			name := CsvProfileList.Name
			nameIdx := GetStoreProfilesPos(name)
			storedProfilesCbx.SetSelected(nameIdx)
		}
	}

	storedProfilesCbx.OnSelected(func(*ui.Combobox) {
		position := storedProfilesCbx.Selected()
		name := profFootBar[position]
		if name != CsvProfileList.Name {
			setProfile(name)
			storedProfilesCbx.SetSelected(GetStoreProfilesPos(name))
			refreshTable()
			//			refreshOption()
			//			refreshSort()
		}
	})
	// Display Foot-Bar
	gridNavProf.Append(onTheFlyEditCheckBox,
		0, 0, 1, 1,
		true, ui.AlignStart, false, ui.AlignCenter)

	gridNavProf.Append(titleLabel,
		0, 0, 1, 1,
		true, ui.AlignCenter, false, ui.AlignCenter)

	gridNavProf.Append(storedProfilesCbx,
		0, 0, 1, 1,
		true, ui.AlignEnd, false, ui.AlignCenter)

	return vbox
}
