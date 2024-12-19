# A Simple Go Server

# Overview

A simple Go web server using Gorilla Mux for managing products.

## Features

- **Endpoints**:
  - `GET /`: Retrieve all products.
  - `PUT /{id}`: Update a product by UUID.
  - `POST /`: Add a new product.

- **Middleware**:
  - Validates JSON payloads for `PUT` and `POST`.
  - Ensures `Content-Type: application/json` in responses.
    
- **Graceful Shutdown**:
  - Handles OS signals (`Interrupt`, `SIGTERM`) for cleanup.
  - Allows 30 seconds for shutdown operations.

## How to Run
1. Clone the repository.
2. Run the server: `go run main.go`.
3. Access endpoints at `http://localhost:9090`.
