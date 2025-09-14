CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    content text NOT NULL,
    author text NOT NULL,
    version integer NOT NULL DEFAULT 1
);
