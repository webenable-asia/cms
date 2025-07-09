import { alovaInstance } from '../alova';

// Enhanced Posts API using Alova
export const alovaPostsApi = {
  // Get all posts using Alova with no caching
  getAll: () => alovaInstance.Get('/api/posts', {
    meta: {
      title: 'Get all posts'
    }
  }),

  // Get single post with no caching
  getById: (id: string) => alovaInstance.Get(`/api/posts/${id}`, {
    meta: {
      title: `Get post ${id}`
    }
  }),

  // Create post
  create: (postData: any) => alovaInstance.Post('/api/posts', postData, {
    meta: {
      title: 'Create post'
    }
  }),

  // Update post  
  update: (id: string, postData: any) => alovaInstance.Put(`/api/posts/${id}`, postData, {
    meta: {
      title: `Update post ${id}`
    }
  }),

  // Delete post
  delete: (id: string) => alovaInstance.Delete(`/api/posts/${id}`, {
    meta: {
      title: `Delete post ${id}`
    }
  })
};
