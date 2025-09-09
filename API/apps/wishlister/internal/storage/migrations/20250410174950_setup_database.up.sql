CREATE TABLE IF NOT EXISTS wishes (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    link TEXT,
    price DECIMAL(10, 2),
    title TEXT NOT NULL,
    description TEXT,
    image_url TEXT,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wishes_user_id ON wishes(user_id);

CREATE TABLE IF NOT EXISTS share_links (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_share_links_user_id ON share_links(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_share_links_token ON share_links(token);
CREATE INDEX IF NOT EXISTS idx_share_links_expires ON share_links(expires_at);