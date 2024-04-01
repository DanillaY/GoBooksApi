
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

CREATE TABLE public.books (
    id bigint NOT NULL,
    current_price bigint,
    old_price bigint,
    title text,
    img_path text,
    page_book_path text,
    vendor text,
    author text,
    translator text,
    production_series text,
    category text,
    publisher text,
    isbn text,
    age_restriction text,
    year_publish text,
    pages_quantity text,
    book_cover text,
    format text,
    weight text,
    book_about text
);


ALTER TABLE public.books OWNER TO postgres;

CREATE SEQUENCE public.books_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.books_id_seq OWNER TO postgres;

ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;

ALTER TABLE ONLY public.books ALTER COLUMN id SET DEFAULT nextval('public.books_id_seq'::regclass);

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (id);
