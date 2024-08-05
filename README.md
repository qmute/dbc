# dbc

db connector

## Usage

由于是私有库，需要配置GIT。

```shell
git config --global url.git@github.com:.insteadOf https://github.com/

cat ~/.gitconfig

[url "git@github.com:"]
    insteadOf = https://github.com/

```

## Getting Started

- [cassandra/scyllaDB](#cassandra/scyllaDB)
- [mySQL](#mySQL)
- [postgreSQL](#postgreSQL)


## cassandra/scyllaDB

### Install

* 下载 
``` 
    go get -u github.com/qmute/dbc/cql
```

* 配合gocqlx使用
```
    go get -u github.com/scylladb/gocqlx
```

### Example

```go
package main

import (
    "time"

    "github.com/qmute/dbc/cql"
) 

var (
	hosts    = []string{}
	username = ""
	password = ""
	keyspace = "dbc_test"
)

func main() {
    opts := []cql.Option{
        cql.WithTimeout(500 * time.Millisecond),
        cql.WithProtoVersion(4),

        // 数据库验证授权, 以下3种
        // 1.普通用户名密码
        // cql.WithPassword(username, password),
        // 2.指定验证器 
        // cql.WithAuthenticator()
        // 3.阿里云cassandra自定义验证器
        // cql.WithAliyunAuth(username, password),
    }

    se, err := cql.Connect(hosts, keyspace, opts...)
    if err != nil {
        // handle error
    }  
    
    // use session doing
    se.Close()
  
}

```

## mySQL

### Install

* 下载 
``` 
    go get -u github.com/qmute/dbc/gdb
```

### Example

```go
package main

import (
    "time"

    "github.com/qmute/dbc/gdb"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
) 

var (
    conn  string
)

func main() {
    opts := []gdb.Option{
		gdb.WithConnMaxLifetime(3 * time.Minute),
		gdb.WithMaxIdleConns(10),
		gdb.WithMaxOpenConns(50),
    }

    cfg := &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
        },
    }
    db, err := gdb.ConnectToMysql(conn, cfg, opts...)
    if err != nil {
        // handle error
    }  
    
    // use db doing
    // ....
  
}

```

## postgreSQL

### Install

* 下载 
``` 
    go get -u github.com/qmute/dbc/gdb
```

### Example

```go
package main

import (
    "time"

    "github.com/qmute/dbc/gdb"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
) 

var (
    conn  string
)

func main() {
    opts := []gdb.Option{
		gdb.WithConnMaxLifetime(3 * time.Minute),
		gdb.WithMaxIdleConns(10),
		gdb.WithMaxOpenConns(50),
    }
    cfg := &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
        },
    }
    db, err := gdb.ConnectToPG(conn, cfg, opts...)
    if err != nil {
        // handle error
    }  
    
    // use db doing
    // ....
}

```

## change log:

- v0.1.0 init project
