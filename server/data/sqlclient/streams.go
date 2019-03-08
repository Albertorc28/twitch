package sqlclient

import (
	"fmt"
	"log"
	m "twitch/server/data/model"
)

//Streams Función que devuelve los streams en directo
func Streams() []m.RStream {
	db := dbConn()
	defer db.Close()

	//Consulta
	query, err := db.Query(
		`SELECT ID, usuario, inicio, expiracion
		FROM stream`)

	if err != nil {
		panic(err.Error())
	}
	defer query.Close()

	//Slice de streams vacío
	var streams = make([]m.RStream, 0)

	//Recorrido de las filas
	for query.Next() {
		var s = m.RStream{}
		if err := query.Scan(&s.ID, &s.Usuario, &s.Inicio, &s.Expira); err != nil {
			log.Fatal(err)
		}
		streams = append(streams, s)
	}

	return streams
}

//CrearStream inserta un nuevo stream en la base de datos
func CrearStream(stream *m.RStream) {
	db := dbConn()
	defer db.Close()

	//Consulta
	res, err := db.Exec(
		`INSERT INTO stream(usuario, inicio, expiracion)
		VALUES(?, ?, ?)`,
		stream.Usuario,
		stream.Inicio,
		stream.Expira)

	if err != nil {
		panic(err.Error())
	}
	stream.ID, err = res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
}

//ActualizarStream mantiene activo el stream con el ID indicado durante 'ttl' minutos más
func ActualizarStream(ID int64, ttl int) {
	db := dbConn()
	defer db.Close()

	//Consulta
	res, err := db.Exec(
		`UPDATE Stream
		SET expiracion = DATE_ADD(expiracion, INTERVAL ? MINUTE)
		WHERE ID = ?`,
		ttl,
		ID)

	if err != nil {
		panic(err.Error())
	}

	filas, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	if filas <= 0 {
		fmt.Printf("La actualización del stream %d no produjo resultados", ID)
	}
}
