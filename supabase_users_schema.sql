-- Supabase SQL schema for users table
-- Run this in your Supabase SQL editor to create the required table

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name TEXT,
    bio TEXT,
    avatar_url TEXT,
    role TEXT DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);

-- Enable Row Level Security (RLS) - recommended for security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Create policies for RLS
-- Allow users to read their own profile
CREATE POLICY "Users can read own profile" ON users
    FOR SELECT 
    USING (auth.uid()::text = id);

-- Allow public read access to basic profile info (for public profiles)
CREATE POLICY "Allow public read access to profiles" ON users
    FOR SELECT 
    USING (true);

-- Allow users to update their own profile (but not their role)
CREATE POLICY "Users can update own profile" ON users
    FOR UPDATE 
    USING (auth.uid()::text = id)
    WITH CHECK (auth.uid()::text = id AND role = OLD.role);

-- Allow admins to update any user (including roles)
CREATE POLICY "Admins can update any user" ON users
    FOR UPDATE
    USING ((SELECT role FROM users WHERE id = auth.uid()::text) = 'admin');

-- Allow anyone to create an account (registration)
CREATE POLICY "Allow public registration" ON users
    FOR INSERT 
    WITH CHECK (true);

-- Alternative: If using anon key without Supabase Auth, uncomment below and comment above policies
-- CREATE POLICY "Allow all operations" ON users 
-- FOR ALL 
-- USING (true) 
-- WITH CHECK (true);

-- Add constraint to ensure username is valid (lowercase alphanumeric and underscores)
ALTER TABLE users ADD CONSTRAINT username_format 
    CHECK (username ~ '^[a-z0-9_]{3,30}$');

-- Add constraint to ensure email is valid
ALTER TABLE users ADD CONSTRAINT email_format 
    CHECK (email ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$');

-- Create a function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update updated_at
CREATE TRIGGER users_updated_at_trigger
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_users_updated_at();

-- Optional: Create a view for public profiles (without sensitive data)
CREATE OR REPLACE VIEW public_user_profiles AS
SELECT 
    id,
    username,
    full_name,
    bio,
    avatar_url,
    role,
    created_at
FROM users;

-- Grant select on the view
GRANT SELECT ON public_user_profiles TO anon, authenticated;


