import { createAlova } from 'alova';
import adapterFetch from 'alova/fetch';
import ReactHook from 'alova/react';

// Define response type for better TypeScript support
interface ApiResponse {
  ok: boolean;
  status: number;
  statusText: string;
  json(): Promise<any>;
}

// Alova instance for API requests
export const alovaInstance = createAlova({
  baseURL: process.env.NODE_ENV === 'production' 
    ? process.env.NEXT_PUBLIC_API_URL || '/api'
    : '/api',
  statesHook: ReactHook,
  requestAdapter: adapterFetch(),
  responded: {
    // Response interceptor for global error handling
    onSuccess: async (response: ApiResponse, method) => {
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Network error' }));
        throw new Error(errorData.error || `HTTP ${response.status}: ${response.statusText}`);
      }
      return response.json();
    },
    onError: (error: Error, method) => {
      console.error('Alova request error:', error.message);
      throw error;
    }
  },
  // Global request timeout
  timeout: 10000,
  // Global request headers
  beforeRequest: method => {
    // Add authentication token if available
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    if (token) {
      method.config.headers = {
        ...method.config.headers,
        'Authorization': `Bearer ${token}`
      };
    }
    
    // Add content-type for non-GET requests
    if (method.type !== 'GET') {
      method.config.headers = {
        ...method.config.headers,
        'Content-Type': 'application/json'
      };
    }
  }
});
export default alovaInstance;
