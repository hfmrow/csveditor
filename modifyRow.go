// modifyRow.go

package main

import "github.com/andlabs/ui"

func tableAddEditRow(spinBoxValue int, flagEdit, dupFlag bool, datas []string) {
	spinBoxValue-- // Adjust value for slice position
	if flagEdit {
		if dupFlag {
			// Duplicate Row
			tableDupInsertRow(ChkRow.Last(), datas)
			mh.checkBoxValue[ChkRow.Last()+1] = 0
			mh.checkBoxValue[ChkRow.Last()] = 1 // Adjust checkbox to last entry
		} else {
			// Edit row
			for idx, field := range datas {
				tableDatas[ChkRow.Last()][idx] = field
				CsvProfileList.Modified = true
				model.RowChanged(ChkRow.Last())
			}
		}
	} else {
		// Insert Row
		tableDupInsertRow(spinBoxValue, datas)
	}
}

func tableDupInsertRow(pos int, datas []string) {
	var tmpValue []string                      // create new empty row
	tableDatas = append(tableDatas, tmpValue)  // Add it to main table
	copy(tableDatas[pos+1:], tableDatas[pos:]) // Duplicate given entry
	tableDatas[pos] = datas                    // Insert row at given position
	// Insert CheckBox state
	var tmpInt int
	mh.checkBoxValue = append(mh.checkBoxValue, tmpInt)    // Add it
	copy(mh.checkBoxValue[pos+1:], mh.checkBoxValue[pos:]) // Duplicate given entry
	mh.checkBoxValue[pos] = 0                              // Insert 0 (for false) at given position

	CsvProfileList.NumberRows++
	CsvProfileList.Modified = true
	model.RowInserted(pos)
}

func tableDeleteRow() {
	for row := len(mh.checkBoxValue) - 1; row >= 0; row-- {
		if mh.checkBoxValue[row] == int(ui.TableTrue) {
			CsvProfileList.NumberRows--
			mh.checkBoxValue = append(mh.checkBoxValue[:row], mh.checkBoxValue[row+1:]...)
			tableDatas = append(tableDatas[:row], tableDatas[row+1:]...)
			CsvProfileList.Modified = true
			model.RowDeleted(row)
			ChkRow.Reset() // Reset stored checked rows history

		}
	}
}
