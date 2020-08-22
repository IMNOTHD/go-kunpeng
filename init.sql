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

GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' IDENTIFIED BY 'canal' WITH GRANT OPTION;
FLUSH PRIVILEGES;