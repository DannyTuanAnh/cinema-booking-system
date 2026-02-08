--
-- PostgreSQL database dump
--


-- Dumped from database version 16.11
-- Dumped by pg_dump version 16.11

-- Started on 2026-02-08 09:45:31

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
-- TOC entry 233 (class 1255 OID 16983)
-- Name: book_seats(bigint, bigint[]); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.book_seats(p_user_id bigint, p_seat_ids bigint[]) RETURNS void
    LANGUAGE plpgsql
    AS $$

-- Khai báo biến để đếm số ghế update được
declare p_updated_count int;

begin

	-- Khóa ghế trước
	perform 1
	from seats
	where seat_id = any(p_seat_ids)
	for update;

	-- Update trạng thái ghế thành booked
	update seats
	set status = 'booked'
	where seat_id = any(p_seat_ids) and status = 'available';

	-- Đếm số ghế update thành công
	get diagnostics p_updated_count = row_count;

	-- So sánh số ghế update thành công với số ghế yêu cầu booked có bằng nhau không
	if p_updated_count <> array_length(p_seat_ids, 1) then
		raise exception 'some seats already booked';
	end if;

	-- cấp nhật số ghế booked thành công vào bảng booking
	insert into bookings(user_id, seat_id)
	select p_user_id, unnest(p_seat_ids);

end;
$$;



SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 226 (class 1259 OID 16483)
-- Name: api_keys; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.api_keys (
    id bigint NOT NULL,
    client_id character varying(50) NOT NULL,
    key_hash character(255) NOT NULL,
    description text,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    revoked_at timestamp with time zone,
    last_used_at timestamp with time zone,
    rate_limit integer,
    rate_window_sec integer DEFAULT 1 NOT NULL
);




--
-- TOC entry 225 (class 1259 OID 16482)
-- Name: api_keys_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.api_keys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



--
-- TOC entry 4970 (class 0 OID 0)
-- Dependencies: 225
-- Name: api_keys_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.api_keys_id_seq OWNED BY public.api_keys.id;


--
-- TOC entry 222 (class 1259 OID 16445)
-- Name: bookings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.bookings (
    booking_id integer NOT NULL,
    user_id integer NOT NULL,
    seat_id integer NOT NULL,
    bookat timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);



--
-- TOC entry 221 (class 1259 OID 16444)
-- Name: bookings_booking_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.bookings_booking_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



--
-- TOC entry 4971 (class 0 OID 0)
-- Dependencies: 221
-- Name: bookings_booking_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.bookings_booking_id_seq OWNED BY public.bookings.booking_id;


--
-- TOC entry 224 (class 1259 OID 16465)
-- Name: movies; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.movies (
    movie_id integer NOT NULL,
    title character varying(100) NOT NULL,
    duration integer,
    description text,
    url_image character varying(255) DEFAULT 'ok.png'::character varying NOT NULL,
    rate numeric(3,1),
    genre character varying(200),
    release_date date,
    director character varying(50),
    cast_list character varying(100),
    CONSTRAINT movies_rate_check CHECK (((rate >= (0)::numeric) AND (rate <= (10)::numeric)))
);


--
-- TOC entry 223 (class 1259 OID 16464)
-- Name: movies_movie_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.movies_movie_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



--
-- TOC entry 4972 (class 0 OID 0)
-- Dependencies: 223
-- Name: movies_movie_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.movies_movie_id_seq OWNED BY public.movies.movie_id;


--
-- TOC entry 228 (class 1259 OID 16505)
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.refresh_tokens (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    revoked boolean DEFAULT false,
    replaced_by_hash text,
    created_at timestamp without time zone DEFAULT now()
);



--
-- TOC entry 227 (class 1259 OID 16504)
-- Name: refresh_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.refresh_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



--
-- TOC entry 4973 (class 0 OID 0)
-- Dependencies: 227
-- Name: refresh_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;


--
-- TOC entry 220 (class 1259 OID 16429)
-- Name: seats; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.seats (
    seat_id integer NOT NULL,
    show_id integer NOT NULL,
    seat_name character varying(10) NOT NULL,
    status character varying(10) NOT NULL,
    CONSTRAINT seats_status_check CHECK (((status)::text = ANY ((ARRAY['available'::character varying, 'booked'::character varying])::text[])))
);



--
-- TOC entry 219 (class 1259 OID 16428)
-- Name: seats_seat_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.seats_seat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



--
-- TOC entry 4974 (class 0 OID 0)
-- Dependencies: 219
-- Name: seats_seat_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.seats_seat_id_seq OWNED BY public.seats.seat_id;


--
-- TOC entry 218 (class 1259 OID 16417)
-- Name: shows; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.shows (
    show_id integer NOT NULL,
    movie_id integer NOT NULL,
    show_time timestamp with time zone NOT NULL
);


--
-- TOC entry 217 (class 1259 OID 16416)
-- Name: shows_show_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.shows_show_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

--
-- TOC entry 4975 (class 0 OID 0)
-- Dependencies: 217
-- Name: shows_show_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.shows_show_id_seq OWNED BY public.shows.show_id;


--
-- TOC entry 216 (class 1259 OID 16399)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    user_id integer NOT NULL,
    email character varying(50) NOT NULL,
    password character varying(255) NOT NULL,
    full_name character(50)
);



--
-- TOC entry 215 (class 1259 OID 16398)
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;




--
-- TOC entry 4976 (class 0 OID 0)
-- Dependencies: 215
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_user_id_seq OWNED BY public.users.user_id;


--
-- TOC entry 4773 (class 2604 OID 16486)
-- Name: api_keys id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_keys ALTER COLUMN id SET DEFAULT nextval('public.api_keys_id_seq'::regclass);


--
-- TOC entry 4769 (class 2604 OID 16448)
-- Name: bookings booking_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings ALTER COLUMN booking_id SET DEFAULT nextval('public.bookings_booking_id_seq'::regclass);


--
-- TOC entry 4771 (class 2604 OID 16468)
-- Name: movies movie_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.movies ALTER COLUMN movie_id SET DEFAULT nextval('public.movies_movie_id_seq'::regclass);


--
-- TOC entry 4777 (class 2604 OID 16508)
-- Name: refresh_tokens id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);


--
-- TOC entry 4768 (class 2604 OID 16432)
-- Name: seats seat_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.seats ALTER COLUMN seat_id SET DEFAULT nextval('public.seats_seat_id_seq'::regclass);


--
-- TOC entry 4767 (class 2604 OID 16420)
-- Name: shows show_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.shows ALTER COLUMN show_id SET DEFAULT nextval('public.shows_show_id_seq'::regclass);


--
-- TOC entry 4766 (class 2604 OID 16402)
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_user_id_seq'::regclass);


--
-- TOC entry 4962 (class 0 OID 16483)
-- Dependencies: 226
-- Data for Name: api_keys; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.api_keys (id, client_id, key_hash, description, is_active, created_at, revoked_at, last_used_at, rate_limit, rate_window_sec) FROM stdin;
12	web	d1d5c2b96b7dc4dfc4a2338dc43df81c2762209280faa7c09690c05830e79a38                                                                                                                                                                                               	\N	t	2026-01-16 21:19:26.892688+07	\N	\N	100	1
\.


--
-- TOC entry 4958 (class 0 OID 16445)
-- Dependencies: 222
-- Data for Name: bookings; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.bookings (booking_id, user_id, seat_id, bookat) FROM stdin;
\.


--
-- TOC entry 4960 (class 0 OID 16465)
-- Dependencies: 224
-- Data for Name: movies; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.movies (movie_id, title, duration, description, url_image, rate, genre, release_date, director, cast_list) FROM stdin;
6	Phi Vụ Động Trời 2	145	Ethan Hunt và đội ngũ của anh tiếp tục cuộc chiến chống lại những kẻ thù nguy hiểm.	/assets/images/film/zootopia.jpg	7.8	Hành động, Phiêu lưu	2024-11-20	Christopher McQuarrie	Tom Cruise, Miles Teller
7	Thế Hệ Kỳ Tích	138	Câu chuyện cảm động về một thế hệ trẻ và những giấc mơ của họ.	/assets/images/film/the-he-ki-tich.jpg	8.2	Tâm lý, Chính kịch	2024-10-10	Various	Vietnamese Actors
8	Chân Trời Rực Rỡ	85	Một cuộc hành trình tài liệu khám phá những kỳ tích của thiên nhiên.	/assets/images/film/ctrr.jpg	8.0	Tài liệu	2024-12-01	Documentary Team	Various
9	Anh Trai Tôi Là Khủng Long	120	Một bộ phim giả tưởng hài hước về anh trai là một chú khủng long.	/assets/images/film/anh-trai-toi-la-khung-long.jpg	7.9	Giả tưởng, Hành động	2024-11-15	Vietnamese Director	Vietnamese Actors
10	Kumanthong Nhật Bản: Vong Nhi Cúp Bế	156	Một bộ phim kinh dị với những yếu tố tâm linh từ các nền văn hóa Á Đông.	/assets/images/film/kumathong-japan.jpg	7.5	Kinh dị, Tâm linh	2024-12-05	Horror Master	Asian Actors
5	Avatar: Lửa Và Tro Tàn	197	Tiếp tục cuộc phiêu lưu trên hành tinh Pandora, Jake Sully và nhóm của anh ta phải đối mặt với những thách thức mới.	/assets/images/film/avatar.jpg	8.5	Giả tưởng, Hành động	2024-12-15	James Cameron	Sam Worthington, Zoe Saldana
\.


--
-- TOC entry 4964 (class 0 OID 16505)
-- Dependencies: 228
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.refresh_tokens (id, user_id, token_hash, expires_at, revoked, replaced_by_hash, created_at) FROM stdin;
50	4	960b8651e8731ffda092d919973cb7d89a5aa775fd69946b5b6a500b7a1fffa1	2026-02-01 13:02:39.364102	t	01990ece547e8606d0357bdb8c29193fda41e65e29872106d96ee9cc90a0ba4d	2026-02-01 13:01:39.368462
51	4	01990ece547e8606d0357bdb8c29193fda41e65e29872106d96ee9cc90a0ba4d	2026-02-01 13:02:54.29041	f	\N	2026-02-01 13:01:54.291438
52	4	f36c9acc6ca6008cd050c74c857a7f8b3fb83177d268e2f9d211d2636c17aacf	2026-02-01 13:04:04.800953	f	\N	2026-02-01 13:03:04.802686
53	6	25591e33097f630ba5799f0f68cbe24158a5a5bcb983c7b11d38b0012c1f14c4	2026-02-01 13:04:27.480513	f	\N	2026-02-01 13:03:27.481269
54	4	6b541ab927fcd832a0d3364ff54264bd908f39f6246e6f8602f402ad624e99de	2026-02-01 13:06:00.716072	t	b7e51ccd367e825beebfb7f12265e9947a441862a6ec3e9ce7210f8f7a6944a4	2026-02-01 13:05:00.716672
56	4	b7e51ccd367e825beebfb7f12265e9947a441862a6ec3e9ce7210f8f7a6944a4	2026-02-01 13:06:36.894929	f	\N	2026-02-01 13:05:36.895388
55	6	27c0530a96554963614d74cff5aded23e6b134b8629319898e9ebd570f90854b	2026-02-01 13:06:20.324224	t	ae1cef3e631e34d8eab2ff59bfbbf928a914d6b4f2b55a2d6ad479ff7ce1e8d4	2026-02-01 13:05:20.32545
57	6	ae1cef3e631e34d8eab2ff59bfbbf928a914d6b4f2b55a2d6ad479ff7ce1e8d4	2026-02-01 13:06:38.703829	f	\N	2026-02-01 13:05:38.704261
58	4	f31238f8d62a8e70a276f465085fc64aeab847d4a2a0276c539f2caeb55c7129	2026-02-01 15:50:15.646965	f	\N	2026-02-01 15:49:15.648987
59	4	9b7b463e7e6678eafc04378536e52e7b860059cb3aeefb6bf16ff344c97f744b	2026-02-01 15:50:29.740963	f	\N	2026-02-01 15:49:29.74209
60	4	c7b15a7ae04bda6b66d889ed66e1b43c932922619bc93a35092ce99c5fbe862f	2026-02-01 16:06:49.49759	f	\N	2026-02-01 16:05:49.498434
61	4	e41f24c68c72906977bb83b371e97b374ace71bed63128075e2e4842c841eb03	2026-02-01 16:25:52.943393	f	\N	2026-02-01 16:24:52.945805
62	4	427648c357c0fd7f5ccfb2a06f528e17ee318d91c7f951c297497b23eeb1a795	2026-02-01 16:31:28.647205	f	\N	2026-02-01 16:30:28.648349
63	4	0ef7d59ea1eaff0f7190b2e1d406e7d1daec1934f2f1cca411d5ac852ccf4bb8	2026-02-01 17:54:26.803202	f	\N	2026-02-01 17:53:26.804106
64	4	d42471271bf448c410985268e7d2ce22a4f1dc55d299a7fc52c1d120fa769244	2026-02-01 17:54:48.197393	f	\N	2026-02-01 17:53:48.200344
65	4	ebd2841d72c7e1325df835adccdad71dceed5e5415658447d79445c0208a8767	2026-02-01 17:56:25.716767	f	\N	2026-02-01 17:55:25.717742
66	4	d04814d73fe89e263bdbe16073e399e53a819701c38785c9fc522cd011e75ac5	2026-02-01 17:57:08.494943	f	\N	2026-02-01 17:56:08.498458
67	4	fb0a7413a2854e34d2c9ba8974a17eb44bc5fc33f93d5639bff27210d1571bb7	2026-02-01 17:57:38.060232	f	\N	2026-02-01 17:56:38.060872
68	4	5fee536b55a1bbb3a967712b58a26d76233ea43965b80478493cf23dd2597130	2026-02-01 18:00:20.894205	f	\N	2026-02-01 17:59:20.895015
69	4	3d5c21824658b105598bae39189721cd47b1ac05230be591e99f5cb32e8f23c2	2026-02-01 18:02:22.700122	f	\N	2026-02-01 18:01:22.701375
70	4	18074961c34dc326f59dbcd4b16b162d761660c230314e193315439e1bb130b8	2026-02-01 18:03:40.815251	f	\N	2026-02-01 18:02:40.816217
71	4	b05f7970a22f5f6d70999e3a894662da2c3b96ec40b09802263639efdb2e7a58	2026-02-01 18:08:56.632857	f	\N	2026-02-01 18:07:56.634327
72	4	45f760076a0ac6022430b78dda1c6203737617e74933baa3fdc1dfb96abdb699	2026-02-01 18:10:47.518486	f	\N	2026-02-01 18:09:47.519504
73	6	10c0527f8bcd554faa54da5ffa91264dafc10db2cae31ecae402337480fd688f	2026-02-01 18:11:33.473454	f	\N	2026-02-01 18:10:33.475179
74	4	f4a5cfd6d6ffb334c9768a87a438fc017ba4c0d22770e117ebceb573fd04585a	2026-02-01 18:12:34.028581	f	\N	2026-02-01 18:11:34.029843
75	4	fc83162c09445599900cd8971f6d5dd0f835c605895269046d861bc525394656	2026-02-01 18:18:02.605741	f	\N	2026-02-01 18:17:02.607883
76	6	1287bfa13859836d046873f8ce68e7a837b609aa9183bd9fdb4beb6a04dc49e7	2026-02-01 18:18:22.840982	f	\N	2026-02-01 18:17:22.84224
77	4	6fb1606f5ec8d129950121e33ffa13fd6e03111cc9ca71c97ef1eed0edc248fb	2026-02-01 18:18:34.259288	f	\N	2026-02-01 18:17:34.260737
78	6	39e774df87c6aa92e6d3344b8fb6886f845df1032182896d42f67bd2c844c01d	2026-02-01 18:26:58.944724	t	76beac3f84d5d0dbeeb76971c25188168986317aa4e735d776bd9d57cb3eb0fa	2026-02-01 18:25:58.946349
80	6	76beac3f84d5d0dbeeb76971c25188168986317aa4e735d776bd9d57cb3eb0fa	2026-02-01 18:27:13.693893	f	\N	2026-02-01 18:26:13.694425
79	4	4df19e1b39bb095a8d45655131c77e7d572bb9e6bc50a8407f16864ba9f93e97	2026-02-01 18:27:11.114646	t	7573ea13438639d241469bbdab4d3caa798d9a5f96b4154638df4c2eccdff559	2026-02-01 18:26:11.115693
81	4	7573ea13438639d241469bbdab4d3caa798d9a5f96b4154638df4c2eccdff559	2026-02-01 18:27:33.562268	t	22a953c2c84cf7c56551192941e40e011ea2990399be072a2e43fd68bff09591	2026-02-01 18:26:33.562739
82	4	22a953c2c84cf7c56551192941e40e011ea2990399be072a2e43fd68bff09591	2026-02-01 18:28:24.484148	f	\N	2026-02-01 18:27:24.484765
83	7	87781ee1e2948171d9ef462ae012ae1c8f75b7e269b5a5b4a7b6400da8b365a4	2026-02-01 18:37:10.207034	f	\N	2026-02-01 18:36:10.207895
84	7	7d23f344a460d21540fc7bc62040fab92020d891007e02d921aa6ba0e43e9d09	2026-02-01 18:39:39.184816	f	\N	2026-02-01 18:38:39.185801
85	7	036525fe207bf05ea22340f5760f09979e716bdb195096e8f761ba60afbe9413	2026-02-01 18:41:31.316826	f	\N	2026-02-01 18:40:31.318608
86	4	b51b5f00af2116554f4a2d498b6c01b6356d215f15d951a1c26fdb7bd8d7482b	2026-02-01 18:45:03.03865	f	\N	2026-02-01 18:44:03.039334
87	4	d9de0cde0563d2e79fe301a93417e4f90eab8870ea844a7fcfdf445b38c8c980	2026-02-04 19:15:03.108873	f	\N	2026-02-04 19:14:03.112114
88	4	ad8388c0124d753a2d7bfc1697086167d8def08c24e81d4a082c424606ce4abc	2026-02-04 19:17:05.730433	f	\N	2026-02-04 19:16:05.731253
89	4	b5b11d9d7af2ec1422858c3e0ee0aeed91b6f42245801833198b265bf6a9df1c	2026-02-04 19:19:01.595583	f	\N	2026-02-04 19:18:01.596155
90	4	6638464490d9ab805e50480d1e7520b9a57f15be26cc5536f7f00882a605d74d	2026-02-04 19:20:02.252209	f	\N	2026-02-04 19:19:02.253526
91	4	2f1bc4fac25659e3241c3603821243866e0cb505c0e933b25b330137a9bd7965	2026-02-04 19:21:27.197293	f	\N	2026-02-04 19:20:27.197859
92	4	ef50a8c7d5c16bdaeda9c7c946c51479757a4222637d667b40afc40ab6b07822	2026-02-04 19:23:50.334417	f	\N	2026-02-04 19:22:50.335653
93	4	f42ff1598ff65b168eb5d446d0a8ee2d7ce3539abea83596253ea43fe6608720	2026-02-04 19:25:01.13236	f	\N	2026-02-04 19:24:01.132853
94	4	eea0108ce3dba1d7bd7eaf59ffeb9ecc4f1c0f9a44612cfb6c49fa0e5e35f006	2026-02-04 19:30:30.867974	f	\N	2026-02-04 19:29:30.869574
95	4	4744fa0365c39061a44017c19034760836d4b7a33396012a8a218b8303aac28e	2026-02-04 19:34:49.291183	f	\N	2026-02-04 19:33:49.293353
96	4	42fc4a05e35874283d717258441c21da1dde500c3aeeddbe6f27271d1d7d900e	2026-02-04 19:40:39.361398	f	\N	2026-02-04 19:39:39.361978
97	4	b1b0a3a3fc49dd3391e7724fe8bb624a1303b13590696fd9b2cbb83c035c3dbc	2026-02-04 19:41:15.351293	f	\N	2026-02-04 19:40:15.352763
98	4	feb788a0ae9b2f37364b5a066e15f9d96560b1ab43d58ab4f588641316a5d6fd	2026-02-04 19:41:59.283159	f	\N	2026-02-04 19:40:59.284272
99	4	b978568560d33d75df6c6ad18c9ffbbdf658595a6a8f385ed470641e6cfaca01	2026-02-04 19:46:42.627927	f	\N	2026-02-04 19:45:42.629814
100	4	103b51d8e1a4d281d8d6540ffe62b090a6b4f0155f297bfcc32c2afca2f1fbfc	2026-02-04 19:52:09.148012	f	\N	2026-02-04 19:51:09.149667
\.


--
-- TOC entry 4956 (class 0 OID 16429)
-- Dependencies: 220
-- Data for Name: seats; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.seats (seat_id, show_id, seat_name, status) FROM stdin;
14	7	A7	available
15	7	B7	available
16	7	A8	available
17	7	B8	available
22	8	A1	available
23	8	B1	available
24	8	A2	available
25	8	B2	available
26	8	A3	available
27	8	B3	available
28	8	A4	available
30	8	A5	available
31	8	B5	available
32	8	A6	available
33	8	B6	available
34	8	A7	available
35	8	B7	available
36	8	A8	available
37	8	B8	available
43	9	B1	available
44	9	A2	available
45	9	B2	available
46	9	A3	available
47	9	B3	available
48	9	A4	available
49	9	B4	available
50	9	A5	available
51	9	B5	available
52	9	A6	available
53	9	B6	available
54	9	A7	available
55	9	B7	available
56	9	A8	available
57	9	B8	available
58	9	A9	available
59	9	B9	available
60	9	A10	available
61	9	B10	available
11	7	B5	available
13	7	B6	available
29	8	B4	available
4	7	A2	available
9	7	B4	available
12	7	A6	available
42	9	A1	available
5	7	B2	available
8	7	A4	available
10	7	A5	available
6	7	A3	available
7	7	B3	available
2	7	A1	available
3	7	B1	available
62	10	A1	available
63	10	A2	available
64	10	A3	available
65	10	A4	available
66	10	A5	available
67	10	A6	available
68	10	A7	available
69	10	A8	available
70	10	B1	available
71	10	B2	available
72	10	B3	available
73	10	B4	available
74	10	B5	available
75	10	B6	available
76	10	B7	available
77	10	B8	available
78	10	C1	available
79	10	C2	available
80	10	C3	available
81	10	C4	available
82	10	C5	available
83	10	C6	available
84	10	C7	available
85	10	C8	available
86	10	D1	available
87	10	D2	available
88	10	D3	available
89	10	D4	available
90	10	D5	available
91	10	D6	available
92	10	D7	available
93	10	D8	available
94	10	E1	available
95	10	E2	available
96	10	E3	available
97	10	E4	available
98	10	E5	available
99	10	E6	available
100	10	E7	available
101	10	E8	available
102	11	A1	available
103	11	A2	available
104	11	A3	available
105	11	A4	available
106	11	A5	available
107	11	A6	available
108	11	A7	available
109	11	A8	available
110	11	B1	available
111	11	B2	available
112	11	B3	available
113	11	B4	available
114	11	B5	available
115	11	B6	available
116	11	B7	available
117	11	B8	available
118	11	C1	available
119	11	C2	available
120	11	C3	available
121	11	C4	available
122	11	C5	available
123	11	C6	available
124	11	C7	available
125	11	C8	available
126	11	D1	available
127	11	D2	available
128	11	D3	available
129	11	D4	available
130	11	D5	available
131	11	D6	available
132	11	D7	available
133	11	D8	available
134	11	E1	available
135	11	E2	available
136	11	E3	available
137	11	E4	available
138	11	E5	available
139	11	E6	available
140	11	E7	available
141	11	E8	available
142	12	A1	available
143	12	A2	available
144	12	A3	available
145	12	A4	available
146	12	A5	available
147	12	A6	available
148	12	A7	available
149	12	A8	available
150	12	B1	available
151	12	B2	available
152	12	B3	available
153	12	B4	available
154	12	B5	available
155	12	B6	available
156	12	B7	available
157	12	B8	available
158	12	C1	available
159	12	C2	available
160	12	C3	available
161	12	C4	available
162	12	C5	available
163	12	C6	available
164	12	C7	available
165	12	C8	available
166	12	D1	available
167	12	D2	available
168	12	D3	available
169	12	D4	available
170	12	D5	available
171	12	D6	available
172	12	D7	available
173	12	D8	available
174	12	E1	available
175	12	E2	available
176	12	E3	available
177	12	E4	available
178	12	E5	available
179	12	E6	available
180	12	E7	available
181	12	E8	available
182	13	A1	available
183	13	A2	available
184	13	A3	available
185	13	A4	available
186	13	A5	available
187	13	A6	available
188	13	A7	available
189	13	A8	available
190	13	B1	available
191	13	B2	available
192	13	B3	available
193	13	B4	available
194	13	B5	available
195	13	B6	available
196	13	B7	available
197	13	B8	available
198	13	C1	available
199	13	C2	available
200	13	C3	available
201	13	C4	available
202	13	C5	available
203	13	C6	available
204	13	C7	available
205	13	C8	available
206	13	D1	available
207	13	D2	available
208	13	D3	available
209	13	D4	available
210	13	D5	available
211	13	D6	available
212	13	D7	available
213	13	D8	available
214	13	E1	available
215	13	E2	available
216	13	E3	available
217	13	E4	available
218	13	E5	available
219	13	E6	available
220	13	E7	available
221	13	E8	available
222	14	A1	available
223	14	A2	available
224	14	A3	available
225	14	A4	available
226	14	A5	available
227	14	A6	available
228	14	A7	available
229	14	A8	available
230	14	B1	available
231	14	B2	available
232	14	B3	available
233	14	B4	available
234	14	B5	available
235	14	B6	available
236	14	B7	available
237	14	B8	available
238	14	C1	available
239	14	C2	available
240	14	C3	available
241	14	C4	available
242	14	C5	available
243	14	C6	available
244	14	C7	available
245	14	C8	available
246	14	D1	available
247	14	D2	available
248	14	D3	available
249	14	D4	available
250	14	D5	available
251	14	D6	available
252	14	D7	available
253	14	D8	available
254	14	E1	available
255	14	E2	available
256	14	E3	available
257	14	E4	available
258	14	E5	available
259	14	E6	available
260	14	E7	available
261	14	E8	available
262	15	A1	available
263	15	A2	available
264	15	A3	available
265	15	A4	available
266	15	A5	available
267	15	A6	available
268	15	A7	available
269	15	A8	available
270	15	B1	available
271	15	B2	available
272	15	B3	available
273	15	B4	available
274	15	B5	available
275	15	B6	available
276	15	B7	available
277	15	B8	available
278	15	C1	available
279	15	C2	available
280	15	C3	available
281	15	C4	available
282	15	C5	available
283	15	C6	available
284	15	C7	available
285	15	C8	available
286	15	D1	available
287	15	D2	available
288	15	D3	available
289	15	D4	available
290	15	D5	available
291	15	D6	available
292	15	D7	available
293	15	D8	available
294	15	E1	available
295	15	E2	available
296	15	E3	available
297	15	E4	available
298	15	E5	available
299	15	E6	available
300	15	E7	available
301	15	E8	available
302	16	A1	available
303	16	A2	available
304	16	A3	available
305	16	A4	available
306	16	A5	available
307	16	A6	available
308	16	A7	available
309	16	A8	available
310	16	B1	available
311	16	B2	available
312	16	B3	available
313	16	B4	available
314	16	B5	available
315	16	B6	available
316	16	B7	available
317	16	B8	available
318	16	C1	available
319	16	C2	available
320	16	C3	available
321	16	C4	available
322	16	C5	available
323	16	C6	available
324	16	C7	available
325	16	C8	available
326	16	D1	available
327	16	D2	available
328	16	D3	available
329	16	D4	available
330	16	D5	available
331	16	D6	available
332	16	D7	available
333	16	D8	available
334	16	E1	available
335	16	E2	available
336	16	E3	available
337	16	E4	available
338	16	E5	available
339	16	E6	available
340	16	E7	available
341	16	E8	available
342	17	A1	available
343	17	A2	available
344	17	A3	available
345	17	A4	available
346	17	A5	available
347	17	A6	available
348	17	A7	available
349	17	A8	available
350	17	B1	available
351	17	B2	available
352	17	B3	available
353	17	B4	available
354	17	B5	available
355	17	B6	available
356	17	B7	available
357	17	B8	available
358	17	C1	available
359	17	C2	available
360	17	C3	available
361	17	C4	available
362	17	C5	available
363	17	C6	available
364	17	C7	available
365	17	C8	available
366	17	D1	available
367	17	D2	available
368	17	D3	available
369	17	D4	available
370	17	D5	available
371	17	D6	available
372	17	D7	available
373	17	D8	available
374	17	E1	available
375	17	E2	available
376	17	E3	available
377	17	E4	available
378	17	E5	available
379	17	E6	available
380	17	E7	available
381	17	E8	available
382	18	A1	available
383	18	A2	available
384	18	A3	available
385	18	A4	available
386	18	A5	available
387	18	A6	available
388	18	A7	available
389	18	A8	available
390	18	B1	available
391	18	B2	available
392	18	B3	available
393	18	B4	available
394	18	B5	available
395	18	B6	available
396	18	B7	available
397	18	B8	available
398	18	C1	available
399	18	C2	available
400	18	C3	available
401	18	C4	available
402	18	C5	available
403	18	C6	available
404	18	C7	available
405	18	C8	available
406	18	D1	available
407	18	D2	available
408	18	D3	available
409	18	D4	available
410	18	D5	available
411	18	D6	available
412	18	D7	available
413	18	D8	available
414	18	E1	available
415	18	E2	available
416	18	E3	available
417	18	E4	available
418	18	E5	available
419	18	E6	available
420	18	E7	available
421	18	E8	available
422	19	A1	available
423	19	A2	available
424	19	A3	available
425	19	A4	available
426	19	A5	available
427	19	A6	available
428	19	A7	available
429	19	A8	available
430	19	B1	available
431	19	B2	available
432	19	B3	available
433	19	B4	available
434	19	B5	available
435	19	B6	available
436	19	B7	available
437	19	B8	available
438	19	C1	available
439	19	C2	available
440	19	C3	available
441	19	C4	available
442	19	C5	available
443	19	C6	available
444	19	C7	available
445	19	C8	available
446	19	D1	available
447	19	D2	available
448	19	D3	available
449	19	D4	available
450	19	D5	available
451	19	D6	available
452	19	D7	available
453	19	D8	available
454	19	E1	available
455	19	E2	available
456	19	E3	available
457	19	E4	available
458	19	E5	available
459	19	E6	available
460	19	E7	available
461	19	E8	available
462	20	A1	available
463	20	A2	available
464	20	A3	available
465	20	A4	available
466	20	A5	available
467	20	A6	available
468	20	A7	available
469	20	A8	available
470	20	B1	available
471	20	B2	available
472	20	B3	available
473	20	B4	available
474	20	B5	available
475	20	B6	available
476	20	B7	available
477	20	B8	available
478	20	C1	available
479	20	C2	available
480	20	C3	available
481	20	C4	available
482	20	C5	available
483	20	C6	available
484	20	C7	available
485	20	C8	available
486	20	D1	available
487	20	D2	available
488	20	D3	available
489	20	D4	available
490	20	D5	available
491	20	D6	available
492	20	D7	available
493	20	D8	available
494	20	E1	available
495	20	E2	available
496	20	E3	available
497	20	E4	available
498	20	E5	available
499	20	E6	available
500	20	E7	available
505	21	A4	available
509	21	A8	available
510	21	B1	available
511	21	B2	available
512	21	B3	available
513	21	B4	available
514	21	B5	available
515	21	B6	available
516	21	B7	available
517	21	B8	available
518	21	C1	available
519	21	C2	available
520	21	C3	available
521	21	C4	available
522	21	C5	available
523	21	C6	available
524	21	C7	available
525	21	C8	available
526	21	D1	available
527	21	D2	available
528	21	D3	available
529	21	D4	available
530	21	D5	available
531	21	D6	available
532	21	D7	available
533	21	D8	available
534	21	E1	available
535	21	E2	available
536	21	E3	available
537	21	E4	available
538	21	E5	available
539	21	E6	available
540	21	E7	available
541	21	E8	available
542	22	A1	available
543	22	A2	available
544	22	A3	available
545	22	A4	available
546	22	A5	available
547	22	A6	available
548	22	A7	available
549	22	A8	available
550	22	B1	available
551	22	B2	available
552	22	B3	available
553	22	B4	available
554	22	B5	available
555	22	B6	available
556	22	B7	available
557	22	B8	available
558	22	C1	available
559	22	C2	available
560	22	C3	available
561	22	C4	available
562	22	C5	available
563	22	C6	available
564	22	C7	available
565	22	C8	available
566	22	D1	available
567	22	D2	available
568	22	D3	available
569	22	D4	available
570	22	D5	available
571	22	D6	available
572	22	D7	available
573	22	D8	available
574	22	E1	available
575	22	E2	available
576	22	E3	available
577	22	E4	available
578	22	E5	available
579	22	E6	available
580	22	E7	available
581	22	E8	available
582	23	A1	available
583	23	A2	available
584	23	A3	available
585	23	A4	available
586	23	A5	available
587	23	A6	available
588	23	A7	available
589	23	A8	available
590	23	B1	available
591	23	B2	available
592	23	B3	available
593	23	B4	available
594	23	B5	available
595	23	B6	available
596	23	B7	available
597	23	B8	available
598	23	C1	available
599	23	C2	available
600	23	C3	available
601	23	C4	available
602	23	C5	available
603	23	C6	available
604	23	C7	available
605	23	C8	available
606	23	D1	available
607	23	D2	available
508	21	A7	available
506	21	A5	available
501	20	E8	available
504	21	A3	available
608	23	D3	available
609	23	D4	available
610	23	D5	available
611	23	D6	available
612	23	D7	available
613	23	D8	available
614	23	E1	available
615	23	E2	available
616	23	E3	available
617	23	E4	available
618	23	E5	available
619	23	E6	available
620	23	E7	available
621	23	E8	available
622	24	A1	available
623	24	A2	available
624	24	A3	available
625	24	A4	available
626	24	A5	available
627	24	A6	available
628	24	A7	available
629	24	A8	available
630	24	B1	available
631	24	B2	available
632	24	B3	available
633	24	B4	available
634	24	B5	available
635	24	B6	available
636	24	B7	available
637	24	B8	available
638	24	C1	available
639	24	C2	available
640	24	C3	available
641	24	C4	available
642	24	C5	available
643	24	C6	available
644	24	C7	available
645	24	C8	available
646	24	D1	available
647	24	D2	available
648	24	D3	available
649	24	D4	available
650	24	D5	available
651	24	D6	available
652	24	D7	available
653	24	D8	available
654	24	E1	available
655	24	E2	available
656	24	E3	available
657	24	E4	available
658	24	E5	available
659	24	E6	available
660	24	E7	available
661	24	E8	available
662	25	A1	available
663	25	A2	available
664	25	A3	available
665	25	A4	available
666	25	A5	available
667	25	A6	available
668	25	A7	available
669	25	A8	available
670	25	B1	available
671	25	B2	available
672	25	B3	available
673	25	B4	available
674	25	B5	available
675	25	B6	available
676	25	B7	available
677	25	B8	available
678	25	C1	available
679	25	C2	available
680	25	C3	available
681	25	C4	available
682	25	C5	available
683	25	C6	available
684	25	C7	available
685	25	C8	available
686	25	D1	available
687	25	D2	available
688	25	D3	available
689	25	D4	available
690	25	D5	available
691	25	D6	available
692	25	D7	available
693	25	D8	available
694	25	E1	available
695	25	E2	available
696	25	E3	available
697	25	E4	available
698	25	E5	available
699	25	E6	available
700	25	E7	available
701	25	E8	available
702	26	A1	available
703	26	A2	available
704	26	A3	available
705	26	A4	available
706	26	A5	available
707	26	A6	available
708	26	A7	available
709	26	A8	available
710	26	B1	available
711	26	B2	available
712	26	B3	available
713	26	B4	available
714	26	B5	available
715	26	B6	available
716	26	B7	available
717	26	B8	available
718	26	C1	available
719	26	C2	available
720	26	C3	available
721	26	C4	available
722	26	C5	available
723	26	C6	available
724	26	C7	available
725	26	C8	available
726	26	D1	available
727	26	D2	available
728	26	D3	available
729	26	D4	available
730	26	D5	available
731	26	D6	available
732	26	D7	available
733	26	D8	available
734	26	E1	available
735	26	E2	available
736	26	E3	available
737	26	E4	available
738	26	E5	available
739	26	E6	available
740	26	E7	available
741	26	E8	available
503	21	A2	available
507	21	A6	available
502	21	A1	available
\.


--
-- TOC entry 4954 (class 0 OID 16417)
-- Dependencies: 218
-- Data for Name: shows; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.shows (show_id, movie_id, show_time) FROM stdin;
8	5	2026-02-01 13:30:00+07
9	5	2026-02-25 17:00:00+07
10	5	2026-02-25 20:00:00+07
11	6	2026-02-25 09:30:00+07
12	6	2026-02-25 12:45:00+07
13	6	2026-02-25 16:15:00+07
14	6	2026-02-25 19:45:00+07
15	6	2026-02-25 22:15:00+07
16	7	2026-02-25 11:00:00+07
17	7	2026-02-25 14:30:00+07
18	7	2026-02-25 18:00:00+07
19	8	2026-02-25 10:30:00+07
20	8	2026-02-25 15:00:00+07
21	9	2026-02-25 10:00:00+07
22	9	2026-02-25 13:00:00+07
23	9	2026-02-25 16:30:00+07
24	9	2026-02-25 19:30:00+07
25	10	2026-02-25 18:00:00+07
26	10	2026-02-25 21:00:00+07
7	5	2026-02-25 17:00:00+07
\.


--
-- TOC entry 4952 (class 0 OID 16399)
-- Dependencies: 216
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (user_id, email, password, full_name) FROM stdin;
1	user1@gmail.com	password123	\N
2	user2@gmail.com	password123	\N
3	user3@gmail.com	password123	\N
4	anh@gmail.com	$2a$10$85FlNmUdH7o2Qpnv4oL0keUWC2nNsDQ2mafSv6pl9GUqnn.h18fxS	Trần Tuấn Anh                                     
5	phuong@gmail.com	$2a$10$7kujkw8ghcOXFLlaC5IJPuqwP97TbvPNt4kV0LsSQV2VljG3SMQIC	Trần Hoàng Phương                                 
6	dang@gmail.com	$2a$10$swl.AE3YrMyf6u0zYmdgQuY0wpwlumUBMjVMUuK/qncTpXs4gAsti	Trần Hải Đăng                                     
7	nhung@gmail.com	$2a$10$/uu6Z03dZEwqcUulL15A.ONl5Jk.GYlJIe/CaKS5exwOd7gKgIMFa	Phương Thị Tuyết Nhung                            
\.


--
-- TOC entry 4977 (class 0 OID 0)
-- Dependencies: 225
-- Name: api_keys_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.api_keys_id_seq', 12, true);


--
-- TOC entry 4978 (class 0 OID 0)
-- Dependencies: 221
-- Name: bookings_booking_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.bookings_booking_id_seq', 173, true);


--
-- TOC entry 4979 (class 0 OID 0)
-- Dependencies: 223
-- Name: movies_movie_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.movies_movie_id_seq', 10, true);


--
-- TOC entry 4980 (class 0 OID 0)
-- Dependencies: 227
-- Name: refresh_tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.refresh_tokens_id_seq', 100, true);


--
-- TOC entry 4981 (class 0 OID 0)
-- Dependencies: 219
-- Name: seats_seat_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.seats_seat_id_seq', 741, true);


--
-- TOC entry 4982 (class 0 OID 0)
-- Dependencies: 217
-- Name: shows_show_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.shows_show_id_seq', 26, true);


--
-- TOC entry 4983 (class 0 OID 0)
-- Dependencies: 215
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_user_id_seq', 7, true);


--
-- TOC entry 4797 (class 2606 OID 16494)
-- Name: api_keys api_keys_key_hash_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_key_hash_key UNIQUE (key_hash);


--
-- TOC entry 4799 (class 2606 OID 16492)
-- Name: api_keys api_keys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_pkey PRIMARY KEY (id);


--
-- TOC entry 4793 (class 2606 OID 16451)
-- Name: bookings bookings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_pkey PRIMARY KEY (booking_id);


--
-- TOC entry 4795 (class 2606 OID 16472)
-- Name: movies movies_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.movies
    ADD CONSTRAINT movies_pkey PRIMARY KEY (movie_id);


--
-- TOC entry 4804 (class 2606 OID 16514)
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 4789 (class 2606 OID 16435)
-- Name: seats seats_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.seats
    ADD CONSTRAINT seats_pkey PRIMARY KEY (seat_id);


--
-- TOC entry 4787 (class 2606 OID 16422)
-- Name: shows shows_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.shows
    ADD CONSTRAINT shows_pkey PRIMARY KEY (show_id);


--
-- TOC entry 4791 (class 2606 OID 16437)
-- Name: seats unique_seat_per_show; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.seats
    ADD CONSTRAINT unique_seat_per_show UNIQUE (show_id, seat_name);


--
-- TOC entry 4783 (class 2606 OID 16406)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 4785 (class 2606 OID 16404)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- TOC entry 4800 (class 1259 OID 16496)
-- Name: idx_api_keys_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_keys_active ON public.api_keys USING btree (is_active);


--
-- TOC entry 4801 (class 1259 OID 16495)
-- Name: idx_api_keys_client_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_keys_client_id ON public.api_keys USING btree (client_id);


--
-- TOC entry 4802 (class 1259 OID 16515)
-- Name: idx_refresh_token_hash; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_refresh_token_hash ON public.refresh_tokens USING btree (token_hash);


--
-- TOC entry 4806 (class 2606 OID 16459)
-- Name: bookings fk_booking_seat; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT fk_booking_seat FOREIGN KEY (seat_id) REFERENCES public.seats(seat_id);


--
-- TOC entry 4807 (class 2606 OID 16454)
-- Name: bookings fk_booking_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT fk_booking_user FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- TOC entry 4805 (class 2606 OID 16438)
-- Name: seats fk_seat_show; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.seats
    ADD CONSTRAINT fk_seat_show FOREIGN KEY (show_id) REFERENCES public.shows(show_id);


-- Completed on 2026-02-08 09:45:32

--
-- PostgreSQL database dump complete
--


