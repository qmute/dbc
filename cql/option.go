package cql

import (
	"time"

	"github.com/gocql/gocql"
)

// Option 是一个连接选项 interface.
type Option interface {
	apply(*gocql.ClusterConfig)
}

type optionFunc func(*gocql.ClusterConfig)

func (f optionFunc) apply(cfg *gocql.ClusterConfig) {
	f(cfg)
}

// WithPassword 用户密码选项.
func WithPassword(u, p string) Option {
	return optionFunc(func(cfg *gocql.ClusterConfig) {
		cfg.Authenticator = &gocql.PasswordAuthenticator{
			Username: u,
			Password: p,
		}
	})
}

// WithTimeout 超时选项.
func WithTimeout(t time.Duration) Option {
	return optionFunc(func(cfg *gocql.ClusterConfig) {
		cfg.Timeout = t
	})
}

// WithConnectTimeout 初始连接超时选项.
func WithConnectTimeout(t time.Duration) Option {
	return optionFunc(func(cfg *gocql.ClusterConfig) {
		cfg.ConnectTimeout = t
	})
}

// WithPort 端口，默认:9042
func WithPort(p int) Option {
	return optionFunc(func(cfg *gocql.ClusterConfig) {
		cfg.Port = p
	})
}

// WithProtoVersion 协议版本选项.
func WithProtoVersion(v int) Option {
	return optionFunc(func(cfg *gocql.ClusterConfig) {
		cfg.ProtoVersion = v
	})
}

// WithAuthenticator 验证器选项
func WithAuthenticator(auth gocql.Authenticator) Option {
	return optionFunc(func(cfg *gocql.ClusterConfig) {
		cfg.Authenticator = auth
	})
}

// WithAliyunAuth 阿里云cassandra验证器
func WithAliyunAuth(u, p string) Option {
	ali := AliAuthenticator{
		Username: u,
		Password: p,
	}

	return WithAuthenticator(ali)
}

// WithNumConns 每个主机的连接数
func WithNumConns(n int) Option {
	return optionFunc(func(cfg *gocql.ClusterConfig) {
		cfg.NumConns = n
	})
}

// WithPageSize Default page size to use for created sessions
func WithPageSize(size int) Option {
	return optionFunc(func(cfg *gocql.ClusterConfig) {
		cfg.PageSize = size
	})
}

// WithSslOpts  ssl 配置
func WithSslOpts(ssl *gocql.SslOptions) Option {
	return optionFunc(func(cfg *gocql.ClusterConfig) {
		cfg.SslOpts = ssl
	})
}
