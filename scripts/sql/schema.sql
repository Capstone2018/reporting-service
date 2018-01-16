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

delimiter $$
create procedure insert_report 
(in `host` varchar(2083), in `path` varchar(2083), 
in `archive_url` varchar(2083), in `query` varchar(1024), 
in `fragment` varchar(320), in `report_type` varchar(100), 
in `description` varchar(1024), in `created_at` datetime,
in `title` varchar(320), in `author_string` varchar(100), 
in `content_summary` varchar(320), in `content_category` varchar(100))
begin
    declare host_id int unsigned default 0;
    declare meta_id int unsigned default 0;
    declare ru_id int unsigned default 0;
    declare exit handler for sqlexception
    begin
        rollback;
    end;

    start transaction;
        if (select exists(select 1 from hostnames where host = `host`)) then
            select id into @host_id from hostnames where host = `host`;
        else
            insert into hostnames(host) values(`host`);
            set @host_id = last_insert_id();
        end if;

        if (select exists(select 1 from url_metadata where query = `query` and fragment = `fragment`)) then
            select id into @meta_id from url_metadata where query = `query` and fragment = `fragment`;
        else
            insert into url_metadata(query, fragment) values(`query`, `fragment`);
            set @meta_id = last_insert_id();
        end if;

        if (select exists(select 1 from report_urls where host_id = @host_id and path = `path` and archive_url = `archive_url` and title = `title` and content_category = `content_category` and content_summary = `content_summary`)) then
            select id into @ru_id from report_urls where host_id = @host_id and path = `path` and archive_url = `archive_url` and title = `title` and content_category = `content_category` and content_summary = `content_summary`;
        else
            insert into report_urls(host_id, path, archive_url, title, author_string, content_summary, content_category) values(@host_id, `path`, `archive_url`, `title`, `author_string`, `content_summary`, `content_category`);
            set @ru_id = last_insert_id();
        end if;
        -- select id into @host_id from hostnames where host = `host`;
        -- if (@host_id = 0) then
        --     insert into hostnames(host) values(`host`);
        --     set @host_id = last_insert_id();
        -- end if;

        -- select id into @meta_id from url_metadata where query = `query` and fragment = `fragment`;
        -- if (@meta_id = 0) then
        --     insert into url_metadata(query, fragment) values(`query`, `fragment`);
        --     set @meta_id = last_insert_id();
        -- end if;

        -- select id into @ru_id from report_urls where host_id = @host_id and path = `path` and archive_url = `archive_url` and title = `title` and content_category = `content_category` and content_summary = `content_summary`; 
        -- if (@ru_id = 0) then
        --     insert into report_urls(host_id, path, archive_url, title, author_string, content_summary, content_category) values(@host_id, `path`, `archive_url`, `title`, `author_string`, `content_summary`, `content_category`);
        --     set @ru_id = last_insert_id();
        -- end if;

        insert into reports(ru_id, meta_id, report_type, description, created_at) values(@ru_id, @meta_id, `report_type`, `description`, `created_at`);
    commit;
end;
$$

delimiter ;