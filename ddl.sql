create table product(
  id serial not null primary key,
  title varchar(32) not null,
  price integer not null,
  date_created bigint not null
);

insert into product(title,price,date_created) values('iPhone12Mini','799000',1618652851556);
insert into product(title,price,date_created) values('iPhone12','899000',1618652851556);
insert into product(title,price,date_created) values('iPhone12Pro','999000',1618652851556);
insert into product(title,price,date_created) values('iPhone12Pro256','1099000',1618652851556);
insert into product(title,price,date_created) values('iPhone12SE','399000',1618652851556);