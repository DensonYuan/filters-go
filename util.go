package filters

import (
	"reflect"
	"strings"
)

type filterTag struct {
	Name   string
	Order  bool
	Search bool
	Match  bool
}

func getFilterTag(field reflect.StructField) *filterTag {
	var res *filterTag
	tag := field.Tag.Get("filter")
	if tag != "" {
		res = &filterTag{}
		tags := strings.Split(tag, ";")
		for _, t := range tags {
			if strings.HasPrefix(t, "name:") {
				res.Name = t[5:]
			} else if t == "order" {
				res.Order = true
			} else if t == "search" {
				res.Search = true
			} else if t == "match" {
				res.Match = true
			}
		}
	}
	return res
}
