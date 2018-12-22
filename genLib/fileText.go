// fileText.go
package genLib

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type File struct {
	Path         string
	Name         string
	NameNoExt    string
	Ext          string
	RealPath     string
	RealName     string
	Output       string
	OutputNewExt string
	Absolute     string
	OsSep        string
	OsListSep    string
}

// Split full filename into path, ext, name, ... optionally add suffix before original extension or change extension
func SplitFilePath(filename string, toAdd ...string) File {
	var outFileName, outFileNameNewExt, absolute string
	realPath, err := filepath.EvalSymlinks(filename)
	Check(err)
	name := filepath.Base(filename)
	ext := filepath.Ext(filename)
	noext := strings.Replace(name, ext, ``, -1)
	if len(toAdd) != 0 {
		outFileName = filepath.Dir(realPath) + string(filepath.Separator) + noext + toAdd[0] + ext
		outFileNameNewExt = filepath.Dir(realPath) + string(filepath.Separator) + noext + "." + toAdd[0]
	}
	if !filepath.IsAbs(filename) {
		absolute, err = filepath.Abs(filename)
		Check(err)
	}
	return File{filepath.Dir(filename), name, noext, ext, filepath.Dir(realPath), filepath.Base(realPath), outFileName, outFileNameNewExt, absolute, string(filepath.Separator), string(filepath.ListSeparator)}
}

// Open file and get (CR, LF, CRLF) > string
func GetEOL(filename string) string {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	textFileBytes, err := ioutil.ReadFile(filename)
	Check(err, `ReadFile: `)
	if bytes.Contains(textFileBytes, bCRLF) {
		return `CRLF`
	} else if bytes.Contains(textFileBytes, bCR) {
		return `CR`
	}
	if bytes.Contains(textFileBytes, bLF) {
		return `LF`
	}
	return `Undefined end of line: ` + filename
}

// Open file and convert EOL (CR, LF, CRLF) then write it back.
func SetEOL(filename, eol string) {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	var outEol []byte
	switch eol {
	case "CR":
		outEol = bCR
	case "LF":
		outEol = bLF
	case "CRLF":
		outEol = bCRLF
	default:
		fmt.Println("EOL convert error: Undefined end of line")
	}
	// Load file LF separated
	textFileBytes, err := ioutil.ReadFile(filename)
	Check(err, `ReadFile: `+filename)
	// Handle end of line
	switch GetEOL(filename) {
	case `CR`:
		textFileBytes = bytes.Replace(textFileBytes, bCR, outEol, -1)
	case `LF`:
		textFileBytes = bytes.Replace(textFileBytes, bLF, outEol, -1)
	case `CRLF`:
		textFileBytes = bytes.Replace(textFileBytes, bCRLF, outEol, -1)
	}
	err = ioutil.WriteFile(filename, textFileBytes, 0644)
	Check(err, "Error writing file in EOL convert operation.")
}

// Write string to file low lvl format
func WriteTextFile(filename, data string, appendIfExist ...bool) {
	var apnd bool
	var file *os.File
	var err error
	if len(appendIfExist) != 0 {
		apnd = appendIfExist[0]
	}

	// open file using READ & WRITE permission
	if apnd { // append
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0660)
		if IsError(err) {
			return
		}
		defer file.Close()
	} else { // Overwrite
		file, err = os.Create(filename)
		if IsError(err) {
			return
		}
		defer file.Close()
		file, err = os.OpenFile(filename, os.O_RDWR, 0644)
		if IsError(err) {
			return
		}
		defer file.Close()
	}
	// write some text to file
	_, err = file.WriteString(data)
	if IsError(err) {
		return
	}
	// save changes
	err = file.Sync()
	if IsError(err) {
		return
	}
}

// Read file low lvl format
func ReadFile(filename string) []byte {
	stats, err := os.Stat(filename)
	if IsError(err) {
		return []byte{}
	}
	// open file
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if IsError(err) {
		return []byte{}
	}
	defer file.Close()
	// read file, line by line
	var datas = make([]byte, stats.Size())
	_, err = file.Read(datas)
	if err != nil && err != io.EOF {
		IsError(err)
		return []byte{}
	}
	return datas
}

// Load text file to slice, opt are same as TrimSpace function "-c, -s, +& or -&" and can be cumulative.
// This function Reconize "CR", "LF", "CRLF", and convert charset to utf-8
func TextFileToLines(filename string, opt ...string) []string {
	bCR := []byte{0x0D}
	bLF := []byte{0x0A}
	bCRLF := []byte{0x0D, 0x0A}
	spliter := ``
	// Load file LF separated
	textFileBytes, err := ioutil.ReadFile(filename)
	Check(err, `ReadFile: `+filename)
	// Handle end of line
	switch GetEOL(filename) {
	case `CR`:
		spliter = string(bCR)
	case `LF`:
		spliter = string(bLF)
	case `CRLF`:
		spliter = string(bCRLF)
	}
	stringsText := strings.Split(CharsetToUtf8(string(textFileBytes)), spliter) // convert to slice
	for idx, line := range stringsText {
		stringsText[idx] = TrimSpace(line, opt...)
	}
	return stringsText
}

// Write slice to file
func LinesToTextFile(filename string, values interface{}) error {
	f, err := os.Create(filename)
	defer f.Close()
	Check(err, `FileWrite: `)
	rv := reflect.ValueOf(values)
	if rv.Kind() != reflect.Slice {
		return errors.New("Not a slice")
	}
	for i := 0; i < rv.Len(); i++ {
		fmt.Fprintln(f, rv.Index(i).Interface())
	}
	return nil
}

// Read data from CSV file
func ReadCsv(filename, comma string, fields, startLine int, endLine ...int) [][]string {
	commaUnEsc, err := strconv.Unquote(`"` + strings.Replace(comma, `"`, ``, -1) + `"`)
	Check(err, `Rune conversion error !`)
	commaR := []rune(commaUnEsc)
	lines := TextFileToLines(filename, `-ct`) // Get file lines entry and trim/remove multi spaces
	if len(endLine) == 0 {                    // There is no value for endline, so give it full length.
		endLine = append(endLine, len(lines))
	}
	newLines := []string{}                    // Array for usedlines
	startLine--                               // Adjuste line number
	for j := startLine; j < endLine[0]; j++ { // Get only needed part of file
		newLines = append(newLines, lines[j])
	}
	csvReader := csv.NewReader(strings.NewReader(strings.Join(newLines, "\n"))) // Convert slice to string then read as csv
	csvReader.Comma = commaR[0]
	csvReader.FieldsPerRecord = fields
	records, err := csvReader.ReadAll()
	Check(err, `Read CSV datas`)
	return records
}

// Write data to CSV file
func WriteCsv(filename, comma string, rows [][]string) {
	commaUnEsc, err := strconv.Unquote(`"` + strings.Replace(comma, `"`, ``, -1) + `"`)
	Check(err, `Rune conversion error !`)
	commaR := []rune(commaUnEsc)
	f, err := os.Create(filename)
	defer f.Close()
	Check(err, " os.Create!")
	w := csv.NewWriter(f)
	w.Comma = commaR[0]
	w.WriteAll(rows)
}

func RemovIfExist(filename string) {
	if _, err := os.Stat(filename); !os.IsNotExist(err) { // File  exist?, delete it
		err := os.Remove(filename)
		Check(err, fmt.Sprintln("Error removing: "+filename))
	}
}
