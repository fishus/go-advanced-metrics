create table IF NOT EXISTS public.metrics_counter (
  name character varying primary key not null,
  value bigint not null
);

create table IF NOT EXISTS public.metrics_gauge (
  name character varying primary key not null,
  value double precision not null
);
