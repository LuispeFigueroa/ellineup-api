# ElLineup API

Backend de ElLineup — plataforma para gestionar ligas de softball. Construido con Go, Gin y PostgreSQL.

**README fue escrito con la ayuda de IA**

##  Links
- **Frontend:** https://luispefigueroa.github.io/ellineup-client
- **API en producción:** https://ellineup-api-production.up.railway.app

## Instrucciones para correr localmente

### Requisitos
- Go 1.21+
- Docker y Docker Compose

### Pasos

1. Clona el repositorio:
```bash
git clone https://github.com/LuispeFigueroa/ellineup-api
cd ellineup-api
```

2. Crea el archivo `.env` en la raíz:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=ellineup
DB_PASSWORD=ellineup123
DB_NAME=ellineup
```

3. Levanta la base de datos:
```bash
docker compose up -d
```

4. Aplica el schema:
```bash
type db\schema.sql | docker exec -i ellineup-db psql -U ellineup -d ellineup
```

5. Corre el servidor:
```bash
go run main.go
```

El servidor corre en `http://localhost:8080`.

## 📡 Endpoints

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/divisiones` | Listar divisiones (soporta `?q=` para buscar) |
| GET | `/divisiones/:id` | Obtener una división |
| POST | `/divisiones` | Crear división |
| PUT | `/divisiones/:id` | Editar división |
| DELETE | `/divisiones/:id` | Eliminar división |
| GET | `/divisiones/:id/equipos` | Listar equipos (soporta `?q=`) |
| GET | `/equipos/:id` | Obtener un equipo |
| POST | `/divisiones/:id/equipos` | Crear equipo |
| PUT | `/equipos/:id` | Editar equipo |
| DELETE | `/equipos/:id` | Eliminar equipo |
| POST | `/equipos/:id/imagen` | Subir logo del equipo |
| GET | `/equipos/:id/jugadores` | Listar jugadores |
| GET | `/jugadores/:id` | Obtener un jugador |
| POST | `/equipos/:id/jugadores` | Crear jugador |
| PUT | `/jugadores/:id` | Editar jugador |
| DELETE | `/jugadores/:id` | Eliminar jugador |
| GET | `/divisiones/:id/partidos` | Listar partidos |
| GET | `/partidos/:id` | Obtener un partido |
| POST | `/divisiones/:id/partidos` | Registrar partido |
| PUT | `/partidos/:id` | Editar partido |
| DELETE | `/partidos/:id` | Eliminar partido |
| GET | `/divisiones/:id/standings` | Tabla de posiciones calculada |

## CORS

CORS (Cross-Origin Resource Sharing) es una política de seguridad del navegador que bloquea requests entre dominios distintos. Como el frontend y el backend corren en orígenes diferentes, configuramos los siguientes headers para permitir la comunicación:

##  Challenges implementados

- **Subir imágenes** — logo de equipos, máximo 1MB 
- **Códigos HTTP correctos** — 201 al crear, 204 al eliminar, 404 si no existe, 400 en input inválido 
- **Validación server-side** — mensajes de error descriptivos en JSON 
- **Búsqueda por nombre** — `?q=` en divisiones y equipos 

## Reflexion
Decidí trabajar con GO para el API para seguir la recomendación de Erick y Dennis, aparte todos mis otros proyectos los habia trabajdo con fast API y quería probar algo diferente pero que mantuviera esa facilidad de trabajar con Fast API. Investigando vi que recomendaban usar Gin cuando se trabaja con Go porque facilitaba mucho el manejo de las rutas y de los middlewares. EN este caso, siento que la velocidad que tiene GO no era tan necesaría para mi proyecto, porque no estoy manejando datos que tienen que modificarse rápidamente, pero de cualquier forma cuando me toque hacer un proyecto donde si pueda aprovechar la velocidad, ya estaré acostumbrado a trabajar una API rest con go. En general, sentí que Go es bastante facil de usar, definiendo los handlers y el router puedo entender rapidamente como esta separado mi proyecto y al ser archivos diferentes para cada model, puedo separar todas las lógicas para trabajar de forma más ordenada y poder entender el código cuando me toque regresar a modificar algo. 

Para la base de datos decidí usar PostgresSQL, porque estoy bastante familiarizado a trabajar tablas relacionales con PostgreSQL por el curso de Base de datos que llevo con Mario Barrientos. Todas las tablas de mi proyecto iban a tener relaciones entre ellas entonces fue la unica opción que consideré. A futuro, seguramente seguiré usando PostgresSQL para manejar base de datos, especialmente si las tablas llegan a tener muchas relaciones. 