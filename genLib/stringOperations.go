// strings.go
package genLib

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Check if string is float
func IsFloat(inString string) bool {
	_, err := strconv.ParseFloat(strings.Replace(strings.Replace(inString, " ", "", -1), ",", ".", -1), 64)
	if err == nil {
		return true
	}
	return false
}

// Check if string is date
func IsDate(inString string) bool {
	dateFormats := NewDateFormat()
	for _, dteFmt := range dateFormats {
		if len(FindDate(inString, dteFmt+" %H:%M:%S")) != 0 {
			return true
		}
	}
	return false
}

// Convert []byte to hexString
func ByteToHexStr(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// Generate a randomized file name
func GenFileName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Md5String(fmt.Sprintln(r.Int63n(time.Now().UnixNano())))
}

// Remove all non alpha-numeric char
func RemoveNonAlNum(inString string) string {
	nonAlNum := regexp.MustCompile(`[[:punct:]]`)
	return nonAlNum.ReplaceAllString(inString, "")
}

// replace all [[:space::]] with underscore "_"
func ReplaceSpace(inString string) string {
	spaceRegex := regexp.MustCompile(`[[:space:]]`)
	return spaceRegex.ReplaceAllString(inString, "_")
}

// Format words text to fit (column/windows with limited width) "max" chars.
func FormatText(str string, max int) string {
	tmpLines := make([]string, 0)
	space := regexp.MustCompile(`[[:space:]]`)
	var outText string
	var countChar, length int

	text := space.Split(str, -1) // Split str at each blank char
	for idxWord := 0; idxWord < len(text); idxWord++ {
		length = len(text[idxWord]) + 1
		if countChar+length <= max {
			tmpLines = append(tmpLines, text[idxWord])
			countChar += length
		} else {
			outText += fmt.Sprintln(strings.Join(tmpLines, " ")) // don't use "\n" in case of windows using
			tmpLines = tmpLines[:0]                              // Clear slice
			countChar = 0
			idxWord--
		}
	} // Get the rest of the text.
	outText += fmt.Sprintln(strings.Join(tmpLines, " ")) // don't use "\n" in case of windows using
	return outText
}

// Reduce string length for display (prefix is separator like: "...", option=0 -> put separator at the begening of output string.
// option=1 -> center, is where separation is placed. option=2 -> line feed, trunc the whole string using LF without shorting it.
// max, is max char length of the output string.
func TruncateString(inString, prefix string, max, option int) string {
	var center, cutAt bool
	switch option {
	case 1:
		center = true
		cutAt = false
		max = max - len(prefix)
	case 2:
		center = false
		cutAt = true
	default:
		center = false
		cutAt = false
		max = max - len(prefix)
	}
	length := len(inString)
	if length > max {
		if cutAt {
			var startCount, endCount int
			var tmpString string
			tmpByteString := make([]byte, 0)
			for currCount := 0; currCount <= (length / max); currCount++ {
				if endCount+max < length-1 {
					endCount += max
					tmpByteString = append(tmpByteString, inString[startCount:(endCount)]...)
					tmpString += fmt.Sprintln(string(tmpByteString))
					tmpByteString = tmpByteString[:0]
					startCount += max
				} else {
					endCount = length
					tmpByteString = append(tmpByteString, inString[startCount:(endCount)]...)
					tmpString += fmt.Sprintln(string(tmpByteString))
				}
			}
			inString = tmpString
		} else if center {
			midLength := max / 2
			inString = inString[:midLength] + prefix + inString[length-midLength-1:]
		} else {
			inString = prefix + inString[length-max:]
		}
	}
	return inString
}

// Some multiple way to trim strings. cmds is optionnal or accept multiples args
func TrimSpace(inputString string, cmds ...string) string {
	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`)    //	to match 2 or more whitespace symbols inside a string
	remInsideNoTab := regexp.MustCompile(`[\p{Zs}]{2,}`) //	(preserve \t) to match 2 or more space symbols inside a string
	newstring := inputString
	if len(cmds) != 0 {
		for _, command := range cmds {
			switch command {
			case "-e": //	Un-Escape specials chars
				tmpString, err := strconv.Unquote(inputString)
				Check(err, `Unquote error on`, inputString)
				newstring = tmpString
			case "+e": //	Escape specials chars
				newstring = fmt.Sprintf("%q", inputString)
			case "-c": //	Trim [[:space:]] and clean multi [[:space:]] inside
				newstring = strings.TrimSpace(remInside.ReplaceAllString(inputString, " "))
			case "-ct": //	Trim [[:space:]] and clean multi [[:space:]] inside (preserve TAB)
				newstring = strings.Trim(remInsideNoTab.ReplaceAllString(inputString, " "), " ")
			case "-s": //	To match 2 or more whitespace leading/ending/inside a string (include \t, \n)
				newstring = strings.Join(strings.Fields(inputString), " ")
			case "-&": //	Replace ampersand CHAR with ampersand HTML code
				newstring = strings.Replace(inputString, "&", "&amp;", -1)
			case "+&": //	Replace ampersand HTML code with ampersand CHAR
				newstring = strings.Replace(inputString, "&amp;", "&", -1)
			default:
				err := errors.New(command + `, does not exist`)
				Check(err, `TrimSpace`)
			}
		}
	}
	return newstring
}
