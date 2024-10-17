package main

import (
	"github.com/christopherhanke/bootdev_gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("datawulf")
	cfg = config.Read()

}
