package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"log"
	"math/rand"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
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

/*
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
			Descripcion:     record[5],
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


	c.String(http.StatusOK, "finalizado")
}*/

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ConnectMySQLDB() *sql.DB {
	// connect to mysql

	Host := getEnv("DB_HOST", "")
	Port := getEnv("DB_PORT", "")
	User := getEnv("DB_USER", "")
	Pass := getEnv("DB_PASS", "")
	DBName := getEnv("DB_NAME", "")

	if Host == "" || Port == "" || User == "" || Pass == "" || DBName == "" {
		panic("No se encontro la variable de entorno para la conexion a la base de datos")
	}

	db, err := sql.Open("mysql", User+":"+Pass+"@tcp("+Host+":"+Port+")/"+DBName+"?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return db
}

func sendEmail(data recoleccionData) bool {

	// select from db configuraciones where key = "email"
	db := ConnectMySQLDB()
	defer db.Close()
	query := "SELECT valor FROM configuraciones WHERE clave = 'correo'"
	row := db.QueryRow(query)
	var email string
	err := row.Scan(&email)
	if err != nil {
		log.Println(err)
		return false
	}

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	message := "From: recycleappbucaramanga@gmail.com\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + "Nueva recoleccion" + "\r\n"
	message += "Content-Type: text/html; charset=\"UTF-8\"\r\n"

	message += "\r\n" + "<b>Se ha registrado una nueva solicitud de recoleccion</b><br/>" + "\r\n" +
		"<b>Nombre: </b>" + data.Nombre + "<br/>" + "\r\n" +
		"<b>Correo: </b>" + data.Correo + "<br/>" + "\r\n" +
		"<b>Telefono: </b>" + data.Telefono + "<br/>" + "\r\n" +
		"<b>Direccion: </b>" + data.Direccion + "<br/>" + "\r\n" +
		"<b>Ciudad: </b>" + data.Ciudad + "<br/>" + "\r\n" +
		"<b>Tipo: </b>" + data.Tipo + "<br/>" + "\r\n" +
		"<b>Dimensiones: </b>" + data.Dimensiones + "<br/>" + "\r\n" +
		"<b>Peso: </b>" + data.Peso + "<br/>" + "\r\n"

	// Send the message
	_, err = srv.Users.Messages.Send("me", &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(message)),
	}).Do()
	if err != nil {
		log.Fatalf("Unable to send message. %v", err)
	}

	return true
}
