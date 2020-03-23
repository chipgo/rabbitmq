package util

import (
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

func BodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func LoadConfig(i interface{}) {
	envconfig.Process("", i)
}
