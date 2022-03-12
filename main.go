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

/* Add methods */
func (s residuo) NombreResiduo() string {
	return s.Nombre
}

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Silence is golden")
	})

	r.GET("/residuos/", getAllResiduos)
	r.GET("/residuos/:clave", findByClave)
	/*	r.POST("/residuos/add", addResiduo)*/

	// Upload route
	r.LoadHTMLFiles("public/upload.html")

	r.GET("/importer", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	//r.POST("/upload", upload)
	r.StaticFS("/files", http.Dir("files"))

	r.Run(":8080")
}
