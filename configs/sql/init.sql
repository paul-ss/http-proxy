CREATE USER proxy WITH PASSWORD 'jw8s0F4';
CREATE DATABASE proxy_db OWNER proxy;
GRANT ALL PRIVILEGES ON DATABASE proxy_db TO proxy;