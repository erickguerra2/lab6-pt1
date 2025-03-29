# La Liga Tracker API

## Descripción

La Liga Tracker API es un servicio RESTful desarrollado en Go que permite gestionar partidos de fútbol. Soporta operaciones CRUD y actualizaciones en tiempo real sobre goles, tarjetas y tiempo extra.

## Características
- Crear, obtener, actualizar y eliminar partidos.
- Registrar goles, tarjetas amarillas, tarjetas rojas y tiempo extra.
- Soporte para CORS.

## Instalación

### Requisitos previos
- Go instalado (versión 1.18 o superior)
- Docker (opcional, para despliegue)

### Clonar el repositorio
```sh
git clone https://github.com/tuusuario/la-liga-tracker.git
cd la-liga-tracker
```

### Ejecutar el servidor
```sh
go run main.go
```
El servidor se iniciará en el puerto `8081`.

## Endpoints

### Obtener todos los partidos
```http
GET /api/matches
```
**Respuesta:**
```json
[
  {
    "id": 1,
    "homeTeam": "Equipo A",
    "awayTeam": "Equipo B",
    "matchDate": "2025-03-28",
    "goalCount": 2,
    "yellowCards": 1,
    "redCards": 0,
    "extraTime": false
  }
]
```

### Obtener un partido por ID
```http
GET /api/matches/{id}
```

### Crear un nuevo partido
```http
POST /api/matches
```
**Cuerpo de la petición:**
```json
{
  "homeTeam": "Equipo A",
  "awayTeam": "Equipo B",
  "matchDate": "2025-03-28"
}
```

### Actualizar un partido
```http
PUT /api/matches/{id}
```

### Eliminar un partido
```http
DELETE /api/matches/{id}
```

### Registrar un gol
```http
PATCH /api/matches/{id}/goals
```

### Registrar tarjeta amarilla
```http
PATCH /api/matches/{id}/yellowcards
```

### Registrar tarjeta roja
```http
PATCH /api/matches/{id}/redcards
```

### Agregar tiempo extra
```http
PATCH /api/matches/{id}/extratime
```

## Despliegue con Docker

```sh
docker build -t la-liga-tracker .
docker run -p 8081:8081 la-liga-tracker
```
