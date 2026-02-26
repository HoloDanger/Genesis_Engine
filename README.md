# Genesis Engine
### Production-Ready Scaffolding CLI

**Genesis** is a CLI tool designed to instantly scaffold high-performance, full-stack applications. It eliminates the hours spent on "setup fatigue"—configuring databases, authentication, and folder structures—allowing you to focus immediately on building features.

Ideal for **Practicums**, **Hackathons**, and **Rapid Prototyping**.

---

## ⚡ The Tech Stack (Modern & Standardized)

Genesis generates projects using an opinionated, industry-standard stack:

- **Frontend:** Next.js 16 (React 19) + Tailwind CSS v4
- **Backend:** Go 1.25+ (High-performance API)
- **Database:** PostgreSQL 16 (via Docker)
- **ORM:** Drizzle ORM (Type-safe database interaction)
- **Authentication:** Better Auth (Pre-configured secure login)
- **Runtime:** Bun (Fast JavaScript runtime & package manager)

---

## 🛠 Prerequisites

Before using Genesis, ensure you have the following installed:

1.  **Go** (v1.23+) - [Download](https://go.dev/dl/)
2.  **Bun** (v1.0+) - [Install](https://bun.sh/)
3.  **Docker** - [Get Docker Desktop](https://www.docker.com/products/docker-personal/)

---

## 🚀 Installation

1.  **Clone the Repository**
    ```bash
    git clone git@github.com:HoloDanger/Genesis_Engine.git
    cd Genesis_Engine
    ```

2.  **Build the CLI Tool**
    ```bash
    go build -o genesis main.go
    ```

3.  **Move to Path (Optional)**
    *Move the binary to your global path to use it from anywhere.*
    ```bash
    mv genesis /usr/local/bin/
    # OR just move it to your projects folder
    mv genesis ~/MyProjects/
    ```

---

## 📖 Usage Guide

Genesis supports three distinct architectural patterns depending on your project needs.

### 1. The Full Stack ("Hybrid")
**Best for:** Capstone Projects, ERPs, SaaS Apps.
Combines a Next.js frontend and a Go backend, both sharing a single Dockerized database.

```bash
# Basic Setup
./genesis -name MyProject -type hybrid

# With OpenAI Integration (Optional)
./genesis -name MyProject -type hybrid -ai=true
```

### 2. Frontend Only ("T3")
**Best for:** Dashboards, simple Web Apps, Demos.
A standalone Next.js application with Database and Auth built-in.

```bash
./genesis -name MyFrontend -type t3
```

### 3. Backend Only ("Go")
**Best for:** High-performance APIs, Microservices, Mobile App Backends.
A pure Go REST API service.

```bash
./genesis -name MyBackend -type go
```

---

## 🔥 Getting Started (After Generation)

Once you have generated your project (e.g., `MyProject`), follow these steps to launch:

1.  **Enter the Directory**
    ```bash
    cd MyProject
    ```

2.  **Start Infrastructure (Database)**
    ```bash
    docker compose up -d
    ```

3.  **Sync Database Schema**
    *This creates your tables in the local database.*
    ```bash
    # For Hybrid/T3 projects:
    cd web && bun db:push
    ```

4.  **Run the App**
    * Open Terminal 1 (Frontend): `cd web && bun dev`
    * Open Terminal 2 (Backend): `cd api && make run`

---

## 🧩 Project Structure

Your generated project will look like this:

```
MyProject/
├── compose.yml       # Docker Database Configuration
├── Makefile          # Quick commands
├── web/              # Next.js Frontend
│   ├── src/app       # Pages & Layouts
│   ├── src/server    # Database Schema
│   └── src/lib       # Auth Configuration
└── api/              # Go Backend
    ├── cmd/          # Entry point
    └── internal/     # Business Logic
```

---

## 💡 Philosophy

- **Zero Config:** Decisions are already made. Don't waste time choosing a linter.
- **Local First:** Designed to run fully offline with local Docker containers.
- **Clean Code:** Generates minimal boilerplate code that is easy to read and extend.

---
*Built with logic and precision.*
