--
-- PostgreSQL database dump
--

-- Dumped from database version 14.4 (Debian 14.4-1.pgdg110+1)
-- Dumped by pg_dump version 15.2

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
INSERT INTO public.goose_db_version (id, version_id, is_applied, tstamp) VALUES (4, 1, true, '2023-05-28 20:50:40.714345');
INSERT INTO public.goose_db_version (id, version_id, is_applied, tstamp) VALUES (6, 2, true, '2023-05-28 20:55:39.297667');


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: boring
--

SELECT pg_catalog.setval('public.goose_db_version_id_seq', 6, true);


--
-- PostgreSQL database dump complete
--

