# StreamVault
<p align="center">
  <a href="#">
    <img src="./assets/logo.svg" alt="Logo de StreamVault" width="150">
  </a>

  <h3 align="center">Una plataforma de streaming de video autohospedada, desarrollada en Go.</h3>


</p>

## Acerca del Proyecto

**StreamVault** nace de la necesidad de tener una plataforma de video simple, eficiente y controlada por el usuario. En lugar de depender de servicios de terceros, este proyecto te da las herramientas para construir tu propio "Netflix" personal o para tu organización. El backend, construido enteramente en **Go**, está diseñado para ser ligero, rápido y fácil de desplegar.

La arquitectura se basa en principios de software moderno, utilizando una capa de datos abstraída mediante **interfaces** para desacoplar la lógica de negocio de la base de datos, lo que lo hace flexible y fácil de mantener.

## 🚀 Características Principales

* **Gestión de Usuarios**: Registro y login directos mediante `username` y `email`.
* **Roles de Usuario (Admin/User)**: Clara distinción entre usuarios normales y administradores con permisos específicos.
* **Panel de Administración Completo**: Una interfaz para que los administradores puedan listar, cambiar el rol y eliminar usuarios, así como gestionar todos los videos subidos.
* **API RESTful Robusta**: **11 endpoints** funcionales que cubren la autenticación, la gestión de contenido y la administración de la plataforma.
* **Arquitectura Desacoplada con Interfaces**: El uso de una capa de datos abstracta (`DataStore`) facilita la testabilidad y la posibilidad de cambiar el motor de base de datos en el futuro.
* **Concurrencia**: Se aprovechan las `goroutines` de Go para tareas en segundo plano (como el procesamiento de video) sin afectar la experiencia del usuario.
* **Configuración Sencilla**: Todo se configura a través de un único archivo `.env`.

## 🛠️ Construido Con

Esta es la tecnología que impulsa StreamVault:

* **Backend**: Go
* **Base de Datos**: PostgreSQL
* **Enrutador**: `gorilla/mux`
* **Autenticación**: `golang-jwt/jwt`
* **Frontend**: HTML, Tailwind CSS, JavaScript (Modular)

## ⚙️ Primeros Pasos

Para tener una copia local funcionando, sigue estos sencillos pasos.

### Prerrequisitos

Asegúrate de tener instalado:
* Go (versión 1.18+)
* Un Servidor de Base de Datos PostgreSQL

### Instalación

1.  **Clona el repositorio**
    ```sh
    git clone https://github.com/tu_usuario/streamvault.git
    ```
2.  **Configura tu base de datos**
    * Asegúrate de que tu servidor PostgreSQL esté corriendo.
    * Crea una base de datos vacía (ej: streaming_db). El nombre debe coincidir con el que pondrás en tu archivo .env.
    * ¡No necesitas ejecutar ningún script SQL! La aplicación creará las tablas necesarias automáticamente en su primer inicio.
3.  **Configura las variables de entorno**
    * Copia `env.example` a un nuevo archivo llamado `.env`.
    * Rellena `.env` con tus credenciales de la base de datos y un secreto para JWT.
4.  **Instala las dependencias de Go**
    ```sh
    go mod tidy
    ```
5.  **Inicia el servidor**
    ```sh
    go run cmd/api/main.go
    ```
6.  **Abre el frontend**
    * Accede a la aplicación a través de `http://localhost:PUERTO`.
  
7.  **Crea un primer usuario y cambia manualmente su rol a admin, en la base de datos **
    ```sql
    UPDATE users SET role = 'admin' WHERE email = 'correo_del_usuario@ejemplo.com'; 
    ```

---

## 📋 Servicios Web

| Método | Ruta                      | Descripción                                 | Protegido (Admin) |
| :----- | :------------------------ | :------------------------------------------ | :---------------: |
| `POST` | `/api/register`           | Registra un nuevo usuario.                  |         No        |
| `POST` | `/api/login`              | Inicia sesión y obtiene un token JWT.       |         No        |
| `GET`  | `/api/videos`             | Obtiene la lista de todos los videos.       |         No        |
| `GET`  | `/api/videos/{id}`        | Obtiene los detalles de un video específico.|         No        |
| `GET`  | `/stream/{filename}`      | Sirve el archivo de video para streaming.   |         No        |
| `POST` | `/api/admin/upload`       | Sube un nuevo archivo de video.             |        **Sí** |
| `PUT`  | `/api/admin/videos/{id}`  | Actualiza los detalles de un video.         |        **Sí** |
| `DELETE`| `/api/admin/videos/{id}`  | Elimina un video y su archivo físico.       |        **Sí** |
| `GET`  | `/api/admin/users`        | Obtiene la lista de todos los usuarios.     |        **Sí** |
| `PUT`  | `/api/admin/users/{id}/role` | Actualiza el rol de un usuario.            |        **Sí** |
| `DELETE`| `/api/admin/users/{id}`   | Elimina un usuario del sistema.             |        **Sí** |

---
### Digrama de clases 
  ```mermaid
classDiagram
    direction BT
    
    class App {
        <<Struct>>
        +Store: DataStore
        +UploadDir: string
        +JwtSecret: string
    }

    class handler {
        <<Struct>>
        -app: App
        +HandleRegisterUser(w, r)
        +HandleLoginUser(w, r)
        +HandleUploadVideo(w, r)
        +HandleDeleteVideo(w, r)
        +...
    }

    class DataStore {
        <<Interface>>
        +Init() error
        +CreateUser(*User) error
        +GetUserByEmail(string) (*User, error)
        +CreateVideo(*Video) error
        +GetAllVideos() ([]*Video, error)
        +...
    }

    class PostgresStore {
        <<Struct>>
        -db: *sql.DB
        +Init() error
        +CreateUser(*User) error
        +GetUserByEmail(string) (*User, error)
        +CreateVideo(*Video) error
        +GetAllVideos() ([]*Video, error)
        +...
    }

    class User {
        <<Model>>
        +ID: int
        +Username: string
        +Email: string
        +Role: string
    }

    class Video {
        <<Model>>
        +ID: int
        +Title: string
        +Category: string
        +FilePath: string
    }

    class main_go["main.go"] {
        <<Punto de Entrada>>
        +main()
    }
    
    main_go --> App : "Creates & Injects"
    main_go --> PostgresStore : "Creates Instance"
    
    App --> handler : "Is used by"
    handler "1" -- "1" App : "Contains"

    App o-- "1" DataStore : "Depends on (DI)"
    
    PostgresStore --|> DataStore : "Implements"
    
    handler ..> DataStore : "Uses"
    PostgresStore ..> User : "Manipulates"
    PostgresStore ..> Video : "Manipulates"

    class Frontend_JS["Frontend (JS)"] {
        <<Conceptual>>
    }

    class app_js["app.js"] {
        +main()
        +router()
    }

    class api_js["api.js"] {
        +request()
        +loginUser()
        +getVideos()
    }

    class ui_js["ui.js"] {
        +showSection()
        +updateAuthUI()
        +renderVideos()
    }

    Frontend_JS o-- "1" app_js : "Orchestrates"
    app_js ..> api_js : "Uses"
    app_js ..> ui_js : "Uses"

    api_js ..> handler : "Makes HTTP requests to"
    
  ```
##  Made with love by missifus <3
