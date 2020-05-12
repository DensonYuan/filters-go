/*
@Time : 2020/4/27 11:55
@Author : yuandunsheng
@File : order_filter
*/

package filters

import (
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
)

type OrderFilter struct {
	Model   interface{}
	OrderBy string
}

func (f *OrderFilter) OrderHandler(db *gorm.DB) *gorm.DB {
	if f.OrderBy != "" {
		var order, field string
		if strings.HasPrefix(f.OrderBy, "-") {
			field = f.OrderBy[1:]
			order = f.OrderBy[1:] + " DESC"
		} else {
			field = f.OrderBy
			order = f.OrderBy
		}
		if f.isOrderFieldValid(field) {
			db = db.Order(order)
		}
	}
	return db
}

func (f *OrderFilter) isOrderFieldValid(field string) bool {
	typeOfModel := reflect.TypeOf(f.Model)
	for i := 0; i < typeOfModel.NumField(); i++ {
		if ft := getFilterTag(typeOfModel.Field(i)); ft != nil && ft.Name == field && ft.Order {
			return true
		}
	}
	return false
}
