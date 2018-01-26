[![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg)](https://github.com/emersion/stability-badges#experimental)

# Table

Table is a Go package that translates .csv files to nice LaTeX2e tables. Currently it is
very idiosyncratic and the API will change as it evolves.

Go is well suited for the development of Command Line Interfaces as well as the manipulation of text files. It's built-in library provides off-the-shelf libraries for parsing encoded files such as .csv or .json files.

LaTeX has packages that can handle csv data files directly, but for larger files they are limited and tend to slow compilation.

With the package one can export from excel to csv and then use a Go preprocessor to build up the tables. The tables are saved to disk and can then be imported to LaTeX with the `input{<tablename.tex>}` command.

## Requirements LaTeX2e

In the examples I have used some `.sty` files from the `phd` package. These can be by-passed with your own styles, with the exception of the `phd-colorpalette.` This file provides color definitions based on the concept of a color palette. Color palettes are set using the `cxset` command and a sinle key.

```[latex]
\cxset{color palette = Black Tulip}
```

The rest of the packages can all be found in standard LaTeX2e distributions.


## Setting-up the go file

The table package, and I know is not a sexy name, can simply be imported using the `import` statement. See
the example files.

A new table is created by:

```[Go]
  r := table.New()
```

This will initialize a set of default properties for the processor.

### Cleaning the data

Files exported to .csv, especially from excel might need a preprocessing stage, where the data is cleaned
and escaped. This can be done in a singlr operation using:

```go
  r.Clean("<filepath>")
```


### Selecting Columns

Selecting the columns to be rendered can be done in a couple of ways. The easiest is to use 

```go
r.ColumnsByName("5-6", "code", "22-25", "short_description", "long_description", 1)
```




```go
func ExampleSmart() {
	r := table.New()

	r.Clean("smartstatus.csv") 
	r.Caption("Smart City", "Current Smart City Cost Commitments")
	r.RefLabel("smartsystems")

	// Skips the first N lines
	r.SkipN = 1

	// Present table sections
	r.HasSections = true

	// A line consisting of only empty lines
	// is translated to either an empty row or 
	// a rule
	r.EmptyToLine = true

	r.Header.M = [][]string{
		{"ITEM", "DESCRIPTION", "VENDOR", "VALUE", "PROJECTED", "WARRANTY", "MAINT."},
		{"No", "", "", "(QAR)", "COST", "PERIOD", ""},
		//{"A", "B", "C"},
	}

	prop := map[string]string{
		"type":                "longtable",
		"table-align":         "c",
		"font-size":           "footnotesize",
		"font-family":         "sffamily",
		"color":               "thetablevrulecolor",
		"thetableheadcolor":   "thetableheadcolor",
		"thetableheadbgcolor": "thetableheadbgcolor",
		"palette":             "black tulip",
		"tabcolumnsep":        "5pt",
		"extrarowheight":      "2.5pt",
		"arraystretch":        "1.3",
		"rowlines":            "false",
	}

	prop["specifier"] = `{|l|% 
                  >{\RaggedRight}p{3.5cm}|% 
                  >{\RaggedRight}p{3.5cm}|%
                  r|r %
                  |c|r|}%`

	r.Columns([]int{0, 1, 2, 3, 4, 5, 6})
	r.SectionCSV("smart.tex", true, prop)

}
```


### The Tabular Specifier

LaTeX tabular require that we provide a specifier.

```latex
  \begin{longtable}{l l l l l}
  \end{longtable}
```

In Go this is provided as part of the property map. Future versions of this package and one basic reason for its development is to avoid the user to have to type the tabular specifier. If an algorithm can be devised for the Go routines to guess a best looks strategy, then one can go back to processing the tabulars with the TeX primitive `\halign`. This in my estimation can speed up compilation by at least two orders of magnitude.

My current thoughts as to the algorithm is as follows:

1.  Iterate through all the columns, determining the dominating type. 
2. If a column is a field with decimal numbers. We have two choices, one is to use an S or D field from the `siunitx` package or the `ddcolumn` or we can use Go and fmt.Sprintf to print the number. In this case for most applications a right justified field is preferable.
3. Cases where we have long alphanumeric strings, will probably need wrapping. In this case we can use a `p{}` or `X` to typeset the cell. 
4. All others center.

Although one can provide a map of properties I do not favour this approach, as it can get extremely verbose. It is fine if you generating your tables programmatically, as it will be one-off, but I still think it is better to spend some more time on the interface.










