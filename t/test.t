=== test: create space sp
--- in
CREATE SPACE sp(partition_num=1024, replica_factor=1);
--- out
=== test: use space sp
--- in
USE sp;
--- out
=== test: empty test
