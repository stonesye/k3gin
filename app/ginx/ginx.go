package ginx

import "github.com/gin-gonic/gin"

func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		// TODO 这里也可以对Response的数据做一下封装
		return err
	}

	return nil
}
