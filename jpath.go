package jpath

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// .something
// ..something
// [sub]
// [union,union]
// [start:end:step]

var TagIdentifier = "jpath"

type filter func(f string, v interface{}) []interface{}

var filterMapping = map[string]filter{
	"(\\.{2}\\w+)": descendingAttributeFilter,
	"(\\.\\w+)":    attributeFilter,
	"(\\[.+\\])":   idxFilter,
}

var matcher *regexp.Regexp
var filters []filter

func init() {
	matchers := make([]string, 0)

	for regex, filter := range filterMapping {
		matchers = append(matchers, regex)
		filters = append(filters, filter)
	}
	matcher = regexp.MustCompile(strings.Join(matchers, "|"))
}

func dumpJson(l string, o interface{}) {
	data, _ := json.MarshalIndent(o, "", "	")
	fmt.Printf("%s: %s\n", l, data)
}

func filterForSegmentResults(segmentResults []string) (filter, string) {
	for i, segmentResult := range segmentResults {
		if segmentResult == "" {
			continue
		}
		return filters[i], segmentResult
	}
	return nil, ""
}

func New(m map[string]interface{}) *Jpath {
	return &Jpath{m}
}

func NewFromBytes(data []byte) (*Jpath, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return New(m), nil
}

type Jpath struct {
	m map[string]interface{}
}

func (j *Jpath) Query(sel string) []interface{} {
	segments := matcher.FindAllStringSubmatch(sel, -1)

	// make the initial nest for the objects to scan
	objs := []interface{}{j.m}

	// loop through all segments in the path
	for _, segmentResults := range segments {
		f, s := filterForSegmentResults(segmentResults[1:])

		// make a place to store the results of the filter
		tempObjs := make([]interface{}, 0)

		// loop through
		for _, o := range objs {
			tempObjs = append(tempObjs, f(s, o)...)
		}
		objs = tempObjs
	}

	return objs
}

func (j *Jpath) Strings(sel string) []string {
	ret := make([]string, 0)
	for _, result := range j.Query(sel) {
		if r, ok := result.(string); ok {
			ret = append(ret, r)
		}
	}
	return ret
}

func (j *Jpath) String(sel string) (string, bool) {
	for _, result := range j.Query(sel) {
		r, ok := result.(string)
		if !ok {
			return "", false
		}
		return r, true
	}
	return "", false
}

func (j *Jpath) Bools(sel string) []bool {
	ret := make([]bool, 0)
	for _, result := range j.Query(sel) {
		if r, ok := result.(bool); ok {
			ret = append(ret, r)
		}
	}
	return ret
}

func (j *Jpath) Bool(sel string) (bool, bool) {
	for _, result := range j.Query(sel) {
		r, ok := result.(bool)
		if !ok {
			return false, false
		}
		return r, true
	}
	return false, false
}

func (j *Jpath) Floats(sel string) []float64 {
	ret := make([]float64, 0)
	for _, result := range j.Query(sel) {
		if r, ok := result.(float64); ok {
			ret = append(ret, r)
		}
	}
	return ret
}

func (j *Jpath) Float(sel string) (float64, bool) {
	for _, result := range j.Query(sel) {
		r, ok := result.(float64)
		if !ok {
			return 0, false
		}
		return r, true
	}
	return 0, false
}

//
func attributeFilter(f string, v interface{}) []interface{} {
	ret := make([]interface{}, 0)

	// shave off the prepended period
	f = f[1:]

	// attribute filter only works on maps
	msi, ok := v.(map[string]interface{})
	if !ok {
		return ret
	}

	// grab the value at the other end of the attribute if it's available
	attr, ok := msi[f]
	if !ok {
		return ret
	}

	return append(ret, attr)
}

//
func idxFilter(f string, v interface{}) []interface{} {
	// shave off the brackets
	f = f[1 : len(f)-1]

	// make sure this is a slice
	slice, ok := v.([]interface{})
	if !ok {
		return make([]interface{}, 0)
	}

	// make sure the index is numeric
	i, err := strconv.Atoi(f)
	if err != nil {
		return make([]interface{}, 0)
	}

	// make sure it's not out of bounds
	if len(slice) <= i {
		return make([]interface{}, 0)
	}

	// grab the entry from the slice
	return slice[i : i+1]
}

//
func descendingAttributeFilter(f string, v interface{}) []interface{} {
	ret := make([]interface{}, 0)

	switch o := v.(type) {
	// if this is a map, then we need to see if one of it's
	// attributes matches our selector.  even if none do,
	// we need to proceed by descending deeper into the object
	// for a match.
	case map[string]interface{}:
		// check to see if a attribute matches the selector, if it
		// does then we just return right now.
		if c, ok := o[f[2:]]; ok {
			return append(ret, c)
		}

		// recursively keep checking for a match
		for _, val := range o {
			ret = append(ret, descendingAttributeFilter(f, val)...)
		}

	// if this a slice, then we need to check each member
	// to see if it matches our selector
	case []interface{}:
		// recursively keep checking for a match
		for _, val := range o {
			ret = append(ret, descendingAttributeFilter(f, val)...)
		}
	}

	return ret
}

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
