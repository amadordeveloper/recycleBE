package main

import (
	"github.com/gin-gonic/gin"
)

func getAllResiduos(c *gin.Context) {

	result := findAllResiduos()
	c.IndentedJSON(200, result)
}

func findByClave(c *gin.Context) {
	result := findResiduosByClave(c.Param("clave"))
	c.IndentedJSON(200, result)
}

/*
func addResiduo(c *gin.Context) {
	return
}
*/
