// Package table provides a datastructure for manipulating
// array-like structures. It also provides utilities to
// read the tables from CSV and producr pdfs using LuaLaTeX.
package table

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"ml/rules"
	"ml/table/caption"
	"os"
	"reflect"
	"stampcircles/is"
	"stampcircles/util"
	"strconv"
	"strings"
)

const (
	cr      = "\\\\"
	nl      = "\n"
	percent = "%"
)

// Errors returned by Table
var (
	errInvalidFieldNames = errors.New("That was a fatal error my friend")
	errInvalidRange      = errors.New("The provided range in SelectColumns() is invalid, my friend")
)

// Field represent a cell of a table.
type Field struct {
	Name string
	t    string
}

// Trigger is used to trigger subheadings in a longtable.
type Trigger struct {
	Names map[int][]string
}

// Stripe is a struct used to add evn and odd column colors.
type Stripe struct {
	EvenColor, OddColor string
	activate            bool
}

// Activate activates a flag for a striped color table. This can also be
// achieved using EveryOddRow or EveryEvenRow.
func (a *Stripe) Activate() {
	a.activate = true
}

// Head is a structure holding strings that describe the head
// of a table.
type Head struct {
	M [][]string
}

// AddHeader adds a slice to the Header. This is a user function.
// to simplify the user interface we allow either a slice
// to be entered or a series of strings. The latter is
// assumed to be easier for TeX users that might need be familiar
// with go.
func (h *Head) AddHeader(s ...string) {
	h.M = append(h.M, s)
}

// Table is a datastructure for tabular data.
// Data: A Table holds
// three types of data, the raw bytes as read from a file,
// An x-y matrix of cells holding values of data,
// an x-y matrix of styles enabling the typesetting of
// nice looking tables with fine grained control down to
// cell level.
type Table struct {
	// the raw data file
	Raw []byte
	// data holds the table as read
	// from the reader and after a clean
	// operation.
	data [][]string
	// number of columns
	ncols int
	// number of rows
	nrows int
	// the source file path
	fields []Field
	inpath string
	// path for saved cleaned file
	outpath string
	// the name of the table compiled to TeX
	texfile  string
	Err      error
	Notes    []string
	Type     string
	property string
	// rendering properties
	prop map[string]string

	// slice with selected cells
	selector []int

	//
	selectedColumns []interface{}

	// number of first lines to skip
	SkipN int
	// Tables can have section names, these are being picked up
	// and used as subtitles. Defaults to false.
	Header          Head //[][]string
	HasHeader       bool
	HasManualHeader bool
	HeaderLines     int

	HasSections bool
	w           *bufio.Writer
	rd          *csv.Reader
	Labels      []string
	// Triggers will use the key of the map to trigger actions in cells or rows
	// for example the word subtotal can provide a signal to the processor
	// that the row is special and should use a multicolumn.
	// They require a column number.
	// A good trick is to have these standardized and include them in a hidden column zero in
	// an excel sheet.
	Trigger

	// Prints a rule if true and we have an empty line
	EmptyToLine bool

	// Captions
	caption.CaptionStyle

	Index     bool
	Landscape bool
	// cell level
	everyCell       bytes.Buffer
	everyCellBefore bytes.Buffer
	everyCellAfter  bytes.Buffer

	// row
	everyFirstRow bytes.Buffer
	everyRow      bytes.Buffer
	everyLastRow  bytes.Buffer
	everyOddRow   bytes.Buffer
	everyEvenRow  bytes.Buffer
	everyNRow     bytes.Buffer
	NRow          int
	//
	everyFirstColumn bytes.Buffer
	everyColumn      bytes.Buffer
	everyEvenColumn  bytes.Buffer
	everyOddColumn   bytes.Buffer
	everyLastColumn  bytes.Buffer

	// for long tables continuation lines
	EndFirstHead bool
	EndHead      string
	endFoot      string

	// for tables
	FloatStart     string
	FloatSpecifier string

	//
	Stripe

	// parser info
	currentline int
	currentcell int
	// reflabel

}

// New creates a new refernce to a Table.
func New() *Table {
	return &Table{fields: []Field{},
		HasSections: false,
		HasHeader:   false,
		HeaderLines: 1,
		SkipN:       0,
		ncols:       0,
		nrows:       0}
}

// Columns selects the columns to be used.
func (t *Table) Columns(ss []int) {
	t.selector = []int(ss)
	t.ncols = len(ss)
}

// checkRange checks if a string is a valid range
// of integers. If it is it returns a slice with
// the integers.
func checkRange(str string) ([]int, bool, error) {
	var intRange []int
	n := strings.Index(str, "-")
	if n > 0 {
		start, err := strconv.Atoi(str[:n])
		end, err := strconv.Atoi(str[n+1:])
		if err != nil {
			return intRange, false, errInvalidFieldNames
		}
		for i := start; i < end+1; i++ {
			intRange = append(intRange, i)
		}
		return intRange, true, nil
	}
	return intRange, false, nil
}

// validateSelectedColumns, validates the
// user provided selection for errors.
// Valid selections are either all integers or
// all strings. User can also provide ranges
// as "3-5", which selects columns 3-5 inclusive.
func (t *Table) validateSelectedColumns(s ...interface{}) error {
	var intRange []int
	strRange := make([]string, 0)

	for _, v := range s {
		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.String:
			str := value.String()
			iRange, ok, _ := checkRange(str)
			if ok {
				intRange = append(intRange, iRange...)
				fmt.Println(intRange, iRange)
			} else {
				strRange = append(strRange, str)
			}

		//return ErrInvalidFieldNames
		case reflect.Int:
			intRange = append(intRange, int(value.Int()))
		default:
			return errInvalidFieldNames
		}
	}

	// invalid conditions or maybe allowed?
	if len(intRange) > 0 && len(strRange) == 0 || len(intRange) == 0 && len(strRange) > 0 {
		return nil
	}
	fmt.Println(intRange)
	fmt.Println(strRange)
	return nil
}

// ColumnsByName selects columns by their name.
// This we need to blend with Columns eventually.
func (t *Table) ColumnsByName(s ...interface{}) {
	if len(s) < 1 {
		panic("Error you need to select at least 1 column")
	}

	// We have input from the user save it n struct.
	// We will validate later, when we start reading the csv
	// file and we know the number of cells.
	t.selectedColumns = append(t.selectedColumns, s...)

}

// GetColumns gets the Columns selected by the user, validates the input
// and returns a slice of column indices.
// todo.
func (t *Table) GetColumns() {
	err := t.validateSelectedColumns(t.selectedColumns)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Ask Donald Knuth")

	}
}

// Vector maps the selected columns.
func (t *Table) Vector(record []string) []string {

	// Check if we have selected user columns by name
	if len(t.selector) == 0 {
		t.GetColumns()
	}

	vector := make([]string, len(t.selector))

	for k, v := range t.selector {
		vector[k] = record[v]
	}
	return vector
}

// SkipFirstN skips the first n lines from the table.
// The index starts at zero.
func (t *Table) SkipFirstN(n int) {
	t.SkipN = n
}

// skip lines
func (t *Table) skiplines() {
	if t.SkipN > 0 {
		for i := 0; i < t.SkipN; i++ {
			t.rd.Read()
		}
	}
}

// ColumnSpecifier describes a table format specification.
// The idea of a head template to specify  table presentational details
// has seen wide adoption in Knuth's TeX and later LaTeX and Carlisle's
// numerous LaTeXe packages.
type ColumnSpecifier struct {
	s  string
	fn func(s string) string
}

// NewColumnType specifies a new column type. A ColumnType is a string of
// on or more characters.
func NewColumnType(cs ColumnSpecifier) string {
	return "l|r|S|p{3cm}|"
}

// RowSpecifier is a datastructure that uses mark-up at the beginning of
// a table row to map properties to all the cells of a row. It is similar
// to applying every row in pgfplotstable macro. For example many excel
// tables in corporate environments have rows totalizing a number of cells
//
// Reserved strings are the following:
// Total, Sum, GrandTotal, S, T, TRANSFORM MAP COL FORMAT
type RowSpecifier struct {
}

// TableHeader adds the header.
func (t *Table) TableHeader(labels []string) string {
	var buf bytes.Buffer
	color, ok := t.prop["thetableheadbgcolor"]
	for _, v := range labels {
		if ok {
			buf.WriteString(RowColor(color))

		}
		buf.WriteString(v)
	}

	return buf.String()

}

// AddVertSpace adds vertical spacing in a longtable
// by inserting single multicolumn commands. The reason
// why multicolumn is used is to not draw any left or right
// borders that might be specified.
func (t *Table) AddVertSpace(w io.Writer, ncells int) {
	s := `\cellcolor{white}&\multicolumn{` + strconv.Itoa(ncells-1) + `}{l}{\color{white}}\\`
	fmt.Fprintln(w, s)
}

// RowColor returns the command for \rowcolor.
func RowColor(colorname string) string {
	return `\rowcolor{` + colorname + `}`
}

// End function handles the closing of the table
// and any landscape or float wrapper.
// end{tabular}
//      endgroup
//       end{landscape}
//          end{table}
func (t *Table) End(w io.Writer) {
	var buf bytes.Buffer
	out := buf.WriteString
	out(`\end{` + t.Type + `}` + "\n" + `\egroup` + "\n")
	if t.FloatStart != "" {
		out("\\end{" + t.FloatStart + "}\n")
	}

	if t.Landscape {
		out("\\end{" + t.Type + "}%\n")
		out("\\egroup\n")
		out("\\end{landscape}%\n")
		fmt.Fprint(w, buf.String())
		return
	}

	fmt.Fprint(w, buf.String())
}

func tableAlign(prop map[string]string) string {
	if val, ok := prop["table-align"]; ok {
		if val == "" {
			return ""
		}
		return `[` + val + `]`
	}
	return ""
}

// Property sets properties for the head of the table. These
// are local to the table.
func (t *Table) Property(prop map[string]string) string {
	var s bytes.Buffer

	if val, ok := prop["font-size"]; ok {
		s.WriteString(`\` + val + " ")
	}

	if val, ok := prop["font-family"]; ok {
		s.WriteString(`\` + val + " ")
	}

	if val, ok := prop["color"]; ok {
		s.WriteString(`\arrayrulecolor{` + val + `}%` + "\n")
	}

	if val, ok := prop["tabcolumnsep"]; ok {
		s.WriteString(`\setlength{\tabcolsep}{` + val + `}%` + "\n")
	}

	if val, ok := prop["extrarowheight"]; ok {
		s.WriteString(`\setlength{\extrarowheight}{` + val + `}%` + "\n")
	}

	if val, ok := prop["arraystretch"]; ok {
		s.WriteString(`\renewcommand{\arraystretch}{` + val + `}%` + "\n")
	}

	return s.String()
}

// Begin sets the table strings based on the specified properties
func (t *Table) Begin(w io.Writer, prop map[string]string) string {
	var buf bytes.Buffer
	out := buf.WriteString
	floatStart := t.FloatStart
	if floatStart != "" {
		out("\\begin{" + floatStart + "}[" + t.FloatSpecifier + "]\n")
		out("\\centering ")
		out(t.CaptionStyle.String() + " ")
	}
	landscape := ""
	columnColor := "thetablehlcolor"
	columnType := "\\newcolumntype{g}{>{\\columncolor{" + columnColor + "}}r}"
	propstr := t.Property(prop)
	t.property = propstr

	// Write the banner on top of the table
	out(Banner())
	if t.Landscape {
		landscape = "\\begin{landscape}"
	}

	out(landscape)
	out("\\bgroup")
	out(`\cxset{palette ` + prop["palette"] + `}%` + "\n")
	out(columnType)
	out(t.property)
	out(`\begin{` + t.Type + `}` + tableAlign(prop) + prop["specifier"] + "\n")
	if floatStart == "" {
		out(t.CaptionStyle.String() + " ")
		out(t.RefLabelCmd())
	}

	return buf.String()
}

// EveryCell prepends and appends strings to
// every table cell.
func (t *Table) EveryCell(prep, app string) {
	t.everyCellBefore.Reset()
	t.everyCellBefore.Write([]byte("" + prep + " "))
	t.everyCellAfter.Reset()
	t.everyCellAfter.Write([]byte("" + app + " "))
}

// ColorAllRows colors rows by prepending a string on the row text.
func (t *Table) ColorAllRows() string {
	return "\\rowcolor{thetableheadbgcolor!0.25!white}"
}

// GetEveryRow  gets the contents of the everyRow
// buffer and returns the contents as a string.
func (t *Table) GetEveryRow() string {
	ev := t.everyRow.String()
	if len(ev) > 0 {
		return ev
	}
	return ""
}

// EveryRow sets the everyRow buffer, with user
// content.
func (t *Table) EveryRow(s string) {
	t.everyRow.WriteString(s)
}

// EveryNRow sets the everyNRow buffer.
func (t *Table) EveryNRow(n interface{}, s string) {
	t.everyNRow.WriteString(s)
	t.NRow = n.(int)
}

// GetEveryNRow gets the value of the everyNRow buffer if
// any or returns an empty string.
func (t *Table) GetEveryNRow() string {
	if t.NRow > 0 && t.currentline%t.NRow == 0 {
		return t.everyNRow.String()
	}
	return ""
}

// ProcessRow processes rows to add row properties or perform calculations
// and filtering. It is responsible to provide any commands prepended or
// appended to rows with EveryRow type of functions as well as process the
// individual records.
//
// If any updates are required to dataframes or database they are done here.
// Any need for multicolumns, should be carried out here.
func (t *Table) ProcessRow(w io.Writer, record []string) {
	s := t.ColorAllRows() // move to table bg
	// Process Records
	s += t.ProcessRecord(w, record)
	// every nth row
	s += t.GetEveryNRow()
	if t.currentline == 6 {
		s += t.GetEveryRow()
	}
	fmt.Fprint(w, s)
}

// ProcessRecord parses a row and applies
// a series of transformations to add the necessary
// tab symbols and decorations. This is done after
// any sorting of records. It writes its contents to
// an io.Writer.
func (t *Table) ProcessRecord(w io.Writer, record []string) string {
	v := ""
	t.EveryCell("", "") // only as example
	sb := t.everyCellBefore.String()
	sa := t.everyCellAfter.String()

	// prepend and append everycell tokens
	s := sb + record[0] + sa

	for _, v = range record[1:] {
		// handle cell first
		v = strings.TrimSpace(v)

		// format negative numbers
		if is.Numeric(v) {
			v = strings.Replace(v, ",", "", -1)
			if strings.HasPrefix(v, "-") {
				v = "\\textcolor{red}{" + `\num{` + v + `}` + "}"
			} else {
				v = `\num{` + v + `}`
				// process non-numeric cell contents
			}
		}
		v = sb + v + sa

		s += " &" + v + " "
	}

	s += " \\\\\n"
	return s
}

// ReadCSV converts a csv file into a .tex file
// containing filtered fields of a t.
// fields []int,
func (t *Table) ReadCSV(fname string, summation bool, prop map[string]string) {
	var vector []string
	//err :=nil
	fields := t.selector
	f1, _ := os.Create(fname)
	t.w = bufio.NewWriter(f1)
	w := t.w
	t.prop = prop
	t.Type = prop["type"]

	fmt.Fprintln(w, t.Begin(w, prop))

	count := 0
	f, _ := os.Open(t.outpath)

	defer f.Close()
	t.rd = csv.NewReader(f)

	//rd := t.rd
	t.csvDefaultSettings()

	t.skiplines()
	t.renderHead()

	for {
		record, err := t.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
		}

		vector = nil
		for _, v := range fields {
			log.Println(v, record[v])
			vector = append(vector, record[v])
		}

		//if record[8] != "" && record[0] != "" {
		count++

		t.ProcessRow(w, vector)
		if prop["rowlines"] == "true" {
			fmt.Fprintln(w, "\\hline")
		}
	}

	t.closeTabular(w)
}

// Read reads the next line.
func (t *Table) Read() ([]string, error) {
	return t.rd.Read()
}

func (t *Table) csvDefaultSettings() {
	t.rd.Comma = ','
	t.rd.LazyQuotes = true
	//t.rd.FieldsPerRecord = -1 //allows missing
}

// renderHead renders the heading of a table.
func (t *Table) renderHead() {
	var buf bytes.Buffer
	w := t.w

	// Since we just started handle any table headings first.
	// we do this if labels is not empty
	if len(t.Header.M) > 0 {
		t.HasManualHeader = true
	}

	hlines := make([]string, 2)

	switch {

	// labels hold strings with & inclusive
	case len(t.Labels) > 0:
		buf.WriteString(t.TableHeader(t.Labels))
		if t.Type == "longtable" {
			buf.WriteString(t.longTableContinuation(hlines))
		}
		// render to io.Writer
		fmt.Fprint(w, buf.String())

	// Manual headers are supplier by the user as []string. Thanks to Knuth and Carlisle
	// we can use a multicolumn	of length 1 to decorate them individually. They
	// default to |c|. These can be customized if an array of decorators is
	// supplied.
	// >{\columncolor{lightgray}}
	// TODO expand customzation
	case t.HasManualHeader:
		for i := 0; i < len(t.Header.M); i++ {
			mc := "\\multicolumn{1}{|>{\\color{" + "thetableheadcolor" + "}\\bfseries}c|}"
			str := ""
			for j := 0; j < len(t.Header.M[i]); j++ {
				str += mc + "{" + t.Header.M[i][j] + "} "
				if j < len(t.Header.M[i])-1 {
					str += " &"
				}
			}

			lbl := str + " \\\\ \n"
			hlines = append(hlines, lbl)

		}
		// Render the header lines including
		// longtable continuation headers and messages
		// similar to freeze panes
		buf.WriteString(t.TableHeader(hlines))

		// if we have a long table we may have continuation
		// lines. In this case render.

		if t.Type == "longtable" {
			buf.WriteString(t.longTableContinuation(hlines))
		}

		// render to io.Writer
		fmt.Fprint(w, buf.String())

	case t.HasHeader:
		record, err := t.Read()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		for i := 0; i < len(t.Header.M); i++ {
			record, err = t.Read()

			// need to filter here
			lbl := strings.Join(record, " &") + "\\\\ \n"

			hlines = append(hlines, lbl)
			fmt.Fprint(w, t.TableHeader(hlines))
			log.Println(record)
		}

	}
}

// renders a TeX comment line
func wcomment(s string) string {
	return fmt.Sprintf("%% %s\n", s)
}

// longTableContinuation uses commands applicable only for long tables.
func (t *Table) longTableContinuation(hlines []string) string {
	var buf bytes.Buffer
	buf.WriteString("\\endfirsthead\n")
	buf.WriteString(`\multicolumn{` + strconv.Itoa(len(t.Header.M[0])) + `}{l}{\ldots continued from previous page}\\` + "\n")
	buf.WriteString(t.TableHeader(hlines) + "\n")
	buf.WriteString(wcomment("ends the head"))
	buf.WriteString("\\endhead" + "\n")
	buf.WriteString(`\multicolumn{` + strconv.Itoa(len(t.Header.M[0])) + `}{r}{continued to next page\ldots}\\` + "\n")
	buf.WriteString(wcomment("ends the footer"))
	buf.WriteString(`\endfoot` + "\n")
	buf.WriteString(wcomment("supresses message on last line"))
	buf.WriteString(`\endlastfoot` + "\n")
	return buf.String()
}

// SectionCSV converts a csv file into a tex file
// It handles longtbales with sections (they look more like documents).
func (t *Table) SectionCSV(fname string, summation bool, prop map[string]string) {

	//fields := t.selector
	f1, _ := os.Create(fname)
	t.w = bufio.NewWriter(f1)
	w := t.w
	t.prop = prop
	t.Type = prop["type"]

	fmt.Fprintln(w, t.Begin(w, prop))
	f, _ := os.Open(t.outpath)
	defer f.Close()

	// csv settings
	t.rd = csv.NewReader(f)
	t.csvDefaultSettings()
	t.skiplines()
	t.renderHead()

	t.currentline = 0
	var inHead = false

	for {

		t.currentline++
		record, err := t.Read()

		// save the number of cells we will need it later
		t.ncols = len(record)

		// skip empty lines
		if len(strings.Join(record, "")) == 0 {
			log.Println("EMPTY RECORD DETECTION")
			if t.EmptyToLine {
				fmt.Fprintln(w, "\\hline")
			}
			record, err = t.Read()
		}

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
		}

		vector := t.Vector(record)
		// needs fixing
		vector[3] = PrintTitleCase(vector[3])

		// detects empty lines find a better method
		condition := false
		if len(record) > 8 {
			if record[8] != "" && record[0] != "" {
				condition = true
			}
		} else {
			condition = true
		}

		if condition {

			// If we have a trigger word we need to take action
			if strings.HasPrefix(strings.TrimSpace(record[1]), "SUBTOTAL") ||
				strings.HasPrefix(strings.TrimSpace(record[1]), "SUBCONTRACTS") ||
				strings.HasPrefix(strings.TrimSpace(record[1]), "MATERIALS") ||
				strings.HasPrefix(strings.TrimSpace(record[1]), "TEC-VAR") ||
				strings.HasPrefix(strings.TrimSpace(record[1]), "CONTRACTS") ||
				strings.HasPrefix(strings.TrimSpace(record[1]), "CINEMA") ||
				strings.HasPrefix(strings.TrimSpace(record[1]), "GRAND") ||
				strings.HasPrefix(strings.TrimSpace(record[1]), "TOTAL") {

				//continue
				// check if we need midrule
				//t.PrintMidrule(w, "1.5pt")
				fmt.Fprintln(w, rules.MidRule("1.5pt"))
				mult := "&\\multicolumn{3}{|p{3.5cm}|}{\\textbf{%s}}"

				tmp := fmt.Sprintf("%s", vector[0]) //ok

				tmp += fmt.Sprintf(mult, vector[1]) //ok

				for i := 4; i < len(vector); i++ {
					tmp += fmt.Sprintf(" &%s", vector[i])
				}
				tmp += " \\\\\n"
				fmt.Fprintln(w, tmp)
				//log.Println(tmp)

				fmt.Fprintln(w, rules.MidRule("1.5pt"))
				t.AddVertSpace(w, len(vector))
				inHead = true

				//fmt.Fprintln(w, "\\bottomrule")

			} else {

				if inHead == true {

					if t.HasSections == true {
						//t.AddVertSpace(w, len(vector))

						sect := GetSectionTitle(record[1])
						fmt.Fprintf(w, "%s\n", AddSection(sect, len(vector)))
						fmt.Fprintln(w, rules.MidRule("1.5pt"))
						//t.TableHeader(w, labels)
						t.ProcessRow(w, vector)
						inHead = false
					}
				}

				// Print non-header, non-summation lines
				// PROCESS ROW FIRST
				t.ProcessRow(w, vector)
				if prop["rowlines"] == "true" {
					fmt.Fprintln(w, "\\hline")
				}
				inHead = false

			}
		}

	}
	// Finally we close the tabular
	t.closeTabular(w)
}

// closes the table environment
func (t *Table) closeTabular(w *bufio.Writer) {
	fmt.Fprintln(w, rules.AddLineSpace("0pt"))
	fmt.Fprintln(w, rules.BottomRule())
	t.End(w)
	w.Flush() // do not forget to flush the buffer
}

// Clean satisfies the Cleaner interface
// It takes a filename or path and performs
// a number of replacements and transformations
// to clean the raw text from common errors.
func (t *Table) Clean(fname string) []byte {
	t.inpath = fname
	outfile := ""
	s, _ := ioutil.ReadFile(fname)
	s1 := strings.Replace(string(s), "\r\n", "\n", -1)
	s1 = strings.Replace(string(s1), "&", "\\&", -1)
	s1 = strings.Replace(string(s1), "%", "\\%", -1)
	s1 = strings.Replace(string(s1), "#DIV/0!", "0.0", -1)
	s1 = strings.Replace(string(s1), "#", "\\#", -1)

	if strings.HasSuffix(fname, ".csv") {
		outfile = strings.TrimSuffix(fname, ".csv") + "a" + ".csv"
	}
	err := ioutil.WriteFile(outfile, []byte(s1), 0666)
	t.outpath = outfile
	if err != nil {
		fmt.Println(err)
	}
	return t.Raw
}

// AddSection adds a section as a Table row in a long
// table by using a multicolumn control sequence.
func AddSection(s string, ncells int) string {
	return `\multicolumn{` + strconv.Itoa(ncells) + `}{l}{\textbf{` + s + `}}\\`
}

// GetSectionTitle is a weird feature but necessary of we are exporting from
// excel for enterprise style spreadsheets. These in many instances
// have numerous rows with subtotals, notes and the like.
// Section titles are input as map[string]string and processed here.
// They are simple replacement routines, but they are also triggers
// that can typeset the row contents differently using multicolumn.
// This is an example, will be replaced by a settings function in the
// next revision.
func GetSectionTitle(s string) string {
	switch {
	case strings.HasPrefix(s, "TEC"):
		return "SECTION THEME EXPERIENCE CENTER"
	case strings.HasPrefix(s, "SUB-SMR"):
		return "SECTION SMART SYSTEMS"
	case strings.HasPrefix(s, "SUB-E"):
		return "SECTION SUBCONTRACTS IN CORE CONTRACT"
	case strings.HasPrefix(s, "SUB-M"):
		return "SECTION SUBCONTRACT PS"
	case strings.HasPrefix(s, "GEN"):
		return "SECTION INDIRECT SITE COSTS"
	case strings.HasPrefix(s, "PC-"):
		return "SECTION PC SUMS (All)"
	case strings.HasPrefix(s, "HVAC"):
		return "SECTION HVAC"
	case strings.HasPrefix(s, "EL-EQ"):
		return "SECTION ELECTRICAL"
	case strings.HasPrefix(s, "SUBCONTRACTS"):
		return "SECTION SUBCONTRACTS"
	default:
		return "SECTION UNKNOWN"
	}

}

// PrintTitleCase prints a string as Title Cased. If the word
// is fully capitalized it will change it to a lower case and
// then if it is not in a lisyt of abbreviations that are normally
// capitalized it will Title Case it.
func PrintTitleCase(s string) string {
	// words that must remain capitalized
	abbr := []string{
		"mep:", "lv", "elv", "tec", "fm-200", "epon", "gpon", "bms",
		"av", "bgm", "cctv", "smatv", "it", "mv", "micc", "gi", "pvc",
		"fp200", "cfd", "ad", "vcd", "fd", "mfd", "sd", "hvac", "chw", "grp",
		"ff", "hv", "edms", "erms", "lpg", "dc", "ac",
		"mdb", "acb", "smdb", "db", "mcc", "pfcu", "fp", "rmu", "ups",
		"co", "dx", "ss", "dne", "ahu", "ahus", "fcu", "fcus", "mcr",
		"hdpe", "pc", "ps", "fas", "pava", "cbs", "ip", "pc", "ps", "it",
	}
	str := strings.ToLower(strings.TrimSpace(s))
	sstr := strings.Split(str, " ")
	temp := ""
	for _, v := range sstr {
		if utils.HasAnySuffix(v, abbr) {
			temp += strings.ToUpper(v) + " "
		} else {
			temp += strings.Title(v) + " "
		}
	}
	str = temp
	return str
}

// Banner prepends a banner to the table
// indicating to the user that the tables have been
// generated automatically.
func Banner() string {
	s := fmt.Sprintf("%%%% This file has been generated from a csv file.\n")
	s += fmt.Sprintf("%%%% automatically by the phd-cli = %s\n", "version 0.20")
	s += fmt.Sprintf("%%%% to get help to regenerate type phd-cli help\n")
	return s
}
