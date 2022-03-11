package main

import (
	"io"
	"math/rand"
	"os"
	"strings"

	"fmt"
	"net/http"

	"encoding/csv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getEnv(key, fallback string) string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func upload(c *gin.Context) {
	//file, header, err := c.Request.FormFile("file")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file load err : %s", err.Error()))
		return
	}

	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ';'

	var LoteResiduos []dataInterface

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("file read err : %s", err.Error()))
			return
		}

		if record[0] == "nombre" {
			continue
		}

		trimClaves := strings.ReplaceAll(record[1], " ", "")
		arrClaves := strings.Split(trimClaves, ",")

		res := residuo{
			Nombre:          record[0],
			Claves:          arrClaves,
			Destino:         record[3],
			Impacto:         record[2],
			Aprovechamiento: record[4],
		}

		LoteResiduos = append(LoteResiduos, res)
	}

	if len(LoteResiduos) == 0 {
		c.String(http.StatusBadRequest, "No hay datos")
		return
	} else {
		sth := saveAll(LoteResiduos, "residuos")
		if sth {
			c.String(http.StatusOK, "Datos cargados")
		} else {
			c.String(http.StatusBadRequest, "Error al cargar datos")
		}
	}

	/* // save file to server
	filename := header.Filename
	out, err := os.Create("files/" + filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		panic(err)
	}
	filepath := "http://localhost:8080/files/" + filename
	c.JSON(http.StatusOK, gin.H{"filepath": filepath})*/

	c.String(http.StatusOK, "finalizado")
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
