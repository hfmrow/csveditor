// tabMakeOption.go

package main

import (
	"strings"

	. "github.com/hfmrow/csveditor/genLib"

	"github.com/andlabs/ui"
)

func tabMakeOptions(vbox *ui.Box) *ui.Box {
	if CsvProfileList.Initialised { // If there is nothing to display, skip it ...

		vbox.Enable() // Used to know if control has been already initialised

		// Create options Form
		entryForm := ui.NewForm()
		entryForm.SetPadded(true)
		vbox.Append(entryForm, false)

		// First part in a group
		gridComma := ui.NewGrid()
		gridComma.SetPadded(true)
		groupCommaQuote := ui.NewGroup("   C.S.V format:  " + TruncateString(CsvProfileList.FileName, "...", 80, 1))
		groupCommaQuote.SetMargined(true)
		groupCommaQuote.SetChild(gridComma)
		vbox.Append(groupCommaQuote, false)

		// Populate fields
		outComma := ui.NewEditableCombobox()
		outComma.Append(`,`)
		outComma.Append(`;`)
		outComma.Append(`\t`)
		outComma.SetText(CsvProfileList.OutComma)
		outComma.OnChanged(func(*ui.EditableCombobox) {
			CsvProfileList.OutComma = outComma.Text()
			CsvProfileList.Modified = true
		})
		outQuote := ui.NewCombobox()
		outQuote.Disable() // TODO using csvtrans
		outQuote.Append(`"`)
		outQuote.Append(`'`)
		outQuote.SetSelected(getValChar("s", CsvProfileList.OutQuote).(int))
		outQuote.OnSelected(func(*ui.Combobox) {
			CsvProfileList.OutQuote = getValChar("s", outQuote.Selected()).(string)
			CsvProfileList.Modified = true
		})
		outCRLF := ui.NewCombobox()
		outCRLF.Append(`CR`)
		outCRLF.Append(`LF`)
		outCRLF.Append(`CRLF`)
		outCRLF.SetSelected(getValChar("c", CsvProfileList.OutCrlf).(int))
		outCRLF.OnSelected(func(*ui.Combobox) {
			CsvProfileList.OutCrlf = getValChar("c", outCRLF.Selected()).(string)
			CsvProfileList.Modified = true
		})

		fmtDate1 := ui.NewCombobox()
		fmtDate2 := ui.NewCombobox()
		fmtDate3 := ui.NewCombobox()

		decimalSep := ui.NewCombobox()
		decimalSep.Append(",")
		decimalSep.Append(".")
		decimalSep.SetSelected(0)
		decimalSep.OnSelected(func(*ui.Combobox) {
			switch decimalSep.Selected() {
			case 0:
				CsvProfileList.DecimalSep = ","
				CsvProfileList.Modified = true
			case 1:
				CsvProfileList.DecimalSep = "."
				CsvProfileList.Modified = true
			}
		})

		// Populate Date format (combobox)
		for idx := 0; idx < len(DatE[0]); idx++ {
			fmtDate1.Append(DatE[0][idx])
			fmtDate2.Append(DatE[0][idx])
			fmtDate3.Append(DatE[0][idx])
		}
		dateFormat := strings.Split(CsvProfileList.DateFormat, "-")
		fmtDate1.SetSelected(getDateFmt(DatE[0][getDateFmt(dateFormat[0])]))
		fmtDate2.SetSelected(getDateFmt(DatE[0][getDateFmt(dateFormat[1])]))
		fmtDate3.SetSelected(getDateFmt(DatE[0][getDateFmt(dateFormat[2])]))
		fmtDate1.OnSelected(func(*ui.Combobox) { storeDateFmt(fmtDate1.Selected(), 0) })
		fmtDate2.OnSelected(func(*ui.Combobox) { storeDateFmt(fmtDate2.Selected(), 1) })
		fmtDate3.OnSelected(func(*ui.Combobox) { storeDateFmt(fmtDate3.Selected(), 2) })

		outCharset := ui.NewCombobox()
		CharsetList := NewCharsetList()

		// Populate charset names (combobox)
		for _, cs := range CharsetList.SimpleCharsetList {
			outCharset.Append(cs)
		}

		_, cShort := CharsetList.GetPos(CsvProfileList.OutCharset)
		outCharset.SetSelected(cShort)
		outCharset.OnSelected(func(*ui.Combobox) {
			_, outcShort := CharsetList.GetCharset(outCharset.Selected())
			CsvProfileList.OutCharset = outcShort
			CsvProfileList.Modified = true
		})
		// Display Date format
		gridComma.Append(ui.NewLabel("Date format: "),
			0, 0, 1, 1,
			true, ui.AlignEnd, false, ui.AlignCenter)
		gridComma.Append(fmtDate1,
			1, 0, 1, 1,
			true, ui.AlignStart, false, ui.AlignCenter)
		gridComma.Append(fmtDate2,
			1, 0, 2, 1,
			true, ui.AlignCenter, false, ui.AlignCenter)
		gridComma.Append(fmtDate3,
			2, 0, 1, 1,
			true, ui.AlignStart, false, ui.AlignCenter)

		// Display decimal separator choice
		gridComma.Append(ui.NewLabel("Decimal separator: "),
			4, 0, 1, 1,
			true, ui.AlignEnd, false, ui.AlignCenter)
		gridComma.Append(decimalSep,
			5, 0, 1, 1,
			true, ui.AlignCenter, false, ui.AlignCenter)

		// Display controls
		gridComma.Append(ui.NewLabel("Out comma: "),
			0, 1, 1, 1,
			true, ui.AlignEnd, false, ui.AlignCenter)
		gridComma.Append(outComma,
			1, 1, 1, 1,
			true, ui.AlignStart, false, ui.AlignEnd)
		gridComma.Append(ui.NewLabel("Out quote: "),
			2, 1, 1, 1,
			true, ui.AlignEnd, false, ui.AlignCenter)
		gridComma.Append(outQuote,
			3, 1, 1, 1,
			true, ui.AlignStart, false, ui.AlignEnd)
		//
		gridComma.Append(ui.NewLabel("Out line-end: "),
			4, 1, 1, 1,
			true, ui.AlignEnd, false, ui.AlignCenter)
		gridComma.Append(outCRLF,
			5, 1, 1, 1,
			true, ui.AlignStart, false, ui.AlignEnd)
		gridComma.Append(ui.NewLabel("Out charset: "),
			6, 1, 1, 1,
			true, ui.AlignEnd, false, ui.AlignCenter)
		gridComma.Append(outCharset,
			7, 1, 1, 1,
			true, ui.AlignStart, false, ui.AlignEnd)

		vbox = tabMakeOptionFields(vbox)
	}
	return vbox
}

// Refresh Option fields and Table tab
func refreshOption() {
	mainTab.Delete(1)
	mainTab.InsertAt("Options", 1, tabOptions()) // Display (Option fields) with new values
}

func storeDateFmt(value, idx int) {
	dateFormat := strings.Split(CsvProfileList.DateFormat, "-")
	switch idx {
	case 0:
		dateFormat[0] = DatE[1][value]
	case 1:
		dateFormat[1] = DatE[1][value]
	case 2:
		dateFormat[2] = DatE[1][value]
	}
	CsvProfileList.DateFormat = strings.Join(dateFormat, "-")
	CsvProfileList.Modified = true
}
