# BACKEND (STILL DEVELOPMENT)

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

   - If you want to use a local image instead of pulling from a registry, you can use the [`./docker-compose-manual.yaml`](./docker-compose-manual.yaml) file. To do this, run the following:

     ```bash
     docker-compose -f docker-compose-manual.yaml up -d
     ```

   - This will run the services using local images.

### Stopping the Services

To stop the services, run:

```bash
docker-compose down
