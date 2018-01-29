package main

import (
	"fmt"
	"golang.org/x/text/message"
	"log"
	"ml/table"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func printCommand(cmd *exec.Cmd) {
	log.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}
func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

// openBrowser tries to open the URL in a browser,
// and returns whether it succeed in doing so.
func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		fmt.Println("opening in browser", url)
		args = []string{"cmd", "/c", "start chrome"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

// Groups
func Example() {
	r := table.New()
	r.Clean("groups.csv") //"j56.csv") //"j56.csv")"sub.csv"
	// Work on settings
	//r.SkipFirstN(1)
	r.SkipN = 1
	r.Landscape = false
	r.Caption("Group Codes")
	r.RefLabel("tbl:groups")
	r.Header.M = [][]string{
		{"CODE", "SHORT DESCRIPTION", "LONG DESCRIPTION"},
		{"A", "B", "C"},
		{"D", "E", "F"},
	}
	r.Columns([]int{0, 1, 2})
	r.ColumnsByName("5-6", "code", "22-25", "short_description", "long_description", 1)
	prop := map[string]string{
		"type":                "longtable",
		"table-align":         "l",
		"font-size":           "footnotesize",
		"font-family":         "sffamily",
		"color":               "thetablevrulecolor",
		"thetableheadcolor":   "thetableheadcolor",
		"thetableheadbgcolor": "thetableheadbgcolor",
		"palette":             "spring onion",
		"tabcolumnsep":        "5pt",
		"extrarowheight":      "2.5pt",
		"arraystretch":        "1.3",
		"rowlines":            "false",
	}

	// Tabular specifier needs to be specified.
	// if left empty we just provide a default
	// for all the strings to be centered.
	prop["specifier"] = `{l%serial 
                  %@{\extracolsep{\fill}}
                  l|% 
                  l|}%`

	// We start procesing here
	r.ReadCSV("groups.tex", true, prop)

}

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

func main() {

	Example()
	ExampleSmart()

	//
	r := table.New()
	r.Clean("j56.csv")
	r.Caption("Material Costs")
	// do load
	r.SkipN = 4

	r.HasSections = true

	r.Header.M = [][]string{
		{"SR", "CAT", "CODE", "DESCRIPTION", "BUDGET", "BUDGET", "CURRENT", "COST AT", "PROJECT."},
		{"", "", "", "", "", "ADJUST.", "BUDGET", "SUBMITTAL", "COST"},
		{"", "", "", "", "", "(QAR)", "(QAR)", "(QAR)", "(QAR)"},
	}

	var tableAlign = "table-align"
	prop := map[string]string{
		"type":                "longtable",
		tableAlign:            "l",
		"font-size":           "footnotesize",
		"font-family":         "sffamily",
		"color":               "thetablevrulecolor",
		"palette":             "black tulip",
		"tabcolumnsep":        "2.5pt",
		"extrarowheight":      "2.5pt",
		"arraystretch":        "1.3",
		"rowlines":            "false",
		"thetableheadcolor":   "thetableheadcolor",
		"thetableheadbgcolor": "thetableheadbgcolor",
		"ignore chars":        "", //lsit of characters to ignore in cells
		"comment chars":       "", //check if in csv
	}

	prop["specifier"] = `{p{.6cm} 
                  %@{\extracolsep{\fill}}
                  p{1cm}|% 
                  r|%
                  >{\RaggedRight}p{2.4cm}|%
                  r|%
                  r|%
                  r|%
                  r?%
                  gV}%`

	r.Trigger.Names = map[int][]string{
		1: {"PREFIX", "GEN", "GENERAL"},
		0: {"CONTAINS", "GEN", "GENERAL"},
	}

	r.Columns([]int{0, 1, 9, 8, 11, 12, 13, 14, 15})
	r.SectionCSV("materials.tex", true, prop)

	r.Header.M = [][]string{
		{"Sr", "Cat", "Code", "Description", "Projected", "Cumulative", "Balance", "Delivered", "Percent"},
		{"", "", "", "", "Cost", "Orders", "Orders", "Orders", ""},
	}
	//r.Labels = []string{s1, s2}
	r.EveryRow("\\hline \n")
	// using the arydshln
	//r.EveryNRow(5, "\\hdashline ")
	prop["palette"] = "oprah"

	prop["specifier"] = `{l:%serial 
                  p{0.8cm}:% 
                  @{}r:%
                  >{\RaggedRight}p{2.2cm}:%
                  r:%
                  r:%
                  r:%
                  r:%
                  l:}%`

	r.Columns([]int{0, 1, 9, 8, 15, 16, 17, 18, 19})

	r.SectionCSV("materials-summary.tex", true, prop)

	prop["specifier"] = `{@{\extracolsep{\fill}}|r|%serial 
                   p{1.5cm}|% 
                   r|%
                   >{\RaggedRight}p{5.2cm}|%
                   }`

	prop["fontsize"] = "small"
	prop["table-align"] = "c"
	prop["palette"] = "black tulip"
	prop["thetableheadbgcolor"] = "thetableheadbgcolor"

	r.Columns([]int{0, 1, 9, 10})
	r.ReadCSV("codes.tex", true, prop)

	r.Clean("test.csv")
	r.Columns([]int{0, 1, 2, 3, 4})
	r.Header.M = [][]string{
		{"Sr", "Cat", "Code", "Description", "Other"},
		{"", "", "", "", "Tests"},
	}
	r.Caption("A Test Table")
	r.FloatStart = "table"
	r.FloatSpecifier = "htbp"
	prop["type"] = "tabularx"
	prop["palette"] = "zealous"
	prop["table-align"] = ""
	prop["specifier"] = "{15cm}[t]{|l|l|l|l|Y|}"
	prop["font-size"] = "Large"

	r.Stripe.Activate()
	r.ReadCSV("test.tex", true, prop)

	//"--interaction=nonstopmode",
	cmd := exec.Command("lualatex", "mat.tex", "&& start chrome.exe", "mat.pdf")
	printCommand(cmd)
	printCommand(cmd)

	p := message.NewPrinter(message.MatchLanguage("en"))
	p.Println(123456.78) // Prints 123,456.78

	p.Printf("%d ducks in a row", 4331) // Prints 4,331 ducks in a row

	p = message.NewPrinter(message.MatchLanguage("nl"))
	p.Println("Hoogte: %f meter", 1244.9) // Prints Hoogte: 1.244,9 meter

	p = message.NewPrinter(message.MatchLanguage("bn"))
	p.Println(123456.78) // Prints ১,২৩,৪৫৬.৭৮

	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)
	openBrowser("mat.pdf")

}
