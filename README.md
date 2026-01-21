# GENESIS ENGINE

### Sovereign Scaffolding CLI

**Doctrine:** Archon T3 (Next.js 16 / React 19 / Tailwind v4)
**Status:** Operational / Void Week Build

---

## I. PURPOSE

Genesis is a high-velocity injection tool designed to spawn production-ready nodes within the Archon ecosystem. It eliminates configuration fatigue by enforcing a strict, opinionated architectural standard.

## II. THE STACK (T3_GAMMA)

- **Next.js 16.1.1** (React 19 / Tailwind v4)
- **Better Auth v1.4.9** (Drizzle Adapter + RBAC)
- **Drizzle ORM** (Postgres)
- **Bun** (Runtime & Package Management)
- **Docker** (Postgres 16-Alpine Instance)
- **Go 1.23+** (The Spear / Backend Logic)

## III. INSTALLATION (THE ARSENAL)

```bash
# 1. Compile the Weapon
go build -o genesis main.go

# 2. Deploy to Global Path
mv genesis ~/Archon/Arsenal/
```

## IV. USAGE

### 1. The Shield (Frontend Only)
```bash
genesis -name <ProjectName> -type t3
```

### 2. The Spear (Backend Only)
```bash
# Standard
genesis -name <ProjectName> -type go

# With AI Modules (OpenAI Integration)
genesis -name <ProjectName> -type go -ai=true
```

### 3. The Archon Hybrid (Twin Engine)
Combines T3 Frontend + Go Backend with shared Database infrastructure.
```bash
genesis -name <ProjectName> -type hybrid -ai=true
```

## V. IGNITION

After spawning a node:

```bash
cd <ProjectName>

# 1. Infrastructure
docker compose up -d  # Spin up the Database

# 2. Synchronization
# For T3/Hybrid:
bun db:push           # Sync Schema to DB

# 3. Execution
bun dev               # Start the Server (T3)
make run              # Start the Server (Go)
```

## VI. PHILOSOPHY

- **Zero Questions:** No interactive prompts. Decisions are already made.
- **Local First:** Integrated Docker provisioning for database sovereignty.
- **Minimalism:** No "slop" code. High density, monospace-first aesthetics.

---

_Logic over Noise. Structure over Chaos._