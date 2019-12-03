=== test: create space sp
--- in
-- DROP SPACE nba;
CREATE SPACE nba(partition_num=128, replica_factor=1);
--- out
=== test: use space nba
--- in
USE nba;
CREATE TAG player (name string, age int);
CREATE TAG team (name string);
CREATE EDGE follow(degree int);
CREATE EDGE serve(start_year int, end_year int);
--- out
=== test: empty test
