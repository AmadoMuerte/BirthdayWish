CREATE TABLE IF NOT EXISTS wishlist (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    link text,
    price text,
    name text,
    created_at timestamp
    with
        time zone NOT NULL DEFAULT now (),
        updated_at timestamp
    with
        time zone NOT NULL DEFAULT now (),
        CONSTRAINT fk_wishlist_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_wishlist_user_id ON wishlist (user_id);
