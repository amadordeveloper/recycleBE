package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/smtp"
	"os"

	_ "github.com/go-sql-driver/mysql"
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

	from := getEnv("EMAIL_USER", "")
	password := getEnv("EMAIL_PASS", "")
	to := []string{email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Subject: Recoleccion de residuos\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// template html email
	html := `<h2>Solicitud de Recolección</h2><br>
				<b>Nombre:</b> ` + data.Nombre + `<br>
				<b>Correo:</b> ` + data.Correo + `<br>
				<b>Telefono:</b> ` + data.Telefono + `<br>
				<b>Direccion:</b> ` + data.Direccion + `<br>
				<b>Ciudad:</b> ` + data.Ciudad + `<br>
				<b>Tipo:</b> ` + data.Tipo + `<br>
				<b>Dimensiones:</b> ` + data.Dimensiones + `<br>
				<b>Peso:</b> ` + data.Peso + `<br>`

	msg := []byte(subject + mime + html)

	// send html email
	err = smtp.SendMail(smtpHost+":"+smtpPort,
		smtp.PlainAuth("", from, password, smtpHost),
		from,
		to,
		msg)

	if err != nil {
		log.Fatal(err)
	}
	return true
}
