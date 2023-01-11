# mysqldump-slice
It's wrap for mysqldump, gzip and /bin/sh. Mysqldump-slice allows to make short dump DB with consistency data.

### Build:
```
go buld -o target/slice cmd/main.go
```

### Change config:
- Create config for slice and change params:
```
cp .conf.yaml conf.yaml
```
- Create config for mysqldump and change params:
```
cp .db.cnf db.cnf
```

### Run:
```
target/slice ./conf.yaml
```


### Description config yaml file:
- Mysql connect for load pk/fk
```
host: "localhost"
database: "test"
user: "admin"
password: "admin"
```
- Path to config for mysqldump
```
default-extra-file: "./db.cnf"
```
- Params for connect mysql
```
max-connect: 10
max-lifetime-connect-minute: 5
max-lifetime-query-second: 3
```
- Flag for show logs
```
log: Yes
```
- Options for setting to save dump file
```
filename:
  path: "./target/"
  prefix: "short"
  gzip: Yes
  date-format: "2006-01-02_15:04:05"
```
- Global limit for each tables
```
tables:
  limit: 100
```
- Table list for full data dump
```
  full:
    - migration_versions
```
- Table list for ignore data dump
```
  ignore:
    - test 
```
- Special setting for some table
```
  specs:
    - name: user
```
- If column is not specified like FK
```
      pk:
        - uuid
```
- Column list for sorting
```
      sort:
        - updated_at
```
- Private limit
```
      limit: 5
```
- Special conditions for to select short data
```
      condition: "updated_at > NOW() - INTERVAL 1 WEEK"
```


### For development:
- Create configs **conf.yaml** and **db.cnf**
- Change network name for block docker-compose.yaml
- Run
```
make watch
```

