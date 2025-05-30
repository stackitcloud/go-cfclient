package client

import (
	"net/url"
	"strings"
	"time"
)

type FilterModifier int

const (
	FilterModifierNone FilterModifier = iota
	FilterModifierGreaterThan
	FilterModifierLessThan
	FilterModifierGreaterThanOrEqual
	FilterModifierLessThanOrEqual
)

func (r FilterModifier) String() string {
	switch r {
	case FilterModifierGreaterThan:
		return "gt"
	case FilterModifierGreaterThanOrEqual:
		return "gte"
	case FilterModifierLessThan:
		return "lt"
	case FilterModifierLessThanOrEqual:
		return "lte"
	default:
		return ""
	}
}

type TimestampFilterList []TimestampFilter

type TimestampFilter struct {
	Timestamp []time.Time
	Operator  FilterModifier
}

func (tl *TimestampFilterList) EqualTo(ts ...time.Time) {
	*tl = append(*tl, TimestampFilter{
		Timestamp: ts,
	})
}

func (tl *TimestampFilterList) Before(ts time.Time) {
	*tl = append(*tl, TimestampFilter{
		Timestamp: []time.Time{
			ts,
		},
		Operator: FilterModifierLessThan,
	})
}

func (tl *TimestampFilterList) BeforeOrEqualTo(ts time.Time) {
	*tl = append(*tl, TimestampFilter{
		Timestamp: []time.Time{
			ts,
		},
		Operator: FilterModifierLessThanOrEqual,
	})
}

func (tl *TimestampFilterList) After(ts time.Time) {
	*tl = append(*tl, TimestampFilter{
		Timestamp: []time.Time{
			ts,
		},
		Operator: FilterModifierGreaterThan,
	})
}

func (tl *TimestampFilterList) AfterOrEqualTo(ts time.Time) {
	*tl = append(*tl, TimestampFilter{
		Timestamp: []time.Time{
			ts,
		},
		Operator: FilterModifierGreaterThanOrEqual,
	})
}

func (tl TimestampFilterList) Serialize(values url.Values, tag string) error {
	for _, t := range tl {
		err := t.Serialize(values, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t TimestampFilter) Serialize(values url.Values, tag string) error {
	if len(t.Timestamp) > 0 {
		if t.Operator != FilterModifierNone {
			tag = tag + "[" + t.Operator.String() + "]"
		}
		timestamps := make([]string, len(t.Timestamp))
		for i, timestamp := range t.Timestamp {
			timestamps[i] = timestamp.Format(time.RFC3339)
		}
		values.Add(tag, strings.Join(timestamps, ","))
	}
	return nil
}

type Filter struct {
	Values []string
}

func (f *Filter) EqualTo(v ...string) {
	f.Values = v
}

func (f Filter) Serialize(values url.Values, tag string) error {
	if len(f.Values) > 0 {
		values.Add(tag, strings.Join(f.Values, ","))
	}
	return nil
}

type ExclusionFilter struct {
	Filter
	Not bool
}

func (e *ExclusionFilter) NotEqualTo(v ...string) {
	e.Values = v
	e.Not = true
}

func (e ExclusionFilter) Serialize(values url.Values, tag string) error {
	if len(e.Values) > 0 {
		if e.Not {
			tag = tag + "[not]"
		}
		values.Add(tag, strings.Join(e.Values, ","))
	}
	return nil
}

type LabelSelector map[string]ExclusionFilter

func (l LabelSelector) Existence(key string) {
	l[key] = ExclusionFilter{}
}

func (l LabelSelector) NotExistence(key string) {
	l[key] = ExclusionFilter{Not: true}
}

func (l LabelSelector) EqualTo(key string, values ...string) {
	l[key] = ExclusionFilter{Filter: Filter{Values: values}}
}

func (l LabelSelector) NotEqualTo(key string, values ...string) {
	l[key] = ExclusionFilter{Filter: Filter{Values: values}, Not: true}
}

func (l LabelSelector) Serialize(values url.Values, tag string) error {
	if len(l) == 0 {
		return nil
	}

	selectors := make([]string, 0, len(l))
	for k, v := range l {
		var ops string
		switch len(v.Values) {
		case 0:
			if v.Not {
				ops = "!"
			}
			selectors = append(selectors, ops+k)
		case 1:
			if v.Not {
				ops = "!="
			} else {
				ops = "="
			}
			selectors = append(selectors, k+ops+v.Values[0])
		default:
			if v.Not {
				ops = " notin ("
			} else {
				ops = " in ("
			}
			selectors = append(selectors, k+ops+strings.Join(v.Values, ",")+")")
		}
	}
	values.Add(tag, strings.Join(selectors, ","))
	return nil
}
