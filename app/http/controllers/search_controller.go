package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/support/facades"
)

type SearchController struct {
}

func (r SearchController) Search(ctx *gin.Context) {
	// facades.Response.Success(ctx, gin.H{
	// 	"Hello": ctx.DefaultQuery("key", "gjc"),
	// })
	// var w http.ResponseWriter
	// err := tplExample.ExecuteWriter(pongo2.Context{"query": r.FormValue("query")}, w)
	// err := tplExample.ExecuteWriter(pongo2.Context{"query": "query"}, w)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
	// gd := gin.Default()
	// r.LoadHTMLFiles("resources/views/index.html")
	// ctx.HTML(200, "index.html", "flysnow_org")
	facades.Response.Success(ctx, gin.H{
		"Hello": "Goravel",
	})
}
