module forum-rpc

go 1.23

require (
	common v0.0.0
	github.com/zeromicro/go-zero v1.9.2
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.30.0
)

replace common => ../common
