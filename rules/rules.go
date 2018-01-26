// Package rules provides datastructures and
// methods for rendering LaTeX rules.
//
// In British Typesetting a "line" is always called a "rule". The
// thickness of the rule is often referred to as its "width". You
// can find more about rules in the LaTeX2e package "booktabs"
// developed by Simon Fear.
package rules

import (
	"fmt"
)

//
type Rule struct {
	Name   string
	Width  string
	Height string
	Cmd    string
	Typ    string
}

func New() *Rule {
	r := Rule{}
	r.Name = "toprule"
	r.Width = "auto"
	r.Height = "auto"
	r.Cmd = "\\toprule "
	r.Typ = "booktabs"
	return &r
}

func render(name string, s ...string) string {
	if len(s) > 0 {
		return fmt.Sprintf("\\%s[%s] ", name, s[0])
	}
	return fmt.Sprintf("\\%s", name)
}

// Rules take an optional agrument denoting their
// width.
func TopRule(s ...string) string {
	return render("toprule", s...)

}

func MidRule(s ...string) string {
	return render("midrule", s...)
}

func BottomRule(s ...string) string {
	return render("bottomrule", s...)
}

func AddLineSpace(s ...string) string {
	return render("addlinespace", s...)
}
