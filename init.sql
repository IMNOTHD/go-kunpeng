/*
运行前必须在原数据库中运行以下sql, 其作用是:
1. 创建一个user, password均为canal的用户, 允许所有网络访问(其他用户也可以, 但是记得在deployer/conf/example/instance.properties中修改)
2. 授予canal超级权限
3. 刷新权限

** 请注意, mysql必须开启binlog, 开启方式在my.ini的[mysqld]段中加入以下三行
log_bin=mysql-bin
binlog-format=ROW #选择row模式
server_id=1 #配置mysql replaction需要定义，不能和canal的slaveId重复
*/

CREATE USER 'canal'@'%' IDENTIFIED BY 'canal';

/*
在Mysql 8.0中, 使用上面的加密方式连接会报错error 2059: Authentication plugin 'caching_sha2_password' cannot be loaded
使用8.0版本的时候, 请额外运行以下两条sql语句
ALTER USER 'canal'@'%' IDENTIFIED BY 'canal' PASSWORD EXPIRE NEVER; #修改加密规则
ALTER USER 'canal'@'%' IDENTIFIED WITH mysql_native_password BY 'canal'; #更新一下用户的密码
*/

GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' IDENTIFIED BY 'canal' WITH GRANT OPTION;
FLUSH PRIVILEGES;