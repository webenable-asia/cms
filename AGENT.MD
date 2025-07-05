# Plan for Developing a Company Website, Blog, and CMS

This document outlines the plan for an AI agent to develop a company website, blog, and a content management system (CMS) using Next.js 15 for the frontend, Go 1.24 (released February 2025) for the backend RESTful API, and CouchDB as the database. This updated plan incorporates Docker Compose for a streamlined local development environment.

## Phase 1: Local Development Setup with Docker Compose

This phase focuses on setting up a containerized local development environment.

1.  **Prerequisites:**
    *   Install Docker and Docker Compose.

2.  **Project Structure:**
    *   Create a main project directory.
    *   Inside, create two subdirectories: `backend` and `frontend`.

3.  **Create `docker-compose.yml`:**
    *   In the root of your project, create a `docker-compose.yml` file. This file will define the `backend`, `frontend`, and `db` services.

    ```yaml


    services:
      db:
        image: couchdb:3
        restart: always
        ports:
          - "5984:5984"
        environment:
          COUCHDB_USER: admin
          COUCHDB_PASSWORD: password
        volumes:
          - couchdb_data:/opt/couchdb/data

      backend:
        build: ./backend
        ports:
          - "8080:8080"
        volumes:
          - ./backend:/app
        depends_on:
          - db
        environment:
          - COUCHDB_URL=http://admin:password@db:5984/

      frontend:
        build: ./frontend
        ports:
          - "3000:3000"
        volumes:
          - ./frontend:/app
          - /app/node_modules
          - /app/.next
        depends_on:
          - backend
        environment:
          - NEXT_PUBLIC_API_URL=http://localhost:8080

    volumes:
      couchdb_data:
    ```

4.  **Create Backend `Dockerfile`:**
    *   Inside the `backend` directory, create a `Dockerfile` for the Go service. This setup uses `air` for live reloading.

    ```dockerfile
    # backend/Dockerfile
    FROM golang:1.24-alpine

    WORKDIR /app

    # Install air for live reloading
    RUN go install github.com/cosmtrek/air@latest

    COPY go.mod go.sum ./
    RUN go mod download

    COPY . .

    EXPOSE 8080

    # Use air for live-reloading
    CMD ["air"]
    ```
    *   You will also need to add a `.air.toml` configuration file in the `backend` directory.

5.  **Create Frontend `Dockerfile`:**
    *   Inside the `frontend` directory, create a `Dockerfile` for the Next.js service.

    ```dockerfile
    # frontend/Dockerfile
    FROM node:20-alpine

    WORKDIR /app

    COPY package*.json ./
    RUN npm install

    COPY . .

    EXPOSE 3000

    CMD ["npm", "run", "dev"]
    ```

6.  **Create Backend Configuration:**
    *   Create `.air.toml` configuration file in the `backend` directory for Go live reloading:
    ```toml
    root = "."
    testdata_dir = "testdata"
    tmp_dir = "tmp"

    [build]
      args_bin = []
      bin = "./tmp/main"
      cmd = "go build -o ./tmp/main ."
      delay = 1000
      exclude_dir = ["assets", "tmp", "vendor", "testdata"]
      exclude_file = []
      exclude_regex = ["_test.go"]
      exclude_unchanged = false
      follow_symlink = false
      full_bin = ""
      include_dir = []
      include_ext = ["go", "tpl", "tmpl", "html"]
      include_file = []
      kill_delay = "0s"
      log = "build-errors.log"
      poll = false
      poll_interval = 0
      rerun = false
      rerun_delay = 500
      send_interrupt = false
      stop_on_root = false

    [color]
      app = ""
      build = "yellow"
      main = "magenta"
      runner = "green"
      watcher = "cyan"

    [log]
      main_only = false
      time = false

    [misc]
      clean_on_exit = false

    [screen]
      clear_on_rebuild = false
      keep_scroll = true
    ```

7.  **Launch the Environment:**
    *   From the project root, run `docker-compose up --build`.
    *   This will build the images and start the containers.
    *   The frontend will be available at `http://localhost:3000`.
    *   The backend will be available at `http://localhost:8080`.
    *   CouchDB will be available at `http://localhost:5984`.

## Phase 2: Backend Development (Go RESTful API)

With the Docker environment running, you can now develop the backend service. The files are mounted as volumes, so changes in your local `backend` directory will trigger a live reload in the container.

1.  **Initialize Go Project:**
    *   `cd backend`
    *   `go mod init cms-backend`
    *   Install required dependencies:
    ```bash
    go get github.com/go-kivik/kivik/v4
    go get github.com/go-kivik/couchdb/v4
    go get github.com/gorilla/mux
    go get github.com/gorilla/handlers
    go get github.com/golang-jwt/jwt/v5
    go get github.com/joho/godotenv
    ```
    *   Create a `main.go` file.

2.  **Structure the Go Application:**
    *   Create packages for `handlers`, `models`, `database`, and `middleware`.
    *   `handlers`: To handle HTTP requests and responses.
    *   `models`: To define data structures (e.g., `Post`, `Page`, `User`).
    *   `database`: To handle the connection and operations with CouchDB.
    *   `middleware`: For authentication, CORS, and logging.

3.  **Implement Database Connection:**
    *   Create `database/connection.go` to handle CouchDB connection:
    ```go
    package database

    import (
        "context"
        "log"
        "os"
        
        "github.com/go-kivik/kivik/v4"
        _ "github.com/go-kivik/couchdb/v4"
    )

    var DB *kivik.DB

    func InitDB() {
        couchdbURL := os.Getenv("COUCHDB_URL")
        if couchdbURL == "" {
            couchdbURL = "http://admin:password@localhost:5984/"
        }

        client, err := kivik.New("couch", couchdbURL)
        if err != nil {
            log.Fatal("Failed to connect to CouchDB:", err)
        }

        DB = client.DB(context.Background(), "cms")
        
        // Create database if it doesn't exist
        if !DB.Exists(context.Background()) {
            if err := client.CreateDB(context.Background(), "cms"); err != nil {
                log.Fatal("Failed to create database:", err)
            }
        }
    }
    ```

4.  **Implement API Endpoints:**
    *   Define CRUD endpoints for blog posts and other content types:
    ```go
    // Example routes in main.go
    r := mux.NewRouter()
    
    // Public routes
    r.HandleFunc("/api/posts", handlers.GetPosts).Methods("GET")
    r.HandleFunc("/api/posts/{id}", handlers.GetPost).Methods("GET")
    r.HandleFunc("/api/pages/{slug}", handlers.GetPage).Methods("GET")
    
    // Protected routes (require authentication)
    protected := r.PathPrefix("/api").Subrouter()
    protected.Use(middleware.AuthMiddleware)
    protected.HandleFunc("/posts", handlers.CreatePost).Methods("POST")
    protected.HandleFunc("/posts/{id}", handlers.UpdatePost).Methods("PUT")
    protected.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")
    ```

5.  **Implement Authentication:**
    *   Add JWT-based authentication middleware for protected routes:
    ```go
    package middleware

    import (
        "context"
        "net/http"
        "strings"
        
        "github.com/golang-jwt/jwt/v5"
    )

    func AuthMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }

            tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
            token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                return []byte("your-secret-key"), nil
            })

            if err != nil || !token.Valid {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
    ```

## Phase 3: Frontend Development (Next.js 15)

Similar to the backend, the `frontend` directory is mounted as a volume, enabling hot-reloading for the Next.js application.

1.  **Initialize Next.js Project:**
    *   `cd frontend`
    *   `npx create-next-app@latest . --typescript --tailwind --eslint --app --src-dir --import-alias "@/*"`
    *   Choose "Yes" for TypeScript, Tailwind CSS, ESLint, App Router, and src/ directory.

2.  **Setup Tailwind CSS v4:**
    *   Install the latest Tailwind CSS v4:
    ```bash
    npm install @tailwindcss/cli@next @tailwindcss/vite@next
    ```
    *   Update `tailwind.config.ts`:
    ```typescript
    import type { Config } from 'tailwindcss'

    const config: Config = {
      content: [
        './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
        './src/components/**/*.{js,ts,jsx,tsx,mdx}',
        './src/app/**/*.{js,ts,jsx,tsx,mdx}',
      ],
      theme: {
        extend: {
          colors: {
            background: 'hsl(var(--background))',
            foreground: 'hsl(var(--foreground))',
            card: {
              DEFAULT: 'hsl(var(--card))',
              foreground: 'hsl(var(--card-foreground))',
            },
            popover: {
              DEFAULT: 'hsl(var(--popover))',
              foreground: 'hsl(var(--popover-foreground))',
            },
            primary: {
              DEFAULT: 'hsl(var(--primary))',
              foreground: 'hsl(var(--primary-foreground))',
            },
            secondary: {
              DEFAULT: 'hsl(var(--secondary))',
              foreground: 'hsl(var(--secondary-foreground))',
            },
            muted: {
              DEFAULT: 'hsl(var(--muted))',
              foreground: 'hsl(var(--muted-foreground))',
            },
            accent: {
              DEFAULT: 'hsl(var(--accent))',
              foreground: 'hsl(var(--accent-foreground))',
            },
            destructive: {
              DEFAULT: 'hsl(var(--destructive))',
              foreground: 'hsl(var(--destructive-foreground))',
            },
            border: 'hsl(var(--border))',
            input: 'hsl(var(--input))',
            ring: 'hsl(var(--ring))',
            chart: {
              '1': 'hsl(var(--chart-1))',
              '2': 'hsl(var(--chart-2))',
              '3': 'hsl(var(--chart-3))',
              '4': 'hsl(var(--chart-4))',
              '5': 'hsl(var(--chart-5))',
            },
          },
          borderRadius: {
            lg: 'var(--radius)',
            md: 'calc(var(--radius) - 2px)',
            sm: 'calc(var(--radius) - 4px)',
          },
        },
      },
      plugins: [require('tailwindcss-animate')],
    }
    export default config
    ```

3.  **Setup Shadcn/UI:**
    *   Initialize Shadcn/UI:
    ```bash
    npx shadcn@latest init
    ```
    *   Configure `components.json`:
    ```json
    {
      "$schema": "https://ui.shadcn.com/schema.json",
      "style": "default",
      "rsc": true,
      "tsx": true,
      "tailwind": {
        "config": "tailwind.config.ts",
        "css": "src/app/globals.css",
        "baseColor": "slate",
        "cssVariables": true,
        "prefix": ""
      },
      "aliases": {
        "components": "@/components",
        "utils": "@/lib/utils"
      }
    }
    ```
    *   Install essential Shadcn/UI components:
    ```bash
    npx shadcn@latest add button card input textarea label form select dropdown-menu navigation-menu sheet dialog alert badge
    ```

4.  **Update Global CSS:**
    *   Update `src/app/globals.css` to include Shadcn/UI variables:
    ```css
    @tailwind base;
    @tailwind components;
    @tailwind utilities;

    @layer base {
      :root {
        --background: 0 0% 100%;
        --foreground: 222.2 84% 4.9%;
        --card: 0 0% 100%;
        --card-foreground: 222.2 84% 4.9%;
        --popover: 0 0% 100%;
        --popover-foreground: 222.2 84% 4.9%;
        --primary: 222.2 47.4% 11.2%;
        --primary-foreground: 210 40% 98%;
        --secondary: 210 40% 96%;
        --secondary-foreground: 222.2 47.4% 11.2%;
        --muted: 210 40% 96%;
        --muted-foreground: 215.4 16.3% 46.9%;
        --accent: 210 40% 96%;
        --accent-foreground: 222.2 47.4% 11.2%;
        --destructive: 0 84.2% 60.2%;
        --destructive-foreground: 210 40% 98%;
        --border: 214.3 31.8% 91.4%;
        --input: 214.3 31.8% 91.4%;
        --ring: 222.2 84% 4.9%;
        --radius: 0.5rem;
        --chart-1: 12 76% 61%;
        --chart-2: 173 58% 39%;
        --chart-3: 197 37% 24%;
        --chart-4: 43 74% 66%;
        --chart-5: 27 87% 67%;
      }

      .dark {
        --background: 222.2 84% 4.9%;
        --foreground: 210 40% 98%;
        --card: 222.2 84% 4.9%;
        --card-foreground: 210 40% 98%;
        --popover: 222.2 84% 4.9%;
        --popover-foreground: 210 40% 98%;
        --primary: 210 40% 98%;
        --primary-foreground: 222.2 47.4% 11.2%;
        --secondary: 217.2 32.6% 17.5%;
        --secondary-foreground: 210 40% 98%;
        --muted: 217.2 32.6% 17.5%;
        --muted-foreground: 215 20.2% 65.1%;
        --accent: 217.2 32.6% 17.5%;
        --accent-foreground: 210 40% 98%;
        --destructive: 0 62.8% 30.6%;
        --destructive-foreground: 210 40% 98%;
        --border: 217.2 32.6% 17.5%;
        --input: 217.2 32.6% 17.5%;
        --ring: 212.7 26.8% 83.9%;
        --chart-1: 220 70% 50%;
        --chart-2: 160 60% 45%;
        --chart-3: 30 80% 55%;
        --chart-4: 280 65% 60%;
        --chart-5: 340 75% 55%;
      }
    }

    @layer base {
      * {
        @apply border-border;
      }
      body {
        @apply bg-background text-foreground;
      }
    }
    ```

5.  **Structure the Next.js Application:**
    *   Use the App Router and create directories for `components`, `app`, `lib`, etc.
    *   Create `src/lib/utils.ts` for utility functions:
    ```typescript
    import { type ClassValue, clsx } from "clsx"
    import { twMerge } from "tailwind-merge"

    export function cn(...inputs: ClassValue[]) {
      return twMerge(clsx(inputs))
    }
    ```

6.  **Build UI Components and Pages:**
    *   Create pages for the website (Home, About, Contact) using Shadcn/UI components.
    *   Create pages for the blog (post list and single post view) with beautiful card layouts.
    *   Create a CMS interface under `/admin` for content management using forms, tables, and dialogs.
    *   Implement dark mode support using next-themes:
    ```bash
    npm install next-themes
    ```
    *   Create a theme provider component:
    ```typescript
    // src/components/theme-provider.tsx
    "use client"

    import * as React from "react"
    import { ThemeProvider as NextThemesProvider } from "next-themes"
    import { type ThemeProviderProps } from "next-themes/dist/types"

    export function ThemeProvider({ children, ...props }: ThemeProviderProps) {
      return <NextThemesProvider {...props}>{children}</NextThemesProvider>
    }
    ```

7.  **Create API Service Layer:**
    *   Create `src/lib/api.ts` for backend communication:
    ```typescript
    const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

    export interface Post {
      id: string
      title: string
      content: string
      slug: string
      published: boolean
      createdAt: string
      updatedAt: string
    }

    export const api = {
      // Public endpoints
      async getPosts(): Promise<Post[]> {
        const response = await fetch(`${API_BASE_URL}/api/posts`)
        if (!response.ok) throw new Error('Failed to fetch posts')
        return response.json()
      },

      async getPost(id: string): Promise<Post> {
        const response = await fetch(`${API_BASE_URL}/api/posts/${id}`)
        if (!response.ok) throw new Error('Failed to fetch post')
        return response.json()
      },

      // Protected endpoints
      async createPost(post: Omit<Post, 'id' | 'createdAt' | 'updatedAt'>, token: string): Promise<Post> {
        const response = await fetch(`${API_BASE_URL}/api/posts`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(post)
        })
        if (!response.ok) throw new Error('Failed to create post')
        return response.json()
      },

      async updatePost(id: string, post: Partial<Post>, token: string): Promise<Post> {
        const response = await fetch(`${API_BASE_URL}/api/posts/${id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(post)
        })
        if (!response.ok) throw new Error('Failed to update post')
        return response.json()
      },

      async deletePost(id: string, token: string): Promise<void> {
        const response = await fetch(`${API_BASE_URL}/api/posts/${id}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        if (!response.ok) throw new Error('Failed to delete post')
      }
    }
    ```

## Phase 4: Integration and Deployment

1.  **Connect Frontend to Backend:**
    *   The `NEXT_PUBLIC_API_URL` environment variable in `docker-compose.yml` is already configured to point to the backend service.
    *   Use the API service layer created in Phase 3 to make requests to the Go backend.
    *   Implement error handling and loading states in your React components.

2.  **CORS Configuration:**
    *   Configure CORS on the Go backend to allow requests from the frontend:
    ```go
    // In main.go
    import "github.com/gorilla/handlers"

    func main() {
        r := mux.NewRouter()
        // ... route definitions ...

        // Configure CORS
        corsOptions := handlers.CORS(
            handlers.AllowedOrigins([]string{"http://localhost:3000"}),
            handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
            handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
        )

        log.Println("Server starting on :8080")
        log.Fatal(http.ListenAndServe(":8080", corsOptions(r)))
    }
    ```

3.  **Environment Variables:**
    *   Create `.env.local` file in frontend directory:
    ```env
    NEXT_PUBLIC_API_URL=http://localhost:8080
    ```
    *   Create `.env` file in backend directory:
    ```env
    COUCHDB_URL=http://admin:password@db:5984/
    JWT_SECRET=your-super-secret-jwt-key
    PORT=8080
    ```

4.  **Build for Production:**
    *   Create production Dockerfiles:
    
    **Backend Production Dockerfile:**
    ```dockerfile
    # backend/Dockerfile.prod
    FROM golang:1.24-alpine AS builder
    WORKDIR /app
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    RUN go build -o main .

    FROM alpine:latest
    RUN apk --no-cache add ca-certificates
    WORKDIR /root/
    COPY --from=builder /app/main .
    EXPOSE 8080
    CMD ["./main"]
    ```

    **Frontend Production Dockerfile:**
    ```dockerfile
    # frontend/Dockerfile.prod
    FROM node:20-alpine AS builder
    WORKDIR /app
    COPY package*.json ./
    RUN npm ci --only=production
    COPY . .
    RUN npm run build

    FROM node:20-alpine AS runner
    WORKDIR /app
    ENV NODE_ENV production
    RUN addgroup --system --gid 1001 nodejs
    RUN adduser --system --uid 1001 nextjs
    COPY --from=builder /app/public ./public
    COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
    COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
    USER nextjs
    EXPOSE 3000
    ENV PORT 3000
    CMD ["node", "server.js"]
    ```

5.  **Production Docker Compose:**
    ```yaml
    # docker-compose.prod.yml
    version: '3.8'

    services:
      db:
        image: couchdb:3
        restart: unless-stopped
        environment:
          COUCHDB_USER: admin
          COUCHDB_PASSWORD: ${COUCHDB_PASSWORD:-password}
        volumes:
          - couchdb_data:/opt/couchdb/data
        networks:
          - app-network

      backend:
        build:
          context: ./backend
          dockerfile: Dockerfile.prod
        restart: unless-stopped
        environment:
          - COUCHDB_URL=http://admin:${COUCHDB_PASSWORD:-password}@db:5984/
          - JWT_SECRET=${JWT_SECRET}
        depends_on:
          - db
        networks:
          - app-network

      frontend:
        build:
          context: ./frontend
          dockerfile: Dockerfile.prod
        restart: unless-stopped
        ports:
          - "80:3000"
        environment:
          - NEXT_PUBLIC_API_URL=http://backend:8080
        depends_on:
          - backend
        networks:
          - app-network

    volumes:
      couchdb_data:

    networks:
      app-network:
        driver: bridge
    ```

6.  **Deployment Options:**
    *   **Cloud Platforms**: Deploy to AWS ECS, Google Cloud Run, or DigitalOcean App Platform.
    *   **VPS Deployment**: Use the production docker-compose on any VPS with Docker installed.
    *   **Frontend Deployment**: Alternatively, deploy the Next.js app to Vercel or Netlify and keep the backend containerized.
    *   **Database**: Use IBM Cloudant (CouchDB-compatible) or self-hosted CouchDB with proper backups.
