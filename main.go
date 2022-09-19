package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type dataInterface interface {
	NombreResiduo() string
}

type clave struct {
	Clave string `json:"clave"`
}
type residuo struct {
	Id              int      `json:"id"`
	Nombre          string   `json:"nombre"`
	Claves          []string `json:"claves"`
	Destino         string   `json:"destino"`
	Impacto         string   `json:"impacto"`
	Aprovechamiento string   `json:"aprovechamiento"`
	Descripcion     string   `json:"descripcion"`
}

type tip struct {
	Id  int    `json:"id"`
	Tip string `json:"tip"`
}

type puntoLimpio struct {
	Id          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Latitud     float64 `json:"latitud"`
	Longitud    float64 `json:"longitud"`
	Descripcion string  `json:"descripcion"`
}

type recoleccionData struct {
	Tipo        string `json:"tipo"`
	Peso        string `json:"peso"`
	Dimensiones string `json:"dimensiones"`
	Direccion   string `json:"direccion"`
	Ciudad      string `json:"ciudad"`
	Nombre      string `json:"nombre"`
	Correo      string `json:"correo"`
	Telefono    string `json:"telefono"`
}

type responseStruct struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

/* Add methods */
func (s residuo) NombreResiduo() string {
	return s.Nombre
}

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Silence is golden")
	})

	// Residuos
	r.GET("/residuos/", getAllResiduos)
	r.GET("/residuos/:clave", findByClave)

	// Tips
	r.GET("/tips/", getAllTips)
	r.GET("/tips/random", getRandomTip)

	// PuntoLimpio
	r.GET("/puntosLimpios/", getAllPuntosLimpios)

	// Recoleccion
	r.POST("/recoleccion/", sendRecoleccionDataToEmail)

	// handler mostrar imagen desde el servidor
	r.GET("/archivos/:folder/:name", func(c *gin.Context) {
		c.File("public/archivos/" + c.Param("folder") + "/" + c.Param("name"))
	})

	// handle error 500 and 404
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route")
	})

	r.NoMethod(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route")
	})

	r.Run(":8080")
}
