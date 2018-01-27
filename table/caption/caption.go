// Package caption provides an API to LaTeX
// captions. Underlying the package is Axel Sommerfeldt's
// "caption" LaTeXe package. The package also maps
// keys to the phd package interface which uses a declarative
// approach for styling documents.
package caption

import (
	"fmt"
	"bytes"
	//"bufio"
)

var (
	labelsep   = []string{"none", "colon", "period", "space", "quad", "newline", "endash"}
	textformat = []string{"empty", "simple", "period"}
)

type Captioner interface {
}


type Margin struct {
	Left string
	Right string
	oneside, twosside string
	// The top and bottom are actually skips
	Top, Bottom string
	buf bytes.Buffer
}

// Setting sets the property value, as well
// as append it to the buffer
func (m *Margin) SetMargin(s string) {
	m.Left = s
	m.buf.WriteString("margin="+s)
}


// CaptionStyle is a datastructure for formatting captions of
// figures and tables.
type CaptionStyle struct {
	listentry     string
	heading       string

	format        string // plain hang etc
	
	// Indent from second line onwards.
	indention     string
	labelformat   string
	labelsep      string
	textformat    string
	justification string
	font          string
	labelfont     string
	textfont      string

	// Caption margin
	Margin       

	// Width 
	Width         string
	style         string
	skips         string
	position      string

	// parskip only useful for
	parskip       string
	// the name of the current environment
	hangindent    string
	// name=Fig.
	name string
	// the \caption command can typeset captions for
	// different types
	typ string
	// for debugging purposes will generate a log file
	// entry showing the caption setup.
	showcaptionsetup bool
	// if there is no caption the label cannot be referenced
	// hence included here
	reflabel string

	buf  bytes.Buffer
}


func (c *CaptionStyle) New() *CaptionStyle {
	c.Margin.Left = "0pt"
	c.Margin.Right = "0pt"
	c.Width = "\\textwidth"
	c.listentry = ""
	c.heading = ""
	c.format = ""
	return c
}



// RefLabel sets the refLabel field. It is used
// to reference a table by a label. See LaTeX
// \ref{} and \label{}
func (c *CaptionStyle) RefLabel(s string) {
	c.reflabel = s
}

// RefLabelCmd returns a "\label" command as a string.
// If the caption is empty it returns an empty string.
// If the user defined it as an empty string we skip.
func (c *CaptionStyle) RefLabelCmd() string {
	// if a label is empty skip
	if c.reflabel != "" {
		return fmt.Sprintf("\\label{%s}\\\\", c.reflabel)
	}
	// if a heading is empty then we skip
	if c.heading != "" {
		return "\\\\" //close caption
	}

	return ""
}

func DeclareCaptionStyle() {

}

// Implement the string interface
func (c *CaptionStyle) String() string {
	// Initialize all Styles by setting them to base
	c.buf.WriteString("\\CaptionStyle{")
	c.buf.WriteString("format=base")
	c.buf.WriteString("}%%\n")
}


// Caption sets the caption of an environment such as that of a figure
// or a table.
func (c *CaptionStyle) Caption(s ...string) string {
	// we do not have any special list entry text
	if len(s) == 1 {
		c.heading = s[0]

	}
	if len(s) == 2 {
		c.listentry = s[0]
		c.heading = s[1]
	}

	return ""
}

// String renders the caption using the stringer interface
// used here to provide idiomatic Go to users
func (c *CaptionStyle) String() string {
	if c.listentry != "" {
		return fmt.Sprintf("\\caption[%s]{%s}", c.listentry, c.heading)
	}

	if c.heading != "" {
		return fmt.Sprintf("\\caption{%s}", c.heading)
	}

	// user did not specify any caption
	return ""
}
