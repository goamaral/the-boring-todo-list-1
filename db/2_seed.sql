--
-- PostgreSQL database dump
--

-- Dumped from database version 14.9 (Debian 14.9-1.pgdg120+1)
-- Dumped by pg_dump version 16.0

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
-- Data for Name: goose_db_version; Type: TABLE DATA; Schema: public; Owner: boring
--

INSERT INTO public.goose_db_version (id, version_id, is_applied, tstamp) VALUES (1, 0, true, '2023-03-06 00:52:32.44439');
INSERT INTO public.goose_db_version (id, version_id, is_applied, tstamp) VALUES (7, 1, true, '2023-09-23 23:33:47.591422');
INSERT INTO public.goose_db_version (id, version_id, is_applied, tstamp) VALUES (8, 2, true, '2023-09-23 23:33:47.622758');


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: boring
--

SELECT pg_catalog.setval('public.goose_db_version_id_seq', 8, true);


--
-- PostgreSQL database dump complete
--

