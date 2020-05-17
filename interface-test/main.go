package main

import (
	"log"
)

type Interface struct {
	Field string
}

type ExpandedInterface struct {
	Field string
	NewField string
}

type ParentInterface struct {
	SimilarInterface ExpandedInterface
	Interface Interface
}

func main() {
	obj := ParentInterface{
		SimilarInterface: ExpandedInterface {
			Field: "similar interface",
			NewField: "similar interface",
		},
		Interface: Interface {
			Field: "interface",
		},
	}

	log.Printf("%#v", obj)
	// This doesn't work
	// obj.SimilarInterface.Echo()

	// This works
	obj.SimilarInterface.Echo()
}

func (Interface)Echo() {
	log.Printf("Inside interface")
}