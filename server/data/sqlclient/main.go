package sqlclient

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Driver mysql
)

//dbConn Función para simplificar la conexión a la base de datos local
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "datamaster"
	dbPass := "metaphase@07"
	dbName := "e03twitch"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	return db
}
