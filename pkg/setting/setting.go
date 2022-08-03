package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret       string
	ExpiredJwt      int
	PageSize        int
	PrefixUrl       string
	RuntimeRootPath string
	ImageSavePath   string
	ImageMaxSize    int
	ImageAllowExts  []string
	ExportSavePath  string
	QrCodeSavePath  string
	FontSavePath    string
	LogSavePath     string
	LogSaveName     string
	LogFileExt      string
	TimeFormat      string
	Issuer          string
	SaltKey         string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	GrpcPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Smtp struct {
	Server   string
	Port     int
	User     string
	Password string
	Identity string
	Sender   string
}

var SmtpSetting = &Smtp{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Port        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type RedisDB struct {
	Host     string
	Port     int
	DB       int
	Key      string
	Password string
}

var RedisDBSetting = &RedisDB{}

type S3 struct {
	SpaceName   string
	SpaceKey    string
	SpaceSecret string
	SpaceBucket string
	Region      string
	EndPoint    string
}

var S3Setting = &S3{}

type Agora struct {
	AppID          string
	AppCertificate string
}

var AgoraSetting = &Agora{}

type ElasticSearch struct {
	ElasticHost         string
	ElasticAuth         bool
	ElasticAuthUsername string
	ElasticAuthPassword string
}

var ElasticSearchSetting = &ElasticSearch{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("config.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'config.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisDBSetting)
	mapTo("smtp", SmtpSetting)
	mapTo("S3", S3Setting)
	mapTo("Agora", AgoraSetting)
	mapTo("Elastic", ElasticSearchSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
