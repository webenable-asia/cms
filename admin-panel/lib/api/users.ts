import { alovaInstance } from '../alova';
import type {
  User,
  UsersResponse,
  QueryParams
} from '../types';

// Users API methods (Admin only)
export const usersApi = {
  // Get all users with pagination (admin only)
  getAll: (params?: QueryParams) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', params.page.toString());
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.search) searchParams.set('search', params.search);
    if (params?.sort) searchParams.set('sort', params.sort);
    
    const queryString = searchParams.toString();
    const url = queryString ? `/users?${queryString}` : '/users';
    
    return alovaInstance.Get<UsersResponse>(url);
  },

  // Get single user by ID (admin only)
  getById: (id: string) => 
    alovaInstance.Get<User>(`/users/${id}`),

  // Delete user (admin only)
  delete: (id: string) => 
    alovaInstance.Delete(`/users/${id}`)
};
