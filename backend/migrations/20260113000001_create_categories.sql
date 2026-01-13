-- Migration: Create categories table
-- Requirements: 5.2, 5.3

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    keywords TEXT[] NOT NULL DEFAULT '{}'
);

-- Insert default categories with keywords for automatic categorization
INSERT INTO categories (name, keywords) VALUES
    ('API', ARRAY['api', 'rest', 'graphql', 'endpoint', 'backend', 'server']),
    ('CLI', ARRAY['cli', 'command', 'terminal', 'shell', 'script', 'console']),
    ('Web App', ARRAY['web', 'frontend', 'react', 'vue', 'angular', 'website', 'webapp']),
    ('Mobile', ARRAY['mobile', 'ios', 'android', 'react native', 'flutter', 'app']),
    ('Other', ARRAY[]::TEXT[])
ON CONFLICT (name) DO NOTHING;
