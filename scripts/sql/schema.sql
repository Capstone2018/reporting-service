-- schema for reporting db

create table reports (
    id int(10) not null auto_increment,
    ru_id int(10) not null,
    report_type varchar(100),
    description varchar(1024) not null,
    created_at datetime not null,

    foreign key (ru_id) references report_urls(id),
    primary key (id)
);


create table report_urls (
    id int(10) not null auto_increment,
    host_id int(10) not null,
    meta_id int(10) not null,
    url varchar(2083) not null,
    archive_url varchar(2083) not null,
    title varchar(320),
    author_string varchar(100),
    content_summary varchar(320),
    content_category varchar(100) not null,

    foreign key (host_id) references hostnames(id),
    foreign key (meta_id) references url_metadata(id),
    primary key (id)
)

create table url_metadata (
    id int(10) not null auto_increment,
    query varchar(1024),
    fragment varchar(320),

    primary key (id)
)

create table hostnames (
    id int(10) not null auto_increment,
    host varchar(2083) not null,

    primary key (id)
);

create table users (
    id int(10) not null auto_increment,
    email varchar(320) not null,
    passhash binary(60) not null,
    username varchar(64) not null,
    created_at datetime not null,

    primary key (id)
);

create table users_websites (
    user_id int(10) not null,
    website_id int(10) not null,

    foreign key (user_id) references users(id),
    foreign key (website_id) references websites(id)
);