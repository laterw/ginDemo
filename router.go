package main

import (
	"flow/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	r.GET("/index", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.GetflowTatal())
	})
	r.GET("/old/:dateTime", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.GetOldTotal(c))
	})
	r.GET("/lfet1", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.GetLeft1(c))
	})
	r.GET("/lfet2", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.GetLeft2(c))
	})
	r.GET("/lfet3/:busssinessKey", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.GetLeft3(c))
	})
	r.GET("/lfet3Init", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.GetLeft3Init(c))
	})
	r.GET("/testWebSocket", func(c *gin.Context) {
		model.TestWebs(c)
	})
	return r
}
