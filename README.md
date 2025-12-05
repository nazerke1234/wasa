# Wasa Text Messenger 

Connect with your friends effortlessly using WASAText! Send and receive messages, whether one-on-one
or in groups, all from the convenience of your PC. Enjoy seamless conversations with text or GIFs and
easily stay in touch through your private chats or group discussions.

## Functionality

You can start private one-on-one chats or create group conversations, share images and GIFs, react to messages with emojis like ❤️, and forward or reply to messages with context. Additionally, you can update your profile by changing your username and profile picture, and easily find other users by searching for their usernames. 

## Tools
 **Backend**: The backend is built with Go, handling the core business logic and serving the API. It uses SQLite as an embedded database, httprouter for lightweight routing, logrus for structured logging, and uuid to generate unique identifiers.
 **Frontend**: The frontend is developed with Vue.js, providing a reactive user interface. Vue Router enables smooth single-page application navigation, while Axios handles API requests. The design is made responsive using Bootstrap alongside custom CSS.
## How to setup and run the application:

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
