// File: internal/db/database.go
package db // El nombre del paquete ahora es 'db'

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Driver de PostgreSQL
)

// InitDB inicializa y devuelve una conexión a la base de datos PostgreSQL.
// Acepta la cadena de conexión (DSN) como argumento.
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la conexión con la base de datos: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error al hacer ping a la base de datos: %w", err)
	}

	log.Println("¡Conexión exitosa a la base de datos!")
	return db, nil
}
