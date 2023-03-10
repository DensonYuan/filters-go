package filters

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// New 创建 ModelFilter, 可通过 gin.Context 初始化
func New(model interface{}, c ...*gin.Context) *ModelFilter {
	mf := &ModelFilter{
		model:     model,
		canOrder:  make(map[string]bool),
		canMatch:  make(map[string]bool),
		canSearch: make(map[string]bool),
	}
	mf.initFunctionalFields()
	if len(c) > 0 {
		mf.initFromGinContext(c[0])
	}
	return mf
}

// SetGlobalConfig 设置全局配置
func SetGlobalConfig(config *Config) {
	globalConfig = config
	if globalConfig.LimitKey == "" {
		globalConfig.LimitKey = defaultLimitKey
	}
	if globalConfig.OffsetKey == "" {
		globalConfig.OffsetKey = defaultOffsetKey
	}
	if globalConfig.OrderKey == "" {
		globalConfig.OrderKey = defaultOrderKey
	}
	if globalConfig.SearchFieldsKey == "" {
		globalConfig.SearchFieldsKey = defaultSearchFieldsKey
	}
	if globalConfig.SearchValueKey == "" {
		globalConfig.SearchValueKey = defaultSearchValueKey
	}
	if globalConfig.FieldsKey == "" {
		globalConfig.FieldsKey = defaultFieldsKey
	}
}

// Query 获取结果集合
func (f *ModelFilter) Query(db *gorm.DB) *gorm.DB {
	db = db.Model(f.model)
	db = f.debugHandler(db)
	db = f.joinHandler(db)
	db = f.orderHandler(db)
	db = f.searchHandler(db)
	db = f.matchHandler(db)
	db = f.clauseHandler(db)
	db = f.paginationHandler(db)
	db = f.selectHandler(db)
	db = f.preloadHandler(db)
	db = f.handleHandler(db)
	return db
}

// Count 获取计数结果
func (f *ModelFilter) Count(db *gorm.DB) (cnt int64, err error) {
	err = f.Query(db).Limit(-1).Offset(-1).Count(&cnt).Error
	return
}

func (f *ModelFilter) Handle(handler Handler) *ModelFilter {
	f.handler = handler
	return f
}

func (f *ModelFilter) Debug() *ModelFilter {
	f.debug = true
	return f
}

// Select 设置查询字段
func (f *ModelFilter) Select(fields string) *ModelFilter {
	f.selectFields = fields
	return f
}

// Where 设置 Where 查询条件
func (f *ModelFilter) Where(query string, args ...interface{}) *ModelFilter {
	f.queries = append(f.queries, queryPair{Query: query, Args: args})
	return f
}

func (f *ModelFilter) Joins(query string, args ...interface{}) *ModelFilter {
	f.joins = append(f.joins, joinPair{Query: query, Args: args})
	return f
}

// Match 设置字段匹配条件
func (f *ModelFilter) Match(field string, value interface{}) *ModelFilter {
	if f.matches == nil {
		f.matches = make(map[string]interface{})
	}
	f.matches[field] = value
	return f
}

// OrderField 返回排序字段
func (f *ModelFilter) OrderField() string {
	return f.orderBy
}

// Order 设置排序字段
func (f *ModelFilter) Order(value string) *ModelFilter {
	f.orderBy = value
	return f
}

// LimitValue 返回分页大小
func (f *ModelFilter) LimitValue() int {
	return f.limit
}

// Limit 设置分页大小
func (f *ModelFilter) Limit(limit int) *ModelFilter {
	f.limit = limit
	return f
}

// OffsetValue 返回分页偏移量
func (f *ModelFilter) OffsetValue() int {
	return f.offset
}

// Offset 设置分页偏移
func (f *ModelFilter) Offset(offset int) *ModelFilter {
	f.offset = offset
	return f
}

// Search 设置搜索字段及值
func (f *ModelFilter) Search(fields string, value string) *ModelFilter {
	f.searchFields = fields
	f.searchValue = value
	return f
}

// Preload 设置预加载条件
func (f *ModelFilter) Preload(query string, args ...interface{}) *ModelFilter {
	if f.preloads == nil {
		f.preloads = make(map[string][]interface{})
	}
	f.preloads[query] = args
	return f
}

// ExtendSearchFields 手动扩充可搜索字段，用于联表等场景
func (f *ModelFilter) ExtendSearchFields(fields ...string) *ModelFilter {
	if f.canSearch == nil {
		f.canSearch = make(map[string]bool)
	}
	for _, field := range fields {
		f.canSearch[field] = true
	}
	return f
}

// ExtendMatchFields 手动扩充可匹配字段，用于联表等场景
func (f *ModelFilter) ExtendMatchFields(fields ...string) *ModelFilter {
	if f.canMatch == nil {
		f.canMatch = make(map[string]bool)
	}
	for _, field := range fields {
		f.canMatch[field] = true
	}
	return f
}

// ExtendOrderFields 手动扩充可排序字段，用于联表等场景
func (f *ModelFilter) ExtendOrderFields(fields ...string) *ModelFilter {
	if f.canOrder == nil {
		f.canOrder = make(map[string]bool)
	}
	for _, field := range fields {
		f.canOrder[field] = true
	}
	return f
}
