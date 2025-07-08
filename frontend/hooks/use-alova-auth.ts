import { useRequest } from 'alova/client';
import ReactHook from 'alova/react';
import { authApi } from '../lib/api/auth';
import type { User, LoginRequest } from '../lib/types';

// Hook for authentication state management using Alova
export function useAlovaAuth() {
  // Use Alova's useRequest hook for getting user profile
  const {
    data: profileData,
    loading: isLoading,
    error,
    send: checkAuth,
  } = useRequest(authApi.profile(), {
    immediate: false, // Don't auto-load
  });

  const user: User | null = profileData || null;
  const isAuthenticated = !!user;

  // Login function using Alova
  const login = async (username: string, password: string) => {
    try {
      const result = await authApi.login({ username, password }).send();
      
      if (result?.user && result.user.role === 'admin') {
        authApi.setToken(result.token);
        await checkAuth(); // Refresh user data
        return result;
      } else {
        authApi.logout();
        throw new Error('Invalid user role or inactive account');
      }
    } catch (error) {
      authApi.logout();
      throw error;
    }
  };

  // Logout function
  const logout = async () => {
    try {
      await authApi.logout();
    } catch (error) {
      console.warn('Logout failed:', error);
    }
  };

  // Initialize auth check on mount if token exists
  const initializeAuth = async () => {
    if (authApi.isAuthenticated()) {
      await checkAuth();
    }
  };

  return {
    user,
    isAuthenticated,
    isLoading,
    login,
    logout,
    checkAuth,
    initializeAuth,
    error
  };
}

/**
 * Returns the authentication API instance
 */
export function useAuthApi() {
  return authApi;
}
