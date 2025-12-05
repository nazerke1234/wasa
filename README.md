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
Open a terminal in the project root and run:

   ```bash
   go run ./cmd/webapi/
   ```
### Frontend
Build dist with:

   ```bash
   yarn run build-prod
   ```
And then run with:

   ```bash
   yarn run preview
   ```
## Build and run container images:

### Build

- **Backend:**

  ```bash
  docker build -t wasatext-backend:latest -f Dockerfile.backend .
  ```
- **Frontend:**
  ```bash
  docker build -t wasatext-frontend:latest -f Dockerfile.frontend .
  ```
  
### Run 

- **Backend:**

  ```bash
  docker run -it --rm -p 3000:3000 wasatext-backend:latest
  ```
- **Frontend:**

  ```bash
  docker run -it --rm -p 8080:80 wasatext-frontend:latest
  ```
