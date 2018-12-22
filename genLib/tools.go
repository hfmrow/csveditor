// tools.go
package genLib

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

// Get OS path separator
func PathSep() string {
	return string(os.PathSeparator)
}

// Make temporary directory
func TempMake(prefix string) string {
	dir, err := ioutil.TempDir("", prefix+"-")
	CheckP(err)
	return dir
}

// Remove directory
func TempRemove(fName string) {
	err := os.RemoveAll(fName)
	Check(err)
}

// Check error function input message is optional or accept multiple arguments.
func Check(err error, message ...string) {
	if err != nil {
		msgs := ``
		if len(message) != 0 { // Make string with messages if exists
			for _, mess := range message {
				msgs += `[` + mess + `]`
			}
		}
		pc, file, line, ok := runtime.Caller(1) //	(pc uintptr, file string, line int, ok bool)
		if ok == false {                        // Remove "== false" if needed
			fName := runtime.FuncForPC(pc).Name()
			fmt.Printf("[%s][%s][File: %s][Func: %s][Line: %d]\n", msgs, err.Error(), file, fName, line)
		} else {
			stack := strings.Split(fmt.Sprintf("%s", debug.Stack()), "\n")
			for idx := 5; idx < len(stack)-1; idx = idx + 2 {
				fmt.Printf("%s[%s][%s]\n", msgs, err.Error(), strings.Join([]string{TrimSpace(stack[idx], "-s"), TrimSpace(stack[idx+1], "-c")}, "]["))
			}
		}
		//	os.Exit(1) //	Break execution
	}
}

// Check error function, this one do panic !!
func CheckP(err error) {
	if err != nil {
		panic(err)
	}
}

// Check error and return true if error exist
func IsError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

// Use function to avoid "Unused variable ..." msgs
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
