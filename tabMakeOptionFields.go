// tabMakeOptionsFields.go

package main

import (
	"github.com/andlabs/ui"
	. "github.com/hfmrow/csveditor/genLib"
)

func tabMakeOptionFields(vbox *ui.Box) *ui.Box {
	gridFields := ui.NewGrid()
	gridFields.SetPadded(true)
	groupFields := ui.NewGroup("")
	groupFields.SetMargined(true)
	groupFields.SetTitle("   Fields format:")
	groupFields.SetChild(gridFields)
	vbox.Append(groupFields, false)
	// If there is nothing to display, skip it ...
	//	if CsvProfileList.Initialised {

	// Add field button with spinbox
	fieldNewButton := ui.NewButton("Add field")
	fieldSpinBox := ui.NewSpinbox(1, CsvProfileList.NumberCols+1)
	fieldSpinBox.SetValue(lastSpinboxOptionFieldsValue)
	fieldNewButton.OnClicked(func(*ui.Button) {
		lastSpinboxOptionFieldsValue = fieldSpinBox.Value()
		CsvProfileList.FieldNames = AppendAt(CsvProfileList.FieldNames, fieldSpinBox.Value()-1, "New field")
		CsvProfileList.FieldDisplay = AppendAt(CsvProfileList.FieldDisplay, fieldSpinBox.Value()-1, getValChar("y", 0).(string)) // "Yes"
		CsvProfileList.FieldOut = AppendAt(CsvProfileList.FieldOut, fieldSpinBox.Value()-1, getValChar("y", 0).(string))         // "Yes"
		CsvProfileList.FieldType = AppendAt(CsvProfileList.FieldType, fieldSpinBox.Value()-1, getValChar("t", 1).(string))       // "String"
		CsvProfileList.NumberCols++
		CsvProfileList.Modified = true
		for row := 0; row < CsvProfileList.NumberRows; row++ {
			tableDatas[row] = AppendAt(tableDatas[row], lastSpinboxOptionFieldsValue-1, "")
		}
		// Reset storeprofile and row checkbox selection
		reloadMainStoreProfile()

		// Refresh tabs
		refreshOptionFields()
	})
	gridFields.Append(fieldNewButton,
		0, 0, 1, 1,
		true, ui.AlignCenter, false, ui.AlignCenter)
	gridFields.Append(ui.NewLabel("Pos :"),
		0, 0, 1, 1,
		true, ui.AlignEnd, false, ui.AlignCenter)
	gridFields.Append(fieldSpinBox,
		1, 0, 1, 1,
		true, ui.AlignStart, false, ui.AlignCenter)

	// Remove field button with spinbox
	fieldDelButton := ui.NewButton("Remove field")
	fieldDelSpinBox := ui.NewSpinbox(1, CsvProfileList.NumberCols)
	fieldDelSpinBox.SetValue(lastDelSpinboxOptionFieldsValue)
	fieldDelButton.OnClicked(func(*ui.Button) {
		lastDelSpinboxOptionFieldsValue = fieldDelSpinBox.Value()
		if CsvProfileList.NumberCols > 0 { // Limit operation to avoid out of range
			CsvProfileList.FieldNames = DeleteSl(CsvProfileList.FieldNames, lastDelSpinboxOptionFieldsValue-1)
			CsvProfileList.FieldDisplay = DeleteSl(CsvProfileList.FieldDisplay, lastDelSpinboxOptionFieldsValue-1)
			CsvProfileList.FieldOut = DeleteSl(CsvProfileList.FieldOut, lastDelSpinboxOptionFieldsValue-1)
			CsvProfileList.FieldType = DeleteSl(CsvProfileList.FieldType, lastDelSpinboxOptionFieldsValue-1)
			CsvProfileList.NumberCols--
			CsvProfileList.Modified = true
			for row := 0; row < CsvProfileList.NumberRows; row++ {
				tableDatas[row] = DeleteSl(tableDatas[row], lastDelSpinboxOptionFieldsValue-1)
			}
		}
		// Reset storeprofile and row checkbox selection
		reloadMainStoreProfile()

		// Refresh tabs
		refreshOptionFields()

	})
	gridFields.Append(fieldDelButton,
		2, 0, 1, 1,
		true, ui.AlignStart, false, ui.AlignCenter)
	gridFields.Append(ui.NewLabel("Pos :"),
		2, 0, 1, 1,
		true, ui.AlignEnd, false, ui.AlignCenter)
	gridFields.Append(fieldDelSpinBox,
		3, 0, 1, 1,
		true, ui.AlignStart, false, ui.AlignCenter)

	// Set LABELS
	labels := []string{"Field name", "Display", "Write out", "Field type"}
	for idx, value := range labels {
		gridFields.Append(ui.NewLabel(value),
			idx, 1, 1, 1,
			true, ui.AlignCenter, false, ui.AlignCenter)
	}
	// Set COLUMN options
	dialogFields = make([]Fields, len(CsvProfileList.FieldNames))

	for idx, value := range CsvProfileList.FieldNames {
		dialogFields[idx].Init(value)

		gridFields.Append(dialogFields[idx].fColumn,
			0, idx+2, 1, 1,
			true, ui.AlignCenter, false, ui.AlignCenter)
		dialogFields[idx].fColumn.OnChanged(func(*ui.Entry) {
			recordFields()
		})
		gridFields.Append(dialogFields[idx].fDisp,
			1, idx+2, 1, 1,
			true, ui.AlignCenter, false, ui.AlignCenter)
		dialogFields[idx].fDisp.SetSelected(getValChar("y", CsvProfileList.FieldDisplay[idx]).(int))
		dialogFields[idx].fDisp.OnSelected(func(*ui.Combobox) { recordFields() })

		gridFields.Append(dialogFields[idx].fOutput,
			2, idx+2, 1, 1,
			true, ui.AlignCenter, false, ui.AlignCenter)
		dialogFields[idx].fOutput.SetSelected(getValChar("y", CsvProfileList.FieldOut[idx]).(int))
		dialogFields[idx].fOutput.OnSelected(func(*ui.Combobox) { recordFields() })

		gridFields.Append(dialogFields[idx].fType,
			3, idx+2, 1, 1,
			true, ui.AlignCenter, false, ui.AlignCenter)
		dialogFields[idx].fType.SetSelected(getValChar("t", CsvProfileList.FieldType[idx]).(int))
		dialogFields[idx].fType.OnSelected(func(*ui.Combobox) { recordFields() })
	}
	//	}
	return vbox
}

// Record fields when modified
func recordFields() {
	for idx, field := range dialogFields {
		CsvProfileList.FieldNames[idx] = field.fColumn.Text()
		CsvProfileList.FieldDisplay[idx] = getValChar("y", field.fDisp.Selected()).(string)
		CsvProfileList.FieldOut[idx] = getValChar("y", field.fOutput.Selected()).(string)
		CsvProfileList.FieldType[idx] = getValChar("t", field.fType.Selected()).(string)
	}
	CsvProfileList.Modified = true
	refreshTable()
	refreshSortFields()
}

// Refresh Option fields and Table tab
func refreshOptionFields() {
	if optionsBox != nil { // Check if initialised or not
		if optionsBox.Enabled() { // Used to know if control has been already initialised
			optionsBox.Delete(3)            // Delete vbox containing (Option fields)
			tabMakeOptionFields(optionsBox) // Display (Option fields) with new values
		} else { // Control has not been already initialised so, do it
			mainTab.Delete(1)                            // Delete Options Tab
			mainTab.InsertAt("Options", 1, tabOptions()) // and refresh it
		}
	}
	refreshTable()
}
