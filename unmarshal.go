package jpath

import (
	"fmt"
	"reflect"
	"strings"
)

// TagIdentifier determines the namespace of the tags accepted for the Unmarshaler.
var TagIdentifier = "jpath"

// Unmarshal functions similarly to json.Unmarshal, except reads the jpath tag
// and marshals results into struct fields based on the results of the queries.
func Unmarshal(data []byte, v interface{}) error {
	d, err := NewFromBytes(data)
	if err != nil {
		return err
	}
	return unmarshal(d, v)
}

func unmarshal(d *Jpath, v interface{}) error {
	vt := reflect.TypeOf(v).Elem()
	vv := reflect.ValueOf(v).Elem()

	if vt.Kind() != reflect.Struct {
		return fmt.Errorf("v must be a struct, got %T", v)
	}

	// we're going to loop through each field in the struct and
	// extract the json path from the tags, query the document with
	// those tags, then set the results on the fields.
	for i := 0; i < vt.NumField(); i++ {
		fieldType := vt.Field(i)
		fieldValue := vv.Field(i)
		fieldKind := fieldValue.Kind()

		// if we can't update this field, we're done with it
		if !fieldValue.CanSet() {
			continue
		}

		// extract the json path from the tag
		tag := fieldType.Tag.Get(TagIdentifier)
		if tag == "" {
			continue
		}

		// eventually we will support marshaling values out, which will be
		// delimited by a space, so this is in preparation for that.
		tagPieces := strings.Split(tag, " ")
		query := tagPieces[0]

		// query the document with the json path
		results := d.Query(query)

		// if this isn't a slice we're unmarshaling into, then just take the
		// first value and set it and we're done.
		if fieldKind != reflect.Slice {
			resultValue := reflect.ValueOf(results[0])

			// ensure that the kinds line up before we set
			if resultValue.Kind() != fieldKind {
				return fmt.Errorf("%s - value of type %s is not assignable to type %s",
					fieldType.Name, resultValue.Type(), fieldType.Type)
			}

			fieldValue.Set(resultValue)
			continue
		}

		// make a slice with the same type as the field
		sl := reflect.MakeSlice(fieldType.Type, 0, len(results))

		// loop through results and append them to the slice we created
		for _, result := range results {
			resultValue := reflect.ValueOf(result)

			// TODO: Find a way to sniff the type of items in the slice and just do a check to see
			// if the kind will fit inside. in the meantime, we're just going to trap panics.
			var err error
			func() {
				defer func() {
					if recover() != nil {
						err = fmt.Errorf("%s - value of type %s will not fit into slice of %s",
							fieldType.Name, resultValue.Type(), fieldType.Type)
					}
				}()

				sl = reflect.Append(sl, resultValue)
			}()

			if err != nil {
				return err
			}
		}

		// finally set the slice
		fieldValue.Set(sl)
	}

	return nil
}
