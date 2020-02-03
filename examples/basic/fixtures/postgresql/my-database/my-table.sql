CREATE TABLE public.sort_attribute (
    id integer NOT NULL,
    name text NOT NULL,
    key text NOT NULL,
    value text NOT NULL,
    input_type text NOT NULL,
    page text NOT NULL,
    "position" integer NOT NULL,
    status integer NOT NULL,
    version integer NOT NULL
);
ALTER TABLE ONLY public.sort_attribute ADD CONSTRAINT unique_sort_attribute_id UNIQUE (id);