CREATE TABLE IF NOT EXISTS orders (
    id serial PRIMARY KEY,
    pet_id serial REFERENCES pets(id) ON DELETE CASCADE,
    UNIQUE(pet_id),
    quantity integer NOT NULL,
    ship_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL,
    complete bool NOT NULL
)