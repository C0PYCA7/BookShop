package main

import (
	"BookShop/book_service/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	fmt.Print(cfg)
}
