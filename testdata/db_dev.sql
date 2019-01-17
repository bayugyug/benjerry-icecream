
create database if not exists benjerry_dev;
create user benjerry_dev;
grant all privileges on benjerry_dev.* to 'benjerry_dev'@'localhost' identified by 'icecream';
grant all privileges on benjerry_dev.* to 'benjerry_dev'@'127.0.0.1' identified by 'icecream';

flush privileges;


