/* Package Preamble provides datastructures
   and methods for handling LaTeX2e preambles
*/   
package preamble

import (
   "bytes"
   )   

// Pack represents a LaTeX Package
type Pack {
	Name string
	Options string
	// 
	Region string
	After   string
	Before  string
	buf bytes.Buffer
}


type Preamble struct {
   Regions []string
   Class string
   Packages []string
}

