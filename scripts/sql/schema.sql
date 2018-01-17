-- schema for reporting db
create table url_metadata (
    id int(10) not null auto_increment,
    query varchar(1024),
    fragment varchar(320),

    primary key (id)
);

create table hostnames (
    id int(10) not null auto_increment,
    host varchar(2083) not null,

    primary key (id)
);

create table report_urls (
    id int(10) not null auto_increment,
    host_id int(10) not null,
    path varchar(2083) not null,
    archive_url varchar(2083) not null,
    title varchar(320),
    author_string varchar(100),
    content_summary varchar(320),
    content_category varchar(100) not null,

    foreign key (host_id) references hostnames(id),
    
    primary key (id)
);

create table reports (
    id int(10) not null auto_increment,
    ru_id int(10) not null,
    meta_id int(10) not null,
    report_type varchar(100),
    description varchar(1024) not null,
    created_at datetime not null,

    foreign key (ru_id) references report_urls(id),
    foreign key (meta_id) references url_metadata(id),
    primary key (id)
);

create table test (
    id int(10) not null auto_increment,
    at_host varchar(100),
    tick_host varchar(100),
    host_id int(10),
    primary key (id)
);

delimiter $$
create procedure insert_report 
(in `host` varchar(2083), in `path` varchar(2083), 
in `archive_url` varchar(2083), in `query` varchar(1024), 
in `fragment` varchar(320), in `report_type` varchar(100), 
in `description` varchar(1024), in `created_at` datetime,
in `title` varchar(320), in `author_string` varchar(100), 
in `content_summary` varchar(320), in `content_category` varchar(100))
begin
    declare host_id int;
    declare meta_id int;
    declare ru_id int;
    declare exit handler for sqlexception
    begin
        rollback;
    end;

    start transaction;
        -- get the host_id if the host already exists
        set @host_id = (select h.id from hostnames h where h.host = `host`);
        if (@host_id is null) then
            insert into hostnames(host) values(host);
            set @host_id = (select last_insert_id());
        end if;
        
        -- get the meta_id if the url_meta already exists
        set @meta_id = (
            select u.id from url_metadata u 
            where u.query = `query` and u.fragment = `fragment`
            );
        if (@meta_id is null) then
            insert into url_metadata(query, fragment) values(`query`, `fragment`);
            set @meta_id = (select last_insert_id());
        end if;

        -- get the ru_id if the report url already exists
        set @ru_id = (
            select r.id from report_urls r 
            where r.host_id = @host_id and r.path = `path`
            and r.archive_url = `archive_url` and r.title = `title` 
            and r.content_category = `content_category` 
            and r.content_summary = `content_summary`
            );
        if (@ru_id is null) then
            insert into report_urls(host_id, path, archive_url, title, author_string, content_summary, content_category) 
                values(@host_id, `path`, `archive_url`, `title`, `author_string`, `content_summary`, `content_category`);
            set @ru_id = (select last_insert_id());
        end if;

        -- finally insert into the reports table
        insert into reports(ru_id, meta_id, report_type, description, created_at) 
            values(@ru_id, @meta_id, `report_type`, `description`, `created_at`);
    commit;
end;
$$

delimiter ;