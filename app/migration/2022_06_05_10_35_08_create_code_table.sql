-- __UP__ 请勿删除该行
create table code
(
    id           bigint unsigned auto_increment primary key,
    copy_from_id bigint unsigned default 0 null comment '从哪个码复制而来的',
    user_id      bigint unsigned           not null,
    name         varchar(50)               not null,
    lng          decimal(10, 6)            not null comment '经度',
    lat          decimal(10, 6)            not null comment '纬度',
    link         varchar(300)              not null comment '二维码的链接',
    times        int unsigned    default 0 not null comment '使用的次数',
    often        tinyint(1)      default 0 not null comment '是否常用',
    share        tinyint(1)      default 1 not null comment '是否其他人可见',
    create_time  datetime                  not null,
    update_time  datetime                  not null
);

create index code_user_id_index
    on code (user_id);

-- __DOWN__ 请勿删除该行
drop table code
