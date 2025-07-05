import { NextRequest, NextResponse } from 'next/server';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:8080';

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const category = searchParams.get('category');
    const limit = searchParams.get('limit');
    const status = searchParams.get('status') || 'published';
    
    // Build query parameters
    const params = new URLSearchParams();
    if (status) params.append('status', status);
    
    // Fetch posts from backend API
    const response = await fetch(`${BACKEND_URL}/api/posts?${params.toString()}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      // Disable caching
      cache: 'no-store',
    });

    if (!response.ok) {
      throw new Error(`Backend responded with ${response.status}`);
    }

    let posts = await response.json();
    
    // Handle null response (empty database)
    if (!posts || posts === null) {
      posts = [];
    }
    
    // Transform backend response to match frontend expectations
    posts = posts.map((post: any) => ({
      _id: post.id || post._id,
      title: post.title,
      content: post.content,
      excerpt: post.excerpt,
      author: post.author,
      created_at: post.created_at || post.createdAt,
      updated_at: post.updated_at || post.updatedAt,
      status: post.status,
      categories: post.categories || [],
      tags: post.tags || [],
      featured_image: post.featured_image ? {
        url: post.featured_image,
        alt: post.image_alt || post.title
      } : undefined,
      meta_title: post.meta_title,
      meta_description: post.meta_description || post.meta_desc,
      reading_time: post.reading_time,
      is_featured: post.is_featured || false,
      view_count: post.view_count || 0,
      scheduled_at: post.scheduled_at,
    }));

    // Filter by category if specified
    if (category) {
      posts = posts.filter((post: any) => 
        post.categories.some((cat: string) => 
          cat.toLowerCase() === category.toLowerCase()
        )
      );
    }

    // Apply limit if specified
    if (limit) {
      const limitNum = parseInt(limit, 10);
      if (!isNaN(limitNum)) {
        posts = posts.slice(0, limitNum);
      }
    }

    // Return in the expected format
    return NextResponse.json({
      posts,
      total: posts.length,
    });
  } catch (error) {
    console.error('Error fetching posts:', error);
    return NextResponse.json(
      { error: 'Failed to fetch posts' },
      { status: 500 }
    );
  }
}
