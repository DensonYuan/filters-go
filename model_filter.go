/*
@Time : 2020/4/8 16:59
@Author : yuandunsheng
@File : model_filter
*/

package filters

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
)

type ModelFilter struct {
	OrderFilter
	SearchFilter
	FieldFilter
	ClauseFilter
	Pagination
}

func (f *ModelFilter) Init(c *gin.Context, model interface{}) {
	f.InitPagination(c)
	f.InitOrderFilter(c, model)
	f.InitSearchFilter(c, model)
	f.InitFieldFilter(c, model)
}

func (f *ModelFilter) InitOrderFilter(c *gin.Context, model interface{}) {
	f.OrderFilter.Model = model
	f.OrderFilter.OrderBy = c.DefaultQuery("_order", "")
}

func (f *ModelFilter) InitSearchFilter(c *gin.Context, model interface{}) {
	f.SearchFilter.Model = model
	f.SearchFilter.SearchFields = c.DefaultQuery("_search_fields", "")
	f.SearchFilter.SearchValue = c.DefaultQuery("_search", "")
}

func (f *ModelFilter) InitFieldFilter(c *gin.Context, model interface{}) {
	f.FieldFilter.Model = model
	m := (map[string][]string)(c.Request.URL.Query())
	for k, v := range m {
		if !strings.HasPrefix(k, "_") && len(v) > 0 {
			f.SetFieldMatch(k, v[0])
		}
	}
}

func (f *ModelFilter) InitPagination(c *gin.Context) {
	f.Pagination.Limit = c.DefaultQuery("_limit", "-1")
	f.Pagination.Offset = c.DefaultQuery("_offset", "0")
}

func (f *ModelFilter) GetQuerySet(db *gorm.DB) *gorm.DB {
	db = f.OrderHandler(db)
	db = f.SearchHandler(db)
	db = f.FieldHandler(db)
	db = f.ClauseHandler(db)
	db = f.PaginationHandler(db)
	return db
}
