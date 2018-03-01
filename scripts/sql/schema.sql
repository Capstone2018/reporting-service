-- schema for reporting db
create table users (
    id serial not null unique,
    email varchar(2083) not null,
    count integer not null,
    primary key (id)
);


create table urls (
    id serial not null unique,
    host varchar(2083) not null,
    path varchar(2083) not null,

    
    primary key (id)
);

create table opengraph (
    id serial not null unique,
    created_at timestamp not null,
    
    url varchar(2083),
    title varchar(40),
    type varchar(7), 
    description varchar(300),
    determiner varchar(5), 
    locale varchar(20),
    locales_alternate text[],
    images jsonb,
    audios jsonb,
    videos jsonb,
    profile jsonb,
    article jsonb,
    book jsonb,

    blob jsonb not null,

    primary key (id)
);

create table query_fragment (
    id serial not null unique,
    query varchar(25),
    fragment varchar(25),

    primary key(id)
);

create table reports (
    id serial not null unique,
    user_id integer not null,

    user_description varchar(1024) not null,
    created_at timestamp not null,

    foreign key (user_id) references users(id),
    
    primary key (id)
);

create table report_types (
    id serial not null unique,
    type varchar(100) not null,

    primary key (id)
);

create table report_types_reports (
    report_id integer references reports(id) ON UPDATE CASCADE ON DELETE CASCADE,
    report_type_id integer references report_types(id) ON UPDATE CASCADE,
    
    constraint report_type_pkey primary key (report_id, report_type_id)
);

create table pages (
    id serial not null unique,
    created_at timestamp not null,
    
    url_id integer not null,
    og_id integer,
    report_id integer,
    query_fragment_id integer,

    wayback varchar(100),
    url_string varchar(2083) not null,

    foreign key (url_id) references urls(id),
    foreign key (og_id) references opengraph(id), 
    foreign key (report_id) references reports(id),
    foreign key (query_fragment_id) references query_fragment(id),

    primary key (id)
);
