// tabTableModelHandler.go

package main

import (
	"strconv"

	"github.com/andlabs/ui"
)

// Table definition
type modelHandler struct {
	rowCount      []int
	checkBoxValue []int
}

func newModelHandler() *modelHandler {
	mh = new(modelHandler)
	mh.checkBoxValue = make([]int, CsvProfileList.NumberRows)
	return mh
}

// Table dimensions
func (mh *modelHandler) NumRows(m *ui.TableModel) int {
	return CsvProfileList.NumberRows // Get nb of row from tDatas
}

// Get nb of columns from tDatas and add aditionalColumns value (for checkbox col, ...)
func (mh *modelHandler) NumCols(m *ui.TableModel) int {
	return CsvProfileList.NumberCols + aditionalVisibleColumns
}

// Add columns type
func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	tValue := make([]ui.TableValue, mh.NumCols(m))
	tValue[0] = ui.TableInt(0)     // Init checkBox column
	tValue[1] = ui.TableString("") // Init rowCount column
	for idx := aditionalVisibleColumns; idx < len(tValue); idx++ {
		tValue[idx] = ui.TableString("") // Init strings columns
	}
	return tValue
}

// Table Values
func (mh *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	if column < mh.NumCols(m) && row < CsvProfileList.NumberRows {
		switch column {
		case 0:
			return ui.TableInt(mh.checkBoxValue[row])
		case 1:
			return ui.TableString(strconv.Itoa(row + 1))
		default:
			return ui.TableString(tableDatas[row][column-aditionalVisibleColumns]) // Get cell Values from tableDatas (csv)
		}
	}
	return nil
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	if column < mh.NumCols(m)+1 && row < CsvProfileList.NumberRows {
		switch column {
		case 0:
			mh.checkBoxValue[row] = int(value.(ui.TableInt)) // checkboxes column
			if mh.checkBoxValue[row] == int(ui.TableTrue) {
				ChkRow.Add(row) // Add last one checked for edition purpose
			} else {
				ChkRow.Delete(row) // Or delete entry if needed
			}
		default:
			if onTheFlyEdit {
				tableDatas[row][column-aditionalVisibleColumns] = string(value.(ui.TableString))
			}
		}
	}
}
