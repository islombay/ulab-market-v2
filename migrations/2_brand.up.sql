create table if not exists brands (
    id uuid primary key ,
    name varchar(30) not null unique 
);