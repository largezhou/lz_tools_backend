-- __UP__ 请勿删除该行
alter table `user`
    add `some` varchar(100) not null;

-- __DOWN__ 请勿删除该行
alter table `user`
    drop column `some`;
