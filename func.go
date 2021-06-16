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
	db = f.preloadHandler(db)
	return db
}

func (f *ModelFilter) GetCount(db *gorm.DB) (cnt int, err error) {
	err = f.GetQuerySet(db).Limit(-1).Offset(0).Count(&cnt).Error
	return
}

func (f *ModelFilter) Delete(db *gorm.DB) (err error) {
	err = f.GetQuerySet(db).Delete(f.model).Error
	return
}

func (f *ModelFilter) Count(db *gorm.DB, value interface{}) *gorm.DB {
	return f.GetQuerySet(db).Limit(-1).Offset(0).Count(value)
}

func (f *ModelFilter) Select(fields string) *ModelFilter {
	f.fields = fields
	return f
}

func (f *ModelFilter) Where(query string, args ...interface{}) *ModelFilter {
	f.queryList = append(f.queryList, query)
	f.argsList = append(f.argsList, args)
	return f
}

func (f *ModelFilter) Match(field string, value interface{}) *ModelFilter {
	if f.mapFieldMatch == nil {
		f.mapFieldMatch = make(map[string]interface{})
	}
	f.mapFieldMatch[field] = value
	return f
}

func (f *ModelFilter) Order(value string) *ModelFilter {
	f.orderBy = value
	return f
}

func (f *ModelFilter) Limit(limit interface{}) *ModelFilter {
	f.limit = limit
	return f
}

func (f *ModelFilter) Offset(offset interface{}) *ModelFilter {
	f.offset = offset
	return f
}

func (f *ModelFilter) Search(fields string, value string) *ModelFilter {
	f.searchFields = fields
	f.searchValue = value
	return f
}

func (f *ModelFilter) Preload(column string, conditions ...interface{}) *ModelFilter {
	f.preloadColumn = column
	f.preloadConditions = conditions
	return f
}
