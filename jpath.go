package jpath

import (
	"encoding/json"
)

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

// Query takes a JSON path and returns a list of objects that matched.
func (j *Jpath) Query(sel string) []interface{} {
	// make the initial nest for the objects to scan
	objs := []interface{}{j.m}

	// loop through all segments in the path
	for _, segmentResults := range segmentsForSelector(sel) {
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

// Strings is a convenience method for returning a JSON path query as a slice of strings.
// Skips results that matched the query, but are not strings.
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

// Bools is a convenience method for returning a JSON path query as a slice of bools.
// Skips results that matched the query, but are not bools.
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

// Floats is a convenience method for returning a JSON path query as a slice of floats.
// Skips results that matched the query, but are not floats.
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
