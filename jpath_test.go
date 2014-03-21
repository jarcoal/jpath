package jpath

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
