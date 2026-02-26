# 🛡️ GENESIS ENGINE: EFFICIENCY BASELINE UPGRADE (PHASE 2)
**Objective:** Optimize Resource Utilization. Establish High-Efficiency Standards. Implement `--resilient` execution.

---

## 🏛️ 1. THE ARCHITECTURAL PIVOT
The current `hybrid` builder utilizes a "Resource-Heavy" stack (Next.js/Postgres). We are implementing the **Integrated SSR Engine** (Go SSR/SQLite) for maximum operational density.

### **A. The `--resilient` Flag**
- **Action:** Add `Resilient bool` to the `Config` struct in `main.go`.
- **Logic:** 
  - If `false`: Generate standard Next.js boilerplate.
  - If `true`: Generate **High-Density SSR** logic.

### **B. The Local Persistence Baseline (SQLite)**
- **Action:** Replace `postgres://` hardcoding in `hybrid/builder.go` with `sqlite.db`.
- **Target:** All generated `.env` files point to local SQLite persistence for deterministic data-locality.

---

## 🏗️ 2. THE HIGH-DENSITY TEMPLATE (HTMX + VANILLA)
Genesis is now capable of generating interfaces that bypass heavy client-side runtimes.

- **The Stack:** Go `html/template` + **HTMX** (for real-time signals) + **Vanilla CSS**.
- **The Deployment:** A single, statically linked binary serving both the API and the UI.
- **The Footprint:** Targeted < 20MB RAM usage.

---

## 🧠 3. INTELLIGENCE INTEGRATION
Update the AI boilerplate generation to match the **Operational Intelligence Framework**.

- **Model ID:** Default to `us.anthropic.claude-sonnet-4-6`.
- **Structure:** Automatically generate the `internal/ai/context.md` file for project-specific grounding.
- **Provider:** Utilize the **Bedrock Runtime Bridge** for high-performance inference.

---

## ⚡ 4. PERFORMANCE OPTIMIZATION (PURGE)
- **Remove:** External authentication dependencies from the `resilient` template.
- **Implement:** Native **JWT + Cookie** logic generated directly into `internal/auth`.
- **Delete:** Unnecessary containerization services for Postgres/Redis in the high-efficiency baseline.

---

## 📊 THE VERDICT
By completing this refactor, Genesis will generate systems that are **Optimized for High-Constraint Environments.** You can deploy mission-critical logic with a near-zero infrastructural footprint.

**Mantra:** Architecture is the ultimate lever for efficiency. **Genesis v2.0 empowers the Architect.**
