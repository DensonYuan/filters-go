/*
@Time : 2020/4/27 11:57
@Author : yuandunsheng
@File : search_filter
*/

package filters

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
)

type SearchFilter struct {
	Model        interface{}
	SearchFields string
	SearchValue  string
}

func (b *SearchFilter) SearchHandler(db *gorm.DB) *gorm.DB {
	clause := ""
	for _, field := range strings.Split(b.SearchFields, ",") {
		if b.isSearchFieldValid(field) {
			format := "`%s` LIKE '%%%s%%'"
			if clause != "" {
				format = " OR " + format
			}
			clause += fmt.Sprintf(format, field, b.SearchValue)
		}
	}
	db = db.Where(clause)
	return db
}

func (b *SearchFilter) isSearchFieldValid(field string) bool {
	if strings.TrimSpace(field) == "" {
		return false
	}
	typeOfModel := reflect.TypeOf(b.Model)
	for i := 0; i < typeOfModel.NumField(); i++ {
		if ft := getFilterTag(typeOfModel.Field(i)); ft != nil && ft.Name == field && ft.Search {
			return true
		}
	}
	return false
}
