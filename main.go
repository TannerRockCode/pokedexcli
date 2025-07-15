package main

import (
	"time"

	"github.com/TannerRockCode/pokedexcli/internal/pokecache"
)

func main() {
	t := 30 * time.Second
	pokecache.NewCache(t)
	repl()
}
