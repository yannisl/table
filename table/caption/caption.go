// Package caption provides an API to LaTeX
// captions. Underlying the package is Axel Sommerfeldt's
// "caption" LaTeXe package. The package also maps
// keys to the phd package interface to use pgf keys
// for styling.
package caption

import (
	"fmt"
)

var (
	labelsep   = []string{"none", "colon", "period", "space", "quad", "newline", "endash"}
	textformat = []string{"empty", "simple", "period"}
)

type Captioner interface {
}

// CaptionStyle is a datastructure for formatting captions of
// figures and tables.
type CaptionStyle struct {
	listentry     string
	heading       string
	format        string // plain hang etc
	indention     string
	labelformat   string
	labelsep      string
	textformat    string
	justification string
	font          string
	labelfont     string
	textfont      string
	margin        string
	style         string
	skips         string
	position      string
	// the name of the current environment
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
}

func (c *CaptionStyle) New() *CaptionStyle {
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

func (c *CaptionStyle) CaptionSetup() {

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
