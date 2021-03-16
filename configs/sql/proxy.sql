CREATE TABLE requests (
    id serial PRIMARY KEY,
    method text,
    host text,
    path text,
    request text
);