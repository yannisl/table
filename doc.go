package main

// Package
//
// Triggers
// Triggers can be specified on the first cell or the first word. Actions
// can then be defined to act on these triggers.
// For example
// Any word starting from a word: will be mapped. Then the map values will be available for
// typesetting either before or after the first run. These can slow the system down by requiring
// two passes within a table they will be mostly translated to multicolumns.
//
//      Project: Doha Oasis
//      Title:....
//      Date:...
//
//
// Selecting Columns
// Particular Columns can be selected by using tbl.Columns(<comma separated list>)
// If this option is empty, all available columns will be provided.
//
// 		column name/.style={key-value-list}
//
// 		hcell[0]= map[sting]string{"align"="l", color="blue",..., }
//
// multicolumn names
// sort
// everyfirstcolumn
// everylastcolumn
// everyevencolumn
// everyoddcolumn
// every column
//
// Configuring Row Styles

// before row
// everyheadrow
// everylastrow
//
// after row
//
// Styling Cells
//
// Besides the possibilities to change column styles and row styles, there are also a
// methods and settings for single cells.
//
// Data Cells can be referenced as cell[x,y] style[row index, column index]=map[string]string{<key values>}
//
// caption
// BeginTableAdd()
//
// This will write the tabular code in a macro rather than a file. YOu will need to provide
// the macro name.
// WriteToMacro()
// Debug()
