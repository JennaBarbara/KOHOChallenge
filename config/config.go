package config

type Config struct {
DB *DBConfig
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
  InputFile: "input.txt",
  OutputFile: "output.txt",
}
}
