--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3
-- Dumped by pg_dump version 15.3

-- Started on 2024-06-10 00:54:46

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

--
-- TOC entry 219 (class 1255 OID 87413)
-- Name: books_trigger(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.books_trigger() RETURNS trigger
    LANGUAGE plpgsql
    AS $$ begin new.search := setweight(to_tsvector('simple',coalesce(new.title,'')), 'A') || ' ' || setweight(to_tsvector('simple',coalesce(new.author,'')), 'B') || ' ' || setweight(to_tsvector('simple',coalesce(new.category,'')), 'C'):: tsvector; return new; end $$;


ALTER FUNCTION public.books_trigger() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 218 (class 1259 OID 87224)
-- Name: book_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.book_users (
    user_id bigint NOT NULL,
    book_id bigint NOT NULL
);


ALTER TABLE public.book_users OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 87207)
-- Name: books; Type: TABLE; Schema: public; Owner: postgres
--

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
    book_about text,
    search tsvector
);


ALTER TABLE public.books OWNER TO postgres;

--
-- TOC entry 214 (class 1259 OID 87206)
-- Name: books_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.books_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.books_id_seq OWNER TO postgres;

--
-- TOC entry 3346 (class 0 OID 0)
-- Dependencies: 214
-- Name: books_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;


--
-- TOC entry 217 (class 1259 OID 87216)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    email text
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 216 (class 1259 OID 87215)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 3347 (class 0 OID 0)
-- Dependencies: 216
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 3183 (class 2604 OID 87210)
-- Name: books id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books ALTER COLUMN id SET DEFAULT nextval('public.books_id_seq'::regclass);


--
-- TOC entry 3184 (class 2604 OID 87219)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 3195 (class 2606 OID 87228)
-- Name: book_users book_users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.book_users
    ADD CONSTRAINT book_users_pkey PRIMARY KEY (user_id, book_id);


--
-- TOC entry 3188 (class 2606 OID 87214)
-- Name: books books_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (id);


--
-- TOC entry 3193 (class 2606 OID 87223)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 3185 (class 1259 OID 87409)
-- Name: books_author; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX books_author ON public.books USING gin (to_tsvector('simple'::regconfig, author));


--
-- TOC entry 3186 (class 1259 OID 87410)
-- Name: books_category; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX books_category ON public.books USING gin (to_tsvector('simple'::regconfig, category));


--
-- TOC entry 3189 (class 1259 OID 87412)
-- Name: books_stock; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX books_stock ON public.books USING gin (to_tsvector('simple'::regconfig, in_stock_text));


--
-- TOC entry 3190 (class 1259 OID 87408)
-- Name: books_title; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX books_title ON public.books USING gin (to_tsvector('simple'::regconfig, title));


--
-- TOC entry 3191 (class 1259 OID 87411)
-- Name: books_vendor; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX books_vendor ON public.books USING gin (to_tsvector('simple'::regconfig, vendor));


--
-- TOC entry 3198 (class 2620 OID 87414)
-- Name: books tsvectorupdate; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE ON public.books FOR EACH ROW EXECUTE FUNCTION public.books_trigger();


--
-- TOC entry 3196 (class 2606 OID 87234)
-- Name: book_users fk_book_users_book; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.book_users
    ADD CONSTRAINT fk_book_users_book FOREIGN KEY (book_id) REFERENCES public.books(id);


--
-- TOC entry 3197 (class 2606 OID 87229)
-- Name: book_users fk_book_users_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.book_users
    ADD CONSTRAINT fk_book_users_user FOREIGN KEY (user_id) REFERENCES public.users(id);


-- Completed on 2024-06-10 00:54:46

--
-- PostgreSQL database dump complete
--

