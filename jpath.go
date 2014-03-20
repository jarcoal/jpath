package jpath

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// .something
// ..something
// [sub]
// [union,union]
// [start:end:step]

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

func Query(sel string, obj map[string]interface{}) []interface{} {
	segments := matcher.FindAllStringSubmatch(sel, -1)

	// make the initial nest for the objects to scan
	objs := []interface{}{obj}

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
