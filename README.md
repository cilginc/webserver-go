## Environment Variables

The server can be configured using the following environment variables:

- **`PORT`**  
  The port on which the server listens for incoming HTTP requests.  
  **Default:** `:8080`  
  **Example:**
  ```bash
  PORT=":3000"
  ```

* **`STATIC_DIR`**
  The absolute path to the directory from which static files are served.
  **Default:** Current working directory (`.`)
  **Example:**

  ```bash
  STATIC_DIR="/app/web"
  ```

* **`WORKER_COUNT`**
  The number of worker goroutines to handle requests.
  **Default:** Number of available CPU cores
  **Example:**

  ```bash
  WORKER_COUNT=16
  ```

---

## Running Locally

1. **Clone the repository:**

   ```bash
   git clone https://github.com/cilginc/webserver-go.git
   cd webserver-go
   ```

2. **Build the server:**

   ```bash
   go build -o bin ./cmd/server
   cd bin
   ```

3. **Start the server:**

   ```bash
   # Using default configuration
   ./server

   # With custom port and static directory
   PORT=":3000" STATIC_DIR="./web/www" ./server
   ```

---

## Running with Docker Compose

Add the following service definition to your `docker-compose.yml` file:

```yaml
services:
  webserver:
    image: ghcr.io/cilginc/server:latest
    container_name: webserver
    ports:
      - "3000:3000"
    volumes:
      - ./www:/www
    environment:
      - PORT=:3000
      - STATIC_DIR=/www
```

Then start the container:

```bash
docker compose up -d
```
