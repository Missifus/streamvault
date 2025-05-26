# StreamVault  
**Sistema de Gestión de Streaming Autohospedado y Seguro**  

StreamVault es una plataforma de streaming **autohospedada** desarrollada en Go, diseñada para almacenar, cifrar y transmitir videos de forma segura. Ideal para proyectos personales o empresariales que requieren control total sobre su contenido.  

---

## 🚀 Características  
- **Autenticación Segura**:  
  - Registro con verificación por email (SMTP).  
  - Contraseñas hasheadas con bcrypt.  
  - Roles de usuario (Admin/Usuario).  
- **Gestión de Contenido**:  
  - Cifrado AES-256 para videos.  
  - Transcodificación a HLS usando FFmpeg.  
  - Almacenamiento local en carpetas estructuradas.  
- **Interfaz Web Mínima**:  
  - Páginas estáticas (login, registro, lista de videos).  
  - Reproductor HTML5 con soporte HLS.  

Made with love by missifus <3
---