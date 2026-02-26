package hybrid

// 1. THE ROOT COMPOSE (The Source of Truth)
const RootCompose = `services:
  # --- THE TRUTH (Database) ---
  postgres:
    image: postgres:16-alpine
    container_name: {{.ProjectName}}-db
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: {{.ProjectName}}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - genesis_net

  # --- THE SHIELD (Frontend) ---
  web:
    build:
      context: ./web
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      DATABASE_URL: postgres://postgres:password@postgres:5432/{{.ProjectName}}
      BETTER_AUTH_URL: http://localhost:3000
    depends_on:
      - postgres
    networks:
      - genesis_net

  # --- THE SPEAR (Backend) ---
  api:
    build:
      context: ./api
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://postgres:password@postgres:5432/{{.ProjectName}}?sslmode=disable
      PORT: 8080
    depends_on:
      - postgres
    networks:
      - genesis_net

volumes:
  postgres_data:

networks:
  genesis_net:
    driver: bridge
`

// 2. THE ROOT MAKEFILE (Orchestrator)
const RootMakefile = `PROJECT := {{.ProjectName}}

.PHONY: all dev up down clean

all: dev

dev:
	@echo "⚔️  [HYBRID] Launching Development Nodes..."
	@echo "   1. Postgres (Background)"
	@docker compose up -d postgres
	@echo "   2. Wait for DB..."
	@sleep 3
	@echo "   3. Launching Terminals (Manual)"
	@echo "      - Web: cd web && bun dev"
	@echo "      - API: cd api && make run"

up:
	@echo "🚀 [PRODUCTION] Spinning up Containers..."
	@docker compose up --build -d

down:
	@echo "🛑 [HALT] Stopping Nodes..."
	@docker compose down

db-push:
	@echo "💾 [SCHEMA] Pushing T3 Schema to DB..."
	@cd web && bun run db:push
`
