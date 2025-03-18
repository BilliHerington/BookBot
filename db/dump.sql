--
-- PostgreSQL database dump
--

-- Dumped from database version 16.4
-- Dumped by pg_dump version 16.4

-- Started on 2025-03-16 15:16:09

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
-- TOC entry 2 (class 3079 OID 16384)
-- Name: adminpack; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS adminpack WITH SCHEMA pg_catalog;


--
-- TOC entry 4876 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION adminpack; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION adminpack IS 'administrative functions for PostgreSQL';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 223 (class 1259 OID 32798)
-- Name: appointments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.appointments (
    appointment_id integer NOT NULL,
    client_name character varying(255) NOT NULL,
    contact_number character varying(20) NOT NULL,
    service_id integer,
    employee_id integer,
    appointment_date date NOT NULL,
    time_start time without time zone NOT NULL,
    time_end time without time zone NOT NULL,
    status character varying(50) NOT NULL
);


ALTER TABLE public.appointments OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 32797)
-- Name: appointments_appointment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.appointments_appointment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.appointments_appointment_id_seq OWNER TO postgres;

--
-- TOC entry 4877 (class 0 OID 0)
-- Dependencies: 222
-- Name: appointments_appointment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.appointments_appointment_id_seq OWNED BY public.appointments.appointment_id;


--
-- TOC entry 217 (class 1259 OID 32770)
-- Name: employees; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.employees (
    employee_id integer NOT NULL,
    name character varying(255) NOT NULL,
    level character varying(20),
    contact_number character varying(20),
    CONSTRAINT employees_level_check CHECK (((level)::text = ANY ((ARRAY['профессионал'::character varying, 'начинающий'::character varying])::text[])))
);


ALTER TABLE public.employees OWNER TO postgres;

--
-- TOC entry 216 (class 1259 OID 32769)
-- Name: employees_employee_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.employees_employee_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.employees_employee_id_seq OWNER TO postgres;

--
-- TOC entry 4878 (class 0 OID 0)
-- Dependencies: 216
-- Name: employees_employee_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.employees_employee_id_seq OWNED BY public.employees.employee_id;


--
-- TOC entry 221 (class 1259 OID 32786)
-- Name: schedule; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schedule (
    schedule_id integer NOT NULL,
    employee_id integer,
    work_date date NOT NULL,
    start_time time without time zone NOT NULL,
    end_time time without time zone NOT NULL
);


ALTER TABLE public.schedule OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 32785)
-- Name: schedule_schedule_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.schedule_schedule_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.schedule_schedule_id_seq OWNER TO postgres;

--
-- TOC entry 4879 (class 0 OID 0)
-- Dependencies: 220
-- Name: schedule_schedule_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.schedule_schedule_id_seq OWNED BY public.schedule.schedule_id;


--
-- TOC entry 219 (class 1259 OID 32777)
-- Name: services; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.services (
    service_id integer NOT NULL,
    name character varying(255) NOT NULL,
    duration interval NOT NULL,
    default_price numeric NOT NULL,
    pro_price numeric NOT NULL,
    description character varying(255)
);


ALTER TABLE public.services OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 32776)
-- Name: services_service_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.services_service_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.services_service_id_seq OWNER TO postgres;

--
-- TOC entry 4880 (class 0 OID 0)
-- Dependencies: 218
-- Name: services_service_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.services_service_id_seq OWNED BY public.services.service_id;


--
-- TOC entry 4707 (class 2604 OID 32801)
-- Name: appointments appointment_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments ALTER COLUMN appointment_id SET DEFAULT nextval('public.appointments_appointment_id_seq'::regclass);


--
-- TOC entry 4704 (class 2604 OID 32773)
-- Name: employees employee_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employees ALTER COLUMN employee_id SET DEFAULT nextval('public.employees_employee_id_seq'::regclass);


--
-- TOC entry 4706 (class 2604 OID 32789)
-- Name: schedule schedule_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schedule ALTER COLUMN schedule_id SET DEFAULT nextval('public.schedule_schedule_id_seq'::regclass);


--
-- TOC entry 4705 (class 2604 OID 32780)
-- Name: services service_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services ALTER COLUMN service_id SET DEFAULT nextval('public.services_service_id_seq'::regclass);


--
-- TOC entry 4870 (class 0 OID 32798)
-- Dependencies: 223
-- Data for Name: appointments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.appointments (appointment_id, client_name, contact_number, service_id, employee_id, appointment_date, time_start, time_end, status) FROM stdin;
7	gg	hh	1	8	2025-03-11	11:00:00	11:30:00	true
8	gg2	hh	1	6	2025-03-11	11:00:00	11:30:00	true
9	gg3	hh	1	8	2025-03-11	15:00:00	15:30:00	true
10	gg4	hh	1	6	2025-03-11	15:00:00	15:30:00	true
\.


--
-- TOC entry 4864 (class 0 OID 32770)
-- Dependencies: 217
-- Data for Name: employees; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.employees (employee_id, name, level, contact_number) FROM stdin;
6	Иван Иванов	профессионал	+79280000006
7	Петр Петров	начинающий	+79280000007
8	Сергей Сергеев	профессионал	+79280000008
9	Анна Аннова	начинающий	+79280000009
10	Елена Еленова	профессионал	+79280000010
\.


--
-- TOC entry 4868 (class 0 OID 32786)
-- Dependencies: 221
-- Data for Name: schedule; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schedule (schedule_id, employee_id, work_date, start_time, end_time) FROM stdin;
12	7	2025-03-11	10:00:00	19:00:00
13	8	2025-03-11	08:00:00	17:00:00
14	9	2025-03-12	09:00:00	18:00:00
11	6	2025-03-11	09:00:00	18:00:00
28	10	2025-05-11	10:00:00	15:00:00
29	10	2024-08-01	09:00:00	20:00:00
30	10	2024-08-03	09:00:00	20:00:00
31	10	2024-08-05	09:00:00	20:00:00
32	10	2024-08-07	09:00:00	20:00:00
33	10	2024-08-09	09:00:00	20:00:00
34	10	2024-08-11	09:00:00	20:00:00
35	10	2024-08-13	09:00:00	20:00:00
36	10	2024-08-15	09:00:00	20:00:00
37	10	2024-08-17	09:00:00	20:00:00
38	10	2024-08-19	09:00:00	20:00:00
39	10	2024-08-21	09:00:00	20:00:00
40	10	2024-08-23	09:00:00	20:00:00
41	10	2024-08-25	09:00:00	20:00:00
42	10	2024-08-27	09:00:00	20:00:00
43	10	2024-08-29	09:00:00	20:00:00
44	10	2024-08-31	09:00:00	20:00:00
45	10	2024-08-02	09:00:00	20:00:00
46	10	2024-08-04	09:00:00	20:00:00
47	10	2024-08-06	09:00:00	20:00:00
48	10	2024-08-08	09:00:00	20:00:00
49	10	2024-08-10	09:00:00	20:00:00
50	10	2024-08-12	09:00:00	20:00:00
51	10	2024-08-14	09:00:00	20:00:00
52	10	2024-08-16	09:00:00	20:00:00
53	10	2024-08-18	09:00:00	20:00:00
54	10	2024-08-20	09:00:00	20:00:00
55	10	2024-08-22	09:00:00	20:00:00
56	10	2024-08-24	09:00:00	20:00:00
57	10	2024-08-26	09:00:00	20:00:00
58	10	2024-08-28	09:00:00	20:00:00
59	10	2024-08-30	09:00:00	20:00:00
60	10	2024-08-01	09:00:00	20:00:00
61	10	2024-08-03	09:00:00	20:00:00
62	10	2024-08-05	09:00:00	20:00:00
63	10	2024-08-07	09:00:00	20:00:00
64	10	2024-08-09	09:00:00	20:00:00
65	10	2024-08-11	09:00:00	20:00:00
66	10	2024-08-13	09:00:00	20:00:00
67	10	2024-08-15	09:00:00	20:00:00
68	10	2024-08-17	09:00:00	20:00:00
69	10	2024-08-19	09:00:00	20:00:00
70	10	2024-08-21	09:00:00	20:00:00
71	10	2024-08-23	09:00:00	20:00:00
72	10	2024-08-25	09:00:00	20:00:00
73	10	2024-08-27	09:00:00	20:00:00
74	10	2024-08-29	09:00:00	20:00:00
75	10	2024-08-31	09:00:00	20:00:00
76	10	2024-08-02	09:00:00	20:00:00
77	10	2024-08-04	09:00:00	20:00:00
78	10	2024-08-06	09:00:00	20:00:00
79	10	2024-08-08	09:00:00	20:00:00
80	10	2024-08-10	09:00:00	20:00:00
81	10	2024-08-12	09:00:00	20:00:00
82	10	2024-08-14	09:00:00	20:00:00
83	10	2024-08-16	09:00:00	20:00:00
84	10	2024-08-18	09:00:00	20:00:00
85	10	2024-08-20	09:00:00	20:00:00
86	10	2024-08-22	09:00:00	20:00:00
87	10	2024-08-24	09:00:00	20:00:00
88	10	2024-08-26	09:00:00	20:00:00
89	10	2024-08-28	09:00:00	20:00:00
90	10	2024-08-30	09:00:00	20:00:00
\.


--
-- TOC entry 4866 (class 0 OID 32777)
-- Dependencies: 219
-- Data for Name: services; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.services (service_id, name, duration, default_price, pro_price, description) FROM stdin;
1	Стрижка	00:30:00	1000	1500	Классическая и модельная стрижка
2	Окрашивание	02:00:00	3000	4500	Профессиональное окрашивание волос
3	Маникюр	01:00:00	2000	3000	Уход за ногтями, покрытие лаком
4	Массаж	01:30:00	2500	3500	Расслабляющий и лечебный массаж
5	Брови	00:45:00	1200	1800	Коррекция и окрашивание бровей
\.


--
-- TOC entry 4881 (class 0 OID 0)
-- Dependencies: 222
-- Name: appointments_appointment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.appointments_appointment_id_seq', 11, true);


--
-- TOC entry 4882 (class 0 OID 0)
-- Dependencies: 216
-- Name: employees_employee_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.employees_employee_id_seq', 12, true);


--
-- TOC entry 4883 (class 0 OID 0)
-- Dependencies: 220
-- Name: schedule_schedule_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.schedule_schedule_id_seq', 90, true);


--
-- TOC entry 4884 (class 0 OID 0)
-- Dependencies: 218
-- Name: services_service_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.services_service_id_seq', 8, true);


--
-- TOC entry 4716 (class 2606 OID 32803)
-- Name: appointments appointments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_pkey PRIMARY KEY (appointment_id);


--
-- TOC entry 4710 (class 2606 OID 32775)
-- Name: employees employees_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (employee_id);


--
-- TOC entry 4714 (class 2606 OID 32791)
-- Name: schedule schedule_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schedule
    ADD CONSTRAINT schedule_pkey PRIMARY KEY (schedule_id);


--
-- TOC entry 4712 (class 2606 OID 32784)
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (service_id);


--
-- TOC entry 4718 (class 2606 OID 32809)
-- Name: appointments appointments_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(employee_id) ON DELETE CASCADE;


--
-- TOC entry 4719 (class 2606 OID 32804)
-- Name: appointments appointments_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(service_id) ON DELETE CASCADE;


--
-- TOC entry 4717 (class 2606 OID 32792)
-- Name: schedule schedule_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schedule
    ADD CONSTRAINT schedule_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(employee_id) ON DELETE CASCADE;


-- Completed on 2025-03-16 15:16:09

--
-- PostgreSQL database dump complete
--

