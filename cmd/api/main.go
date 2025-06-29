package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"streamvault/internal/api"
	"streamvault/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: No se encontró el archivo .env")
	}

	// ... (código para cargar variables de entorno es el mismo)
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	jwtSecret := os.Getenv("JWT_SECRET")
	uploadDir := os.Getenv("UPLOAD_DIR")

	psqlInfo := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName)

	// Inicializar la implementación concreta del store
	store, err := storage.NewPostgresStore(psqlInfo)
	if err != nil {
		log.Fatalf("Error fatal al inicializar el data store: %v", err)
	}

	log.Println("Conexión a la base de datos y store inicializado exitosamente.")

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	// Crear la App, inyectando la INTERFAZ, no la implementación concreta
	app := &api.App{
		Store:     store,
		UploadDir: uploadDir,
		JwtSecret: jwtSecret,
	}

	r := api.NewRouter(app)

	port := "8080"
	fmt.Printf("Servidor escuchando en el puerto :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
