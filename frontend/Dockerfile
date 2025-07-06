# Use Node.js 22 Alpine for smaller image size
FROM node:22-alpine

# Install pnpm globally
RUN npm install -g pnpm

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json pnpm-lock.yaml* ./

# Install dependencies
RUN pnpm install

# Copy source code
COPY . .

# Expose port 3000
EXPOSE 3000

# Start development server
CMD ["pnpm", "run", "dev"]
