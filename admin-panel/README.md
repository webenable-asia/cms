# WebEnable Admin Panel

This is the separate admin panel for the WebEnable CMS system. It provides a dedicated interface for content management, user administration, and system configuration.

## Features

- **Separate from main frontend** - Enhanced security and maintainability
- **Markdown Editor** - Rich content editing with live preview for blog posts
- **User Management** - Admin user administration
- **Dashboard** - Analytics and system overview
- **Content Management** - Blog posts, pages, and media management
- **Dark Mode Support** - Professional admin interface with theme switching

## Tech Stack

- **Next.js 15.3.5** - React framework with App Router
- **TypeScript** - Type-safe development
- **Tailwind CSS** - Utility-first styling
- **Radix UI** - Accessible component primitives
- **Markdown Editor** - @uiw/react-md-editor for content creation
- **Authentication** - JWT-based auth with session management

## Development

```bash
# Install dependencies
pnpm install

# Start development server (runs on port 3001)
pnpm dev

# Build for production
pnpm build

# Start production server
pnpm start
```

## Environment Variables

The admin panel shares environment variables with the main application through the `.env` file in the project root.

Key variables:
- `NEXT_PUBLIC_API_URL` - Backend API URL
- `BACKEND_URL` - Internal backend URL for SSR
- `JWT_SECRET` - JWT signing secret
- `ADMIN_USERNAME` - Default admin username
- `ADMIN_PASSWORD` - Default admin password

## Architecture

```
admin-panel/
├── app/                    # Next.js App Router pages
│   ├── admin/             # Admin routes
│   │   ├── dashboard/     # Main dashboard
│   │   ├── login/         # Authentication
│   │   ├── posts/         # Content management
│   │   ├── users/         # User management
│   │   └── contacts/      # Contact management
│   ├── globals.css        # Global styles with Markdown editor CSS
│   ├── layout.tsx         # Root layout with providers
│   └── page.tsx           # Root redirect page
├── components/            # Reusable UI components
│   ├── ui/               # Base UI components (Radix UI)
│   ├── auth/             # Authentication components
│   ├── admin/            # Admin-specific components
│   ├── post-editor.tsx   # Markdown post editor
│   └── theme-provider.tsx # Theme context provider
├── hooks/                # Custom React hooks
│   ├── use-auth.ts       # Authentication hook
│   ├── use-api.ts        # API interaction hooks
│   └── use-posts.ts      # Posts management hooks
├── lib/                  # Utility libraries
│   ├── api.ts            # API client configuration
│   ├── utils.ts          # General utilities
│   └── types.ts          # TypeScript type definitions
└── types/                # Type definitions
    └── api.ts            # API response types
```

## Deployment

The admin panel is deployed as a separate container in the Docker Compose stack:

- **Service name**: `admin-panel`
- **Port**: 3001
- **Routing**: Caddy proxy routes `/admin*` to admin panel
- **Public frontend**: All other routes go to main frontend

## Security Features

- **Separate deployment** - Admin panel isolated from public frontend
- **JWT Authentication** - Secure token-based authentication
- **Session Management** - Secure session handling
- **CSRF Protection** - Cross-site request forgery protection
- **Security Headers** - Comprehensive security headers via Caddy
- **Input Validation** - Server-side and client-side validation
- **Rate Limiting** - API rate limiting for admin endpoints

## Admin Panel vs Frontend

| Feature | Admin Panel | Main Frontend |
|---------|-------------|---------------|
| Purpose | Content management | Public website |
| Authentication | Required | Optional |
| Content | Admin interface | Public content |
| Port | 3001 | 3000 |
| Routes | `/admin/*` | All other routes |
| Users | Admin users only | Public visitors |
| Performance | Admin optimized | Public optimized |

## Content Management

### Blog Posts
- **Markdown Editor** - Rich text editing with live preview
- **SEO Optimization** - Meta titles, descriptions, and social media tags
- **Category Management** - Organize posts by categories
- **Tag System** - Flexible tagging for content organization
- **Featured Images** - Hero images for posts
- **Publishing Workflow** - Draft, review, and publish states
- **Reading Time Calculation** - Automatic reading time estimation

### User Management
- **Admin Users** - Create and manage admin accounts
- **Role-based Access** - Different permission levels
- **Activity Tracking** - Log admin activities
- **Password Management** - Secure password policies

### Analytics Dashboard
- **Content Statistics** - Post views, engagement metrics
- **User Analytics** - Admin activity tracking
- **System Health** - Performance monitoring
- **Database Metrics** - Storage and performance stats

## API Integration

The admin panel communicates with the Go backend API:

- **Authentication**: `/api/auth/login`, `/api/auth/logout`
- **Posts**: `/api/posts` (CRUD operations)
- **Users**: `/api/users` (User management)
- **Analytics**: `/api/analytics` (Dashboard metrics)
- **Health**: `/api/health` (System status)

All API calls are automatically proxied through Caddy for security and performance.

## Contributing

1. Follow the existing code structure and patterns
2. Use TypeScript for all new code
3. Ensure responsive design for all admin interfaces
4. Add proper error handling and loading states
5. Include JSDoc comments for complex functions
6. Test admin functionality thoroughly before deployment

## Production Deployment

The admin panel is production-ready with:

- **Docker containerization** - Consistent deployment
- **Health checks** - Automated monitoring
- **Resource limits** - Memory and CPU constraints
- **Graceful shutdowns** - Proper container lifecycle
- **Logging** - Structured JSON logging
- **Monitoring** - Application and system metrics

## Support

For issues related to the admin panel:

1. Check the application logs: `podman compose logs admin-panel`
2. Verify backend connectivity: `curl http://localhost/api/health`
3. Check authentication: Ensure JWT tokens are valid
4. Review Caddy routing: Verify `/admin*` routes are proxying correctly

The admin panel is designed to be a professional, secure, and efficient content management interface for the WebEnable CMS system.
