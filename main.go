package main

import (
	"github.com/brittandeyoung/ckia/cmd"
	_ "github.com/brittandeyoung/ckia/cmd/aws"
)

func main() {
	cmd.Execute()
}
