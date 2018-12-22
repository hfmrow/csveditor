// tabMakeSort.go

package main

import (
	"fmt"

	. "github.com/hfmrow/csveditor/genLib"

	"github.com/andlabs/ui"
)

func tabMakeSort(vbox *ui.Box) *ui.Box {
	if CsvProfileList.Initialised { // If there is nothing to display, skip it ...

		vbox.Enable() // Used to know if control has been already initialised

		// Create options Form
		entryForm := ui.NewForm()
		entryForm.SetPadded(true)
		vbox.Append(entryForm, false)

		// First part in a group
		gridSort := ui.NewGrid()
		gridSort.SetPadded(true)
		groupSort := ui.NewGroup("   C.S.V sorting:  " + TruncateString(CsvProfileList.FileName, "...", 80, 1))
		groupSort.SetMargined(true)
		groupSort.SetChild(gridSort)
		vbox.Append(groupSort, false)

		// Set COLUMN sort storage
		sortFields = make([]int, 0)

		// Add controls
		addColBtn := ui.NewButton("Add field")
		clrColBtn := ui.NewButton("Reset all")
		caseSvChk := ui.NewCheckbox("Case sensitive")
		doSortBtn := ui.NewButton("Sort")
		addColCbx := ui.NewCombobox()
		for _, col := range CsvProfileList.FieldNames {
			addColCbx.Append(col)
		}
		addColCbx.SetSelected(0)
		// Add OnClicked
		addColBtn.OnClicked(func(*ui.Button) { // Add button OnClicked
			var addOk bool = true
			for _, fieldIdx := range sortFields {
				if fieldIdx == addColCbx.Selected() {
					addOk = false
					break
				}
			}
			if addOk {
				sortFields = append(sortFields, addColCbx.Selected())
				refreshSortFields()
			}
		})
		// Clear OnClicked
		clrColBtn.OnClicked(func(*ui.Button) { // Clear button OnClicked
			resetSortFields()
		})
		// Case sensitive OnToggled
		caseSvChk.OnToggled(func(*ui.Checkbox) { caseSensitiveSortFlag = caseSvChk.Checked() })

		// Sort OnClicked
		doSortBtn.OnClicked(func(*ui.Button) { // Sort button OnClicked
			if len(sortFields) != 0 {
				// Get column names for the new table
				var tmpDatas [][]string
				tmpDatas = append(tmpDatas, CsvProfileList.FieldNames)
				// Duplicate table to make independant one
				sortedDatas := make([][]string, len(tableDatas))
				copy(sortedDatas, tableDatas)
				// Create temp filename to store result ...
				filename := tempDir + GenFileName() + ".tmp"
				for _, fieldIdx := range sortFields {
					switch CsvProfileList.FieldType[fieldIdx] {
					case TypE[0]: // Date
						sortedDatas = SliceSortDate(sortedDatas, CsvProfileList.DateFormat+" %H:%M:%S", fieldIdx, -1, sortAscDscCheckbox[fieldIdx].fChkbox.Checked())
					case TypE[1]: // String
						SliceSortString(sortedDatas, fieldIdx, sortAscDscCheckbox[fieldIdx].fChkbox.Checked(), caseSensitiveSortFlag)
					case TypE[2]: // Numeric
						SliceSortFloat(sortedDatas, fieldIdx, sortAscDscCheckbox[fieldIdx].fChkbox.Checked(), CsvProfileList.DecimalSep)
					default:
						fmt.Println("Sort Error: Bad field type !")
					}
				}
				CsvProfileList.Modified = true // Something changed, store information
				sortedDatas = append(tmpDatas, sortedDatas...)
				// Write it
				WriteCsv(filename, ",", sortedDatas)
				// Load, Store and set CSV
				addAndSetCsv(filename, "Sort")
				ChkRow.Reset() // Reset stored checked rows history
				// Refresh tabs
				refreshTable()
			}
		})

		// Display controls
		gridSort.Append(addColBtn,
			0, 0, 1, 1,
			true, ui.AlignEnd, false, ui.AlignCenter)

		gridSort.Append(addColCbx,
			1, 0, 1, 1,
			true, ui.AlignCenter, false, ui.AlignCenter)

		gridSort.Append(clrColBtn,
			2, 0, 1, 1,
			true, ui.AlignStart, false, ui.AlignCenter)

		gridSort.Append(caseSvChk,
			3, 0, 1, 1,
			true, ui.AlignStart, false, ui.AlignCenter)

		gridSort.Append(doSortBtn,
			4, 0, 1, 1,
			true, ui.AlignStart, false, ui.AlignCenter)

		sortBox = tabMakeSortFields(vbox)
	}
	return vbox
}

func resetSortFields() {
	if sortFieldsFlag {
		sortBox.Delete(3)
		sortFieldsFlag = false
		sortFields = sortFields[:0]
		sortAscDscCheckbox = sortAscDscCheckbox[:0]
	}
}

func refreshSort() {
	mainTab.Delete(2)
	mainTab.InsertAt("Sort", 2, tabSort())
}
