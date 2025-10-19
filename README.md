# 🖼️ ImgAI — REST API de Galería de Imágenes Inteligente

**ImgAI** es una API REST desarrollada en **Go (Golang)** utilizando **Gin**, **GORM**, **MySQL** y **OpenAI API**, que permite a los usuarios:
- Subir, listar y eliminar fotos.
- Organizar fotos en álbumes.
- Marcar favoritos.
- Analizar automáticamente las imágenes usando IA para generar descripciones y etiquetas.
<img width="400" height="768" alt="image" src="https://github.com/user-attachments/assets/862985dc-d1b1-41f4-b2a7-51d49fb836ba" />

---

## 🚀 Tecnologías principales

| Componente | Descripción |
|-------------|--------------|
| **Go 1.22+** | Lenguaje principal del backend |
| **Gin Gonic** | Framework HTTP rápido y minimalista |
| **GORM** | ORM para manejo de base de datos relacional |
| **MySQL / MariaDB** | Motor de base de datos |
| **OpenAI API** | Genera descripciones y etiquetas de imágenes |
| **bcrypt** | Hash seguro de contraseñas |
| **JWT (JSON Web Tokens)** | Autenticación basada en tokens |
| **Docker (opcional)** | Despliegue y pruebas locales |
| **Postman** | Colección de pruebas de endpoints incluida |
| **MINIO** | Almacenamiento en buckets (s3-like)
---
## Guia de Instalacion

### 1. Clonar el repositorio

### 2. Crea un archivo .env con los sig valores:
DATABASE_URL=root:123456@tcp(db:3306)/imagAI?parseTime=true
DB_PASSWORD=123456
DB_NAME=imagAI
S3_ENDPOINT=http://s3:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_REGION=us-east-1
S3_BUCKET=imagai
SECRET_JWT=<some random value>
OPENAI_API_KEY=<keyvalue>

### 3. Haz la imagen el servicio api
Ejecuta el comando: 
```bash
  docker build -t imgai-api
```

### 4. Levanta todos los servicios con Docker Compose

```
  docker-compose up -d --build
```

### 5. Importa la coleccion del servicio en Postman

<img width="572" height="884" alt="image" src="https://github.com/user-attachments/assets/e1138b3b-8df0-49e5-a0be-8f38ff6cbabe" />













