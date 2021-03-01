package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Node interface{}

type Chardata string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func (e *Element) String() string {
	// attrs := ""
	// for _, v := range e.Attr {
	// 	attrs += fmt.Sprintf("%s=%q  ", v.Name.Local, v.Value)
	// }
	// return fmt.Sprintf("%s %s\n %+s", e.Type.Local, attrs, e.Children)

	b := new(bytes.Buffer)
	prettyPrintNode(b, e, 1)
	return b.String()
}

func prettyPrintNode(w io.Writer, n Node, depth int) {
	switch n := n.(type) {
	case Chardata:
		fmt.Fprintf(w, "%*s%s", depth*2, "", n)
	case *Element:
		attrs := ""
		for _, v := range n.Attr {
			attrs += fmt.Sprintf("%s=%q  ", v.Name.Local, v.Value)
		}
		fmt.Fprintf(w, "%*s%s %s\n", depth*2, "", n.Type.Local, attrs)
		for _, c := range n.Children {
			prettyPrintNode(w, c, depth+1)
		}
	}
}

func main() {

	f, _ := os.Open("d.html")

	dec := xml.NewDecoder(f)

	stack := []*Element{}
	var root Node

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "token: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			child := &Element{
				Type: tok.Name,
				Attr: tok.Copy().Attr,
			}
			if len(stack) == 0 {
				root = child
			} else {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, child)
			}
			stack = append(stack, child)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if len(stack) == 0 {
				fmt.Println("empty stack, Chardata: ", tok)
				continue
			}
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, Chardata(tok))
		}
	}

	fmt.Println(root)
}

// func main() {

// 	f, _ := os.Open("d.html")

// 	// dec := xml.NewDecoder(os.Stdin)
// 	dec := xml.NewDecoder(f)
// 	var depth int
// 	for {
// 		tok, err := dec.Token()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			fmt.Fprintf(os.Stderr, "token: %v\n", err)
// 			os.Exit(1)
// 		}
// 		switch tok := tok.(type) {
// 		case xml.StartElement:
// 			depth++
// 			fmt.Printf("% *s%s\n", depth+1, " ", tok.Name.Local)
// 			if len(tok.Attr) == 0 {
// 				continue
// 			}
// 			fmt.Printf("% *s %s\n", depth+1+1, " ", "Attribs:")

// 			for _, a := range tok.Attr {
// 				fmt.Printf("% *s%s = %q\n", depth+1+1+1, " ", a.Name.Local, a.Value)
// 			}
// 		case xml.EndElement:
// 			depth--
// 		}

// 	}
// }
