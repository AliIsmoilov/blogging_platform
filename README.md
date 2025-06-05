# Blogging Platform Documentation

## 📘 Overview

This project is a blogging platform developed in Go (Golang), designed to provide a scalable 
solution for creating and managing blogs.

---

## 🧱 System Architecture

### High-Level Architecture

The platform follows a client-server architecture, where:

- **Frontend**: A web interface for users to interact with the platform.
- **Backend**: A Go-based server handling business logic and API requests.
- **Database**: Stores user data, user blog posts.

### Components

- **API Server**: Handles HTTP requests and serves the frontend.
- **Database**: Manages data persistence.

---

## ⚙️ Technologies Used

- **Programming Language**: Go (Golang)
- **Database**: [PostgreSQL, MongoDB, Neo4j]
- **Web Framework**: [ Gin ]
- **Authentication**: [Not developed yet]

---

## 🛠️ Installation and Usage

### Prerequisites

- Go 1.18 or higher
- MongoDB
- Neo4j

### Installation Steps

#### 1. Clone the repository:

```bash
    git clone https://github.com/AliIsmoilov/blogging_platform.git
```

--

#### 2. Create the `.env` File
Create a `.env` file in the root directory of the project by copying .env.example

--

####  3. Install dependencies:

```bash
    go mod tidy
```

--

#### 4. Configure the Databases
To connect the application to the required databases, you'll need to create a `.env` file in the root directory of the project. This file will store your database connection settings.

P.S Use .env.example file

The application supports **three types of databases**:

- **PostgreSQL** – for relational data
- **MongoDB** – for document-oriented data
- **Neo4j** – for graph-based data

--

#### 5. Run the Application
```bash
    go run .\cmd\main.go
```