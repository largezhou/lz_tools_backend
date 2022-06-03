-- __UP__ 请勿删除该行
alter table `user`
    add tmp varchar(100) not null;

-- __DOWN__ 请勿删除该行
alter table `user`
    drop column `tmp`;
