package mysql

const (
	// Username Example: root
	Username = "haetae"
	// Password Example: 123456
	Password = "4ZeARKHFMy5GStf3"
	// protocol Example: tcp
	Protocol = "tcp"
	// Address Example: 127.0.0.1
	Address = "121.36.28.205"
	// Port Example: 3306
	Port = 3306
	// Dbname Example: test
	Dbname = "haetae"
	// Addition Example: param1=value1&...&paramN=valueN
	// 务必加上parseTime=true, 否则查询datetime会报错
	// unsupported Scan, storing driver.Value type []uint8 into type *time.Time
	Addition = "charset=utf8mb4&parseTime=true"
)
