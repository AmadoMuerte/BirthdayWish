CREATE TABLE users (
    id text PRIMARY KEY,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone,
    email text NOT NULL,
    name text NOT NULL,
    age integer,
    gender text
);
