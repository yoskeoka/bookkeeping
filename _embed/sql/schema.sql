-- SQLite3

drop table if exists accounts;
create table accounts(
    code integer not null primary key,
    name text,
    is_bs boolean,
    is_left boolean
);

drop table if exists general_ledger;
create table general_ledger(
    id integer not null primary key,
    date date,
    code integer,
    description text,
    left integer DEFAULT 0,
    right integer DEFAULT 0
);


