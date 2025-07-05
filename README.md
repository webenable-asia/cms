# WebEnable CMS

A modern content management system built with Next.js 15, Go 1.24, and CouchDB.

## Features

- **Modern Stack**: Next.js 15 frontend with Go 1.24 RESTful API backend (released Feb 2025)
- **Database**: CouchDB for flexible document storage
- **Authentication**: JWT-based authentication system
- **Content Management**: Create, edit, and publish blog posts
- **Admin Interface**: Clean and intuitive admin dashboard
- **Responsive Design**: Mobile-friendly interface with Tailwind CSS
- **Theme Toggle**: Light/Dark/System theme support with animated transitions
- **UI Components**: Radix UI components with shadcn/ui design system
- **Docker Support**: Containerized development environment

## Prerequisites

- Docker and Docker Compose
- Node.js 20+ (for local development)
- Go 1.24+ (for local development - released February 2025)

## Quick Start

1. **Clone and navigate to the project:**
   ```bash
   cd /Users/tsaa/workspace/projects/webenable/cms
   ```

2. **Start the development environment:**
   ```bash
   docker-compose up --build
   ```

3. **Access the applications:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - CouchDB: http://localhost:5984 (admin/password)

4. **Default admin credentials:**
   - Username: `admin`
   - Password: `password`

## Project Structure

```
├── docker-compose.yml          # Docker Compose configuration
├── backend/                    # Go backend application
│   ├── Dockerfile             # Backend Docker configuration
│   ├── .air.toml              # Air live reload configuration
│   ├── main.go                # Main application entry point
│   ├── go.mod                 # Go module dependencies
│   ├── handlers/              # HTTP request handlers
│   │   ├── posts.go           # Post-related endpoints
│   │   ├── posts_protected.go # Protected post operations
│   │   └── auth.go            # Authentication endpoints
│   ├── models/                # Data models
│   │   └── models.go          # Post and User models
│   ├── database/              # Database connection and setup
│   │   └── database.go        # CouchDB initialization
│   └── middleware/            # HTTP middleware
│       └── auth.go            # JWT authentication middleware
└── frontend/                  # Next.js frontend application
    ├── Dockerfile             # Frontend Docker configuration
    ├── package.json           # Node.js dependencies
    ├── next.config.js         # Next.js configuration
    ├── tailwind.config.js     # Tailwind CSS configuration
    ├── app/                   # Next.js App Router
    │   ├── layout.tsx         # Root layout
    │   ├── page.tsx           # Home page
    │   ├── globals.css        # Global styles
    │   ├── blog/              # Blog section
    │   │   ├── page.tsx       # Blog listing
    │   │   └── [id]/          # Individual post pages
    │   └── admin/             # Admin section
    │       ├── page.tsx       # Admin login
    │       ├── dashboard/     # Admin dashboard
    │       └── posts/         # Post management
    ├── components/            # Reusable React components
    │   ├── navigation.tsx     # Main navigation with theme toggle
    │   ├── theme-provider.tsx # Theme context provider
    │   └── ui/                # UI component library
    │       ├── button.tsx     # Button component
    │       ├── dropdown-menu.tsx # Dropdown menu component
    │       └── theme-toggle-reference.tsx # Theme toggle component
    └── lib/                   # Utility libraries
        └── api.ts             # API client configuration
```

## API Endpoints

### Public Endpoints
- `GET /api/posts` - Get all published posts
- `GET /api/posts/{id}` - Get a specific post
- `POST /api/auth/login` - User authentication

### Protected Endpoints (require JWT token)
- `POST /api/posts` - Create a new post
- `PUT /api/posts/{id}` - Update an existing post
- `DELETE /api/posts/{id}` - Delete a post

## Development

### Theme System

The application includes a comprehensive theme system with:

- **Light Mode**: Clean, bright interface for daytime use
- **Dark Mode**: Dark theme for reduced eye strain
- **System Mode**: Automatically follows OS theme preference
- **Animated Transitions**: Smooth icon animations and theme switching
- **Persistent Storage**: Theme preference saved in localStorage

The theme toggle is located in the navigation bar and uses:
- **next-themes**: Theme management library
- **Radix UI**: Accessible dropdown components
- **Lucide Icons**: Sun/Moon/Monitor icons with CSS transitions
- **CSS Variables**: Comprehensive color system for both themes

### Backend Development

The backend uses Air for live reloading. Any changes to Go files will automatically restart the server.

To run the backend locally without Docker:
```bash
cd backend
go mod download
go install github.com/cosmtrek/air@latest
air
```

### Frontend Development

The frontend uses Next.js with hot reloading enabled. Changes to React components will be reflected immediately.

To run the frontend locally without Docker:
```bash
cd frontend
npm install
npm run dev
```

### Database Management

CouchDB is accessible at http://localhost:5984/_utils with admin credentials (admin/password).

The application automatically creates the following databases:
- `posts` - Stores blog posts
- `users` - Stores user information

## Environment Variables

### Backend
- `COUCHDB_URL` - CouchDB connection string (default: http://admin:password@db:5984/)
- `PORT` - Server port (default: 8080)

### Frontend
- `NEXT_PUBLIC_API_URL` - Backend API URL (default: http://localhost:8080)

## Production Deployment

For production deployment, you'll need to:

1. **Update Docker configurations** for production builds
2. **Set secure environment variables** (JWT secret, database credentials)
3. **Configure CORS** for your production domain
4. **Set up SSL/HTTPS** termination
5. **Use a production-ready database setup**

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License.
