// Package caption provides an API to LaTeX
// captions. Underlying the package is Axel Sommerfeldt's
// "caption" LaTeXe package. The package also maps
// keys to the phd package interface which uses a declarative
// approach for styling documents.
//
// A caption can be rendered in many different ways.
//
//     Table 1.1 This is a caption
//     figure 2:  This is a caption for
//        my scientific figure.
//
// The long term intention of the package is to provide
// renderers for html, tex, pdf, docx as well as create a
// component for React or Vue javascript renderers.
// document independent 
package caption

import (
	"bytes"
	"fmt"
	//"bufio"
	//"github.com/golang/net/html/atom"
)

const (
	NONE  = "none"
	COLON = "colon"
)

// A TokenType is the type of a Token.
type TokenType uint32

const (
	// ErrorToken means that an error occurred during tokenization.
	ErrorToken TokenType = iota
	// TextToken means a text node.
	TextToken
	// A StartTagToken looks like <a>. \bgroup
	StartTagToken
	// An EndTagToken looks like </a>. \egroup
	EndTagToken
	// A SelfClosingTagToken tag looks like <br/>. \section[]{}...{} 
	SelfClosingTagToken
	// A CommentToken looks like <!--x-->. in HTML or % for TeX
	CommentToken
	// A DoctypeToken looks like <!DOCTYPE x> this is class in LaTeX not used for TeX
	DoctypeToken
)

// A Token consists of a TokenType and some Data (tag name for start and end
// tags, content for text, comments and doctypes). A tag Token may also contain
// a slice of Attributes. Data is unescaped for all Tokens (it looks like "a<b"
// rather than "a&lt;b"). For tag Tokens, DataAtom is the atom for Data, or
// zero if Data is not a known tag name.
type Token struct {
	Type     TokenType
	//DataAtom atom.Atom
	Data     string
	Attr     []Attribute
}

type Element struct {
	toks []Token
}

var (
	labelsep   = []string{NONE, COLON, "period", "space", "quad", "newline", "endash"}
	textformat = []string{"empty", "simple", "period"}
)

// Renderer provides the interace for all captioning rendering. I will extend it
// to a full interface for the document at a later stage.
//
//
type Renderer interface {
	// adds the setup to the preamble or relevant css strings
	AddToPreamble(out *bytes.Buffer, options map[string]string)
	Style()
	SetCommand()
}

// html renderers
type html struct {
	// inherit all methods. We will override later.
	Renderer
}

func HtmlRenderer(out *bytes.Buffer, options map[string]string) Renderer {
	return &html{}
}

func TeXRenderer(out *bytes.Buffer, options map[string]string) Renderer {
	cs := New()
	return cs
}

// AddToPreamble implements the Renderer interface. It adds the package to the
// LaTeX preamble. The bytes.Buffer refers to region 2 of the preamble template.
func (c *CaptionStyle) AddToPreamble(out *bytes.Buffer, options map[string]string) {
	out.WriteString("\\usepackage{caption}\n")
	out.WriteString("\\CaptionSetUp{}\n")
}

func (c *CaptionStyle) Style() {

}

func (c *CaptionStyle) SetCommand() {

}

// Margin is a struct representing the margins of the
// container.
type Margin struct {
	Left              string
	Right             string
	oneside, twosside string
	// The top and bottom are actually skips
	Top, Bottom string
	buf         bytes.Buffer
}

// Setting sets the property value, as well
// as append it to the buffer
// TODO
func (m *Margin) SetMargin(s ...string) string {
	if len(s) < 1 {
		return ""
	}
	m.Left = s[0]
	m.buf.WriteString("margin=" + s[0])
	return m.buf.String()
}

// An Attribute is an attribute namespace-key-value triple. Namespace is
// non-empty for foreign attributes like xlink, Key is alphabetic (and hence
// does not contain escapable characters like '&', '<' or '>'), and Val is
// unescaped (it looks like "a<b" rather than "a&lt;b").
//
// Namespace is only used by the parser, not the tokenizer.
// Not used for all elements or I can use for the element
type Attribute struct {
	Namespace, Key, Val string
}


// Not used at the moment
type Font struct {
	Name        string
	Series      string
	Shape       string
	FontFace    string
	FontFamily  string
	FontWeight  string
	FontVariant string
	FontSize    string
	Typ         string
}

// use atoms to refer to attributes
var dict = map[string]string{"font-size": "10pt",
							 "font-weight": "normal",
						}
	
// SetProperty sets any property
// f.SetProperty("font-size", "12pt")
func (f *Font) SetProperty(s1, s2 string) {
	key:=s1
	if _, ok := dict[key]; ok {
       dict[key]=s2
	}
}

func (f *Font) GetProperty(key string) string {
	if val, ok := dict[key]; ok {
       return val
	}
	return ""
}

// How to handle document wise properties
// doc:=NewDocument()
// doc.Caption.SetProperty()
// doc.Head1.SetProperty()
// doc.Head2.SetProperty()
// doc.Paragraph.Properties("parindent": "0pt", etc)

// CaptionStyle is a datastructure for formatting captions of
// figures and tables.
type CaptionStyle struct {
	listentry string
	heading   string

	format string // plain hang etc

	// Indent from second line onwards.
	indention     string
	labelformat   string
	labelsep      string
	textformat    string
	justification string

	// fonts and text decorations
	Font
	labelfont string
	textfont  string

	// Caption margin
	Margin

	// Width
	Width    string
	style    string
	skips    string
	position string

	// parskip only useful for
	parskip string
	// the name of the current environment
	hangindent string
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

	buf bytes.Buffer
}

func New() *CaptionStyle {
	c := &CaptionStyle{}
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

// DeclareCaptionStyle creates a new named caption style. Think of it as
// a class in css.
func DeclareCaptionStyle() {

}

// Caption sets the caption of an environment such as that of a figure
// or a table. This is the user entry command.
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
