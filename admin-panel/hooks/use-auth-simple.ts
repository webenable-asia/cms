'use client';

import { useState, useEffect } from 'react';
import { authApi, tokenManager } from '@/lib/api';

export const useAuth = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // Check for existing token on mount
  useEffect(() => {
    const token = tokenManager.getToken();
    if (token) {
      setIsAuthenticated(true);
    }
  }, []);

  const login = async (username: string, password: string) => {
    try {
      const result = await authApi.login({ username, password });
      if (result.token) {
        setIsAuthenticated(true);
        return;
      }
      throw new Error('Login failed');
    } catch (error) {
      console.error('Login error:', error);
      throw new Error('Invalid credentials');
    }
  };

  const logout = () => {
    tokenManager.removeToken();
    setIsAuthenticated(false);
  };

  return {
    isAuthenticated,
    login,
    logout,
  };
};
