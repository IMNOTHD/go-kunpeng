package canal

const (
	// canal连接相关设置
	Address     = "127.0.0.1"
	Port        = 11111
	Username    = ""
	Password    = ""
	Destination = "example"
	SoTimeOut   = 60000
	IdleTimeOut = 60 * 60 * 1000
)

const (
	// 轮询间隔时间, 单位ms
	PollingInterval = 300
)
