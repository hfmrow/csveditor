// varStructDeclaration.go

package main

import (
	"fmt"

	"os"

	"github.com/andlabs/ui"
	"github.com/hfmrow/csveditor/genLib"
)

// App infos
const appName = "Csv Editor"
const appVers = "v1.6.1"
const appCreat = "H.F.M"
const copyRight = "©2018 MIT licence"

// Store UI pointers
var model *ui.TableModel
var mh *modelHandler
var mainwin *ui.Window // Andlabs ui main window
var mainTab *ui.Tab
var tableBox *ui.Box   // Store current Table boxed
var optionsBox *ui.Box // Store current Options boxed
var sortBox *ui.Box    // Store current SortCompare boxed

// Other vars
var argFilename string // Used to pass arguments filename
var tempDir string     // Temp dir
// var wg = &sync.WaitGroup{} // Used to wait until end of a function

var CsvProfileList ProfileCsv       // Store actual CSV profile
var tableDatas [][]string           // Store CSV datas
var aditionalVisibleColumns int = 2 // Column 0 checkbox, column 1 for checkBox text

// Some other usefull declarations
var charsetList genLib.CharsetList // Store Different charset names
var lastSpinboxOptionFieldsValue int = 1
var lastDelSpinboxOptionFieldsValue int = 1
var sortFields []int // Used to store fields in sortCompare tab

var sortFieldsFlag bool        // Used to know if we need to refresh fields box in sort tab
var caseSensitiveSortFlag bool // Define sort type (case sensitive)
var onTheFlyEdit bool          // Flag for direct cell edition

// Store sort type for each fields (ascending/descending)
var sortAscDscCheckbox []ascDscCheckbox

type ascDscCheckbox struct {
	fChkbox *ui.Checkbox
}

var SeP = []string{`"`, `'`}
var CrlF = []string{`CR`, `LF`, `CRLF`}
var YesnO = []string{`Yes`, `No`}
var TypE = []string{`Date`, `String`, `Numeric`}
var DatE = [][]string{{`Day`, `Month`, `Year`}, {`%d`, `%m`, `%y`}}

// Retreive (Value < in > Pos) of input. val=(s:sep, c:crlf, y:yesno, t:type)
func getValChar(val string, in interface{}) interface{} {
	switch v := in.(type) {
	case int:
		switch val {
		case "s":
			return SeP[in.(int)]
		case "c":
			return CrlF[in.(int)]
		case "y":
			return YesnO[in.(int)]
		case "t":
			return TypE[in.(int)]
		default:
			return fmt.Errorf("Unexpected data type: %s", val)
		}
	case string:
		switch val {
		case "s":
			return genLib.GetStrIndex(SeP, in.(string))
		case "c":
			return genLib.GetStrIndex(CrlF, in.(string))
		case "y":
			return genLib.GetStrIndex(YesnO, in.(string))
		case "t":
			return genLib.GetStrIndex(TypE, in.(string))
		default:
			return fmt.Errorf("Unexpected data type: %s", val)
		}
	default:
		return fmt.Errorf("Unexpected variable type: %v", v)
	}
	return fmt.Errorf("Matching value not found: %v", in)
}

// Used to store multiples Tables/Profiles
var storeProfiles []ProfileTable
var profName = []string{"Main", "Search", "Sort", "Compare"}
var profFootBar []string

type ProfileTable struct {
	Prof        ProfileCsv
	Table       [][]string
	Name        string
	Initialised bool
}

func NewStoreProfiles() []ProfileTable {
	profFootBar = profFootBar[:0]
	profs := make([]ProfileTable, 4)
	for idx, _ := range profs {
		profs[idx].Name = profName[idx]
		profs[idx].Prof = ProfileCsv{}
		profs[idx].Table = [][]string{}
		profs[idx].Initialised = false
	}
	return profs
}

// Get word position in "ProfName" slice
func GetStoreProfilesPos(name string) int {
	for idx, nme := range profFootBar {
		if name == nme {
			return idx
		}
	}
	return -1
}

// Used to store field values
var dialogFields []Fields // Used to store field values in options tab
type Fields struct {
	fColumn *ui.Entry
	fDisp   *ui.Combobox
	fOutput *ui.Combobox
	fType   *ui.Combobox
	fIdx    int
}

func (f *Fields) Init(name string) {
	f.fColumn = ui.NewEntry()
	f.fDisp = ui.NewCombobox()
	f.fOutput = ui.NewCombobox()
	f.fType = ui.NewCombobox()
	f.fColumn.SetText(name)
	f.fDisp.Append(getValChar("y", 0).(string)) // Yes
	f.fDisp.Append(getValChar("y", 1).(string)) // No
	f.fDisp.SetSelected(getValChar("y", "Yes").(int))
	f.fOutput.Append(getValChar("y", 0).(string)) // Yes
	f.fOutput.Append(getValChar("y", 1).(string)) // No
	f.fOutput.SetSelected(getValChar("y", "Yes").(int))
	f.fType.Append(getValChar("t", 0).(string)) // Date
	f.fType.Append(getValChar("t", 1).(string)) // String
	f.fType.Append(getValChar("t", 2).(string)) // Numeric
	f.fType.SetSelected(getValChar("t", "String").(int))
	f.fIdx = -1
}

// Used to store checked rows, lifo stack to get last selected row for edition purpose
var ChkRow CheckedRow // Checkbox history
type CheckedRow struct {
	Checked []int
}

func NewCheckedRow() CheckedRow { // Initialise struct
	return CheckedRow{}
}
func (cr *CheckedRow) Add(val int) { // Add row entry
	cr.Checked = append(cr.Checked, val)
}
func (cr *CheckedRow) Reset() { // Reset rows entry
	cr.Checked = cr.Checked[:0]
}
func (cr *CheckedRow) Last() int { // Get entry for last checked row.
	return cr.Checked[len(cr.Checked)-1]
}
func (cr *CheckedRow) Delete(val int) { // Remove value
	for idx, value := range cr.Checked { // Search fo value
		if value == val { // delete if exist
			cr.Checked = append(cr.Checked[:idx], cr.Checked[idx+1:]...)
		}
	}
}
func (cr *CheckedRow) NotEmpty() bool { // Check if struct is empty or not
	if len(cr.Checked) != 0 {
		return true
	}
	return false
}

// Change current dir ... (debug purpose)
func chgToSpecificDir(dir string) {
	err := os.Chdir(dir)
	genLib.Check(err, "os.Chdir")
}

// Convert date format, input "Day" or "%d", return 0,  "Month" or "%m", return 1 ...
func getDateFmt(in string) int {
	switch in {
	case DatE[0][0]:
		return 0
	case DatE[0][1]:
		return 1
	case DatE[0][2]:
		return 2
	case DatE[1][0]:
		return 0
	case DatE[1][1]:
		return 1
	case DatE[1][2]:
		return 2
	}
	return -1
}
