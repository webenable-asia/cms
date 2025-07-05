'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { CalendarIcon, ClockIcon, EyeIcon, ShareIcon, ArrowLeftIcon } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import { format } from 'date-fns';

interface Post {
  _id: string;
  title: string;
  content: string;
  excerpt: string;
  author: string;
  created_at: string;
  updated_at: string;
  status: 'draft' | 'published' | 'scheduled';
  categories: string[];
  tags: string[];
  featured_image?: {
    url: string;
    alt: string;
  };
  meta_title?: string;
  meta_description?: string;
  reading_time?: number;
  is_featured: boolean;
  view_count: number;
  scheduled_at?: string;
}

export default function BlogPostPage() {
  const params = useParams();
  const [post, setPost] = useState<Post | null>(null);
  const [relatedPosts, setRelatedPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPost = async () => {
      try {
        const response = await fetch(`/api/posts/${params.id}`);
        if (!response.ok) {
          throw new Error('Post not found');
        }
        const data = await response.json();
        setPost(data);

        // Increment view count
        await fetch(`/api/posts/${params.id}/view`, { method: 'POST' });

        // Fetch related posts
        if (data.categories.length > 0) {
          const relatedResponse = await fetch(`/api/posts?category=${data.categories[0]}&limit=3`);
          if (relatedResponse.ok) {
            const relatedData = await relatedResponse.json();
            setRelatedPosts(relatedData.posts.filter((p: Post) => p._id !== data._id));
          }
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setLoading(false);
      }
    };

    if (params.id) {
      fetchPost();
    }
  }, [params.id]);

  const sharePost = () => {
    if (navigator.share && post) {
      navigator.share({
        title: post.title,
        text: post.excerpt,
        url: window.location.href,
      });
    } else {
      // Fallback to copy to clipboard
      navigator.clipboard.writeText(window.location.href);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen py-12">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto">
            <div className="animate-pulse">
              <div className="h-8 bg-gray-200 rounded w-3/4 mb-4"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2 mb-8"></div>
              <div className="h-64 bg-gray-200 rounded mb-8"></div>
              <div className="space-y-4">
                <div className="h-4 bg-gray-200 rounded"></div>
                <div className="h-4 bg-gray-200 rounded w-5/6"></div>
                <div className="h-4 bg-gray-200 rounded w-4/6"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !post) {
    return (
      <div className="min-h-screen py-12">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto text-center">
            <h1 className="text-4xl font-bold mb-4">Post Not Found</h1>
            <p className="text-muted-foreground mb-8">
              {error || 'The post you are looking for does not exist.'}
            </p>
            <Button asChild>
              <Link href="/blog">
                <ArrowLeftIcon className="w-4 h-4 mr-2" />
                Back to Blog
              </Link>
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen py-12">
      <div className="container mx-auto px-4">
        <div className="max-w-4xl mx-auto">
          {/* Back Button */}
          <Button variant="ghost" asChild className="mb-8">
            <Link href="/blog">
              <ArrowLeftIcon className="w-4 h-4 mr-2" />
              Back to Blog
            </Link>
          </Button>

          {/* Post Header */}
          <div className="mb-8">
            {post.is_featured && (
              <Badge variant="secondary" className="mb-4">
                Featured Post
              </Badge>
            )}
            
            <h1 className="text-4xl font-bold mb-4">{post.title}</h1>
            
            <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground mb-6">
              <div className="flex items-center">
                <CalendarIcon className="w-4 h-4 mr-2" />
                {format(new Date(post.created_at), 'MMMM d, yyyy')}
              </div>
              
              {post.reading_time && (
                <div className="flex items-center">
                  <ClockIcon className="w-4 h-4 mr-2" />
                  {post.reading_time} min read
                </div>
              )}
              
              <div className="flex items-center">
                <EyeIcon className="w-4 h-4 mr-2" />
                {post.view_count} views
              </div>
              
              <Button variant="ghost" size="sm" onClick={sharePost}>
                <ShareIcon className="w-4 h-4 mr-2" />
                Share
              </Button>
            </div>

            {/* Categories and Tags */}
            <div className="flex flex-wrap gap-2 mb-6">
              {post.categories.map((category) => (
                <Badge key={category} variant="default">
                  {category}
                </Badge>
              ))}
              {post.tags.map((tag) => (
                <Badge key={tag} variant="outline">
                  #{tag}
                </Badge>
              ))}
            </div>

            {post.excerpt && (
              <p className="text-xl text-muted-foreground leading-relaxed">
                {post.excerpt}
              </p>
            )}
          </div>

          {/* Featured Image */}
          {post.featured_image && (
            <div className="mb-8">
              <Image
                src={post.featured_image.url}
                alt={post.featured_image.alt || post.title}
                width={800}
                height={400}
                className="w-full h-auto rounded-lg"
                priority
              />
            </div>
          )}

          {/* Post Content */}
          <div className="prose prose-lg max-w-none mb-12">
            <ReactMarkdown>{post.content}</ReactMarkdown>
          </div>

          {/* Author Info */}
          <div className="border-t pt-8 mb-12">
            <div className="flex items-center mb-4">
              <div className="w-12 h-12 bg-primary rounded-full flex items-center justify-center text-primary-foreground font-semibold mr-4">
                {post.author.charAt(0).toUpperCase()}
              </div>
              <div>
                <h3 className="font-semibold">{post.author}</h3>
                <p className="text-sm text-muted-foreground">
                  Published on {format(new Date(post.created_at), 'MMMM d, yyyy')}
                </p>
              </div>
            </div>
          </div>

          {/* Related Posts */}
          {relatedPosts.length > 0 && (
            <div>
              <h2 className="text-2xl font-bold mb-6">Related Posts</h2>
              <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
                {relatedPosts.map((relatedPost) => (
                  <Card key={relatedPost._id} className="hover:shadow-md transition-shadow">
                    <CardHeader>
                      <CardTitle className="text-lg line-clamp-2">
                        <Link href={`/blog/${relatedPost._id}`} className="hover:text-primary">
                          {relatedPost.title}
                        </Link>
                      </CardTitle>
                    </CardHeader>
                    <CardContent>
                      <p className="text-sm text-muted-foreground line-clamp-3 mb-4">
                        {relatedPost.excerpt}
                      </p>
                      <div className="flex items-center justify-between text-xs text-muted-foreground">
                        <span>{format(new Date(relatedPost.created_at), 'MMM d, yyyy')}</span>
                        {relatedPost.reading_time && (
                          <span>{relatedPost.reading_time} min read</span>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
