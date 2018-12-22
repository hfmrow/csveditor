// tabMakeSortFields.go

/// +build OMIT
package main

import (
	"strconv"

	"github.com/andlabs/ui"
)

func tabMakeSortFields(vbox *ui.Box) *ui.Box { // Second part, display selected fields
	// Be sure there is something to display
	if len(sortFields) != 0 {
		gridSortFields := ui.NewGrid()
		gridSortFields.SetPadded(true)
		groupSortFields := ui.NewGroup("   Selected field(s):  ")
		groupSortFields.SetMargined(true)
		groupSortFields.SetChild(gridSortFields)
		vbox.Append(groupSortFields, false)

		// Make slice of checkbox for sorting direction (ascending or descending)
		sortAscDscCheckbox = make([]ascDscCheckbox, len(CsvProfileList.FieldNames))
		for idx, _ := range CsvProfileList.FieldNames {
			sortAscDscCheckbox[idx].fChkbox = ui.NewCheckbox("")
		}

		// Display labels
		gridSortFields.Append(ui.NewLabel("Sort order"),
			0, 0, 1, 1,
			true, ui.AlignEnd, false, ui.AlignCenter)
		gridSortFields.Append(ui.NewLabel("Field name"),
			1, 0, 1, 1,
			true, ui.AlignCenter, false, ui.AlignCenter)
		gridSortFields.Append(ui.NewLabel("Field type"),
			2, 0, 1, 1,
			true, ui.AlignStart, false, ui.AlignCenter)
		gridSortFields.Append(ui.NewLabel("Ascending sort"),
			3, 0, 1, 1,
			true, ui.AlignStart, false, ui.AlignCenter)

		// Display columns
		for idx, fieldIdx := range sortFields {
			entry := ui.NewEntry()
			entry.SetText(CsvProfileList.FieldNames[fieldIdx])
			entry.SetReadOnly(true)

			gridSortFields.Append(ui.NewLabel(strconv.Itoa(idx+1)),
				0, idx+1, 1, 1,
				true, ui.AlignEnd, false, ui.AlignCenter)
			gridSortFields.Append(entry,
				1, idx+1, 1, 1,
				true, ui.AlignCenter, false, ui.AlignCenter)
			gridSortFields.Append(ui.NewLabel(CsvProfileList.FieldType[fieldIdx]),
				2, idx+1, 1, 1,
				true, ui.AlignStart, false, ui.AlignCenter)
			gridSortFields.Append(sortAscDscCheckbox[fieldIdx].fChkbox,
				3, idx+1, 1, 1,
				true, ui.AlignStart, false, ui.AlignCenter)
		}
		sortFieldsFlag = true
	}
	return vbox
}

func refreshSortFields() {
	if sortFieldsFlag {
		sortBox.Delete(3)
	}
	sortBox = tabMakeSortFields(sortBox)
}
