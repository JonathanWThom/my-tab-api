CREATE TABLE drinks (
  id serial PRIMARY KEY,
  oz NUMERIC NOT NULL,
  percent NUMERIC NOT NULL,
  stddrink NUMERIC NOT NULL,
  imbibed_on TIMESTAMPTZ NOT NULL
  --- NEEDS USER FOREIGN KEY REFERENCE ---
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username text NOT NULL UNIQUE,
    password text NOT NULL,
    uuid uuid NOT NULL DEFAULT uuid_generate_v4()
);

CREATE UNIQUE INDEX users_pkey ON users(id int4_ops);
CREATE UNIQUE INDEX users_username_key ON users(username text_ops);