# Storage API

Esta es una API construida con Golang que interactúa con el almacenamiento de Cloudflare R2. La API permite gestionar archivos mediante endpoints REST y se puede ejecutar fácilmente dentro de un contenedor Docker.

## Requisitos

Antes de empezar, asegúrate de tener instalado lo siguiente:

- Docker
- Docker Compose (opcional pero recomendado)

## Configuración de Docker

Para crear y ejecutar la API dentro de un contenedor Docker, utiliza `docker-compose` con la configuración incluida o ejecuta directamente la imagen desde Docker Hub.

### Comando de CLI

Para ejecutar la API directamente desde Docker Hub, usa el siguiente comando:

```bash
docker run --name storage-api -p 4003:4003 -p 4004:4004 --env-file .env tutitoos/storage-api
```

Este comando hará lo siguiente:
- Descargará la imagen `tutitoos/storage-api` desde Docker Hub.
- Iniciará un contenedor con los puertos 4003 y 4004 mapeados a tu máquina local.

### Docker Compose

Si prefieres usar `docker-compose`, asegúrate de tener la siguiente configuración en tu archivo `docker-compose.yml`:

```yaml
version: '3.8'

services:
  storage-api:
    image: tutitoos/storage-api
    container_name: storage-api
    ports:
      - "4003:4003"
      - "4004:4004"
    env_file:
      - .env
    deploy:
      resources:
        limits:
          cpus: "0.5"
          memory: 500m
```

Luego, ejecuta:

```bash
docker-compose up --build
```

### Variables de entorno

En el archivo `.env`, debes configurar las siguientes variables para que la API funcione correctamente:

```env
PORT="4003"                        # Puerto en el que la API escucha.
TRY_PORT="4004"                    # Puerto alternativo si PORT no está disponible.
API_URL="http://localhost:{PORT}"  # URL base de la API con placeholder para el puerto.
TOKEN=""                           # Token de autorización para acceder a la API.

# Cloudflare
CLOUDFLARE_ACCOUNT_ID=""           # ID de la cuenta de Cloudflare.
CLOUDFLARE_ACCESS_KEY_ID=""        # ID de la clave de acceso de Cloudflare.
CLOUDFLARE_ACCESS_KEY_SECRET=""    # Clave secreta de acceso de Cloudflare.

# Cloudflare Storage
BUCKET_NAME=""                     # Nombre del bucket en Cloudflare.
BUCKET_REGION="eeur"               # Región del bucket en Cloudflare. Ejemplo:  "eeur"

# Rules Explorer
EXCLUDE_FOLDER=""                  # Carpetas que no deben ser incluidas en la respuesta de la API. Ejemplo: "my-folder,my-other-folder"
EXCLUDE_FILE=""                    # Archivos que no deben ser incluidos en la respuesta de la API. Ejemplo: "my-file.txt,my-other-file.txt"

# Rules Whitelist IP
WHITELIST_IPS="127.0.0.1,::1"      # IPs que no deben ser incluidas en la respuesta de la API. Ejemplo: "127.0.0.1,::1"
```

### Verificar la API

Una vez que el contenedor esté en funcionamiento, la API estará disponible en `http://localhost:4003` (o en `http://localhost:4004` si el puerto principal no está disponible).

Puedes verificar que la API está funcionando haciendo una solicitud GET a la raíz:

```bash
curl http://localhost:4003/v1
```

## Endpoints de la API

La API expone varios endpoints para interactuar con el almacenamiento de Cloudflare R2.

### 1. `GET /v1`

Este endpoint devuelve una respuesta básica para verificar que el servicio está activo.

**Ejemplo:**

```bash
curl http://localhost:4003/v1
```

### 2. `GET /v1/files/*`

Devuelve una lista de archivos en una carpeta específica dentro del almacenamiento. Si no se proporciona la carpeta, devuelve los archivos en la raíz.

**Ejemplo:**

```bash
curl http://localhost:4003/v1/files/my-folder
```

### 3. `GET /v1/file/*`

Devuelve el contenido de un archivo específico.

**Ejemplo:**

```bash
curl http://localhost:4003/v1/file/my-folder/file.txt
```

### 4. `DELETE /v1/file/*`

Elimina un archivo específico.

**Ejemplo:**

```bash
curl -X DELETE http://localhost:4003/v1/file/my-folder/file.txt
```

### 5. `POST /v1/file`

Sube uno o más archivos al almacenamiento. Los archivos deben ser enviados como parte de una solicitud `form-data`.

**Cuerpo de la solicitud:**
- `files`: Los archivos a subir (clave del `form-data`).

**Ejemplo:**

```bash
curl -X POST http://localhost:4003/v1/file \
  -F "files=@/path/to/local/file1.txt" \
  -F "files=@/path/to/local/file2.jpg"
```