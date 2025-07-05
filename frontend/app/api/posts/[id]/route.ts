import { NextRequest, NextResponse } from 'next/server';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:8080';

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;
    
    // Fetch post from backend API
    const response = await fetch(`${BACKEND_URL}/api/posts/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      // Disable caching
      cache: 'no-store',
    });

    if (!response.ok) {
      if (response.status === 404) {
        return NextResponse.json(
          { error: 'Post not found' },
          { status: 404 }
        );
      }
      throw new Error(`Backend responded with ${response.status}`);
    }

    const post = await response.json();
    
    // Transform backend response to match frontend expectations
    const transformedPost = {
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
    };

    return NextResponse.json(transformedPost);
  } catch (error) {
    console.error('Error fetching post:', error);
    return NextResponse.json(
      { error: 'Failed to fetch post' },
      { status: 500 }
    );
  }
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;
    const url = new URL(request.url);
    
    // Handle view count increment
    if (url.pathname.endsWith('/view')) {
      const response = await fetch(`${BACKEND_URL}/api/posts/${id}/view`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`Backend responded with ${response.status}`);
      }

      return NextResponse.json({ success: true });
    }

    return NextResponse.json(
      { error: 'Invalid request' },
      { status: 400 }
    );
  } catch (error) {
    console.error('Error processing request:', error);
    return NextResponse.json(
      { error: 'Failed to process request' },
      { status: 500 }
    );
  }
}
