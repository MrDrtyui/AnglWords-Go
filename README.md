# Anglwords Project

## Project Description

Anglwords is a comprehensive assistant designed to help users learn new words. It features a robust Go backend for word management and authentication, with future plans for a Telegram bot and a web interface. The core functionality revolves around allowing users to store words, retrieve their translations, and get AI-generated example sentences.

## Features

### Backend
*   **User Authentication**: JWT access tokens and secure refresh tokens.
*   **Word Management**: Users can create, view their own words, view public words, and get translations with difficulty levels.
*   **AI Integration**: Uses a neuro-service (Gemini) for word translation and difficulty assessment.
*   **Database**: PostgreSQL with Ent ORM for clean and efficient data access.
*   **API**: RESTful API built with Chi router.
*   **Configuration**: Flexible configuration management using YAML.
*   **Docker Support**: Easy local development setup with Docker Compose.

### Planned Features
*   **Telegram Bot**: Interact with the word learning system directly via Telegram.
*   **Web Interface**: A user-friendly web application for managing and learning words.

## Project Structure

```
.
├── backend/                  # Go Backend application
│   ├── cmd/api/              # Application entry point
│   ├── domain/               # Domain models and DTOs
│   ├── internal/             # Internal packages (auth, user, word, router, config, db, neuronet)
│   ├── ent/                  # Ent ORM generated code and schemas
│   └── config/               # Configuration files
├── docker/                   # Docker Compose setup for services (e.g., PostgreSQL)
├── GEMINI.md                 # Project overview and requirements (in Russian)
└── README.md                 # This README file
```

## Technologies Used

*   **Go**: Programming Language (1.20+)
*   **Chi**: HTTP Router (v5)
*   **Ent ORM**: Entity Framework for Go (v0.14.5)
*   **PostgreSQL**: Relational Database (16)
*   **JWT**: JSON Web Tokens for authentication
*   **Bcrypt**: Password hashing
*   **Docker**: Containerization
*   **Gemini API**: For AI-powered translations and example sentences
*   **Swagger**: API Documentation

## Getting Started

Follow these steps to set up and run the Anglwords backend locally.

### Prerequisites

*   [Go](https://golang.org/doc/install) (version 1.20 or higher)
*   [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)

### 1. Start PostgreSQL

Navigate to the `docker` directory and start the PostgreSQL container:

```bash
cd docker
docker compose up -d
```
This will start a PostgreSQL 16 instance on `localhost:5432` using credentials defined in `.env.postgres`.

### 2. Configure the Application

Create a configuration file `config/dev.yaml` in the project root (if it doesn't already exist) and populate it with your settings. Make sure to replace placeholders like `your-secret-key` and `your-gemini-api-key`.

```yaml
env: dev
http:
  port: ":9000"
database:
  url: "postgres://user:pass@localhost:5432/postgres?sslmode=disable"
jwt:
  secret: "my-super-secret-jwt-key-change-in-production"
  accessTtlHours: 1h
  refreshTtlHours: 720h
gemini:
  api-key: "YOUR_GEMINI_API_KEY" # Replace with your actual Gemini API Key
```

### 3. Generate Ent ORM Code & Swagger Docs

Navigate to the `backend` directory and run the Go generate commands:

```bash
cd backend
go generate ./ent
go generate ./...
```
This will generate the necessary Ent ORM code based on the schemas and update the Swagger API documentation.

### 4. Run the Application

From the `backend` directory, run the application, ensuring you provide the `CONFIG_PATH` environment variable:

```bash
CONFIG_PATH=../config/dev.yaml go run cmd/api/main.go
```
The server will start on `http://localhost:9000` (or the port specified in your `dev.yaml`).

### API Documentation (Swagger)
Access the API documentation at `http://localhost:9000/swagger/index.html`.

## API Endpoints

All endpoints respond with JSON. Error responses are typically in the format: `{"error": "error message"}`.

### Authentication Endpoints

These endpoints are used for user registration, login, and token management.

#### `POST /auth/register`
*   **Description**: Creates a new user account.
*   **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "securepassword",
      "username": "myusername"
    }
    ```
*   **Success Response (201 Created)**: Returns user details and authentication tokens.
    ```json
    {
      "user": {
        "id": 1,
        "email": "user@example.com",
        "username": "myusername"
      },
      "access_token": "eyJ...",
      "refresh_token": "Mz8F..."
    }
    ```

#### `POST /auth/login`
*   **Description**: Authenticates an existing user.
*   **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "securepassword"
    }
    ```
*   **Success Response (200 OK)**: Returns user details and authentication tokens.
    ```json
    {
      "user": {
        "id": 1,
        "email": "user@example.com",
        "username": "myusername"
      },
      "access_token": "eyJ...",
      "refresh_token": "Mz8F..."
    }
    ```

#### `POST /auth/refresh`
*   **Description**: Rotates the refresh token and generates new access and refresh tokens.
*   **Request Body**:
    ```json
    {
      "refresh_token": "Mz8F..."
    }
    ```
*   **Success Response (200 OK)**: Returns new authentication tokens.
    ```json
    {
      "user": {
        "id": 1,
        "email": "user@example.com",
        "username": "myusername"
      },
      "access_token": "eyJ...",
      "refresh_token": "oTQm..."
    }
    ```

#### `POST /auth/logout`
*   **Description**: Revokes the refresh token, effectively logging out the user.
*   **Request Body**:
    ```json
    {
      "refresh_token": "oTQm..."
    }
    ```
*   **Success Response (204 No Content)**

### Word Management Endpoints

These endpoints are for managing words, some require authentication.

#### `POST /words`
*   **Description**: Creates a new word for the authenticated user.
*   **Authentication**: Required (Bearer Token)
*   **Request Body**:
    ```json
    {
      "word": "hello"
    }
    ```
*   **Success Response (201 Created)**: Returns the newly created word details.
    ```json
    {
      "id": 1,
      "word": "hello",
      "ruWord": "привет",
      "level": "A1",
      "createdAt": "2025-01-01T12:00:00Z"
    }
    ```

#### `GET /words/my`
*   **Description**: Retrieves all words associated with the authenticated user.
*   **Authentication**: Required (Bearer Token)
*   **Success Response (200 OK)**: Returns a list of words.
    ```json
    [
      {
        "id": 1,
        "word": "hello",
        "ruWord": "привет",
        "level": "A1",
        "createdAt": "2025-01-01T12:00:00Z"
      }
    ]
    ```

#### `GET /words/all`
*   **Description**: Retrieves all words in the database (publicly accessible).
*   **Authentication**: None
*   **Success Response (200 OK)**: Returns a list of all words.
    ```json
    [
      {
        "id": 1,
        "word": "hello",
        "ruWord": "привет",
        "level": "A1",
        "createdAt": "2025-01-01T12:00:00Z"
      },
      {
        "id": 2,
        "word": "world",
        "ruWord": "мир",
        "level": "A1",
        "createdAt": "2025-01-01T12:05:00Z"
      }
    ]
    ```

#### `GET /word/{word}`
*   **Description**: Retrieves a single word by its string representation (publicly accessible).
*   **Authentication**: None
*   **Path Parameters**: `word` (string) - The word to retrieve.
*   **Success Response (200 OK)**: Returns the word details.
    ```json
    {
      "id": 1,
      "word": "hello",
      "ruWord": "привет",
      "level": "A1",
      "createdAt": "2025-01-01T12:00:00Z"
    }
    ```

## Development

### Generate Ent ORM Code

After modifying any schema files in `backend/ent/schema/`:

```bash
cd backend
go generate ./ent
```

### Build the Application

From the `backend` directory:

```bash
cd backend
make build
# The binary will be created at backend/bin/app
```

### Run Tests

From the `backend` directory:

```bash
cd backend
make test
```

### Lint Code

From the `backend` directory (requires `golangci-lint`):

```bash
cd backend
make lint
```

## License

This project is a template and is open for use and modification as per your needs.

## Contributing

Feel free to fork this repository and customize it for your own projects.