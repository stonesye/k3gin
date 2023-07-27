package redisx

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisConfig struct {
	Addr     Addr     // 地址(IP:Port)
	DB       DB       // 数据库编号
	Password Password // 密码
	Prefix   Prefix   // 存储前缀
}

type Addr string
type Password string
type DB int
type Prefix string

func WihtMustRedisConfig(addr Addr, password Password) func(*RedisConfig) {
	return func(rcf *RedisConfig) {
		rcf.Addr = addr
		rcf.Password = password
	}
}

func WithDB(DB DB) func(*RedisConfig) {
	return func(rcf *RedisConfig) {
		rcf.DB = DB
	}
}

func WihtPrefix(prefix Prefix) func(*RedisConfig) {
	return func(rcf *RedisConfig) {
		rcf.Prefix = prefix
	}
}

func NewRedisStore(options ...func(*RedisConfig)) (*Store, func(), error) {

	var redisCfg = &RedisConfig{}

	for _, option := range options {
		option(redisCfg)
	}

	cli := redis.NewClient(&redis.Options{
		Addr:     string(redisCfg.Addr),
		Password: string(redisCfg.Password),

		DialTimeout:  time.Second * 5, //闲置重新建立连接数时间 默认5s
		ReadTimeout:  time.Second * 3, //设置读超时时间,默认3s
		WriteTimeout: time.Second * 3, //设置写超时时间，默认3S

		PoolSize:     10,              //每个CPU上的默认最大连接总数, 默认是10
		MinIdleConns: 5,               //最小空闲连接数
		PoolTimeout:  time.Second * 4, //所有的连接打满后，请求超时时间 Default is ReadTimeout + 1 second
		IdleTimeout:  time.Minute * 5, //闲置请求超时时间， Default is 5 minutes. -1 disables idle timeout check.
	})
	return &Store{cli: cli, prefix: string(redisCfg.Prefix)}, func() {
		cli.Close()
	}, nil
}

type Store struct {
	cli    *redis.Client // Redis Client
	prefix string        // Wire and read data to redis store's prefix
}

// Storer 令牌存储接口
type Storer interface {
	// Get 获取
	Get(context.Context, string) (string, error)
	// Set 存储令牌数据，并指定到期时间
	Set(context.Context, string, interface{}, time.Duration) error
	// Check 检查令牌是否存在
	Check(context.Context, string) (bool, error)

	Delete(context.Context, string) (bool, error)
	// Close 关闭存储
	Close() error
}

func (s *Store) wrapperKey(k string) string {
	return fmt.Sprintf("%s%s", s.prefix, k)
}

func (s *Store) Get(ctx context.Context, k string) (string, error) {
	cmd := s.cli.Get(s.wrapperKey(k))
	return cmd.Result()
}

func (s *Store) Set(ctx context.Context, k string, v interface{}, expiration time.Duration) error {
	cmd := s.cli.Set(s.wrapperKey(k), v, expiration)
	return cmd.Err()
}

func (s *Store) Delete(ctx context.Context, k string) (bool, error) {
	cmd := s.cli.Del(s.wrapperKey(k))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *Store) Check(ctx context.Context, k string) (bool, error) {
	cmd := s.cli.Exists(s.wrapperKey(k))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *Store) Close() error {
	return s.cli.Close()
}
