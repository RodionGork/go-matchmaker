package main

import (
	"os"

	"github.com/rodiongork/go-matchmaker/pkg/matcher"
	"github.com/rodiongork/go-matchmaker/pkg/network"
)

var m *matcher.Matcher

var usersRequestFields = map[string]string {
    "name": "string",
    "skill": "float64",
    "latency": "float64",
}

func usersRequest(body map[string]any) string {
    m.Enqueue(body["name"].(string), body["skill"].(float64), body["latency"].(float64))
    return ""
}

func main() {
	m = matcher.New()
	network.HandleJSON("/users", usersRequest, usersRequestFields)
	network.Start(os.Getenv("TCP_PORT"))
}
