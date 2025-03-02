package config

import (
	"os"
	"path"
)

const Version = "0.1.0"

var BasePath = os.ExpandEnv("$HOME/.loggernaut-cli")
var OutboxPath = path.Join(BasePath, "outbox")
