# Multi-stage build for optimized production and development
ARG NODE_ENV=development

# ===== Dependencies Stage =====
FROM docker.io/node:22-alpine AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files
COPY package*.json pnpm-lock.yaml* ./

# Install dependencies based on environment
RUN if [ "$NODE_ENV" = "production" ]; then \
    pnpm install --frozen-lockfile --prod; \
  else \
    pnpm install --frozen-lockfile; \
  fi

# ===== Builder Stage (Production only) =====
FROM docker.io/node:22-alpine AS builder
WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Accept build arguments for environment variables
ARG NEXT_PUBLIC_API_URL
ARG BACKEND_URL

# Set environment variables for build
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV BACKEND_URL=$BACKEND_URL

# Build the application
RUN pnpm run build

# ===== Production Stage =====
FROM node:22-alpine AS production
WORKDIR /app

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy built application
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000
ENV PORT 3000
ENV NODE_ENV production

CMD ["node", "server.js"]

# ===== Development Stage =====
FROM docker.io/node:22-alpine AS development
WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules
COPY package*.json pnpm-lock.yaml* ./

# Copy source code
COPY . .

EXPOSE 3000
ENV NODE_ENV development

CMD ["pnpm", "run", "dev"]

# ===== Final Stage =====
FROM ${NODE_ENV} AS final
