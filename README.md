# BACKEND (STILL DEVELOPMENT)

### Key Features
1. **Role-Based Access Control (RBAC)**
Provides fine-grained access control for different user roles (e.g., Admin, User).
Ensures that only authorized users can perform specific operations on the system.
2. **Authentication and Authorization**
Utilizes JWT (JSON Web Tokens) for secure and stateless authentication.
Access tokens are used to validate and authorize user requests across services.
3. **gRPC for Service Communication**
Implements gRPC for efficient and reliable communication between microservices.
Ensures high performance with protocol buffers (protobuf) for data serialization.
Services communicate seamlessly using RPC methods for actions like creating books, authors, and user management.
4. **REST API with Swagger Documentation**
Provides a RESTful API for external interaction with the system.
The API is documented using Swagger for easy exploration and testing.
API documentation is automatically generated and accessible through a web interface.

## Setup
This project consists of four services: `auth_service`, `book_service`, `category_service`, and `author_service`. Below are the steps to set up and run the services.

### Prerequisites

- Ensure you have Docker and Docker Compose installed on your machine.

### Steps to Run the Project

1. **Set up Environment Files**

   - In the root project directory, create a `.env` file by copying the contents of [`.env.example`](./.env.example) and fill it up:

     ```bash
     cp .env.example .env
     ```

   - In each service directory (`auth_service`, `book_service`, `category_service`, `author_service`), create a `.env` file by copying the contents of the corresponding `.env.example` file and fill it up:

     ```bash
     cp <service_name>/.env.example <service_name>/.env
     ```

     Replace `service_name` with the name of each service (e.g., `auth_service`, `book_service`, etc.).

2. **Run Docker Compose**

   - In the root project directory, run the following command to start the services in detached mode:

     ```bash
     docker-compose up -d
     ```

   - This will start all the services defined in the [`./docker-compose.yaml`](./docker-compose.yaml) file.

3. **Use Local Image (Optional)**

   - If you want to use a local image instead of pulling from a registry, you can leave empty the `*_SERVICE_IMAGE` in the `.env` file (refer to [`./env.example`](./env.example)) and run docker compose with `--build` flag:

     ```bash
     docker-compose up -d --build
     ```

   - This will run the services using local images.

## Swagger API Documentation
### Available Swagger Endpoints
- auth_service:
  `http://{host}:8001/swagger/index.html`

- book_service:
  `http://{host}:8003/swagger/index.html`

- category_service: Not implemented yet

- author_service:
  `http://{host}:8002/swagger/index.html`

## gRPC Ports
- auth_service:
 `{host}:7001`

- auth_service:
 `{host}:7002`

- category_service:
  Not implemented yet

- author_service:
  Not implemented yet
