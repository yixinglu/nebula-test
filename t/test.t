=== test: create space sp
--- in
-- DROP SPACE nba;
CREATE SPACE IF NOT EXISTS nba(partition_num=128, replica_factor=1);
--- out
{ "error_code": 0 }


=== test: use space nba
--- in
USE nba;
CREATE TAG IF NOT EXISTS player (name string, age int);
CREATE TAG IF NOT EXISTS team (name string);
CREATE EDGE IF NOT EXISTS follow(degree int);
CREATE EDGE IF NOT EXISTS serve(start_year int, end_year int);
--- out
{ "error_code": 0 }


=== test: empty test

=== test: insert vertices
More description,
New line test.
--- in: wait=10s
INSERT VERTEX player(name, age) VALUES 100:("Tim Duncan", 42), \
 101:("Tony Parker", 36), \
 102:("LaMarcus Aldridge", 33);
INSERT VERTEX team(name) VALUES 200:("Warriors");
INSERT VERTEX team(name) VALUES 201:("Nuggets");
INSERT VERTEX player(name, age) VALUES 121:("Useless", 60);
--- out
{
  "error_code": 0,
  "space_name": "nba"
}


=== test: insert edges
--- in:
INSERT EDGE follow(degree) VALUES 100 -> 101:(95), \
 100 -> 102:(90), \
 102 -> 101:(75);
INSERT EDGE serve(start_year, end_year) VALUES 100 -> 200:(1997, 2016);
INSERT EDGE serve(start_year, end_year) VALUES 101 -> 201:(1999, 2018);
--- out
{ "error_code": 0 }


=== test: fetch vertex props
--- in
FETCH PROP ON player 100;
--- out
{
	"error_code": 0,
	"column_names": ["player.name", "player.age"],
	"space_name": "nba",
	"rows": [
    {
		  "columns": [
        { "str": "Tim Duncan" },
        { "integer": 42 }
      ]
    }
  ]
}


=== test: fetch edge props
--- in
FETCH PROP ON serve 100 -> 200;
--- out
{
	"error_code": 0,
	"column_names": ["serve.start_year", "serve.end_year"],
	"space_name": "nba",
	"rows": [
    {
		  "columns": [
        { "integer": 1997 },
        { "integer": 2016 }
      ]
    }
  ]
}


=== test: find the vertex that VID 100 follows, whose age is greater than 35
--- in
GO FROM 100 OVER follow WHERE $$.player.age >= 35 \
YIELD $$.player.name AS Teammate, $$.player.age AS Age;
--- out: type=json, order=false
{
	"error_code": 0,
	"column_names": ["Teammate", "Age"],
	"space_name": "nba",
	"rows": [
    {
		  "columns": [
        { "str": "Tony Parker" },
        { "integer": 36 }
      ]
    }
  ]
}


=== test: Yield
--- in
YIELD 10
--- out: type=json, order=false
{
  "error_code": 0,
  "error_msg": ""
}
