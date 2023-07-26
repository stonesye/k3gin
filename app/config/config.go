package config

import (
	"encoding/json"
	"fmt"
	"github.com/koding/multiconfig"
	"os"
	"strings"
	"sync"
)

var (
	C    = new(Config)
	once sync.Once
)

// MustLoad 批量加载所有类型的配置文件内容
func MustLoad(fpaths ...string) {
	once.Do(func() {
		// 加载load加载器配置
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		// 补充其他类型的加载器
		for _, fpath := range fpaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
		}

		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}

		m.MustLoad(C)
	})
}

// PrintWithJSON 将配置文件用Json打印出来
func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", "")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error : " + err.Error())
			return
		}

		os.Stdout.WriteString(string(b) + "\n")
	}
}

type Config struct {
	RunMode     string
	WWW         string
	Swagger     bool
	Pprof       bool
	PrintConfig bool
	HTTP        HTTP
	Log         Log
	Redis       Redis
	CORS        CORS
	GZIP        GZIP
	SESSION     SESSION
	Gorm        Gorm
	RMySQL      MySQL
	WMySQL      MySQL
}

type SESSION struct {
	Enable bool
	Secret string
}

type HTTP struct {
	Host             string
	Port             int
	CertFile         string
	KeyFile          string
	ShutdownTimeout  int
	MaxContentLength int64
}

type Log struct {
	Level         int
	Format        string
	Output        string
	OutputFile    string
	RotationCount int
	RotationTime  int
}

type Redis struct {
	Addr     string
	Password string
}

type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

type GZIP struct {
	Enable             bool
	ExcludedExtentions []string
	ExcludedPaths      []string
}

type Gorm struct {
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	TablePrefix       string
	EnableAutoMigrate bool
}

type MySQL struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	Parameters string
}

// DSN 将参数拼接成数据库打开链接
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}
