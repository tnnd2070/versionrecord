# versionrecord
# TODO说明

1. 读取csv文件到二维表中，data[应用编号][部署单元][ip][版本]
2. 选出一个部署单元中最大的版本号
3. 根据应用编号，部署单元，查找另个csv文件的版本
4. 生成一个对比文档 应用编号，部署单元，csv1的版本，csv2的版本

### 测试
go run .\main.go -file1 test/hxpre_pkgversion_2021121921.csv -file2 test/hxpro_pkgversion_2021121921.csv

### recordpg 编译
`sh build.sh`

### postgresql 初始化

#### pg启动

```dockerfile
docker  volume create pgdata
docker run -it --name cl_postgress \
        --restart always \
        -e TZ='Asia/Shanghai' \
        -e POSTGRES_PASSWORD='abc123' \
        -e ALLOW_IP_RANGE=0.0.0.0/0 \
        -v pgdata:/var/lib/postgresql \
        -p 55435:5432 -d postgres
```
 - 用户，授权
```sql
CREATE USER vc WITH PASSWORD '1qaz@WSX';

create database vcdb;

GRANT ALL PRIVILEGES ON DATABASE vcdb TO vc;
```
  - 创建表

```sql
CREATE TABLE version_history(
   id serializable PRIMARY KEY,
   sysid varchar(32),
   unitid varchar(32),
   pkgname text,
   ip char(50),
   vsion varchar(32),
   timestamp timestamp
);
```
