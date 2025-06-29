# StreamVault
<p align="center">
  <a href="#">
    <img src="https://placehold.co/150x150/4f46e5/ffffff?text=SV" alt="Logo de StreamVault">
  </a>

  <h3 align="center">Una plataforma de streaming de video autohospedada, desarrollada en Go.</h3>

  <p align="center">
    Control total sobre tu contenido, desde la subida hasta el streaming.
    <br />
    <a href="#"><strong>Explorar la documentación (próximamente)</strong></a>
    ·
    <a href="#">Reportar un Bug</a>
    ·
    <a href="#">Solicitar una Característica</a>
  </p>
</p>

## Acerca del Proyecto

**StreamVault** nace de la necesidad de tener una plataforma de video simple, eficiente y controlada por el usuario. En lugar de depender de servicios de terceros, este proyecto te da las herramientas para construir tu propio "Netflix" personal o para tu organización. El backend, construido enteramente en **Go**, está diseñado para ser ligero, rápido y fácil de desplegar.

La arquitectura se basa en principios de software moderno, utilizando una capa de datos abstraída mediante **interfaces** para desacoplar la lógica de negocio de la base de datos, lo que lo hace flexible y fácil de mantener.

## 🚀 Características Principales

* **Gestión de Usuarios Simplificada**: Registro y login directos mediante `username` y `email`.
* **Roles de Usuario (Admin/User)**: Clara distinción entre usuarios normales y administradores con permisos específicos.
* **Panel de Administración Completo**: Una interfaz para que los administradores puedan listar, cambiar el rol y eliminar usuarios, así como gestionar todos los videos subidos.
* **API RESTful Robusta**: **11 endpoints** funcionales que cubren la autenticación, la gestión de contenido y la administración de la plataforma.
* **Arquitectura Desacoplada con Interfaces**: El uso de una capa de datos abstracta (`DataStore`) facilita la testabilidad y la posibilidad de cambiar el motor de base de datos en el futuro.
* **Demostración de Concurrencia**: Se aprovechan las `goroutines` de Go para simular tareas en segundo plano (como el procesamiento de video) sin afectar la experiencia del usuario.
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
* PostgreSQL

### Instalación

1.  **Clona el repositorio**
    ```sh
    git clone git clone https://github.com/tu_usuario/streamvault.git
    ```
2.  **Configura tu base de datos**
    * Crea una base de datos en PostgreSQL (ej: `streaming_db`).
    * Ejecuta el script SQL del proyecto para crear las tablas `users` y `videos`.
3.  **Configura las variables de entorno**
    * Copia `env.example` (si existe) a un nuevo archivo llamado `.env`.
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
    * Usa una extensión como "Live Server" en VS Code sobre el archivo `web/index.html` o sirve la carpeta `web/` con un servidor local (`python -m http.server`).
    * Accede a la aplicación a través de `http://localhost:PUERTO`.

---

## 📋 Servicios Web Implementados (11 en total)

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

## Made with love by missifus <3