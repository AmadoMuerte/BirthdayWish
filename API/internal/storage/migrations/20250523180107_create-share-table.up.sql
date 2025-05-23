CREATE TABLE share_wishlist_access (
    id SERIAL PRIMARY KEY,
    created_at timestamp
    with
        time zone NOT NULL DEFAULT now (),
        updated_at timestamp
    with
        time zone,
        user_id integer NOT NULL,
        access_token text NOT NULL
);
