CREATE TABLE public.current_media (
    id integer NOT NULL,
    title public.citext,
    type integer,
    weight numeric DEFAULT 1,
    date_added date DEFAULT now(),
    referrer character varying,
    removed date,
    started date,
    priority boolean DEFAULT false,
    genre character varying
);
