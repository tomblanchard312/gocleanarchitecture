-- Supabase SQL schema for comments table
-- Run this in your Supabase SQL editor to create the required table

CREATE TABLE IF NOT EXISTS comments (
    id TEXT PRIMARY KEY,
    blog_post_id TEXT NOT NULL REFERENCES blog_posts(id) ON DELETE CASCADE,
    author_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    parent_id TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_comments_blog_post_id ON comments(blog_post_id);
CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments(author_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);

-- Enable Row Level Security (RLS)
ALTER TABLE comments ENABLE ROW LEVEL SECURITY;

-- Create policies for RLS

-- Allow public read access to comments
CREATE POLICY "Allow public read access" ON comments
    FOR SELECT 
    USING (true);

-- Users can insert their own comments
CREATE POLICY "Users can insert their own comments" ON comments
    FOR INSERT 
    WITH CHECK (auth.uid()::text = author_id);

-- Users can update their own comments
CREATE POLICY "Users can update their own comments" ON comments
    FOR UPDATE 
    USING (auth.uid()::text = author_id)
    WITH CHECK (auth.uid()::text = author_id);

-- Users can delete their own comments
CREATE POLICY "Users can delete their own comments" ON comments
    FOR DELETE 
    USING (auth.uid()::text = author_id);

-- Admins can delete any comment
CREATE POLICY "Admins can delete any comment" ON comments
    FOR DELETE
    USING ((SELECT role FROM users WHERE id = auth.uid()::text) = 'admin');

-- Create a function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_comments_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update updated_at
CREATE TRIGGER comments_updated_at_trigger
    BEFORE UPDATE ON comments
    FOR EACH ROW
    EXECUTE FUNCTION update_comments_updated_at();

-- Alternative: If using anon key without Supabase Auth, uncomment below and comment above policies
-- CREATE POLICY "Allow all operations" ON comments 
-- FOR ALL 
-- USING (true) 
-- WITH CHECK (true);

