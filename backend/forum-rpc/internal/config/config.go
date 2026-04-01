package config

import "github.com/zeromicro/go-zero/zrpc"

type DatabaseConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConf struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Config struct {
	zrpc.RpcServerConf
	Database DatabaseConf
	Redis    RedisConf
}
