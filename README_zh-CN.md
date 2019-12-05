# 集成测试

Nebula Test 工具用来读取自定义的[测试文件](t/test.t)，然后逐个解析每段的测试用例的输入，向 [Nebula](https://github.com/vesoft-inc/nebula) 服务端发起请求，然后比对请求后的响应和测试用例中的输出是否一致，进而完成集成测试的目的。

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Nebula Test 使用方式](#nebula-test-使用方式)
- [用例](#用例)
    - [1. 测试用例 Title 和 Description](#1-测试用例-title-和-description)
        - [1.1 前缀标识](#11-前缀标识)
        - [1.2 解释](#12-解释)
        - [1.3 完整示例](#13-完整示例)
    - [2. 测试用例输入](#2-测试用例输入)
        - [2.1 前缀标识](#21-前缀标识)
        - [2.2 解释](#22-解释)
        - [2.3 完整示例](#23-完整示例)
    - [3. 测试用例输出](#3-测试用例输出)
        - [3.1 前缀标识](#31-前缀标识)
        - [3.2 解释](#32-解释)
        - [3.3 完整示例](#33-完整示例)
        - [3.4 Column Type](#34-column-type)
- [测试用例的编写](#测试用例的编写)
- [参考](#参考)
- [TODO](#todo)

<!-- markdown-toc end -->


## Nebula Test 使用方式

Nebula Test 设计为独立执行的程序，可以直接运行如下的命令进行测试：

```bash
$ nebula-test --user user --password password \
     --address 127.0.0.1:3699 \
     --file test/test.t
```

可以采用起多个进程的方式来独立的并发测试多个文件。

## 用例

测试文件格式：

```text
=== test: insert vertices
--- in: wait=30s
USE nba;
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
--- in
INSERT EDGE follow(degree) VALUES 100 -> 101:(95), \
 100 -> 102:(90), \
 102 -> 101:(75);
INSERT EDGE serve(start_year, end_year) VALUES 100 -> 200:(1997, 2016);
INSERT EDGE serve(start_year, end_year) VALUES 101 -> 201:(1999, 2018);
--- out
{ "error_code": 0 }


=== test: Another: find the vertex that VID 100 follows, whose age is greater than 35
Firstly
Secondly
More description is ok.
--- in
GO FROM 100 OVER follow WHERE $$.player.age >= 35 \
YIELD $$.player.name AS Teammate, $$.player.age AS Age;
--- out: order=false
{
  "error_code": 0,
  "error_msg": "",
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
Computation.
--- in
YIELD 10
--- out: type=json, order=false
{
  "error_code": 0
}
```

上述中的格式说明如下：

### 1. 测试用例 Title 和 Description

#### 1.1 前缀标识

```text
=== test
```

#### 1.2 解释

**title**: 在前缀标识后的部分作为一个测试用例的标题部分，可以用来区分不同的测试用例，在测试报告中也会给出这部分。需要注意：标题部分只能写在一行之中。

示例：

```text
=== test: Go with filter condition
```

**Description**: 在 title 下面的紧邻行中还可以用来更详细的描述当前测试用例的更多信息，这部分内容没有行数的限制，但是也不会体现在测试报告中。

示例：

```text
find the vertex that VID 100 follows,
whose age is greater than 35
```

#### 1.3 完整示例

```text
=== test: Go with filter conditions
find the vertex that VID 100 follows,
whose age is greater than 35
```

### 2. 测试用例输入

#### 2.1 前缀标识

```text
--- in
```

#### 2.2 解释

**输入**: 前缀标识下面的部分用来输入合法的 nGQL 语句，如果一条 nGQL 需要多行输入，必须在行尾使用 `\` 连接多行。每条 nGQL 需使用 `;` 结尾。

示例：

```text
INSERT VERTEX player(name, age) VALUES 100:("Tim Duncan", 42), \
 101:("Tony Parker", 36), \
 102:("LaMarcus Aldridge", 33);
INSERT VERTEX team(name) VALUES 200:("Warriors");
INSERT VERTEX team(name) VALUES 201:("Nuggets");
INSERT VERTEX player(name, age) VALUES 121:("Useless", 60);
```

**选项**：目前还提供以下选项来定制输入部分的语句执行约束。

- `wait`: 用来指定上个测试用例执行结束以后，该测试用例等待多久再开始执行。比如 schema 配置后，等待 Meta Server 同步数据。

示例：

```text
--- in: wait=1m10s
```

时间的单位：

- `h`: hours
- `m`: minutes
- `s`: seconds
- `ms`: milliseconds
- `us`: microseconds
- `ns`: nanoseconds

#### 2.3 完整示例

```text
--- in: wait=1m10s
INSERT VERTEX player(name, age) VALUES 100:("Tim Duncan", 42), \
 101:("Tony Parker", 36), \
 102:("LaMarcus Aldridge", 33);
INSERT VERTEX team(name) VALUES 200:("Warriors");
INSERT VERTEX team(name) VALUES 201:("Nuggets");
INSERT VERTEX player(name, age) VALUES 121:("Useless", 60);
```

### 3. 测试用例输出

#### 3.1 前缀标识

```text
--- out
```

#### 3.2 解释

**输出**： 前缀标识下面用来放置上述输入中**最后一条 nGQL 输出的结果**。目前支持三种格式的结果比较，详细表述在下面的**选项**中解释。

**选项**： 用来配置输出结果的比较方式。

- `type`: 表示下面结果的格式，有如下取值：
  - `json`: JSON 格式，其中的字段参考 [graph.thrift](https://github.com/vesoft-inc/nebula/blob/master/src/interface/graph.thrift#L107-L114) 定义，为**默认值**。

    如下示例所示，其中除却 `error_code` 之外的其他字段都是**可选**，即如果结果中没有给出，则会忽略该字段。

    示例：

    ```json
    --- out
    {
      "error_code": 0,
      "error_msg": "",
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
    ```
<!--
  - `table`: 表格形式，必须跟 nebula console 中的输出一致。

    示例：

    ```text
    --- out: type=table
    ============================
    | player.name | player.age |
    ============================
    | Tim Duncan  | 42         |
    ----------------------------
    ```

  - `row`: 按行排列，跟上述 `table` 格式类似，只是没有表格格式，通过 TAB 来对齐不同列。

    示例：

    ```text
    --- out: type=row
    player.name	player.age
    Tim Duncan	42
    ```
-->

- `order`: 表示输出行是否需要经过排序，如果是排序后的结果，则严格按照给出的顺序比较。只有两种取值。
  - `false`: 结果不排序，默认值。
  - `true`：结果排序。

  示例：

  ```text
  --- out: order=true
  ```

#### 3.3 完整示例

```text
--- out: order=true
{
  "error_code": 0,
  "error_msg": "",
  "column_names": ["player.name", "player.age"],
  "space_name": "nba",
  "rows": [
    {
      "columns": [
        { "str": "Tim Duncan" },
        { "integer": 42 }
      ]
    },
    {
      "columns": [
        { "str": "Tony Parker" },
        { "integer": 36 }
      ]
    }
  ]
}
```

#### 3.4 Column Type

在上述 JSON 输出格式中，每行里的每列都需要一种类型的 key 来指定对应的值，key 的值可参看 [graph.thrift](https://github.com/vesoft-inc/nebula/blob/master/src/interface/graph.thrift#L75-L99)，具体如下所示：

- str: string
- integer
- bool_val
- id
- single_precision
- double_precision
- timestamp
- year
- month
- date
- datetime
- path

## 测试用例的编写

在编写测试用例时比较建议：

- 不同的模块按测试文件（以 `.t` 结尾）划分，作为一个 test suits。
- 不同测试文件尽量自包含，即可拿该文件独立运行，不依赖其他测试文件的执行顺序。数据也可以尽量的自包含（在测试的开始自行创建独立的 SPACE 和 INSERT 需要的数据）。

## 参考

1. [twitter/mysql-test](https://github.com/twitter-forks/mysql/tree/master/mysql-test) 
1. [openresty/test-nginx](https://openresty.gitbooks.io/programming-openresty/content/testing/test-file-layout.html) 

## TODO

- [X] Sleep some time between adjacent test cases
- [ ] Pretty Print Test Report
- [ ] Support Webhook
- [X] Support comparation in order for output results
