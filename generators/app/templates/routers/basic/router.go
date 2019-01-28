package basic

import (
	"<%= myAppPath %>/ginserver"

	"<%= myAppPath %>/controllers"

	"github.com/gin-gonic/gin"
)

func CreateBasicRouter(eng *gin.Engine) (err error) {
	c := controllers.NewController()
	v1 := eng.Group("/v1")
	{
		Brand := v1.Group("/Brand")
		{
			Brand.Use(ginserver.HandleTokenVerify())
			Brand.GET("GetList", c.BrandGetList)
		}
	}
	return nil
}
