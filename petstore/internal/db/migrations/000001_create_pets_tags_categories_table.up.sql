CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name  text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    deleted bool NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS tags (
    id serial PRIMARY KEY,
    name VARCHAR(255)  NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
    id serial PRIMARY KEY,
    name VARCHAR(255)  NOT NULL
);

CREATE TABLE IF NOT EXISTS pets (
    id serial PRIMARY KEY,
    category_id serial REFERENCES categories(id),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    photo_urls TEXT[]
);

CREATE TABLE IF NOT EXISTS pet_tags (
  pet_id serial REFERENCES pets(id) ON DELETE CASCADE,
  tag_id serial REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (pet_id, tag_id)
);


