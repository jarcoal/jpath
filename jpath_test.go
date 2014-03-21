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
	// simple index access
	{
		title, ok := document.String("$.store.book[0].title")
		if !ok {
			t.Fatal("expected ok to be true")
		}

		if title != "Sayings of the Century" {
			t.Fatalf("expected value to be 'Sayings of the Century', got %s", title)
		}
	}

	// invalid index access
	{
		if _, ok := document.String("$.store.book[10].title"); ok {
			t.Fatal("expected ok to be false")
		}
	}

	// reverse access
	{
		lastAuthor, ok := document.String("$.store.book[-2].author")
		if !ok {
			t.Fatal("expected ok to be true")
		}

		if lastAuthor != "Herman Melville" {
			t.Fatalf("expected lastAuthor to be 'Herman Melville', got %s", lastAuthor)
		}
	}

	// invalid reverse access
	{
		if _, ok := document.String("$.store.book[-10].author"); ok {
			t.Fatal("expected ok to be false")
		}
	}

	// wildcard all indexes
	{
		authors := document.Strings("$.store.book[*].author")

		if len(authors) != 4 {
			t.Fatalf("expected 4 authors, got %v", len(authors))
		}
	}
}

func TestSliceSelector(t *testing.T) {
	// slice access
	{
		middleTitles := document.Strings("$.store.book[1:3].title")

		if len(middleTitles) != 2 {
			t.Fatalf("expected 2 titles, got %v", len(middleTitles))
		}

		if middleTitles[0] != "Sword of Honour" {
			t.Fatalf("expected first title to be 'Sword of Honour', got %v", middleTitles[0])
		}

		if middleTitles[1] != "Moby Dick" {
			t.Fatalf("expected first title to be 'Moby Dick', got %v", middleTitles[1])
		}
	}

	// slice access with empty start
	{
		firstTwoTitles := document.Strings("$.store.book[:2].title")

		if len(firstTwoTitles) != 2 {
			t.Fatalf("expected 2 titles, got %v", len(firstTwoTitles))
		}

		if firstTwoTitles[0] != "Sayings of the Century" {
			t.Fatalf("expected second title to be 'Sayings of the Century', got %v",
				firstTwoTitles[0])
		}

		if firstTwoTitles[1] != "Sword of Honour" {
			t.Fatalf("expected second title to be 'Sword of Honour', got %v", firstTwoTitles[1])
		}
	}

	// slice access with empty end
	{
		lastTwoTitles := document.Strings("$.store.book[2:].title")

		if len(lastTwoTitles) != 2 {
			t.Fatalf("expected 2 titles, got %v", len(lastTwoTitles))
		}

		if lastTwoTitles[0] != "Moby Dick" {
			t.Fatalf("expected first title to be 'Moby Dick', got %v", lastTwoTitles[0])
		}

		if lastTwoTitles[1] != "The Lord of the Rings" {
			t.Fatalf("expected first title to be 'The Lord of the Rings', got %v", lastTwoTitles[1])
		}
	}

	// reverse slice access
	{
		middleTitles := document.Strings("$.store.book[-3:-1].title")

		if len(middleTitles) != 2 {
			t.Fatalf("expected 2 titles, got %v", len(middleTitles))
		}

		if middleTitles[0] != "Sword of Honour" {
			t.Fatalf("expected first title to be 'Sword of Honour', got %v", middleTitles[0])
		}

		if middleTitles[1] != "Moby Dick" {
			t.Fatalf("expected first title to be 'Moby Dick', got %v", middleTitles[1])
		}
	}

	// reverse slice access with empty start
	{
		firstTwoTitles := document.Strings("$.store.book[:-2].title")

		if len(firstTwoTitles) != 2 {
			t.Fatalf("expected 2 titles, got %v", len(firstTwoTitles))
		}

		if firstTwoTitles[0] != "Sayings of the Century" {
			t.Fatalf("expected second title to be 'Sayings of the Century', got %v",
				firstTwoTitles[0])
		}

		if firstTwoTitles[1] != "Sword of Honour" {
			t.Fatalf("expected second title to be 'Sword of Honour', got %v", firstTwoTitles[1])
		}
	}

	// reverse slice access with empty end
	{
		lastTwoTitles := document.Strings("$.store.book[-2:].title")

		if len(lastTwoTitles) != 2 {
			t.Fatalf("expected 2 titles, got %v", len(lastTwoTitles))
		}

		if lastTwoTitles[0] != "Moby Dick" {
			t.Fatalf("expected first title to be 'Moby Dick', got %v", lastTwoTitles[0])
		}

		if lastTwoTitles[1] != "The Lord of the Rings" {
			t.Fatalf("expected first title to be 'The Lord of the Rings', got %v", lastTwoTitles[1])
		}
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
