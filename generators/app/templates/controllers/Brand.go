package controllers

import (
	"net/http"

	"<%= myAppPath %>/httputil"
	"<%= myAppPath %>/models"

	"github.com/gin-gonic/gin"
)

// 获取品牌 godoc
// @Summary 获取品牌
// @Description 获取品牌
// @Tags Brand 品牌
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Security OAuth2Application[admin]
// @router /Brand/GetList [get]
func (c *Controller) BrandGetList(ctx *gin.Context) {
	data, err := models.GetBrandList()
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}
