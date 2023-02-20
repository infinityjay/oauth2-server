package config

import (
    "io/ioutil"
    // "fmt"
    // "strings"
    "log"

    "gopkg.in/yaml.v2"
)

var Cfg App

func Setup() {
    content, err := ioutil.ReadFile("app.yaml")
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    err = yaml.Unmarshal(content, &Cfg)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
}