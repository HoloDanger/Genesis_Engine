package goservice

// 1. IDENTITY (go.mod)
const GoMod = `module {{.Name}}

go 1.26.0

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/go-playground/validator/v10 v10.30.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/joho/godotenv v1.5.1
{{if .WithAI}}	github.com/sashabaranov/go-openai v1.41.2{{end}}
)
`

// 2. ENTRYPOINT (cmd/api/main.go)
const MainGo = `package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"{{.Name}}/internal/config"
	"{{.Name}}/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// 1. FAIL-FAST CONFIGURATION
	cfg := config.Load()

	// 2. SERVER INITIALIZATION
	app := server.NewServer(cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      app.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second, 
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("[SYSTEM] {{.Name}} Online on :%s\n", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("[FATAL] Server failed: %v", err)
	}
}
`

// 3. CONFIG STRUCT (internal/config/config.go)
const ConfigGo = `package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port         string ` + "`" + `env:"PORT" envDefault:"8080"` + "`" + `
	DatabaseURL  string ` + "`" + `env:"DATABASE_URL,required"` + "`" + `
	AdminSecret  string ` + "`" + `env:"ADMIN_SECRET,required"` + "`" + `
	{{if .WithAI}}OpenAIApiKey string ` + "`" + `env:"OPENAI_API_KEY,required"` + "`" + `
	OpenAIModel  string ` + "`" + `env:"OPENAI_MODEL" envDefault:"gpt-4o-mini"` + "`" + `{{end}}
}

func Load() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("[FATAL] Configuration Failure: %v", err)
	}
	return &cfg
}
`

// 4. SERVER STRUCT (internal/server/server.go)
const ServerGo = `package server

import (
	"database/sql"
	"fmt"
	"log"

	"{{.Name}}/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct {
	db         *sql.DB
	config     *config.Config
	auditQueue chan AuditEntry
}

func NewServer(cfg *config.Config) *Server {
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	fmt.Println("🔌 [DATABASE] Connected to Postgres")

	s := &Server{
		db:         db,
		config:     cfg,
		auditQueue: make(chan AuditEntry, 100),
	}

	// Launch the Deterministic Ledger Engine
	go s.StartAuditWorker()

	return s
}
`

// 5. ROUTER (internal/server/routes.go)
const RoutesGo = `package server

import (
	"encoding/json"
	"net/http"
{{if .WithAI}}
	"{{.Name}}/internal/ai"
{{end}}
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// 1. PUBLIC ROUTES (Health & Baseline)
	mux.HandleFunc("GET /health", s.healthHandler)

{{if .WithAI}}
	// 2. AI ORCHESTRATION (The Sword)
	aiService := ai.NewService(s.db, s.config.OpenAIApiKey, s.config.OpenAIModel)
	aiHandler := ai.NewHandler(aiService)
{{end}}

	// 3. PROTECTED ROUTES (The Spear)
	protected := http.NewServeMux()
	protected.HandleFunc("GET /me", s.meHandler)

{{if .WithAI}}
	// AI Features restricted to Admin by default in boilerplate
	protected.Handle("POST /ai/generate", s.RBACMiddleware("ADMIN")(http.HandlerFunc(aiHandler.HandleGenerate)))
	protected.Handle("POST /ai/chat", s.RBACMiddleware("ADMIN")(http.HandlerFunc(aiHandler.HandleChat)))
{{end}}

	// Mount protected routes under /api/ with Global Auth Guard
	mux.Handle("/api/", http.StripPrefix("/api", s.AuthMiddleware(protected)))

	return s.StandardStack(mux)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "active",
		"system": "{{.Name}}",
	})
}

func (s *Server) meHandler(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value(UserIDKey)
	userID := "unknown"
	if val != nil {
		userID = val.(string)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "authenticated",
		"userID": userID,
	})
}
`

// 6. MIDDLEWARE (internal/server/middleware.go)
const MiddlewareGo = `package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"
)

type UserContextKey string

const (
	UserIDKey   UserContextKey = "userID"
	UserRoleKey UserContextKey = "userRole"
)

// --- AUTH MIDDLEWARE ---

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// 1. Extraction: Cookie > Header
		cookie, err := r.Cookie("better-auth.session_token")
		if err == nil {
			token = cookie.Value
		} else {
			authHeader := r.Header.Get("Authorization")
			if t, found := strings.CutPrefix(authHeader, "Bearer "); found {
				token = t
			}
		}

		if token == "" {
			http.Error(w, "Unauthorized: No Signal", http.StatusUnauthorized)
			return
		}

		// 2. Formatting: Handle Better-Auth token.signature
		if strings.Contains(token, ".") {
			if before, _, found := strings.Cut(token, "."); found {
				token = before
			}
		}

		// 3. Verification: Query the Iron Link (The Database)
		// Assuming 'user' table has a 'role' column (ADMIN, CLERK, SYSTEM)
		query := ` + "`" + `
			SELECT s.user_id, u.role
			FROM public.session s
			JOIN public.user u ON s.user_id = u.id
			WHERE s.token = $1 
			AND s.expires_at > NOW()
		` + "`" + `

		var userID, role string
		err = s.db.QueryRow(query, token).Scan(&userID, &role)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized: Invalid Token", http.StatusUnauthorized)
				return
			}
			log.Printf("❌ [AUTH] DB Error: %v", err)
			http.Error(w, "System Failure", http.StatusInternalServerError)
			return
		}

		// 4. Identification: Inject into Context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// --- RBAC MIDDLEWARE ---

func (s *Server) RBACMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(UserRoleKey).(string)

			authorized := false
			for _, r := range allowedRoles {
				if r == role {
					authorized = true
					break
				}
			}

			if !authorized {
				log.Printf("⚠️ [RBAC] Access Denied for role: %s on path: %s", role, r.URL.Path)
				http.Error(w, "Forbidden: Insufficient Authority", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// --- STANDARD STACK ---

func (s *Server) StandardStack(next http.Handler) http.Handler {
	return s.Logger(s.CORS(s.Recoverer(next)))
}

func (s *Server) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// Hardened to localhost:3000 by default for hybrid local dev
		if origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
`

// 7. AUDIT LOGIC (internal/server/audit.go)
const AuditGo = `package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type AuditEntry struct {
	UserID   string
	Action   string
	EntityID string
	Payload  any
}

func (s *Server) StartAuditWorker() {
	fmt.Println("📝 [AUDIT] Serialized Worker Online")
	for entry := range s.auditQueue {
		s.processAudit(entry)
	}
}

func (s *Server) LogAudit(ctx context.Context, action string, entityID string, payload any) {
	userID, _ := ctx.Value(UserIDKey).(string)
	if userID == "" {
		userID = "SYSTEM"
	}

	s.auditQueue <- AuditEntry{
		UserID:   userID,
		Action:   action,
		EntityID: entityID,
		Payload:  payload,
	}
}

func (s *Server) processAudit(entry AuditEntry) {
	payloadBytes, _ := json.Marshal(entry.Payload)

	query := ` + "`" + `
		INSERT INTO audit_logs (id, user_id, action, entity_id, payload, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	` + "`" + `

	id := uuid.New().String()
	_, err := s.db.Exec(query, id, entry.UserID, entry.Action, entry.EntityID, string(payloadBytes))
	
	if err != nil {
		log.Printf("🔥 [AUDIT FAILURE] %v", err)
	} else {
		fmt.Printf("📝 [AUDIT] %s | %s\n", entry.Action, entry.UserID)
	}
}
`

// 8. SQLC (sqlc.yaml)
const SQLCConfig = `version: "2"
sql:
  - schema: "internal/db/schema.sql"
    queries: "internal/db/query.sql"
    engine: "postgresql"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
`

// 9. DATABASE SKELETON (internal/db/schema.sql)
const SchemaSQL = `CREATE TABLE IF NOT EXISTS public.user (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    role TEXT NOT NULL DEFAULT 'CLERK', -- ADMIN, CLERK, SYSTEM
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.session (
    token TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES public.user(id),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.audit_logs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES public.user(id),
    action TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    payload TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
`

// 10. QUERY SKELETON (internal/db/query.sql)
const QuerySQL = `-- name: GetUserBySession :one
SELECT u.id, u.email, u.role
FROM public.session s
JOIN public.user u ON s.user_id = u.id
WHERE s.token = $1 
AND s.expires_at > NOW();

-- name: CreateAuditLog :exec
INSERT INTO public.audit_logs (id, user_id, action, entity_id, payload, created_at)
VALUES ($1, $2, $3, $4, $5, NOW());
`

// 11. AI SERVICE (internal/ai/service.go)
const AIServiceGo = `package ai

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Service struct {
	db     *sql.DB
	client *openai.Client
	model  string
}

func NewService(db *sql.DB, apiKey string, model string) *Service {
	if model == "" {
		model = openai.GPT4oMini // Default fallback
	}
	client := openai.NewClient(apiKey)
	return &Service{
		db:     db,
		client: client,
		model:  model,
	}
}

func (s *Service) GenerateDescription(ctx context.Context, specs string) (string, error) {
	prompt := fmt.Sprintf("You are a technical copywriter. Convert these specs into a professional 2-sentence description.\n\nSPECS: %s", specs)

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: s.model,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)

	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (s *Service) ChatWithInventory(ctx context.Context, query string) (string, error) {
	// 1. RAG Context (Simple Dump Strategy)
	// In production, use Vector Search. For now, we dump the last 50 items.
	rows, err := s.db.QueryContext(ctx, "SELECT name, price, description FROM product LIMIT 50")
	var inventoryJson []byte
	
	if err == nil {
		defer rows.Close()
		var inventory []map[string]any
		for rows.Next() {
			var name, desc string
			var price int
			if err := rows.Scan(&name, &price, &desc); err == nil {
				inventory = append(inventory, map[string]any{
					"name": name, "price": price, "desc": desc,
				})
			}
		}
		inventoryJson, _ = json.Marshal(inventory)
	} else {
		inventoryJson = []byte("[]") // Empty if no DB or table yet
	}

	// 2. Inference
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: s.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: fmt.Sprintf("You are a Strategic AI. DATA: %s", string(inventoryJson)),
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: query,
				},
			},
			// Give reasoning models room to breathe
			MaxCompletionTokens: 5000, 
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
`

// 12. AI HANDLER (internal/ai/handler.go)
const AIHandlerGo = `package ai

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

type GenerateRequest struct {
	Specs string ` + "`" + `json:"specs"` + "`" + `
}

type ChatRequest struct {
	Query string ` + "`" + `json:"query"` + "`" + `
}

func (h *Handler) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}

	desc, err := h.service.GenerateDescription(r.Context(), req.Specs)
	if err != nil {
		http.Error(w, "AI Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"description": desc})
}

func (h *Handler) HandleChat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}

	response, err := h.service.ChatWithInventory(r.Context(), req.Query)
	if err != nil {
		http.Error(w, "AI Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"response": response})
}
`

// 13. CONFIG (.env)
const EnvFile = `PORT=8080
DATABASE_URL=postgres://postgres:password@localhost:5432/{{.Name}}?sslmode=disable
ADMIN_SECRET=supersecret
{{if .WithAI}}OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o-mini{{end}}
`

// 14. IGNORE (.gitignore)
const GitIgnore = `bin/
main
.DS_Store
.env
`

// 15. MAKEFILE
const Makefile = `PROJECT_NAME := {{.Name}}

.PHONY: all build run test lint clean docker sqlc

all: build

build:
	@echo "⚙️  Building $(PROJECT_NAME)..."
	@go build -o bin/api cmd/api/main.go

run:
	@echo "🚀 Launching $(PROJECT_NAME)..."
	@go run cmd/api/main.go

tidy:
	@go mod tidy

test:
	@go test -v ./...

lint:
	@echo "🔍 Linting $(PROJECT_NAME)..."
	@golangci-lint run ./...

sqlc:
	@echo "🧬 Generating SQL Spears..."
	@sqlc generate

clean:
	@rm -rf bin/

docker:
	@docker compose up -d
`

// 16. LINTING (.golangci.yml)
const LintConfig = `run:
  timeout: 5m
  tests: true

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gosec
    - revive
    - misspell

linters-settings:
  govet:
    check-shadowing: true
  revive:
    rules:
      - name: exported
        arguments:
          - "disable"

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
`

// 17. DOCKER (compose.yml)
const DockerCompose = `services:
  postgres:
    image: postgres:16-alpine
    container_name: {{.Name}}-db
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: {{.Name}}
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
`
