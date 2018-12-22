// jsonHandler.go

/// +build OMIT

package genLib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// read json datas from file to given interface / structure
// ReadJson(filename, &person)
func ReadJson(filename string, interf interface{}) {
	textFileBytes, err := ioutil.ReadFile(filename)
	Check(err, "ioutil.ReadFile!")

	err = json.Unmarshal(textFileBytes, &interf)
	Check(err, "json.Unmarshal!")
}

// Write json datas to file from given interface / structure
// i.e: WriteJson(filename, &person)
func WriteJson(filename string, interf interface{}) {
	var out bytes.Buffer
	jsonData, err := json.Marshal(&interf)
	Check(err, "json.Marshal!")

	err = json.Indent(&out, jsonData, "", "\t")
	Check(err, `json.Indent!`)

	f, err := os.Create(filename)
	defer f.Close()
	Check(err, `Create opt file!`)

	_, err = fmt.Fprintln(f, string(out.Bytes()))
	Check(err, `fmt.Fprintln!`)
}
