
create database if not exists benjerry;
create user benjerry;
grant all privileges on benjerry.* to 'benjerry'@'localhost' identified by 'icecream';
grant all privileges on benjerry.* to 'benjerry'@'127.0.0.1' identified by 'icecream';

flush privileges;

