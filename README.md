# mysqldump-slice
It's wrap for mysqldump, gzip and /bin/sh. Mysqldump-slice allows to make short dump DB with consistency data.

Build:
```
go buld -o target/slice cmd/main.go
```

Run:
```
target/slice ./.slice.yaml
```

Example config yaml file:
```
host: "localhost"
database: "test"
user: "admin"
password: "admin"
max-connect: 10
log: Yes

filename:
  path: "./target/"
  prefix: "short"
  gzip: Yes
  date-format: "2006-01-02_15:04:05"

tables:
  limit: 100
  full:
    - migration_versions
  ignore:
    - test 
  specs:
    - name: user
      pk:
        - uuid
      sort:
        - updated_at
      limit: 5
      condition: "updated_at > NOW() - INTERVAL 1 WEEK"
```
