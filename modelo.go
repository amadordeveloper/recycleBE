package main

func findAllResiduos() []interface{} {
	// use ConnectMySQLDB to connect to the database and query all data from table residuos

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
	// iterate over the rows using struct residuo
	var result []interface{}
	for rows.Next() {
		var r residuo
		// scan the current row into the struct transforming claves into an array of strings
		err = rows.Scan(&r.Id, &r.Nombre, &r.Destino, &r.Impacto, &r.Aprovechamiento, &r.Descripcion)
		if err != nil {
			panic(err.Error())
		}
		// query the claves table and get the claves for the current residuo
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
		// add the claves array to the residuo struct
		r.Claves = claves
		// add the residuo struct to the result array
		result = append(result, r)
	}
	return result
}

func findResiduosByClave(keyword string) interface{} {
	// use ConnectMySQLDB to connect to the database and query all data from table residuos where clave LIKE %keyword%

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
							AND C.clave LIKE ?;`
	rows, err := conn.Query(sqlQuery, "%"+keyword+"%")

	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	// iterate over the rows using struct residuo
	var result []interface{}
	for rows.Next() {
		var r residuo
		// scan the current row into the struct transforming claves into an array of strings
		err = rows.Scan(&r.Id, &r.Nombre, &r.Destino, &r.Impacto, &r.Aprovechamiento, &r.Descripcion)
		if err != nil {
			panic(err.Error())
		}
		// query the claves table and get the claves for the current residuo
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
		// add the claves array to the residuo struct
		r.Claves = claves
		// add the residuo struct to the result array
		result = append(result, r)
	}
	return result
}
