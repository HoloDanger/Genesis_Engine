package goservice

// 1. IDENTITY (go.mod)
const GoMod = `module {{.Name}}

go 1.25.5

require (
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
{{if .WithAI}}	github.com/sashabaranov/go-openai v1.32.0{{end}}
)
`

// 2. ENTRYPOINT (cmd/api/main.go)
const MainGo = `package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"{{.Name}}/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := server.NewServer()

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      app.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second, // Extended for AI latency
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("[SYSTEM] {{.Name}} Online on :%s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("[FATAL] Server failed: %v", err)
	}
}
`

// 3. SERVER STRUCT (internal/server/server.go)
const ServerGo = `package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Server struct {
	db *sql.DB
}

func NewServer() *Server {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is required in .env")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	fmt.Println("🔌 [DATABASE] Connected to Postgres")

	return &Server{
		db: db,
	}
}
`

// 4. ROUTER (internal/server/routes.go)
const RoutesGo = `package server

import (
	"encoding/json"
	"net/http"
{{if .WithAI}}
	"fmt"
	"os"
	"{{.Name}}/internal/ai"
{{end}}
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// 1. Public Routes
	mux.HandleFunc("GET /health", s.healthHandler)

{{if .WithAI}}
	// 2. AI Initialization
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️ [WARNING] OPENAI_API_KEY is missing! AI features will fail.")
	}
	model := os.Getenv("OPENAI_MODEL")
	aiService := ai.NewService(s.db, apiKey, model)
	aiHandler := ai.NewHandler(aiService)
{{end}}

	// 3. Protected Routes
	protected := http.NewServeMux()
	protected.HandleFunc("GET /me", s.meHandler)

{{if .WithAI}}
	protected.HandleFunc("POST /ai/generate", aiHandler.HandleGenerate)
	protected.HandleFunc("POST /ai/chat", aiHandler.HandleChat)
{{end}}

	// Mount protected routes under /api/
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

// 5. MIDDLEWARE (internal/server/middleware.go)
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
const UserIDKey UserContextKey = "userID"

// --- AUTH MIDDLEWARE ---

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// 1. Try Standard Better-Auth Cookie
		cookie, err := r.Cookie("better-auth.session_token")
		if err == nil {
			token = cookie.Value
		} else {
			// 2. Try Raw Session Token
			cookie, err = r.Cookie("session_token")
			if err == nil {
				token = cookie.Value
			} else {
				// 3. Try Authorization Header
				authHeader := r.Header.Get("Authorization")
				if t, found := strings.CutPrefix(authHeader, "Bearer "); found {
					token = t
				}
			}
		}

		if token == "" {
			log.Println("🔑 [AUTH] No Token found in Cookies or Header")
			http.Error(w, "Unauthorized: No Signal", http.StatusUnauthorized)
			return
		}

		// Strip signature (Better-Auth format: token.signature)
		// Note: We only strip if a dot is present
		if strings.Contains(token, ".") {
			if before, _, found := strings.Cut(token, "."); found {
				token = before
			}
		}

		s.verifySession(w, r, next, token)
	})
}

func (s *Server) verifySession(w http.ResponseWriter, r *http.Request, next http.Handler, token string) {
	query := ` + "`" + `
		SELECT user_id 
		FROM public.session 
		WHERE token = $1 
		AND expires_at > NOW()
	` + "`" + `

	var userID string
	err := s.db.QueryRow(query, token).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("🔑 [AUTH] Invalid or Expired Token: %s", token)
			http.Error(w, "Unauthorized: Invalid Token", http.StatusUnauthorized)
			return
		}
		log.Printf("❌ [AUTH] DB Error: %v", err)
		http.Error(w, "System Failure", http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(r.Context(), UserIDKey, userID)
	next.ServeHTTP(w, r.WithContext(ctx))
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
		
		// In a hybrid local environment, we allow localhost:3000
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

// 6. AI SERVICE (internal/ai/service.go)
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

// 7. AI HANDLER (internal/ai/handler.go)
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

// 8. CONFIG (.env)
const EnvFile = `PORT=8080
DATABASE_URL=postgres://postgres:password@localhost:5432/{{.Name}}?sslmode=disable
{{if .WithAI}}OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-5-nano{{end}}
`

// 9. IGNORE (.gitignore)
const GitIgnore = `bin/
main
.DS_Store
.env
`

// 10. MAKEFILE
const Makefile = `PROJECT_NAME := {{.Name}}

.PHONY: all build run test clean docker

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

clean:
	@rm -rf bin/

docker:
	@docker compose up -d
`

// 11. DOCKER (compose.yml)
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
