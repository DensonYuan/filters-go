/*
@Time : 2020/4/20 11:45
@Author : yuandunsheng
@File : field_filter
*/

package filters

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
)

type FieldFilter struct {
	Model         interface{}
	mapFieldMatch map[string]interface{}
}

func (f *FieldFilter) SetFieldMatch(field string, value interface{}) {
	if f.mapFieldMatch == nil {
		f.mapFieldMatch = make(map[string]interface{})
	}
	f.mapFieldMatch[field] = value
}

func (f *FieldFilter) FieldHandler(db *gorm.DB) *gorm.DB {
	for k, v := range f.mapFieldMatch {
		if f.isFilterFieldValid(k) {
			db = db.Where(fmt.Sprintf("`%s` = ?", k), v)
		}
	}
	return db
}

func (f *FieldFilter) isFilterFieldValid(field string) bool {
	typeOfModel := reflect.TypeOf(f.Model)
	for i := 0; i < typeOfModel.NumField(); i++ {
		if ft := getFilterTag(typeOfModel.Field(i)); ft != nil && ft.Name == field && ft.Match {
			return true
		}
	}
	return false
}
