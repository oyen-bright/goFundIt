package main

import (
	"github.com/oyen-bright/goFundIt/config"
)

func main() {

	_, err := config.LoadConfig()
	if err != nil {
		panic(err)

	}

}
