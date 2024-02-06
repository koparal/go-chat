# Chat Application

This application allows users to chat in real-time using websockets.

Technologies : Golang, React JS, PostgreSQL, Redis

## Features

- [x] Real-time chat using websockets
- [x] User authentication (register,login) and authorization
- [x] Ability to create chat rooms and join existing ones
- [x] Ability to create/update/delete topics

## Design

The project structure is organized as follows:

```plaintext
- chat
   - backend
      - cmd
      - configs
      - db
      - docs
      - internal
         - cache
         - chat
         - config
         - router
         - topic
         - user
         - utils
      - Makefile
   - frontend
   - ...
   - README.md
```

## Getting Started

1. Clone the repository.

   ```
   git clone https://github.com/koparal/go-chat
   ```

2. Navigate to the project directory:

   ```
   cd go-chat
   ```

3. Install dependencies.

   ```
   cd frontend
   npm install
   
   cd ../backend
   go mod download
   ```

4. Set up the PostgreSQL database and Redis server.

   ```
   cd backend
   make create-db
   make create-table
   make migrate
   make create-redis (optional) if you have redis server, you can configure it from config file.
   ```

5. Configuration Settings
    ```
      nano backend/configs/dev|prod.json
      
      ## Database
      - **Host:** localhost
      - **Port:** 5433
      - **User:** root
      - **Password:** 123
      - **Table:** chat
      
      ## Router
      - **Server:** localhost
      - **Port:** 8080
      - **Mode:** dev
      
      ## Cache
      - **Host:** localhost
      - **Port:** 6379
   
    ```
6. Run seeder to create admin user

   ```
   cd backend
   go run cmd/seeder/main.go
   ```

7. Run the frontend and backend apps:

   ```
   cd frontend
   npm run dev or yarn dev
   
   cd ../backend
   go run cmd/server/main.go
   ```

8. Access the application at `http://localhost:3000` in your web browser.
9. See the api swagger documentation at `http://localhost:8080/swagger/index.html`