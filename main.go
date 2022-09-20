package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
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

	r.Run(":80")

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	re, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(re.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range re.Labels {
		fmt.Printf("- %s\n", l.Name)
	}

}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
