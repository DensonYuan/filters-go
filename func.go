package filters

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// 通过 gin.Context 初始化 ModelFilter
func InitModelFilter(c *gin.Context, model interface{}) *ModelFilter {
	mf := &ModelFilter{}
	mf.init(c, model)
	return mf
}

// 创建 ModelFilter，传入 model 对象
func NewModelFilter(model interface{}) *ModelFilter {
	mf := &ModelFilter{}
	mf.model = model
	return mf
}

func (f *ModelFilter) GetQuerySet(db *gorm.DB) *gorm.DB {
	db = db.Model(f.model)
	db = f.orderHandler(db)
	db = f.searchHandler(db)
	db = f.matchHandler(db)
	db = f.clauseHandler(db)
	db = f.paginationHandler(db)
	db = f.selectHandler(db)
	return db
}

func (f *ModelFilter) GetCount(db *gorm.DB) (cnt int, err error) {
	err = f.GetQuerySet(db).Limit(-1).Offset(0).Count(&cnt).Error
	return
}

func (f *ModelFilter) Count(db *gorm.DB, value interface{}) *gorm.DB {
	return f.GetQuerySet(db).Limit(-1).Offset(0).Count(value)
}

func (f *ModelFilter) Select(fields string) {
	f.fields = fields
}

func (f *ModelFilter) Where(query string, args ...interface{}) {
	if query != "" {
		f.query = query
		f.args = args
	}
}

func (f *ModelFilter) Match(field string, value interface{}) {
	if f.mapFieldMatch == nil {
		f.mapFieldMatch = make(map[string]interface{})
	}
	f.mapFieldMatch[field] = value
}

func (f *ModelFilter) Order(value string) {
	f.orderBy = value
}

func (f *ModelFilter) Limit(limit interface{}) {
	f.limit = limit
}

func (f *ModelFilter) Offset(offset interface{}) {
	f.offset = offset
}

func (f *ModelFilter) Search(fields string, value string) {
	f.searchFields = fields
	f.searchValue = value
}
