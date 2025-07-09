// API Types for the CMS application

export interface User {
  id: string;
  username: string;
  email: string;
  role: 'user'; // Frontend only handles regular users, admin roles handled in admin-panel
  created_at: string;
  updated_at: string;
}

export interface Post {
  id: string;
  title: string;
  content: string;
  author_id: string;
  author?: User;
  published: boolean;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
  message: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface RegisterResponse {
  message: string;
  user: User;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  published?: boolean;
}

export interface UpdatePostRequest {
  title?: string;
  content?: string;
  published?: boolean;
}

export interface ApiError {
  error: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface PostsResponse extends PaginatedResponse<Post> {}
export interface UsersResponse extends PaginatedResponse<User> {}

export interface QueryParams {
  page?: number;
  limit?: number;
  search?: string;
  sort?: string;
}
