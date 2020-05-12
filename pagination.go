/*
@Time : 2020/4/7 15:03
@Author : yuandunsheng
@File : filters
*/

package filters

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type Pagination struct {
	Limit  string
	Offset string
}

func (p *Pagination) PaginationHandler(db *gorm.DB) *gorm.DB {
	limit, err := strconv.Atoi(p.Limit)
	if err != nil {
		limit = -1
	}
	offset, _ := strconv.Atoi(p.Offset)
	db = db.Limit(limit)
	db = db.Offset(offset)
	return db
}
