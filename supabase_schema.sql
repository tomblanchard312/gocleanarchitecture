-- Supabase SQL schema for blog_posts table
-- Run this in your Supabase SQL editor to create the required table

CREATE TABLE IF NOT EXISTS blog_posts (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_blog_posts_created_at ON blog_posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_blog_posts_author_id ON blog_posts(author_id);

-- Enable Row Level Security (RLS)
ALTER TABLE blog_posts ENABLE ROW LEVEL SECURITY;

-- RLS Policies: Public can read, but only authors can modify their own posts
CREATE POLICY "Allow public read access" ON blog_posts
    FOR SELECT USING (true);

CREATE POLICY "Users can insert their own posts" ON blog_posts
    FOR INSERT WITH CHECK (auth.uid()::text = author_id);

CREATE POLICY "Users can update their own posts" ON blog_posts
    FOR UPDATE USING (auth.uid()::text = author_id)
    WITH CHECK (auth.uid()::text = author_id);

CREATE POLICY "Users can delete their own posts" ON blog_posts
    FOR DELETE USING (auth.uid()::text = author_id);