# nebula-test

[中文](README_zh-CN.md)

Test framework for [Nebula Graph](https://github.com/vesoft-inc/nebula.git) project.

## Format

```text
=== test: create space sp
--- in
CREATE SPACE sp(partition_num=1024, replica_factor=1);
--- out
```

## TODO

- [ ] Pretty report
- [ ] Deploy in k8s cluster
