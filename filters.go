package jpath

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type filter func(f string, v interface{}) []interface{}

var filterMapping = map[string]filter{
	"\\.{2}\\w+": descendingAttributeFilter,
	"\\.\\w+":    attributeFilter,
	"\\[.+\\]":   idxFilter,
}

var matcher *regexp.Regexp
var filters []filter

// build the complete regex from all of the registered filters
func init() {
	matchers := make([]string, 0)

	for regex, filter := range filterMapping {
		matchers = append(matchers, fmt.Sprintf("(%s)", regex))
		filters = append(filters, filter)
	}
	matcher = regexp.MustCompile(strings.Join(matchers, "|"))
}

// filterForSegmentResults takes segment result pieces and finds the filter that handles
// that segment type
func filterForSegmentResults(segmentResults []string) (filter, string) {
	for i, segmentResult := range segmentResults {
		if segmentResult == "" {
			continue
		}
		return filters[i], segmentResult
	}
	return nil, ""
}

// segmentsForSelector parses a selector into filterable chunks
func segmentsForSelector(sel string) [][]string {
	return matcher.FindAllStringSubmatch(sel, -1)
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

	// wildcard, so we return everything inside of the slice
	if f == "*" {
		return slice
	}

	// check to see if this is a simple index access
	if i, err := strconv.Atoi(f); err == nil {
		// make sure it's not out of range
		if outOfRange(i, slice) {
			return make([]interface{}, 0)
		}

		start := i

		// for reverse access, add the length of the slice
		if i < 0 {
			start += len(slice)
		}

		return slice[start : start+1]
	}

	// check if this a slice access
	if sl := strings.Split(f, ":"); len(sl) > 0 {
		var start, end int
		var err error

		// get the initial value for the start of the slice
		if sl[0] == "" {
			start = 0
		} else if start, err = strconv.Atoi(sl[0]); err != nil || outOfRange(start, slice) {
			return make([]interface{}, 0)
		}

		// if the start is a reverse access, put it into range
		if start < 0 {
			start += len(slice)
		}

		// get the initial value for the end of the slice
		if sl[1] == "" {
			end = len(slice)
		} else if end, err = strconv.Atoi(sl[1]); err != nil || outOfRange(end, slice) {
			return make([]interface{}, 0)
		}

		// if the end is a reverse access, put it into range
		if end < 0 {
			end += len(slice)
		}

		return slice[start:end]
	}

	// some invalid filter
	return make([]interface{}, 0)
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

func outOfRange(idx int, sl []interface{}) bool {
	return math.Abs(float64(idx)) >= float64(len(sl))
}
