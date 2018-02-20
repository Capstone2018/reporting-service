-- schema for reporting db
create table users (
    id int(10) not null auto_increment,

    primary key (id)
);

create table hostnames (
    id int(10) not null auto_increment,
    host varchar(2083) not null,

    primary key (id)
);

create table urls (
    id int(10) not null auto_increment,
    host_id int(10) not null,
    path varchar(2083) not null,

    foreign key (host_id) references hostnames(id),
    
    primary key (id)
);

create table images (
    id int(10) not null auto_increment,
    og_url int(10) not null,
    url_id int(10) not null,
    type varchar(7),
    width int(5),
    height int(5),

    foreign key (og_url) references opengraph(id),
    foreign key (url_id) references urls(id),

    primary key (id)
);

create table audios (
    id int(10) not null auto_increment,
    og_url int(10) not null,
    url_id int(10) not null,
    type varchar(7),
    width int(5),
    height int(5),

    foreign key (og_url) references opengraph(id),
    foreign key (url_id) references urls(id),
    
    primary key (id)
);

create table videos (
    id int(10) not null auto_increment,
    og_url int(10) not null,
    url_id int(10) not null,
    type varchar(7),
    width int(5),
    height int(5),
    
    foreign key (og_url) references opengraph(id),
    foreign key (url_id) references urls(id),
    
    primary key (id)
);

create table profiles (
    id int(10) not null auto_increment,
    url_id int(10) not null,
    firstname varchar(100),
    lastname varchar(100),
    username varchar(100),
    gender varchar(10),
    foreign key (url_id) references urls(id),
    
    primary key (id)
);

create table articles (
    id int(10) not null auto_increment,
    url_id int(10) not null,
    
    published_time datetime,
    modified_time datetime,
    expiration_time datetime,
    section varchar(100),
    tags varchar(300),

    foreign key (url_id) references urls(id),
    foreign key (author_id) references authors(id),
    
    primary key (id)
);

create table books (
    id int(10) not null auto_increment,
    url_id int(10) not null,
    author_id int(10),

    isbn varchar(100),
    release_date datetime,
    tags varchar(300),
    foreign key (url_id) references urls(id),
    
    primary key (id)
);

create table opengraph (
    id int(10) not null auto_increment,
    url_id int(10) not null,

    article_id int(10),
    book_id int(10),
    profile_id int(10),
    
    title varchar(40),
    type varchar(7), 
    description varchar(300),
    determiner varchar(5), 
    locale varchar(20),
    locales_alternate varchar(100),

    foreign key (url_id) references urls(id),

    primary key (id)
);


create table reports (
    id int(10) not null auto_increment,
    user_id int(10) not null,
    og_id int(10) not null,
    report_type varchar(100) not null,
    user_description varchar(1024) not null,
    created_at datetime not null,

    foreign key (user_id) references users(id),
    foreign key (og_id) references opengraph(id),

    primary key (id)
);