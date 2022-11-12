package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func findAllResiduos() []interface{} {
	conn := ConnectMySQLDB()
	sqlQuery := `SELECT
								R.id,
								R.nombre,
								D.nombre,
								R.impacto,
								R.aprovechamiento,
								R.descripcion
							FROM
								residuos R,
								destinos D
							WHERE
								R.id_destino = D.id;`

	rows, err := conn.Query(sqlQuery)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	var result []interface{}
	for rows.Next() {
		var r residuo
		err = rows.Scan(&r.Id, &r.Nombre, &r.Destino, &r.Impacto, &r.Aprovechamiento, &r.Descripcion)
		if err != nil {
			panic(err.Error())
		}

		sqlQuery := `SELECT
										C.clave as claves
									FROM
										residuos R,
										claves C,
										clave_residuo CR
									WHERE
										CR.id_residuo = R.id
									AND
										CR.id_clave = C.id
									AND R.id = ?;`
		rowsClaves, err := conn.Query(sqlQuery, r.Id)
		if err != nil {
			panic(err.Error())
		}
		var claves []string
		for rowsClaves.Next() {
			var c clave
			err = rowsClaves.Scan(&c.Clave)
			if err != nil {
				panic(err.Error())
			}
			claves = append(claves, c.Clave)
		}

		r.Claves = claves

		result = append(result, r)
	}
	return result
}

func findResiduosByClave(keyword string) interface{} {
	conn := ConnectMySQLDB()
	sqlQuery := `SELECT
								R.id,
								R.nombre,
								D.nombre,
								R.impacto,
								R.aprovechamiento,
								R.descripcion
							FROM
								residuos R,
								destinos D,
								clave_residuo CR,
								claves C
							WHERE
								CR.id_residuo = R.id
							AND
								CR.id_clave = C.id
							AND R.id_destino = D.id
							AND C.clave LIKE ?
							GROUP BY R.id
							ORDER BY R.nombre;
							`
	rows, err := conn.Query(sqlQuery, "%"+keyword+"%")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	var result []interface{}

	for rows.Next() {
		var r residuo
		err = rows.Scan(&r.Id, &r.Nombre, &r.Destino, &r.Impacto, &r.Aprovechamiento, &r.Descripcion)
		if err != nil {
			panic(err.Error())
		}

		sqlQuery := `SELECT
										C.clave as claves
									FROM
										residuos R,
										claves C,
										clave_residuo CR
									WHERE
										CR.id_residuo = R.id
									AND
										CR.id_clave = C.id
									AND R.id = ?;`
		rowsClaves, err := conn.Query(sqlQuery, r.Id)
		if err != nil {
			panic(err.Error())
		}
		var claves []string
		for rowsClaves.Next() {
			var c clave
			err = rowsClaves.Scan(&c.Clave)
			if err != nil {
				panic(err.Error())
			}
			claves = append(claves, c.Clave)
		}

		r.Claves = claves
		result = append(result, r)
	}
	return result
}

func findAllTips() []interface{} {
	conn := ConnectMySQLDB()
	sqlQuery := `SELECT
								T.id,
								T.tip
							FROM
								tips T;`
	rows, err := conn.Query(sqlQuery)
	if err != nil {
		log.Panic(err.Error())
		handlers := gin.New()
		handlers.Use(gin.Logger())
		handlers.Use(gin.Recovery())
	}
	defer conn.Close()
	var result []interface{}
	for rows.Next() {
		var t tip
		err = rows.Scan(&t.Id, &t.Tip)
		if err != nil {
			log.Panic(err.Error())
			handlers := gin.New()
			handlers.Use(gin.Logger())
			handlers.Use(gin.Recovery())
		}

		result = append(result, t)
	}
	return result
}

func findRandomTip() []interface{} {

	conn := ConnectMySQLDB()
	sqlQuery := `SELECT
								T.id,
								T.tip
							FROM
								tips T
							ORDER BY RAND()
							LIMIT 1;`
	rows, err := conn.Query(sqlQuery)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	var result []interface{}
	for rows.Next() {
		var t tip
		err = rows.Scan(&t.Id, &t.Tip)
		if err != nil {
			panic(err.Error())
		}

		result = append(result, t)
	}

	return result
}

func findAllPuntosLimpios() []interface{} {
	conn := ConnectMySQLDB()
	sqlQuery := `SELECT
								P.id,
								P.nombre,
								P.latitud,
								P.longitud,
								P.descripcion
							FROM
								puntos_limpios P;`
	rows, err := conn.Query(sqlQuery)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	var result []interface{}
	for rows.Next() {
		var p puntoLimpio
		err = rows.Scan(&p.Id, &p.Nombre, &p.Latitud, &p.Longitud, &p.Descripcion)
		if err != nil {
			panic(err.Error())
		}

		result = append(result, p)
	}
	return result
}

func sendRecoleccionData(data recoleccionData) responseStruct {
	if data.Tipo == "" ||
		data.Peso == "" ||
		data.Dimensiones == "" ||
		data.Direccion == "" ||
		data.Ciudad == "" ||
		data.Nombre == "" ||
		data.Telefono == "" ||
		data.Correo == "" {
		return responseStruct{Status: false, Message: "Todos los campos son obligatorios"}
	}

	conn := ConnectMySQLDB()
	sqlQuery := `INSERT INTO
								recolecciones (
									tipo,
									peso,
									dimensiones,
									direccion,
									ciudad,
									nombre,
									telefono,
									correo)
							VALUES (
								?, ?, ?, ?, ?, ?, ?, ?);`
	stmt, err := conn.Prepare(sqlQuery)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	_, err = stmt.Exec(data.Tipo,
		data.Peso,
		data.Dimensiones,
		data.Direccion,
		data.Ciudad,
		data.Nombre,
		data.Telefono,
		data.Correo)
	if err != nil {
		panic(err.Error())
	}

	if sendEmail(data) {
		return responseStruct{
			Status:  true,
			Message: "Solicitud de recolección registrada con éxito",
		}
	} else {
		return responseStruct{
			Status:  false,
			Message: "Error al enviar el correo",
		}
	}

}
