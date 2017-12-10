-- schema for reporting db
create table users (
    id int(10) not null auto_increment,
    email varchar(320) not null,
    passhash binary(60) not null,
    username varchar(64) not null,
    createdAt datetime not null,

    primary key (id)
);

create table users_websites (
    userID int(10) not null,
    websiteID int(10) not null,

    foreign key (userID) references users(id),
    foreign key (websiteID) references websites(id)
);

create table websites (
    id int(10) not null auto_increment,
    url varchar(2083) not null,
    host varchar(2083) not null,
    --description varchar(1024),

    primary key (id)
);

create table reports (
    id int(10) not null auto_increment,
    --user_id int(10) not null,
    websiteID int(10) not null,
    description varchar(1024) not null,
    createdAt datetime not null,
    --foreign key (user_id) references users(id),
    foreign key (websiteID) references websites(id),
    primary key (id)
);