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

```[Go]
func ExampleSmart() {
	r := table.New()

	r.Clean("smartstatus.csv") //"j56.csv") //"j56.csv")"sub.csv"
	r.Caption("Smart City", "Current Smart City Cost Commitments")
	r.RefLabel("smartsystems")
	// Work on settings.csv.
	//r.SkipFirstN(1).csv
	r.SkipN = 1
	r.HasSections = true
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