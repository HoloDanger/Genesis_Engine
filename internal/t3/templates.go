package t3

// 1. IDENTITY (The Package File)
// Note: We use specific versions based on T3_Blueprint
const PackageJSON = `{
  "name": "{{.Name}}",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "eslint",
    "db:push": "drizzle-kit push",
    "db:studio": "drizzle-kit studio"
  },
  "dependencies": {
    "next": "16.1.1",
    "react": "^19.2.3",
    "react-dom": "^19.2.3",
    "better-auth": "^1.4.9",
    "drizzle-orm": "^0.45.1",
    "postgres": "^3.4.7",
    "dotenv": "^17.2.3"
  },
  "devDependencies": {
    "@tailwindcss/postcss": "^4",
    "@types/node": "^20",
    "@types/react": "^19",
    "@types/react-dom": "^19",
    "eslint": "^9",
    "eslint-config-next": "16.1.1",
    "tailwindcss": "^4",
    "typescript": "^5",
    "drizzle-kit": "^0.31.8",
    "tsx": "^4.21.0"
    },
  "ignoreScripts": [
    "sharp",
    "unrs-resolver"
  ],
  "trustedDependencies": [
    "sharp",
    "unrs-resolver"
  ]
}`

// 2. THE STYLE (Tailwind v4 CSS)
const GlobalCSS = `@import "tailwindcss";
@config "../../tailwind.config.ts";

@theme {
  --font-sans: var(--font-geist-sans), ui-sans-serif, system-ui, sans-serif,
    "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
}

body {
  background-color: #0a0a0a;
  color: #e5e5e5;
  font-family: var(--font-sans); /* Force the clean font */
  -webkit-font-smoothing: antialiased; /* Make it crisp */
}
`

// 3. THE CONFIG (Tailwind and PostCSS Config)
const PostCSSConfig = `const config = {
  plugins: {
    "@tailwindcss/postcss": {},
  },
};
export default config;
`

const TailwindConfig = `import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};
export default config;
`

// 4. THE INTERFACE (Page.tsx)
const MainPage = `"use client";

import Link from "next/link";
{{if .IsHybrid}}import { useEffect, useState } from "react";{{end}}
import { authClient } from "@/src/lib/auth-client";
import { useRouter } from "next/navigation";

export default function Home() {
{{if .IsHybrid}}  const [status, setStatus] = useState("SCANNING...");
  const [userId, setUserId] = useState<string | null>(null);{{end}}
  const router = useRouter();

{{if .IsHybrid}}  useEffect(() => {
    // Ping the Go Backend via the Next.js Proxy
    fetch("/go-api/me")
      .then((res) => {
        if (res.ok) return res.json();
        throw new Error("Unauthorized");
      })
      .then((data) => {
        setStatus("LINK ESTABLISHED");
        setUserId(data.userID);
      })
      .catch(() => {
        setStatus("LINK SEVERED (Unauthorized)");
        setUserId(null);
      });
  }, []);{{end}}

  const handleLogout = async () => {
    await authClient.signOut({
      fetchOptions: {
        onSuccess: () => {
{{if .IsHybrid}}          setStatus("LINK SEVERED (Unauthorized)");
          setUserId(null);{{end}}
          router.refresh();
        },
      },
    });
  };

  return (
    <main className="min-h-screen flex items-center justify-center space-y-4 flex-col font-sans bg-black text-white">
      <h1 className="text-4xl font-bold tracking-tighter">
        GENESIS NODE: {{.Name}}
      </h1>
      
{{if .IsHybrid}}      <div className="flex flex-col items-center space-y-2">
        <div className={"px-4 py-2 border rounded text-xs font-mono transition-colors " + (
          status.includes("ESTABLISHED") 
            ? "border-green-900 bg-green-950 text-green-400" 
            : status.includes("SCANNING")
            ? "border-yellow-900 bg-yellow-950 text-yellow-400"
            : "border-red-900 bg-red-950 text-red-400"
        )}>
          NEURAL LINK: {status}
        </div>
        
        {userId && (
           <div className="text-xs font-mono text-neutral-500">
             IDENTITY: {userId}
           </div>
        )}
      </div>{{end}}

      <div className="flex gap-4">
        <Link 
          href="/auth" 
          className="px-4 py-2 text-sm font-medium text-black bg-white rounded hover:bg-neutral-200 transition-colors"
        >
          Auth Portal
        </Link>
        
        <button
          onClick={handleLogout}
          className="px-4 py-2 text-sm font-medium text-white border border-neutral-800 bg-neutral-900 rounded hover:bg-red-900 hover:border-red-800 transition-colors"
        >
          Disconnect
        </button>
      </div>
    </main>
  );
}
`

// 5. The Nervous System (Layout.tsx)
const RootLayout = `import "./globals.css";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
`

// 6. TypeScript Configuration (tsconfig.json)
const TSConfig = `{
  "compilerOptions": {
    "target": "ES2017",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "react-jsx",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": [
    "next-env.d.ts",
    "**/*.ts",
    "**/*.tsx",
    ".next/types/**/*.ts",
    ".next/dev/types/**/*.ts",
    "**/*.mts"
  ],
  "exclude": ["node_modules"]
}
`

// 7. THE SHIELD (.gitignore)
const GitIgnore = `# System Files
.DS_Store
Thumbs.db

# Dependencies
node_modules/
.pnp
.pnp.js

# Build Output
.next/
out/
build/
dist/

# Secrets
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# Logs
npm-debug.log*
yarn-debug.log*
yarn-error.log*
`

// 8. THE CONTROL PLANE (next.config.ts)
const NextConfig = `import type { NextConfig } from "next";

const nextConfig: NextConfig = {
{{if .IsHybrid}}  async rewrites() {
    return [
      {
        source: "/go-api/:path*",
        destination: "http://localhost:8080/api/:path*", // Proxy to Go Backend
      },
    ];
  },{{end}}
};

export default nextConfig;
`

// 9. AUTH SERVER CONFIG (src/lib/auth.ts)
const AuthServer = `import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from "@/src/server/db";
import * as schema from "@/src/server/schema";

export const auth = betterAuth({
  database: drizzleAdapter(db, {
    provider: "pg",
    schema: schema,
  }),
  emailAndPassword: {
    enabled: true,
  }
});
`

// 10. AUTH API ROUTE (src/app/api/auth/[...all]/route.ts)
const AuthRoute = `import { auth } from "@/src/lib/auth";
import { toNextJsHandler } from "better-auth/next-js";

export const { GET, POST } = toNextJsHandler(auth);
`

// 11. AUTH CLIENT (src/lib/auth-client.ts)
const AuthClient = `import { createAuthClient } from "better-auth/react";

export const authClient = createAuthClient({
  baseURL: process.env.BETTER_AUTH_URL,
});
`

// 12. THE ENGINE (compose.yml)
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

// 13. AUTH UI (src/app/auth/page.tsx)
const LoginPage = `"use client";

import { authClient } from "@/src/lib/auth-client";
import { useState } from "react";

export default function AuthPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [isLogin, setIsLogin] = useState(true);

  const handleSubmit = async () => {
    if (isLogin) {
      await authClient.signIn.email({
        email,
        password,
        callbackURL: "/",
      });
    } else {
      await authClient.signUp.email({
        email,
        password,
        name,
        callbackURL: "/",
      });
    }
  };

  return (
    <main className="min-h-screen flex items-center justify-center flex-col font-sans bg-black text-white">
      <div className="w-full max-w-md p-8 space-y-6 border border-neutral-800 rounded-lg bg-neutral-900">
        <h1 className="text-2xl font-bold tracking-tight text-center">
          {isLogin ? "Sign In" : "Create Account"}
        </h1>
        
        <div className="space-y-4">
          {!isLogin && (
            <div className="space-y-2">
              <label className="text-sm font-medium text-neutral-400">Name</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-3 py-2 bg-black border border-neutral-800 rounded focus:border-white focus:outline-none"
                placeholder="Agent Name"
              />
            </div>
          )}
          
          <div className="space-y-2">
            <label className="text-sm font-medium text-neutral-400">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-3 py-2 bg-black border border-neutral-800 rounded focus:border-white focus:outline-none"
              placeholder="contact@holodanger.dev"
            />
          </div>
          
          <div className="space-y-2">
            <label className="text-sm font-medium text-neutral-400">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 bg-black border border-neutral-800 rounded focus:border-white focus:outline-none"
              placeholder="••••••••"
            />
          </div>

          <button
            onClick={handleSubmit}
            className="w-full py-2 font-medium bg-white text-black rounded hover:bg-neutral-200 transition-colors"
          >
            {isLogin ? "Enter System" : "Initialize Node"}
          </button>
        </div>

        <div className="text-center text-sm text-neutral-500">
          <button 
            onClick={() => setIsLogin(!isLogin)}
            className="hover:text-white underline"
          >
            {isLogin ? "Need access? Initialize" : "Already active? Enter"}
          </button>
        </div>
      </div>
    </main>
  );
}
`
