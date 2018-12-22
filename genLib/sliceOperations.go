// sliceOperations.go

package genLib

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Search in 2d string slice. cs=case sensitive, ww=whole word, rx=regex
func SearchSl(find string, table [][]string, cs, ww, rx bool) ([][]string, error) {
	if len(table) != 0 {
		if len(find) != 0 {
			var outTable [][]string
			if !cs {
				find = strings.ToLower(find)
			}

			regX, err := regexp.Compile(`(` + find + `)`)
			if err != nil {
				return [][]string{}, err
			}
			if ww {
				regX, err = regexp.Compile(`(` + find + `\b)`)
				if err != nil {
					return [][]string{}, err
				}
			}
			if rx {
				regX, err = regexp.Compile(find)
				if err != nil {
					return [][]string{}, err
				}
			}

			for idxRow := 0; idxRow < len(table); idxRow++ {
				for _, col := range table[idxRow] {
					if !cs {
						col = strings.ToLower(col)
					}
					if regX.MatchString(col) {
						outTable = append(outTable, table[idxRow])
						break // Avoid duplicate when element found twice in same row
					}
				}
			}
			if len(outTable) == 0 {
				return [][]string{}, errors.New(find + "\n\nNot found ...")
			}
			return outTable, nil // Result found then return it
		}
	}
	return [][]string{}, errors.New("Nothing to search ...")
}

// Get index of a string in a slice
func GetStrIndex(slice []string, item string) int {
	for i, _ := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

// Search in 2d string slice if a row exist.
func IsExist(slice [][]string, cmpRow []string) bool {
	for _, mainRow := range slice {
		if reflect.DeepEqual(mainRow, cmpRow) {
			return true
		}
	}
	return false
}

// Add data at the begining of a string slice
func Preppend(slice []string, prepend ...string) []string {
	return append(prepend, slice...)
}

// Add data at a specified position in slice of a string
func AppendAt(slice []string, pos int, insert ...string) []string {
	if pos > len(slice) {
		pos = len(slice)
	} else if pos < 0 {
		pos = 0
	}
	return append(slice[:pos], append(insert, slice[pos:]...)...)
}

// Delete value at specified position in string slice
func DeleteSl(slice []string, pos int) []string {
	return append(slice[:pos], slice[pos+1:]...)
}

// Delete value at specified position in string slice
func DeleteSl1(slice []string, pos int) []string {
	copy(slice[pos:], slice[pos+1:])
	slice[len(slice)-1] = "" // or the zero value of T
	return slice[:len(slice)-1]
}

// Remove duplicat entry in a string slice
func RemoveDupSl(slice []string) []string {
	var isExist bool
	tmpSlice := make([]string, 0)
	for _, inValue := range slice {
		isExist = false
		inValue = RemoveNonAlNum(inValue)
		for _, outValue := range tmpSlice {
			if outValue == inValue {
				isExist = true
				break
			}
		}
		if !isExist {
			tmpSlice = append(tmpSlice, inValue)
		}
	}
	fmt.Println(len(tmpSlice))
	return tmpSlice
}

// Sort 2d string slice with date inside
func SliceSortDate(slice [][]string, fmtDate string, dateCol, secDateCol int, ascendant bool) [][]string {
	fieldsCount := len(slice[0]) // Get nb of columns
	var firstLine int
	var previous, after string
	var positiveidx, negativeidx int
	// compute unix date using given column numbers
	for idx := firstLine; idx < len(slice); idx++ {
		dateStr := FindDate(slice[idx][dateCol], fmtDate)
		if dateStr != nil { // search for 1st column
			slice[idx] = append(slice[idx], fmt.Sprintf("%d", FormatDate(fmtDate, dateStr[0]).Unix()))
		} else if secDateCol != -1 { // Check for second column if it was given
			dateStr = FindDate(slice[idx][secDateCol], fmtDate)
			if dateStr != nil { // If date was not found in 1st column, search for 2nd column
				slice[idx] = append(slice[idx], fmt.Sprintf("%d", FormatDate(fmtDate, slice[idx][secDateCol]).Unix()))
			} else { //  in case where none of the columns given contain date field, put null string if there is no way to find a date
				slice[idx] = append(slice[idx], ``)
			}
		} else { // put null string if there is no way to find a date
			slice[idx] = append(slice[idx], ``)
		}
	}
	// Ensure we always have a value in sorting field (get previous or next closer)
	for idx := firstLine; idx < len(slice); idx++ {
		if slice[idx][fieldsCount] == `` {
			for idxFind := firstLine + 1; idxFind < len(slice); idxFind++ {
				positiveidx = idx + idxFind
				negativeidx = idx - idxFind
				if positiveidx >= len(slice) { // Check index to avoiding 'out of range'
					positiveidx = len(slice) - 1
				}
				if negativeidx <= 0 {
					negativeidx = 0
				}
				after = slice[positiveidx][fieldsCount] // Get previous or next value
				previous = slice[negativeidx][fieldsCount]
				if previous != `` { // Set value, prioritise the previous one.
					slice[idx][fieldsCount] = previous
					break
				}
				if after != `` {
					slice[idx][fieldsCount] = after
					break
				}
			}
		}
	}
	tmpLines := make([][]string, 0)
	if ascendant != true {
		// Sort by unix date preserving order descendant
		sort.SliceStable(slice, func(i, j int) bool { return slice[i][len(slice[i])-1] > slice[j][len(slice[i])-1] })
		for idx := firstLine; idx < len(slice); idx++ { // Store row count elements - 1
			tmpLines = append(tmpLines, slice[idx][:len(slice[idx])-1])
		}
	} else {
		// Sort by unix date preserving order ascendant
		sort.SliceStable(slice, func(i, j int) bool { return slice[i][len(slice[i])-1] < slice[j][len(slice[i])-1] })
		for idx := firstLine; idx < len(slice); idx++ { // Store row count elements - 1
			tmpLines = append(tmpLines, slice[idx][:len(slice[idx])-1])
		}
	}
	return tmpLines
}

// Sort 2d string slice
func SliceSortString(slice [][]string, col int, ascendant, caseSensitive bool) {
	transform := func(inString string) string {
		return inString
	}
	if !caseSensitive {
		transform = func(inString string) string { return strings.ToLower(inString) }
	}

	if ascendant != true {
		// Sort string preserving order descendant
		sort.SliceStable(slice, func(i, j int) bool { return transform(slice[i][col]) > transform(slice[j][col]) })
	} else {
		// Sort string preserving order ascendant
		sort.SliceStable(slice, func(i, j int) bool { return transform(slice[i][col]) < transform(slice[j][col]) })
	}
}

// Sort 2d string with float value
func SliceSortFloat(slice [][]string, col int, ascendant bool, decimalChar string) {
	if ascendant != true {
		// Sort string (float) preserving order descendant
		sort.SliceStable(slice, func(i, j int) bool {
			return StringDecimalSwitchFloat(decimalChar, slice[i][col]) > StringDecimalSwitchFloat(decimalChar, slice[j][col])
		})
	} else {
		// Sort string (float) preserving order ascendant
		sort.SliceStable(slice, func(i, j int) bool {
			return StringDecimalSwitchFloat(decimalChar, slice[i][col]) < StringDecimalSwitchFloat(decimalChar, slice[j][col])
		})
	}
}

// Convert comma to dot if needed and return 0 if input string is empty.
func StringDecimalSwitchFloat(decimalChar, inString string) float64 {
	if inString == "" {
		inString = "0"
	}
	switch decimalChar {
	case ",":
		f, _ := strconv.ParseFloat(strings.Replace(inString, ",", ".", 1), 64)
		return f
	case ".":
		f, _ := strconv.ParseFloat(inString, 64)
		return f
	}
	return -1
}
