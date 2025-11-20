# Authentication Microservice

Microservicio de autenticación para Bikes2Road que gestiona la creación y validación de tokens JWT. Este servicio se comunica con el microservicio de usuarios para validar credenciales y generar tokens de acceso.

## Arquitectura

Este proyecto sigue una **arquitectura hexagonal** (también conocida como arquitectura de puertos y adaptadores) con la siguiente estructura:

```
authentication/
├── cmd/api/                    # Punto de entrada de la aplicación
│   ├── config/                # Configuración
│   ├── container/             # Inyección de dependencias
│   └── main.go               # Punto de entrada principal
├── internal/
│   ├── adapters/             # Adaptadores (capa de infraestructura)
│   │   ├── http/            # Adaptador HTTP (entrada)
│   │   └── httpclient/      # Cliente HTTP (salida)
│   ├── domain/              # Entidades y lógica de negocio
│   ├── ports/               # Interfaces (contratos)
│   └── services/            # Servicios de aplicación
└── docs/                    # Documentación Swagger generada
```

## Características

- ✅ Autenticación basada en JWT
- ✅ Tokens de acceso y refresh
- ✅ Validación de tokens
- ✅ Integración con microservicio de usuarios
- ✅ Documentación Swagger/OpenAPI
- ✅ Arquitectura hexagonal
- ✅ Health check endpoint
- ✅ Dockerizado

## Tecnologías

- **Go 1.21+**
- **Gin** - Framework web
- **JWT** - golang-jwt/jwt/v5
- **Swagger** - Documentación API
- **Docker** - Containerización

## Endpoints

### Autenticación

#### POST /api/v1/auth/login
Autentica un usuario y retorna tokens JWT.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe"
  },
  "tokens": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

#### POST /api/v1/auth/validate
Valida un token JWT.

**Request:**
```json
{
  "token": "eyJhbGc..."
}
```

**Response:**
```json
{
  "valid": true,
  "claims": {
    "user_id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "exp": 1234567890,
    "iat": 1234567890
  }
}
```

#### POST /api/v1/auth/refresh
Refresca un token JWT usando el refresh token.

**Request:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response:**
```json
{
  "tokens": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

### Health Check

#### GET /health
Verifica que el servicio esté funcionando.

**Response:**
```json
{
  "status": "OK"
}
```

## Variables de Entorno

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `PORT` | Puerto del servidor | `8080` |
| `HOST` | Host del servidor | `0.0.0.0` |
| `JWT_SECRET_KEY` | Clave secreta para firmar JWT | **Requerido** |
| `JWT_ACCESS_TOKEN_EXPIRATION` | Expiración del access token (horas) | `24` |
| `JWT_REFRESH_TOKEN_EXPIRATION` | Expiración del refresh token (horas) | `24` |
| `USERS_SERVICE_URL` | URL del microservicio de usuarios | `http://localhost:8083` |

## Instalación y Ejecución

### Requisitos Previos

- Go 1.21 o superior
- Microservicio de usuarios ejecutándose

### Configuración

1. Clonar el repositorio:
```bash
git clone <repository-url>
cd authentication
```

2. Copiar el archivo de ejemplo de variables de entorno:
```bash
cp .env.example .env
```

3. Editar `.env` y configurar las variables necesarias:
```bash
# IMPORTANTE: Cambiar JWT_SECRET_KEY en producción
JWT_SECRET_KEY=your-super-secret-key-change-this-in-production
USERS_SERVICE_URL=http://localhost:8083
```

### Ejecución Local

1. Instalar dependencias:
```bash
go mod download
```

2. Generar documentación Swagger:
```bash
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

3. Ejecutar la aplicación:
```bash
go run cmd/api/main.go
```

O compilar y ejecutar:
```bash
go build -o bin/auth-service cmd/api/main.go
./bin/auth-service
```

La aplicación estará disponible en `http://localhost:8080`

### Ejecución con Docker

1. Construir la imagen:
```bash
docker build -t bikes2road/authentication:latest .
```

2. Ejecutar el contenedor:
```bash
docker run -p 8080:8080 \
  -e JWT_SECRET_KEY=your-secret-key \
  -e USERS_SERVICE_URL=http://users-service:8083 \
  bikes2road/authentication:latest
```

## Documentación API

Una vez que la aplicación esté ejecutándose, la documentación Swagger estará disponible en:

```
http://localhost:8080/swagger/index.html
```

## Desarrollo

### Estructura del Proyecto

- **cmd/api**: Punto de entrada y configuración de la aplicación
- **internal/domain**: Entidades de dominio y lógica de negocio
- **internal/ports**: Interfaces que definen los contratos
- **internal/adapters**: Implementaciones de las interfaces (HTTP, clientes externos)
- **internal/services**: Servicios de aplicación que orquestan la lógica

### Generar Documentación Swagger

Después de modificar los comentarios de Swagger en los handlers:

```bash
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

### Ejecutar Tests

```bash
go test ./...
```

## Seguridad

- Los tokens JWT están firmados con HS256
- El secret key debe ser una cadena aleatoria y segura en producción
- Los tokens tienen expiración configurable
- Las contraseñas se validan usando bcrypt en el microservicio de usuarios

## Licencia

Apache 2.0
