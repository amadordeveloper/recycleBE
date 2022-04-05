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
	result := findResiduosByClave(c.Param("clave"))
	c.IndentedJSON(200, result)
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
	// print to the console json data from POST request
	/*jsonData, err := c.GetRawData()
	if err != nil {
		fmt.Println(err)
	}*/
	var data recoleccionData
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
	}
	resp := sendRecoleccionData(data)
	var response gin.H
	if resp == true {
		response = gin.H{
			"message": "Email sent successfully",
			"status":  "success",
		}
	} else {
		response = gin.H{
			"message": "Email not sent",
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
