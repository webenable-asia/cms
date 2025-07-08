import { alovaInstance } from '../alova';
import type {
  Post,
  PostsResponse,
  CreatePostRequest,
  UpdatePostRequest,
  QueryParams
} from '../types';

// Posts API methods
export const postsApi = {
  // Get all posts with pagination
  getAll: (params?: QueryParams) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', params.page.toString());
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.search) searchParams.set('search', params.search);
    if (params?.sort) searchParams.set('sort', params.sort);
    
    const queryString = searchParams.toString();
    const url = queryString ? `/posts?${queryString}` : '/posts';
    
    return alovaInstance.Get<PostsResponse>(url, {
      cacheFor: 60 * 1000, // Cache for 1 minute
      meta: {
        authRequired: false // Public endpoint
      }
    });
  },

  // Get single post by ID
  getById: (id: string) => 
    alovaInstance.Get<Post>(`/posts/${id}`, {
      cacheFor: 5 * 60 * 1000, // Cache for 5 minutes
      meta: {
        authRequired: false
      }
    }),

  // Create new post (requires auth)
  create: (postData: CreatePostRequest) => 
    alovaInstance.Post<Post>('/posts', postData),

  // Update existing post (requires auth)
  update: (id: string, postData: UpdatePostRequest) => 
    alovaInstance.Put<Post>(`/posts/${id}`, postData),

  // Delete post (requires auth)
  delete: (id: string) => 
    alovaInstance.Delete(`/posts/${id}`),

  // Get posts by current user (requires auth)
  getMyPosts: (params?: QueryParams) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set('page', params.page.toString());
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.search) searchParams.set('search', params.search);
    
    const queryString = searchParams.toString();
    const url = queryString ? `/posts/my?${queryString}` : '/posts/my';
    
    return alovaInstance.Get<PostsResponse>(url, {
      cacheFor: 30 * 1000 // Cache for 30 seconds
    });
  }
};
