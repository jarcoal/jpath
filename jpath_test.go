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
	            "title": "Sayings of the Century"
	        },
	        {
	            "author": "Evelyn Waugh",
	            "category": "fiction",
	            "price": 12.99,
	            "title": "Sword of Honour"
	        },
	        {
	            "author": "Herman Melville",
	            "category": "fiction",
	            "isbn": "0-553-21311-3",
	            "price": 8.99,
	            "title": "Moby Dick"
	        },
	        {
	            "author": "J. R. R. Tolkien",
	            "category": "fiction",
	            "isbn": "0-395-19395-8",
	            "price": 22.99,
	            "title": "The Lord of the Rings"
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

func TestAttributeSelector(t *testing.T) {
	color, ok := document.String("$.store.bicycle.color")
	if !ok {
		t.Fatal("expected ok to be true")
	}

	if color != "red" {
		t.Fatal("expected value to be 'red', got %s", color)
	}
}

func TestDescendentSelector(t *testing.T) {
	// all prices
	{
		results := document.Floats("$..price")

		if len(results) != 5 {
			t.Fatalf("expected 5 results, got %v", len(results))
		}
	}

	// just book prices
	{
		results := document.Floats("$..book..price")

		if len(results) != 4 {
			t.Fatalf("expected 4 results, got %v", len(results))
		}
	}
}

func TestIndexSelector(t *testing.T) {
	title, ok := document.String("$.store.book[0].title")
	if !ok {
		t.Fatal("expected ok to be true")
	}

	if title != "Sayings of the Century" {
		t.Fatalf("expected value to be 'Sayings of the Century', got %s", title)
	}
}

func TestUnmarshal(t *testing.T) {
	st := new(struct {
		Color       string    `jpath:"$.store.bicycle.color"`
		Prices      []float64 `jpath:"$..price"`
		FirstAuthor string    `jpath:"$..author"`
	})

	if err := Unmarshal(documentBytes, st); err != nil {
		t.Fatal(err.Error())
	}

	if st.Color != "red" {
		t.Fatalf("expected Color to be '%s', got '%s'", st.Color)
	}

	if len(st.Prices) != 5 {
		t.Fatalf("expected Prices to have 5 entries, got %v", len(st.Prices))
	}

	if st.FirstAuthor != "Nigel Rees" {
		t.Fatalf("expected FirstAuthor to be '%s', got '%s'", st.FirstAuthor)
	}
}

func TestInvalidUnmarshal(t *testing.T) {
	// needs to be a struct
	{
		m := make(map[string]interface{})
		if err := Unmarshal(documentBytes, &m); err == nil {
			t.Fatal("expected an error to be returned")
		}
	}

	// price is a float64, not a string
	{
		st := new(struct {
			Price string `jpath:"$..price"`
		})

		if err := Unmarshal(documentBytes, st); err == nil {
			t.Fatal("expected an error to be returned")
		}
	}

	// authors are strings, not floats
	{
		st := new(struct {
			Authors []float64 `jpath:"$..author"`
		})

		if err := Unmarshal(documentBytes, st); err == nil {
			t.Fatal("expected an error to be returned")
		}
	}
}
