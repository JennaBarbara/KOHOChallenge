package config

type Config struct {
  DB *DBConfig
  VL *VelocityLimits
  InputFile string
  OutputFile string

}

type DBConfig struct {
  Dialect  string
  Host     string
  Port     int
  Username string
  Password string
  Name     string
  Charset  string
}

type VelocityLimits struct {
  DailyLoadLimit int64
  DailyAmountLimit float64
  WeeklyAmountLimit float64
}

func GetConfig() *Config {
return &Config{
  DB: &DBConfig{
    Dialect:  "mysql",
    Host:     "127.0.0.1",
    Port:     3306,
    Username: "root",
    Password: "1234",
    Name:     "loadfundsapp",
    Charset:  "utf8",
  },
  VL: &VelocityLimits{
    DailyLoadLimit: 3,
    DailyAmountLimit: 5000.00,
    WeeklyAmountLimit: 20000.00,
  },
  InputFile: "input.txt",
  OutputFile: "output.txt",
}
}
