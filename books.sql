--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3
-- Dumped by pg_dump version 15.3

-- Started on 2024-06-05 15:34:05

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

CREATE TABLE public.book_users (
    user_id bigint NOT NULL,
    book_id bigint NOT NULL
);

CREATE TABLE public.books (
    id bigint NOT NULL,
    current_price bigint,
    old_price bigint,
    title text,
    img_path text,
    page_book_path text,
    vendor_url text,
    vendor text,
    author text,
    translator text,
    production_series text,
    category text,
    publisher text,
    isbn text,
    age_restriction text,
    year_publish bigint,
    pages_quantity text,
    book_cover text,
    format text,
    weight text,
    in_stock_text text,
    book_about text
);

CREATE SEQUENCE public.books_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;



CREATE TABLE public.users (
    id bigint NOT NULL,
    email text
);

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

ALTER TABLE ONLY public.books ALTER COLUMN id SET DEFAULT nextval('public.books_id_seq'::regclass);

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


ALTER TABLE ONLY public.book_users
    ADD CONSTRAINT book_users_pkey PRIMARY KEY (user_id, book_id);

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.book_users
    ADD CONSTRAINT fk_book_users_book FOREIGN KEY (book_id) REFERENCES public.books(id);

ALTER TABLE ONLY public.book_users
    ADD CONSTRAINT fk_book_users_user FOREIGN KEY (user_id) REFERENCES public.users(id);
