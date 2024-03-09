package main

import (
	"MCOnebot/pkg"
	"MCOnebot/pkg/minecraft"
	"fmt"
)

func main() {
	err := pkg.Init()
	if err != nil {
		return
	}
	msauth, err := minecraft.GetMCcredentials("data/msauth.json", "88650e7e-efee-4857-b9a9-cf580a00ef43")
	if err != nil {
		panic(err)
	}
	fmt.Println(msauth.Name)
}
