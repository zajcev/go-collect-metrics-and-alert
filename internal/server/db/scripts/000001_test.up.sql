BEGIN;

CREATE TABLE IF NOT EXISTS public.metrics (id varchar NOT NULL, type varchar NOT NULL,delta bigint NULL,value double precision NULL,CONSTRAINT id UNIQUE (id));

COMMIT;