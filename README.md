# ImgAI
## Objetivo
Crear una aplicación web donde los usuarios puedan subir, visualizar, organizar y gestionar imágenes de manera segura y eficiente, con un backend robusto en Go y un frontend interactivo en React.

<img width="400" height="768" alt="image" src="https://github.com/user-attachments/assets/862985dc-d1b1-41f4-b2a7-51d49fb836ba" />

### Arquitectura General
- Frontend (React)
- Interfaz de usuario intuitiva para:
- Subir imágenes (drag & drop o selector de archivos)
- Visualizar galerías
- Buscar imágenes por nombre o tags

### Estado y lógica:
- React + React Router (para rutas internas)
- Context API
- Comunicación con backend mediante API REST

Backend (Go)
- Servidor HTTP usando Gin

Funcionalidades:
- CRUD de imágenes (crear, leer, actualizar, eliminar)
- Autenticación y autorización
- Gestión de metadatos (nombre, etiquetas, fecha, tamaño, formato)

Almacenamiento:
- En la nube (S3<MinIO>)
- Base de datos para metadatos (PostgreSQL o MySQL)
- Validación de archivos (tipo, tamaño)

Base de datos
- Tabla de usuarios
- Tabla de imágenes
- Tabla de álbumes o categorías (opcional)
- Almacenamiento de archivos

Opciones:
- Local: Carpeta en servidor con rutas en DB
- Cloud: S3, MinIO (local para pruebas)

