package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"streamvault/internal/api"
	"streamvault/internal/db"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Cargar las variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: No se encontró el archivo .env, se usarán las variables de entorno del sistema.")
	}

	// Leer la configuración de la base de datos desde las variables de entorno
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	jwtSecret := os.Getenv("JWT_SECRET")
	uploadDir := os.Getenv("UPLOAD_DIR")

	// Validar que las variables esenciales estén presentes
	if dbUser == "" || dbPassword == "" || dbName == "" || dbHost == "" || jwtSecret == "" || uploadDir == "" {
		log.Fatal("Error: Una o más variables de entorno requeridas no están definidas (DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, JWT_SECRET, UPLOAD_DIR).")
	}

	// Construir el DSN (Data Source Name) para la conexión a la base de datos
	psqlInfo := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName)

	// Inicializar la conexión con la base de datos desde el paquete 'db'
	database, err := db.InitDB(psqlInfo) // <-- FUNCIÓN LLAMADA DESDE EL NUEVO PAQUETE
	if err != nil {
		log.Fatalf("Error fatal al inicializar la base de datos: %v", err)
	}
	defer database.Close()

	// Crear el directorio de subidas si no existe
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		log.Println("Creando directorio de subidas:", uploadDir)
		os.Mkdir(uploadDir, 0755)
	}

	// Crear una instancia de la aplicación con la configuración cargada
	app := &api.App{
		DB:        database, // Usamos la variable 'database' para evitar conflictos con el nombre del paquete 'db'
		UploadDir: uploadDir,
		JwtSecret: jwtSecret,
	}

	// Configurar el enrutador
	r := api.NewRouter(app)

	// Iniciar el servidor
	port := "8080"
	fmt.Printf("Servidor escuchando en el puerto :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
