import { createAlova } from 'alova';
import adapterFetch from 'alova/fetch';
import ReactHook from 'alova/react';
import { tokenManager } from './api';

// Define response type for better TypeScript support
interface ApiResponse {
  ok: boolean;
  status: number;
  statusText: string;
  json(): Promise<any>;
  text(): Promise<string>;
}

// Alova instance for API requests
export const alovaInstance = createAlova({
  baseURL: typeof window !== 'undefined' && window.location.hostname === 'localhost' 
    ? '/api'
    : '/api',
  statesHook: ReactHook,
  requestAdapter: adapterFetch(),
  
  // Disable caching globally for admin panel
  cacheFor: 0,
  
  // Global request headers and token management
  beforeRequest: (method: any) => {
    // Add authentication token if available using shared token manager
    const token = tokenManager.getToken();
    console.log('Alova beforeRequest - URL:', method.url);
    console.log('Alova beforeRequest - Method:', method.type);
    console.log('Alova beforeRequest - Token:', token ? `${token.substring(0, 20)}...` : 'No token');
    
    if (token) {
      method.config.headers = {
        ...method.config.headers,
        'Authorization': `Bearer ${token}`
      };
      console.log('Alova beforeRequest - Authorization header added');
    } else {
      console.log('Alova beforeRequest - No token available, skipping Authorization header');
    }
    
    // Add content-type for non-GET requests
    if (method.type !== 'GET') {
      method.config.headers = {
        ...method.config.headers,
        'Content-Type': 'application/json'
      };
    }
    
    console.log('Alova beforeRequest - Final headers:', method.config.headers);
  },

  responded: {
    // Response interceptor for global error handling
    onSuccess: async (response: ApiResponse, method: any) => {
      console.log('Alova response success - URL:', method.url);
      console.log('Alova response success - Status:', response.status);
      console.log('Alova response success - OK:', response.ok);
      
      if (!response.ok) {
        let errorData;
        try {
          errorData = await response.json();
        } catch (e) {
          errorData = { error: await response.text() || 'Network error' };
        }
        
        console.log('Alova response success - Error data:', errorData);
        
        // If 401, clear the token and redirect to login
        if (response.status === 401) {
          console.log('Alova 401 error - clearing token and redirecting');
          tokenManager.removeToken();
          if (typeof window !== 'undefined') {
            window.location.href = '/admin/auth/login';
          }
        }
        
        const errorMessage = errorData.error || errorData.message || `HTTP ${response.status}: ${response.statusText}`;
        console.log('Alova response success - Throwing error:', errorMessage);
        throw new Error(errorMessage);
      }
      
      // Handle 204 No Content responses
      if (response.status === 204) {
        console.log('Alova response success - 204 No Content, returning empty object');
        return {};
      }
      
      try {
        const jsonData = await response.json();
        console.log('Alova response success - JSON data:', jsonData);
        return jsonData;
      } catch (e) {
        console.log('Alova response success - Failed to parse JSON, returning empty object');
        // If response isn't JSON, return empty object
        return {};
      }
    },
    onError: (error: Error, method: any) => {
      console.error('Alova request error - URL:', method.url);
      console.error('Alova request error - Method:', method.type);
      console.error('Alova request error - Error:', error.message);
      console.error('Alova request error - Full error:', error);
      
      // Check if error contains 401 status
      if (error.message.includes('401') || error.message.includes('Unauthorized')) {
        console.log('Alova error handler - 401 detected, clearing token');
        tokenManager.removeToken();
        if (typeof window !== 'undefined') {
          window.location.href = '/admin/auth/login';
        }
      }
      throw error;
    }
  },
  
  // Global request timeout
  timeout: 10000
});
export default alovaInstance;
