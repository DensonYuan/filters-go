package filters

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
)

type ModelFilter struct {
	model             interface{}
	orderBy           string
	searchFields      string
	searchValue       string
	mapFieldMatch     map[string]interface{}
	query             string
	args              []interface{}
	limit             interface{}
	offset            interface{}
	fields            string
	preloadColumn     string
	preloadConditions []interface{}
}

//////////////////////////////////////////////////////////////////////////////////////

func (f *ModelFilter) init(c *gin.Context, model interface{}) {
	f.model = model
	f.limit = c.DefaultQuery("_limit", "-1")
	f.offset = c.DefaultQuery("_offset", "0")
	f.orderBy = c.DefaultQuery("_order", "")
	f.searchFields = c.DefaultQuery("_search_fields", "")
	f.searchValue = c.DefaultQuery("_search", "")
	f.fields = c.DefaultQuery("_fields", "")

	m := (map[string][]string)(c.Request.URL.Query())
	for k, v := range m {
		if !strings.HasPrefix(k, "_") && len(v) > 0 {
			f.Match(k, v[0])
		}
	}
}

//////////////////////////////////////////////////////////////////////////////////////

func (f *ModelFilter) orderHandler(db *gorm.DB) *gorm.DB {
	if f.orderBy != "" {
		var order, field string
		if strings.HasPrefix(f.orderBy, "-") {
			field = f.orderBy[1:]
			order = "`" + f.orderBy[1:] + "` DESC"
		} else {
			field = f.orderBy
			order = "`" + f.orderBy + "`"
		}
		if f.isOrderFieldValid(field) {
			db = db.Order(order)
		}
	}
	return db
}

func (f *ModelFilter) isOrderFieldValid(field string) bool {
	typeOfModel := reflect.TypeOf(f.model)
	for i := 0; i < typeOfModel.NumField(); i++ {
		if ft := getFilterTag(typeOfModel.Field(i)); ft != nil && ft.Name == field && ft.Order {
			return true
		}
	}
	return false
}

//////////////////////////////////////////////////////////////////////////////////////

func (f *ModelFilter) paginationHandler(db *gorm.DB) *gorm.DB {
	db = db.Limit(f.limit)
	db = db.Offset(f.offset)
	return db
}

//////////////////////////////////////////////////////////////////////////////////////

func (f *ModelFilter) searchHandler(db *gorm.DB) *gorm.DB {
	if f.searchValue == "" {
		return db
	}
	var fields []string
	if f.searchFields != "" {
		fields = strings.Split(f.searchFields, ",")
	} else {
		typeOfModel := reflect.TypeOf(f.model)
		for i := 0; i < typeOfModel.NumField(); i++ {
			if ft := getFilterTag(typeOfModel.Field(i)); ft != nil && ft.Search {
				fields = append(fields, ft.Name)
			}
		}

	}
	clause := ""
	for _, field := range fields {
		if f.searchFields == "" || f.isSearchFieldValid(field) {
			format := "`%s` LIKE '%%%s%%'"
			if clause != "" {
				format = " OR " + format
			}
			clause += fmt.Sprintf(format, field, f.searchValue)
		}
	}
	db = db.Where(clause)
	return db
}

func (f *ModelFilter) isSearchFieldValid(field string) bool {
	if strings.TrimSpace(field) == "" {
		return false
	}
	typeOfModel := reflect.TypeOf(f.model)
	for i := 0; i < typeOfModel.NumField(); i++ {
		if ft := getFilterTag(typeOfModel.Field(i)); ft != nil && ft.Name == field && ft.Search {
			return true
		}
	}
	return false
}

//////////////////////////////////////////////////////////////////////////////////////

func (f *ModelFilter) matchHandler(db *gorm.DB) *gorm.DB {
	for k, v := range f.mapFieldMatch {
		if f.isMatchFieldValid(k) {
			db = db.Where(fmt.Sprintf("`%s` = ?", k), v)
		}
	}
	return db
}

func (f *ModelFilter) isMatchFieldValid(field string) bool {
	typeOfModel := reflect.TypeOf(f.model)
	for i := 0; i < typeOfModel.NumField(); i++ {
		if ft := getFilterTag(typeOfModel.Field(i)); ft != nil && ft.Name == field && ft.Match {
			return true
		}
	}
	return false
}

//////////////////////////////////////////////////////////////////////////////////////

func (f *ModelFilter) clauseHandler(db *gorm.DB) *gorm.DB {
	if f.query != "" {
		return db.Where(f.query, f.args...)
	}
	return db
}

//////////////////////////////////////////////////////////////////////////////////////

func (f *ModelFilter) selectHandler(db *gorm.DB) *gorm.DB {
	if f.fields != "" {
		return db.Select(f.fields)
	}
	return db
}

//////////////////////////////////////////////////////////////////////////////////////

func (f *ModelFilter) preloadHandler(db *gorm.DB) *gorm.DB {
	return db.Preload(f.preloadColumn, f.preloadConditions...)
}
