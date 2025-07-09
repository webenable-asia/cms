import { alovaInstance } from '../alova';
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  User
} from '../types';

// Authentication API methods
export const authApi = {
  // Login user
  login: (credentials: LoginRequest) => 
    alovaInstance.Post<LoginResponse>('/auth/login', credentials, {
      meta: {
        authRole: 'login' // Mark as login request for Alova token authentication
      }
    }),

  // Register new user
  register: (userData: RegisterRequest) => 
    alovaInstance.Post<RegisterResponse>('/auth/register', userData, {
      meta: {
        authRequired: false
      }
    }),

  // Get current user profile
  profile: () => 
    alovaInstance.Get<User>('/auth/profile'),

  // Logout (server-side call + client-side cleanup)
  logout: () => 
    alovaInstance.Post('/auth/logout', {}, {
      meta: {
        authRole: 'logout' // Mark as logout request for Alova token authentication
      }
    }),

  // Check if user is authenticated
  isAuthenticated: () => {
    if (typeof window === 'undefined') return false;
    return !!localStorage.getItem('token');
  },

  // Get stored token
  getToken: () => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('token');
  },

  // Set token in localStorage
  setToken: (token: string) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', token);
    }
  }
};
