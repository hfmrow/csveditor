// rwProfileCsv.go

package main

import (
	"os"
	"strings"

	"github.com/hfmrow/csveditor/genLib"
)

// Update opt profile
func updOptProfile(profileList *ProfileCsv) {
	optFilename := getOptFileName(profileList.FileName)
	profileList.Modified = false
	genLib.WriteJson(optFilename, &profileList)
	CsvProfileList = *profileList
}

// Check for saving modified table
func doCheckAndSaveOnExit(profileList *ProfileCsv) {
	var restructuredCsv [][]string
	// Make CSV datas
	restructuredCsv = append(restructuredCsv, profileList.FieldNames)
	restructuredCsv = append(restructuredCsv, tableDatas...)
	// Count nb of deleted column(s)
	count := 0
	for _, act := range profileList.FieldOut {
		if act == getValChar("y", 1).(string) { // "No"
			count++
		}
	}
	tmpProfileList := ProfileCsv{}
	// Remove unwanted columns before writing csv
	if count != 0 {
		tmpDatas := make([][]string, len(restructuredCsv))
		// Some temp var
		var FieldDisplay, FieldType []string
		for idx := 0; idx < len(restructuredCsv[0]); idx++ {
			if profileList.FieldOut[idx] == getValChar("y", 0).(string) { // "Yes"
				for row := 0; row < len(restructuredCsv); row++ {
					tmpDatas[row] = append(tmpDatas[row], restructuredCsv[row][idx])
				}
				FieldDisplay = append(FieldDisplay, profileList.FieldDisplay[idx])
				FieldType = append(FieldType, profileList.FieldType[idx])
			}
		}
		// Write to CSV file
		restructuredCsv = tmpDatas
		genLib.WriteCsv(profileList.FileName, strings.Replace(storeProfiles[0].Prof.OutComma, `"`, ``, -1), restructuredCsv) // Remove `"`
		tmpProfileList = NewProfileCsv(profileList.FileName)                                                                 // Re-create ProfileList
		tmpProfileList.FieldDisplay = FieldDisplay
		tmpProfileList.FieldType = FieldType
	} else {
		// Write to CSV file without modification
		genLib.WriteCsv(profileList.FileName, strings.Replace(storeProfiles[0].Prof.OutComma, `"`, ``, -1), restructuredCsv) // Remove `"`
		tmpProfileList = NewProfileCsv(profileList.FileName)                                                                 // Re-create ProfileList
		tmpProfileList.FieldDisplay = storeProfiles[0].Prof.FieldDisplay
		tmpProfileList.FieldType = storeProfiles[0].Prof.FieldType
	}
	// Preserve some profile informations
	tmpProfileList.DateFormat = storeProfiles[0].Prof.DateFormat
	tmpProfileList.OutComma = storeProfiles[0].Prof.OutComma
	tmpProfileList.OutCrlf = storeProfiles[0].Prof.OutCrlf
	tmpProfileList.OutQuote = storeProfiles[0].Prof.OutQuote
	tmpProfileList.OutCharset = storeProfiles[0].Prof.OutCharset
	tmpProfileList.DecimalSep = storeProfiles[0].Prof.DecimalSep
	// Convert charset if needed
	actCharset := genLib.DetectCharsetFile(profileList.FileName)
	if actCharset != storeProfiles[0].Prof.OutCharset {
		genLib.FileCharsetSwitch(actCharset, storeProfiles[0].Prof.OutCharset, profileList.FileName)
	}
	// Convert Eol if needed
	actEOL := genLib.GetEOL(profileList.FileName)
	if actEOL != storeProfiles[0].Prof.OutCrlf {
		genLib.SetEOL(profileList.FileName, storeProfiles[0].Prof.OutCrlf)
	}
	// Save it
	updOptProfile(&tmpProfileList)
}

// Make option file name (i.e replace .ext by .opt)
func getOptFileName(filename string) string {
	return genLib.SplitFilePath(filename, "opt").OutputNewExt
}

// Check for how to proceed where optfile exist or not ...
func rwProfileCsv(filename string) ProfileCsv {
	optFilename := getOptFileName(filename)
	profileList := NewProfileCsv(filename) // Create ProfileList

	if _, err := os.Stat(optFilename); os.IsNotExist(err) { // Opt file does not exist
		genLib.WriteJson(optFilename, &profileList) // then save it ...
	} else { // Opt file exist
		genLib.ReadJson(optFilename, &profileList) // load it
	}
	return profileList
}

// add csv to ProfileTables and set as default
func addAndSetCsv(filename, name string) {
	profile := rwProfileCsv(filename)
	table := genLib.ReadCsv(filename, profile.Comma, profile.NumberCols, profile.FirstLine+1)
	for idx := 0; idx < len(storeProfiles); idx++ {
		if storeProfiles[idx].Name == name {
			storeProfiles[idx].Prof = profile
			storeProfiles[idx].Prof.Name = storeProfiles[idx].Name
			storeProfiles[idx].Table = table
			storeProfiles[idx].Initialised = true
		}
	}
	setProfile(name)
}

// Set profile as mainProfile
func setProfile(name string) {
	for idx := 0; idx < len(storeProfiles); idx++ {
		if storeProfiles[idx].Name == name {
			CsvProfileList = storeProfiles[idx].Prof
			CsvProfileList.Name = name
			tableDatas = storeProfiles[idx].Table
			if optionsBox != nil {
				// Disabling Option Tab if focused data is not Main table
				if name != profName[0] { // as "Main"
					CsvProfileList.FileName = name
					optionsBox.Disable()
				} else {
					optionsBox.Enable()
				}
			}
			mainwin.SetTitle(genLib.TruncateString(CsvProfileList.FileName, "...", 80, 1))
			break // ADDED 25-12
		}
	}
}

// Reset stored profile to main
func reloadMainStoreProfile() {
	ChkRow.Reset() // Reset stored checked rows history
	storeProfiles = NewStoreProfiles()
	storeProfiles[0].Prof = CsvProfileList
	storeProfiles[0].Prof.Name = profName[0] // as "Main"
	storeProfiles[0].Table = tableDatas
	storeProfiles[0].Initialised = CsvProfileList.Initialised
	setProfile(profName[0]) // as "Main"
	refreshTable()
}
