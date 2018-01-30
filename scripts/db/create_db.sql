-- Database: dfoodie

-- DROP DATABASE dfoodie;

CREATE DATABASE dfoodie
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'C'
    LC_CTYPE = 'C'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE dfoodie
    IS 'Database to store all microservices data';