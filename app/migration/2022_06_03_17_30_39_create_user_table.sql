-- __UP__ 请勿删除该行
create table user
(
    id          bigint unsigned auto_increment primary key,
    uuid        varchar(60)             not null,
    union_id    varchar(60)             not null,
    nickname    varchar(100) default '' null,
    avatar      varchar(300) default '' null,
    create_time datetime                not null,
    update_time datetime                not null,
    constraint user_union_id_uindex unique (union_id),
    constraint user_uuid_uindex unique (uuid)
);

-- __DOWN__ 请勿删除该行
drop table user;
