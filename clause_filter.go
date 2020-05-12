/*
@Time : 2020/4/28 14:32
@Author : yuandunsheng
@File : clause_filter
*/

package filters

import "github.com/jinzhu/gorm"

type ClauseFilter struct {
	query string
	args  []interface{}
}

func (f *ClauseFilter) SetClause(query string, args ...interface{}) {
	if query != "" {
		f.query = query
		f.args = args
	}
}

func (f *ClauseFilter) ClauseHandler(db *gorm.DB) *gorm.DB {
	if f.query != "" {
		return db.Where(f.query, f.args...)
	}
	return db
}
