create table product(
  id bigint not null primary key auto_increment,
  title varchar(32) not null,
  price int not null,
  date_created bigint not null
);

insert into product(title,price,date_created) value('iPhone12Mini','799000',1618652851556);
insert into product(title,price,date_created) value('iPhone12','899000',1618652851556);
insert into product(title,price,date_created) value('iPhone12Pro','999000',1618652851556);
insert into product(title,price,date_created) value('iPhone12Pro256','1099000',1618652851556);
insert into product(title,price,date_created) value('iPhone12SE','399000',1618652851556);