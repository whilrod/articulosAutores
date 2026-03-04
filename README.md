# Microservicio de Autores y Artículos

API REST para gestión de autores y artículos de un medio digital.

### Requisitos del proyecto
- [x] Modelos de dominio (Autor, Artículo, Score)
- [x] Migraciones MySQL
- [x] Repositorios (capa de datos)
- [x] Servicios (lógica de aplicación)
- [x] Handlers HTTP con Gin
- [x] Docker Compose (MySQL + API)
- [x] Documentación básica
- [x] Tests unitarios (pendiente implementación)
- [x] Tests de integración (pendiente implementación)

## Entrega

Repositorio: [https://github.com/whilrod/articulosAutores](https://github.com/whilrod/articulosAutores)

## Estructura del proyecto

```
📁 articulosAutores/
├── 📁 cmd/                    # 📦 Ejecutables (punto de entrada)
│   └── 📁 api/                
│       └── main.go            # 🚀 Inicio de la aplicación
├── 📁 internal/                # 🧩 Código privado
│   ├── 📁 domain/              # 🧠 Modelos y reglas de negocio
│   │   ├── autor.go
│   │   └── articulo.go
│   ├── 📁 application/         # ⚙️ Casos de uso (servicios)
│   └── 📁 infrastructure/       # 🔌 Adaptadores externos
│       ├── 📁 handlers/        # Capa HTTP (Gin)
│       └── 📁 repositories/    # Capa de datos (MySQL)
├── 📁 migrations/               # 🗄️ Scripts SQL versionados
├── 📁 test/                     # 🧪 Tests de integración
├── Dockerfile                    # 🐳 Imagen de la API
├── docker-compose.yml            # 🐙 Orquestación completa
└── README.md                     # 📚 Documentación
```

### Notas adicionales
- Los IDs son UUID v4
- Los artículos nuevos siempre se crean en estado "borrador"
- La fecha de publicación se asigna automáticamente al publicar
- El score se calcula dinámicamente, no se almacena en BD
- La API expone el endpoint `/health` para verificar estado
![alt text](image.png)
## Requisitos previos

- [Go](https://golang.org/dl/) 1.21 o superior
- [Docker](https://www.docker.com/products/docker-desktop/) y Docker Compose
- [Git](https://git-scm.com/)

## Ejecución del proyecto

### 1. Clonar el repositorio

- git clone https://github.com/whilrod/articulosAutores.git
- cd articulosAutores

### 2. Levantar automáticamente (API + MySQL)
- docker-compose up -d
#### Verificar que MySQL esté corriendo:
- docker ps
##### Debería aparecer "articulos-mysql" con estado "healthy"

### 3. Ejecutar migraciones

#### Instalar herramienta de migraciones (si no la tienes)
- go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

#### Ejecutar migraciones
- migrate -path migrations -database "mysql://articulos_user:articulos_pass@tcp(localhost:3306)/articulos_db" up

### 4. Ejecutar la API
- go run cmd/api/main.go
###### La API estará disponible en http://localhost:8080

#### Endpoints

##### Autores
- POST /api/v1/autores - Crear un nuevo autor
- GET /api/v1/autores/{id} - Obtener autor por ID

##### Artículos
- POST /api/v1/articulos - Crear artículo en estado BORRADOR
- POST /api/v1/articulos/{id}/publicar - Publicar un artículo
- GET /api/v1/articulos?estado=publicado&pagina=1&limite=10 - Listar artículos publicados
- GET /api/v1/autores/{id}/articulos?estado=publicado - Listar artículos por autor
- GET /api/v1/autores/{id}/resumen - Resumen del autor

## 🧪 Testing

### Tests unitarios (rápidos, sin dependencias)
```bash
# Tests de dominio (modelos, validaciones, score)
go test ./internal/domain/... -v

# Tests de aplicación (servicios)
go test ./internal/application/... -v
```

### Tests de integración (requieren MySQL)
```bash
# 1. Levantar MySQL (si no está corriendo)
docker-compose up -d mysql

# 2. Ejecutar tests de integración
go test ./test/integration/... -v

# 3. Opcional: detener MySQL
docker-compose down
```

### Tests específicos por requerimiento

#### Prueba unitaria cálculo de score
```bash
go test ./internal/domain/score_test.go -v
```

#### Prueba unitaria validación antes de publicar
```bash
go test ./internal/domain/articulo_test.go -v -run TestValidarParaPublicar
```

#### Prueba unitaria endpoint Top autores
```bash
go test ./internal/application/top_autores_service_test.go -v
```

#### Prueba de integración (publicar artículo + verificar BD)
```bash
go test ./test/integration/publicar_articulo_test.go -v
```

### Cobertura de código
```bash
go test -cover ./...
```

### Tests dentro del contenedor Docker
```bash
# Ejecutar todos los tests dentro del contenedor
docker exec articulos-api go test ./...
```

## EXTRA 
#### Comandos útiles de Docker
##### Ver logs de MySQL
- docker logs articulos-mysql
##### Entrar a MySQL
- docker exec -it articulos-mysql mysql -u articulos_user -p
##### Detener servicios
- docker-compose down
##### Detener y eliminar volúmenes (borra datos)
- docker-compose down -v

