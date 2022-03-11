package main

import (
	"github.com/gin-gonic/gin"
)

func getAllResiduos(c *gin.Context) {

	result := findAll("residuos")
	c.IndentedJSON(200, result)
}

func findByClave(c *gin.Context) {
	result := findBy("residuos", "Claves", c.Param("clave"), true)
	c.IndentedJSON(200, result)
}

func addResiduo(c *gin.Context) {
	return
}
