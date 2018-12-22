// csvprofile.go
package main

import (
	"encoding/csv"

	"strconv"
	"time"

	"strings"

	. "github.com/hfmrow/csveditor/genLib"
)

// Try to find informations about CSV file format like comma separator, quote separator, number of fields, first line with column names, endline type (crlf)
type ProfileCsv struct {
	Comma, Quote, Crlf, OutComma, OutQuote, OutCrlf, DecimalSep string
	NumberCols, FirstLine, NumberRows                           int
	Initialised, Modified                                       bool
	Charset, OutCharset, FileName, DateFormat, Name             string
	LastUsage                                                   int64
	FieldNames, FieldDisplay, FieldOut, FieldType               []string
}

// Create new struct
func NewProfileCsv(filename string) ProfileCsv {
	p := ProfileCsv{}
	p.Init(filename)
	return p
}

// Initialise it
func (p *ProfileCsv) Init(filename string) {
	p.FileName = filename

	// Retrieve line end type
	p.Crlf = GetEOL(filename)
	if p.OutCrlf == "" { // Set value for output if it is not set
		p.OutCrlf = p.Crlf
	}

	p.LastUsage = time.Now().Unix()
	p.DateFormat = `%d-%m-%y`
	p.DecimalSep = ","
	p.Initialised = true
	colCount := []RowStore{}
	// Preparing input file
	lines := TextFileToLines(filename, "-ct") // Load file in slice
	newString := strings.Join(lines, " ")     //	Join lines to make unique string

	// Get Encoding
	p.Charset = DetectCharsetFile(filename)
	p.OutCharset = p.Charset
	// Get separator and quote chars
	tmpSep := GetSep(newString) // Get comma and quote
	p.Comma, _ = strconv.Unquote(tmpSep[0])
	p.OutComma = p.Comma
	p.Quote, _ = strconv.Unquote(tmpSep[1])
	p.OutQuote = p.Quote

	//Get number of columns & Firstline with column names
	sepRne := []rune(p.Comma) // Put comma to rune
	for lineNb := 0; lineNb < len(lines); lineNb++ {
		r := csv.NewReader(strings.NewReader(lines[lineNb])) // Make newreader for csv
		r.Comma = sepRne[0]
		records, err := r.Read() // Read line from csv reader to get fields count
		if err == nil {          //	Store in slice column number for each lines. If error occure, just jump next line
			colCount = AppendIfMissingC(colCount, RowStore{lineNb, len(records), 0, lines[lineNb]})
		}
	}
	// Find first line with same number of columns (mostly repeated)
	lastTot := 0
	highVal := 0
	for lineNb := 0; lineNb < len(colCount); lineNb++ {
		if colCount[lineNb].Tot > lastTot {
			lastTot = colCount[lineNb].Tot
			highVal = lineNb
		}
	}
	p.NumberCols = colCount[highVal].Cnt
	p.FirstLine = colCount[highVal].Idx + 1
	p.NumberRows = lastTot
	// Get column names
	r := csv.NewReader(strings.NewReader(lines[colCount[highVal].Idx])) // Make newreader for csv (column names)
	r.Comma = sepRne[0]
	recordsFieldNames, err := r.Read()
	Check(err, `Read col Names profile error !`)
	p.FieldNames = recordsFieldNames
	// Set Default values for display, output and type of columns
	for idx := 0; idx < len(p.FieldNames); idx++ {
		p.FieldDisplay = append(p.FieldDisplay, "Yes")
		p.FieldOut = append(p.FieldOut, "Yes")
	}
	// Check for column types. Not really accurate but do the job most of time.
	tmpCsv := ReadCsv(p.FileName, p.Comma, p.NumberCols, p.FirstLine, p.NumberRows)
	var dateCount, floatCount, totalRow, totalRowCount int
	var value string
	if p.NumberRows > 10 {
		totalRow = 10
	} else {
		totalRow = len(tmpCsv) - 1
	}
	for idxCol, _ := range tmpCsv[0] {
		for idxRow := 0; idxRow < totalRow; idxRow++ {
			value = tmpCsv[idxRow][idxCol]
			if value != "" {
				totalRowCount++
				if IsDate(value) {
					dateCount++
				} else if IsFloat(value) {
					floatCount++
				}
			}
		}
		if totalRowCount == 0 { // To avoid divide by zero ...
			dateCount = 0
			floatCount = 0
		} else {
			dateCount = int(float64((dateCount * 100) / totalRowCount))
			floatCount = int(float64((floatCount * 100) / totalRowCount))
		}
		totalRowCount = 0
		switch {
		case dateCount > 60:
			p.FieldType = append(p.FieldType, "Date")
			dateCount = 0
		case floatCount > 60:
			p.FieldType = append(p.FieldType, "Numeric")
			floatCount = 0
		default:
			p.FieldType = append(p.FieldType, "String")
		}
	}
}

// Update it
func UpdateProfile(prof *ProfileCsv) *ProfileCsv {
	prof.Crlf = prof.OutCrlf
	prof.Charset = prof.OutCharset

	tmpOutComma, err := strconv.Unquote(`"` + strings.Replace(prof.OutComma, `"`, ``, -1) + `"`)
	Check(err, "Uodate comma profile error !")
	prof.OutComma = tmpOutComma
	prof.Comma = prof.OutComma

	prof.OutQuote, err = strconv.Unquote(`"` + strings.Replace(prof.OutQuote, `"`, ``, -1) + `"`)
	Check(err, "Uodate quote profile error !")
	prof.Quote = prof.OutQuote
	prof.Crlf = prof.OutCrlf

	return prof
}
