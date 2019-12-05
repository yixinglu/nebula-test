=== test: create space sp
--- in
-- DROP SPACE nba;
CREATE SPACE IF NOT EXISTS nba(partition_num=128, replica_factor=1);
--- out
=== test: use space nba
--- in
USE nba;
CREATE IF NOT EXISTS TAG player (name string, age int);
CREATE IF NOT EXISTS TAG team (name string);
CREATE IF NOT EXISTS EDGE follow(degree int);
CREATE IF NOT EXISTS EDGE serve(start_year int, end_year int);
--- out
