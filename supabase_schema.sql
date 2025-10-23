-- Supabase SQL schema for blog_posts table
-- Run this in your Supabase SQL editor to create the required table

CREATE TABLE IF NOT EXISTS blog_posts (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create an index on created_at for better performance when ordering
CREATE INDEX IF NOT EXISTS idx_blog_posts_created_at ON blog_posts(created_at DESC);

-- Enable Row Level Security (RLS) - optional but recommended
ALTER TABLE blog_posts ENABLE ROW LEVEL SECURITY;

-- Create policies for RLS (adjust based on your authentication needs)
-- This example allows all operations for authenticated users
CREATE POLICY "Allow all operations for authenticated users" ON blog_posts
    FOR ALL USING (auth.role() = 'authenticated');

-- Alternative: Allow public read access and authenticated write access
-- CREATE POLICY "Allow public read access" ON blog_posts
--     FOR SELECT USING (true);
-- 
-- CREATE POLICY "Allow authenticated write access" ON blog_posts
--     FOR ALL USING (auth.role() = 'authenticated');