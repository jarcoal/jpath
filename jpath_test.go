package jpath

import (
	"testing"
)

var documentBytes = []byte(`{
	"store": {
	    "bicycle": {
	        "color": "red",
	        "price": 19.95
	    },
	    "book": [
	        {
	            "author": "Nigel Rees",
	            "category": "reference",
	            "price": 8.95,
	            "title": "Sayings of the Century",
	            "available": true
	        },
	        {
	            "author": "Evelyn Waugh",
	            "category": "fiction",
	            "price": 12.99,
	            "title": "Sword of Honour",
	            "available": true
	        },
	        {
	            "author": "Herman Melville",
	            "category": "fiction",
	            "isbn": "0-553-21311-3",
	            "price": 8.99,
	            "title": "Moby Dick",
	            "available": false
	        },
	        {
	            "author": "J. R. R. Tolkien",
	            "category": "fiction",
	            "isbn": "0-395-19395-8",
	            "price": 22.99,
	            "title": "The Lord of the Rings",
	            "available": false
	        }
	    ]
	}
}`)

var document *Jpath

func init() {
	var err error
	if document, err = NewFromBytes(documentBytes); err != nil {
		panic(err.Error())
	}
}

func TestStrings(t *testing.T) {
	if len(document.Strings("$..title")) != 4 {
		t.FailNow()
	}

	if len(document.Strings("$..price")) != 0 {
		t.FailNow()
	}
}

func TestString(t *testing.T) {
	if _, ok := document.String("$..title"); !ok {
		t.FailNow()
	}

	if _, ok := document.String("$..price"); ok {
		t.FailNow()
	}
}

func TestBools(t *testing.T) {
	if len(document.Bools("$..available")) != 4 {
		t.FailNow()
	}

	if len(document.Bools("$..author")) != 0 {
		t.FailNow()
	}
}

func TestBool(t *testing.T) {
	if _, ok := document.Bool("$..available"); !ok {
		t.FailNow()
	}

	if _, ok := document.Bool("$..isbn"); ok {
		t.FailNow()
	}
}

func TestFloats(t *testing.T) {
	if len(document.Floats("$..price")) != 5 {
		t.FailNow()
	}

	if len(document.Floats("$..available")) != 0 {
		t.FailNow()
	}
}

func TestFloat(t *testing.T) {
	if _, ok := document.Float("$..price"); !ok {
		t.FailNow()
	}

	if _, ok := document.Float("$..available"); ok {
		t.FailNow()
	}
}
