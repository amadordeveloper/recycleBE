package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func getAllResiduos(c *gin.Context) {
	result := findAllResiduos()
	c.IndentedJSON(200, result)
}

func findByClave(c *gin.Context) {
	if c.Param("clave") != "" {
		result := findResiduosByClave(c.Param("clave"))
		c.IndentedJSON(200, result)
	} else {
		result := findAllResiduos()
		c.IndentedJSON(200, result)
	}
}

func getAllTips(c *gin.Context) {
	result := findAllTips()
	c.IndentedJSON(200, result)
}

func getAllPuntosLimpios(c *gin.Context) {
	result := findAllPuntosLimpios()
	c.IndentedJSON(200, result)
}

func getRandomTip(c *gin.Context) {
	result := findRandomTip()
	c.IndentedJSON(200, result)
}

func sendRecoleccionDataToEmail(c *gin.Context) {
	var data recoleccionData
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
	}
	resp := sendRecoleccionData(data)
	var response gin.H
	if resp.Status == true {
		response = gin.H{
			"message": resp.Message,
			"status":  "success",
		}
	} else {
		response = gin.H{
			"message": resp.Message,
			"status":  "error",
		}
	}
	c.IndentedJSON(200, response)
}

/*
func addResiduo(c *gin.Context) {
	return
}
*/
