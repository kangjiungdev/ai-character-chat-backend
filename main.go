package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env 파일을 로드하지 못했습니다. (무시하고 계속 진행)")
	}

	r := SetupRouter()
	r.Run(":4000")
}
