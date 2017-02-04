package config

import (
    "encoding/json"
    "io/ioutil"
)

// MyConfig struct
// This is the struct that the config.json must have
type MyConfig struct {
    Domain string // forum domain (if doesnt have dmain use IP)
    Port string // listen on port

    SITE_KEY string // SITE_KEY is the key to generate session and other important stuff, please change it on production

    TopicsPerPage int
    PostsPerPage int
    MaxTitleLength int
    MinTitleLength int
    MaxContentLength int
    MinContentLength int
    MaxNameLength int
    MaxDescriptionLength int
    MinDescriptionLength int

    // DB info
    DbConfig MyDBConfig

    // SMTP info
    SMTPConfig MySMTPConfig
    UserActivationEmailTemplate string // HTML template for the user activation email
}

// DB Config
type MyDBConfig struct {
    Host string
    User string
    Pass string
    Name string // DB Collection name
}

// SMTP server config to send emails
type MySMTPConfig struct {
    Host string
    Port string
    User string
    Pass string
    SenderAddress string
}

var instance *MyConfig = nil

func CreateInstance(filename string) *MyConfig {
    var err error
    instance, err = loadConfig(filename)
    if err != nil {
        // use defaults
        instance = &MyConfig{
            Domain: "localhost",
            Port: "3000",
            SITE_KEY: "Change_Me",
            TopicsPerPage: 20,
            PostsPerPage: 20,
            MaxTitleLength: 140,
            MinTitleLength: 5,
            MaxContentLength: 1000,
            MinContentLength: 5,
            MaxNameLength: 50,
            MaxDescriptionLength: 140,
            MinDescriptionLength: 5,
            DbConfig: MyDBConfig{
                Host: "localhost",
                User: "",
                Pass: "",
                Name: "GoBBit",
            },
            SMTPConfig: MySMTPConfig{
                Host: "",
                Port: "",
                User: "",
                Pass: "",
                SenderAddress: "",
            },
        }
    }

    return instance
}

func GetInstance() *MyConfig {
    return instance
}

func loadConfig(filename string) (*MyConfig, error){
    var s *MyConfig

    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return s, err
    }
    // Unmarshal json
    err = json.Unmarshal(bytes, &s)
    return s, err
}

