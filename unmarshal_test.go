package jpath

import (
	"testing"
)

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
