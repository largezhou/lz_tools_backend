-- __UP__ 请勿删除该行
create table user
(
    id          bigint unsigned auto_increment
        primary key,
    username    varchar(20) not null,
    create_time datetime    not null,
    update_time datetime    not null,
    constraint user_username_uindex
        unique (username)
);

-- __DOWN__ 请勿删除该行
drop table user
