import { useState, useEffect } from 'react';
import { postsApi } from '../lib/api/posts';
import type { Post, PostsResponse, CreatePostRequest, UpdatePostRequest, QueryParams } from '../lib/types';

// Hook for managing posts using Alova
export function usePosts(params?: QueryParams) {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState(0);

  const fetchPosts = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await postsApi.getAll(params).send();
      setPosts(response.items || []);
      setTotal(response.total || 0);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch posts');
      setPosts([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, [params?.page, params?.limit, params?.search, params?.sort]);

  const createPost = async (postData: CreatePostRequest) => {
    try {
      const newPost = await postsApi.create(postData).send();
      setPosts(prev => [newPost, ...prev]);
      return newPost;
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to create post');
    }
  };

  const updatePost = async (id: string, postData: UpdatePostRequest) => {
    try {
      const updatedPost = await postsApi.update(id, postData).send();
      setPosts(prev => prev.map(post => post.id === id ? updatedPost : post));
      return updatedPost;
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to update post');
    }
  };

  const deletePost = async (id: string) => {
    try {
      console.log('Attempting to delete post:', id);
      await postsApi.delete(id).send();
      console.log('Post deleted successfully:', id);
      setPosts(prev => prev.filter(post => post.id !== id));
    } catch (err) {
      console.error('Delete post error:', err);
      throw new Error(err instanceof Error ? err.message : 'Failed to delete post');
    }
  };

  const refreshPosts = () => {
    fetchPosts();
  };

  return {
    posts,
    loading,
    error,
    total,
    createPost,
    updatePost,
    deletePost,
    refreshPosts,
  };
}

// Hook for a single post
export function usePost(id: string) {
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    const fetchPost = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await postsApi.getById(id).send();
        setPost(response);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch post');
        setPost(null);
      } finally {
        setLoading(false);
      }
    };

    fetchPost();
  }, [id]);

  return {
    post,
    loading,
    error,
  };
}
