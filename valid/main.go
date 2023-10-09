package main

import (
	"fmt"
	"gdk-test/wallet"
)

func main() {
	wallet := wallet.Wallet{}
	err := wallet.Init()
	if err != nil {
		fmt.Println(err)
	}
	mnemonic := "donkey vacuum you canoe tooth today toss brisk quick inherit faint wing lesson monitor staff host wish drift exist anchor wagon scorpion cage subway"
	err = wallet.Login(mnemonic)
	if err != nil {
		fmt.Println(err)
	}
}
