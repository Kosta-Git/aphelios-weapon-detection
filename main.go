package main

import (
	"aphelios/pkg"
	"context"
	"fmt"
)

func main() {
	weaponWatchContext, weaponWatchCtxCancel := context.WithCancel(context.Background())
	defer weaponWatchCtxCancel()
	weapons := pkg.WatchApheliosWeapons("./assets/spells/", 1, weaponWatchContext)
	for {
		select {
		case weapon := <-weapons:
			fmt.Printf("%v", weapon)
		}
	}
}
