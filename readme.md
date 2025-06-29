<<<<<<< HEAD
# StreamVault

## Acerca del Proyecto

**StreamVault** es una plataforma de video simple, eficiente y controlada por el usuario. El backend, construido enteramente en **Go**, está diseñado para ser ligero, rápido y fácil de desplegar.

La arquitectura se basa en principios de software moderno, utilizando una capa de datos abstraída mediante **interfaces** para desacoplar la lógica de negocio de la base de datos, lo que lo hace flexible y fácil de mantener.

## 🚀 Características Principales

* **Gestión de Usuarios Segura**: Registro, login con JWT y verificación obligatoria por correo electrónico.
* **Roles de Usuario**: Distinción entre usuarios normales y administradores con permisos específicos.
* **API RESTful Completa**: Nueve endpoints para gestionar usuarios y videos (CRUD completo para videos por parte del admin).
* **Arquitectura Desacoplada**: Uso de interfaces para la capa de datos, facilitando la testabilidad y el cambio de motor de base de datos.
* **Rendimiento y Concurrencia**: Aprovecha las `goroutines` de Go para tareas asíncronas como el envío de correos, sin afectar la experiencia del usuario.
* **Configuración Sencilla**: Todo se configura a través de variables de entorno.

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
* Un servicio SMTP (como [Brevo](https://www.brevo.com/)) para el envío de correos.

### Instalación

1.  **Clona el repositorio**
    ```sh
    git clone [https://github.com/Missifus/streamvault.git]
    ```
2.  **Configura tu base de datos**
    * Crea una base de datos en PostgreSQL.
    * Ejecuta el script SQL del proyecto para crear las tablas.
3.  **Configura las variables de entorno**
    * Copia `env.example` a un nuevo archivo llamado `.env`.
    * Rellena `.env` con tus credenciales de la base de datos, el secreto JWT y las credenciales SMTP.
4.  **Instala las dependencias de Go**
    ```sh
    go mod tidy
    ```
5.  **Inicia el servidor**
    ```sh
    go run cmd/api/main.go
    ```
6.  **Abre el frontend**
    * Navega a la carpeta `web/` y abre `index.html` en tu navegador.

    Made with love by missifus <3
=======
StreamVault
Acerca del Proyecto
StreamVault nace de la necesidad de tener una plataforma de video simple, eficiente y controlada por el usuario. En lugar de depender de servicios de terceros, este proyecto te da las herramientas para construir tu propio "Netflix" personal o para tu organización. El backend, construido enteramente en Go, está diseñado para ser ligero, rápido y fácil de desplegar.

La arquitectura se basa en principios de software moderno, utilizando una capa de datos abstraída mediante interfaces para desacoplar la lógica de negocio de la base de datos, lo que lo hace flexible y fácil de mantener.

🚀 Características Principales
Gestión de Usuarios Segura: Registro, login con JWT y verificación obligatoria por correo electrónico.

Roles de Usuario: Distinción entre usuarios normales y administradores con permisos específicos.

API RESTful Completa: Nueve endpoints para gestionar usuarios y videos (CRUD completo para videos por parte del admin).

Arquitectura Desacoplada: Uso de interfaces para la capa de datos, facilitando la testabilidad y el cambio de motor de base de datos.

Rendimiento y Concurrencia: Aprovecha las goroutines de Go para tareas asíncronas como el envío de correos, sin afectar la experiencia del usuario.

Configuración Sencilla: Todo se configura a través de variables de entorno.

🛠️ Construido Con
Esta es la tecnología que impulsa StreamVault:

Backend: Go

Base de Datos: PostgreSQL

Enrutador: gorilla/mux

Autenticación: golang-jwt/jwt

Frontend: HTML, Tailwind CSS, JavaScript (Modular)

⚙️ Primeros Pasos
Para tener una copia local funcionando, sigue estos sencillos pasos.

Prerrequisitos
Asegúrate de tener instalado:

Go (versión 1.18+)

PostgreSQL

Un servicio SMTP (como Brevo) para el envío de correos.

Instalación
Clona el repositorio

git clone https://github.com/tu_usuario/streamvault.git

Configura tu base de datos

Crea una base de datos en PostgreSQL.

Ejecuta el script SQL del proyecto para crear las tablas.

Configura las variables de entorno

Copia env.example a un nuevo archivo llamado .env.

Rellena .env con tus credenciales de la base de datos, el secreto JWT y las credenciales SMTP.

Instala las dependencias de Go

go mod tidy

Inicia el servidor

go run cmd/api/main.go

Abre el frontend

Navega a la carpeta web/ y abre index.html en tu navegador.

Made with love by missifus <3
---
>>>>>>> 56248dec56c5c90427fbce668945fca80a8045e0
