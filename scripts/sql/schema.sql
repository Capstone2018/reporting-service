-- schema for reporting db
-- create table users (
--     id serial not null unique,

--     primary key (id)
-- );

create table hostnames (
    id serial not null unique,
    host varchar(2083) not null,

    primary key (id)
);

create table urls (
    id serial not null unique,
    host_id integer not null,
    path varchar(2083) not null,

    foreign key (host_id) references hostnames(id),
    
    primary key (id)
);

create table opengraph (
    id serial not null unique,
    created_at timestamp not null,

    url_id integer not null,
    article_id integer,
    book_id integer,
    profile_id integer,
    
    title varchar(40),
    type varchar(7), 
    description varchar(300),
    determiner varchar(5), 
    locale varchar(20),
    locales_alternate varchar(100),
    images jsonb,
    audios jsonb,
    videos jsonb,
    profiles jsonb,
    articles jsonb,
    books jsonb,

    foreign key (url_id) references urls(id),

    primary key (id)
);

create table query (
    id serial not null unique,
);

create table page (
    id serial not null unique,
    created_at timestamp not null, 

    primary key (id)
);


create table reports (
    id serial not null unique,
    user_id integer not null,
    og_id integer not null,
    report_type varchar(100) not null,
    user_description varchar(1024) not null,
    created_at timestamp not null,

    foreign key (user_id) references users(id),
    foreign key (og_id) references opengraph(id),

    primary key (id)
);