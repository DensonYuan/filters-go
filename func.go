/*
@Time : 2020/5/12 16:12
@Author : yuandunsheng
@File : func
*/

package filters

import "github.com/gin-gonic/gin"

// 通过 gin.Context 初始化 ModelFilter
func InitModelFilter(c *gin.Context, model interface{}) *ModelFilter {
	mf := &ModelFilter{}
	mf.Init(c, model)
	return mf
}

// 创建 ModelFilter，传入 model 对象
func NewModelFilter(model interface{}) *ModelFilter {
	mf := &ModelFilter{}
	mf.OrderFilter.Model = model
	mf.FieldFilter.Model = model
	mf.SearchFilter.Model = model
	return mf
}
