create table if not exists fetchers
(
    id       serial  not null
        constraint fetchers_pk
            primary key,
    url      text    not null,
    interval integer not null,
    job_id   integer
);

alter table fetchers
    owner to postgres;

create unique index if not exists fetchers_id_uindex
    on fetchers (id);

create table if not exists histories
(
    fetcher_id integer
        constraint responses_fetchers_id_fk
            references fetchers
            on delete cascade,
    response   text,
    duration   double precision,
    created_at double precision
);

alter table histories
    owner to postgres;