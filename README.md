# Woody's Backend API

![Go Version](https://img.shields.io/badge/Go-1.24.4-00ADD8?logo=go&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-1.30.0-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14%2B-336791?logo=postgresql&logoColor=white)
![gorilla/mux](https://img.shields.io/badge/gorilla%2Fmux-1.8.1-green)
![License](https://img.shields.io/badge/License-MIT-yellow)

API REST para una plataforma de compartición de proyectos de carpintería. Permite gestionar usuarios (autenticados vía Firebase), proyectos, comentarios con respuestas, likes/dislikes en comentarios, valoraciones de proyectos y listas de proyectos de usuarios.

## Tech Stack

- **Go 1.24.4**
- **gorilla/mux** — Router HTTP
- **GORM + gorm.io/driver/postgres** — ORM y driver de PostgreSQL
- **joho/godotenv** — Carga de variables de entorno
- **PostgreSQL** — Base de datos

## Estructura del proyecto

```
woodys-backend/
├── main.go             # Entry point: conexión DB, router y registro de rutas
├── config/             # Carga de variables de entorno y DSN
├── db/                 # Conexión GORM a PostgreSQL
├── middlewares/        # CORS, Content-Type JSON, recálculo de ratings
├── models/             # Modelos GORM (User, Project, Comment, CommentLike, Rating, ...)
├── routes/             # Handlers HTTP organizados por dominio
├── Makefile            # Comandos de build, run, lint
├── .golangci.yml       # Configuración de linters
├── .env.example        # Plantilla de variables de entorno
└── go.mod / go.sum     # Dependencias del módulo
```

## Requisitos

- Go 1.24.4 o superior
- PostgreSQL

## Instalación y ejecución

1. **Clonar el repositorio**

   ```bash
   git clone <repository-url>
   cd woodys-backend
   ```

2. **Instalar dependencias**

   ```bash
   go mod download
   ```

3. **Configurar variables de entorno**

   ```bash
   cp .env.example .env
   # Editar .env con los valores reales
   ```

4. **Correr el servidor**

   ```bash
   go run main.go
   # o con Make:
   make run
   ```

   El servidor escucha en `:8080`.

> **Nota:** el esquema de la base de datos se gestiona manualmente; el proyecto no ejecuta migraciones automáticas al iniciar.

## Variables de entorno

| Variable       | Descripción                              | Default     | Requerida |
| -------------- | ---------------------------------------- | ----------- | --------- |
| `DB_HOST`      | Host de la base de datos                 | `localhost` | Sí        |
| `DB_USER`      | Usuario de la base de datos              | `postgres`  | Sí        |
| `DB_PASSWORD`  | Contraseña de la base de datos           | —           | Sí        |
| `DB_NAME`      | Nombre de la base de datos               | `postgres`  | Sí        |
| `DB_PORT`      | Puerto de la base de datos               | `5432`      | No        |
| `DB_SSL_MODE`  | Modo SSL (`require` / `disable`)         | `disable`   | No        |
| `SERVER_PORT`  | Puerto del servidor HTTP                 | `8080`      | No        |

## Comandos Make

| Comando              | Descripción                                          |
| -------------------- | ---------------------------------------------------- |
| `make run`           | Corre la aplicación                                  |
| `make build`         | Compila el binario en `./bin/woodys-backend`         |
| `make run-bin`       | Compila y corre el binario                           |
| `make lint`          | Corre los linters (golangci-lint si está instalado)  |
| `make fmt`           | Formatea el código (gofmt + goimports)               |
| `make vet`           | Corre `go vet`                                       |
| `make check`         | Corre fmt + lint + vet (pre-commit)                  |
| `make deps`          | Descarga y ordena dependencias                       |
| `make clean`         | Limpia artefactos de build                           |
| `make install-tools` | Instala herramientas de desarrollo                   |

## Endpoints

### Stats

| Método | Ruta     | Descripción                                        |
| ------ | -------- | -------------------------------------------------- |
| GET    | `/stats` | Estadísticas de la landing (proyectos, usuarios, valoraciones, promedio) |

### Profile Pictures

| Método | Ruta                    | Descripción                    |
| ------ | ----------------------- | ------------------------------ |
| GET    | `/profile-pictures`     | Lista todas las fotos de perfil |
| GET    | `/profile-picture/{id}` | Obtiene una foto de perfil     |

### Users

| Método | Ruta                        | Descripción                       |
| ------ | --------------------------- | --------------------------------- |
| GET    | `/users/{id}`               | Obtiene un usuario por ID         |
| GET    | `/users/uid/{firebase_uid}` | Obtiene el ID de un usuario por su Firebase UID |
| GET    | `/users/{id}/projects`      | Lista los proyectos de un usuario |
| POST   | `/users`                    | Crea un usuario                   |
| PUT    | `/users/{id}`               | Actualiza un usuario              |
| DELETE | `/users/{id}`               | Elimina un usuario                |

### Projects

| Método | Ruta               | Descripción              |
| ------ | ------------------ | ------------------------ |
| GET    | `/projects/search` | Busca proyectos          |
| GET    | `/projects/{id}`   | Obtiene un proyecto      |
| POST   | `/projects`        | Crea un proyecto         |
| PUT    | `/projects/{id}`   | Actualiza un proyecto    |
| DELETE | `/projects/{id}`   | Elimina un proyecto      |

### Comments

| Método | Ruta                          | Descripción                        |
| ------ | ----------------------------- | ---------------------------------- |
| GET    | `/projects/{id}/comments`     | Lista los comentarios de un proyecto |
| POST   | `/projects/{id}/comments`     | Crea un comentario                 |
| POST   | `/comments/{id}/reply`        | Responde a un comentario           |
| GET    | `/comments/{id}/replies`      | Lista las respuestas de un comentario |
| DELETE | `/comments/{id}`              | Elimina (o marca) un comentario    |

### Comment Likes

| Método | Ruta                                | Descripción                                    |
| ------ | ----------------------------------- | ---------------------------------------------- |
| POST   | `/comments/{id}/likes`              | Crea una reacción (body: `user_id`, `value`: `"like"` o `"dislike"`) |
| GET    | `/comments/{id}/likes`              | Conteos agregados (`{"likes": N, "dislikes": M}`) |
| GET    | `/comments/{id}/likes/{user_id}`    | Reacción de un usuario a un comentario         |
| PUT    | `/comment-likes/{id}`               | Invierte la reacción (toggle like ↔ dislike)   |
| DELETE | `/comment-likes/{id}`               | Elimina la reacción                            |

### Ratings

| Método | Ruta                                  | Descripción                                  |
| ------ | ------------------------------------- | -------------------------------------------- |
| POST   | `/projects/{id}/ratings`              | Crea una valoración (body: `user_id`, `value`) |
| GET    | `/projects/{id}/ratings`              | Lista las valoraciones de un proyecto        |
| GET    | `/projects/{id}/ratings/{user_id}`    | Valoración de un usuario a un proyecto       |
| PUT    | `/projects/{id}/ratings`              | Actualiza la valoración (body: `user_id`, `value`) |
| DELETE | `/projects/{id}/ratings/{user_id}`    | Elimina la valoración de un usuario          |

### Project Lists

| Método | Ruta                                          | Descripción                     |
| ------ | --------------------------------------------- | ------------------------------- |
| GET    | `/users/{id}/project-lists`                   | Lista las listas de un usuario  |
| GET    | `/project-lists/{id}`                         | Obtiene una lista               |
| POST   | `/project-lists`                              | Crea una lista                  |
| PUT    | `/project-lists/{id}`                         | Actualiza una lista             |
| DELETE | `/project-lists/{id}`                         | Elimina una lista               |
| POST   | `/project-lists/{id}/projects`                | Agrega un proyecto a una lista  |
| GET    | `/project-lists/{id}/projects`                | Lista los proyectos de una lista |
| DELETE | `/project-lists/{list_id}/projects/{project_id}` | Quita un proyecto de una lista |

## Reglas de negocio

- **Borrado de comentarios con respuestas:** `DELETE /comments/{id}` no elimina físicamente un comentario que tiene respuestas; en su lugar, su contenido pasa a ser `"Comentario Eliminado"`, preservando el hilo de respuestas. Si no tiene respuestas, se borra físicamente.
- **Toggle de reacciones:** `PUT /comment-likes/{id}` no recibe body; simplemente invierte el valor actual (`true` → `false`, like ↔ dislike).
- **Unicidad por usuario:** un usuario solo puede tener una valoración por proyecto y una reacción por comentario. Intentar crear una duplicada devuelve `409 Conflict`.
- **Ratings agregados:** los campos `average_rating` y `rating_count` de cada proyecto se recalculan automáticamente ante cada creación, actualización o eliminación de valoraciones.
- **Respuestas de borrado:** todos los handlers de eliminación responden `200 OK` con un mensaje JSON de confirmación (`{"message": "... deleted successfully"}`).

## Linting

El proyecto usa [golangci-lint](https://golangci-lint.run/) configurado en `.golangci.yml`. Antes de commitear:

```bash
make check   # fmt + lint + vet
```
