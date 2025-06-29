### Backend

1. **Prerequisites:**
   - [Go1.24.4](https://golang.org/dl/)
   - SQLite

2. **Build and Run**

   Open a terminal in the project root and run:

   ```bash
   CGO_ENABLED=1 go build -o webapi ./cmd/webapi
   ./webapi
   ```

   By default, the server listens on port `3000`. You can adjust settings (such as API host, database file, and timeouts) via command-line flags or by editing the configuration file (default location: `/conf/config.yml`).

## Docker

The project includes Dockerfiles for both the backend and frontend. Use the following commands to build and run container images:

### Build Container Images

- **Backend:**

  ```bash
  docker build -t wasatext-backend:latest -f Dockerfile.backend .
  ```
### Run Container Images

- **Backend:**

  ```bash
  docker run -it --rm -p 3000:3000 wasatext-backend:latest
  ```
