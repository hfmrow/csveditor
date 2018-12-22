# CsvEditor
*A simple software to Edit, create, modify, search, sort, save csv files.*
- It is not designed to use with large csv files (over 10k entries), at the risk of longer processing times.
- An issue exist with columns list, in option tab or edition window when you have more than 16 fields (columns) and screen size of 1080 pixel height. The window haven't vertical scroll bar so you get some columns and "ok", "cancel" buttons to be invisible. Simply use *[alt]+left mouse button to move whole window*.
- ~~Date format actually can't handle single char. i.e: 4/12/18,  2018/12/4, cannot be used correctly. I will work on a future solution.~~ Resolved.

## How it's made
- Programed with go language: [golang](https://golang.org/doc/) 
- GUI provided by [andlabs go-libui](https://github.com/andlabs/ui), Platform-native GUI library for Go. 

## Fonctionnalities
- Edition, on the fly or using fields window.
- Adding/Removing/Duplicating row.
- Adding/Removing fields.
- Create new csv file.
- Auto detection of field names row.
- Auto set of fields type (string, date, numeric).
- Searching on whole document.
- Sorting by string, date, numeric type. (option must be set for date format and numeric values for decimal separator if needed). The Sort By Date feature allows you to automatically recognize the date in a string.
- Selection of displayed fields.
- Selection of saved fields.
- Option for saving csv: charset type, line-end type, comma character.
- Load csv via command line.

## Some pictures and explanations  

*This is the main screen.*  
![Main](/images/main.png  "Main")  

*New row window. If no line is checked, you get this window. If one or more lines are checked, you get an edit window with a duplicate option for the last checked line. You can check [Cell editing on the fly] at the bottom left to edit a row directly by double-clicking a cell.*  
![NewEntry](/images/newentry.png  "NewEntry")  

*Search option. Everything is said looking at the picture*  
![Search option](/images/search.png  "Search option")  

*Selecting the table, Where you can change the target table, you have 3 modes, main, search and sort. Only the "main" tables have possibility to modify the options, "search" and "sort" can be sorted, searched, saved, edited. If you want to change the column options, you must select the main view, otherwise the option is grayed out.*  
![Table selection](/images/tabsel.png  "Table selection")  

*Options tab. Output file depend on selected options, field types are used in sort tab to define type of sorting, decimal character and date format are also used in it.*  
![Options tab](/images/options.png  "Options tab")  

*Sort tab. The sort order must be taken into consideration to obtain the desired result*  
![Sort tab](/images/sort.png  "Sort tab")  

## How to compile
- Be sure you have golang installed in right way. [Go installation](https://golang.org/doc/install).
- Open terminal window and at command prompt, type: `go get github.com/hfmrow/csveditor`
- If you have a problem with GTK, please check at [andlabs go-libui](https://github.com/andlabs/ui) for gui installation instruction.
	
        Debian, Ubuntu, etc.: sudo apt-get install libgtk-3-dev
        Red Hat/Fedora, etc.: sudo dnf install gtk3-devel


### Misc informations
- I'm working on linuxmint 18.3 (more informations available under release tab).
- I haven't tested compilation under Windows or Mac OS, but all file access functions, line-end manipulations or charset implementation are made with OS portability in mind.  

## You got an issue ?
- Give information about used plateform and OS version.
- Provide a method to reproduce the problem (possibly a sample with problem csv file).
