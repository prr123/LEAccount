// createLEAccount.go
// program that generates Lets encrypt Account and saves keys
// author: prr azul software
// date: 21 Jan 2024
// copyright 2024 prr, azulsoftware
//
// refactor of creLEAcnt
// 29 Dec 2024 changes:
//  - change LE Account name file <name>LE<Test|Prod>.yaml

package main

import (
	"log"
	"fmt"
	"os"
	"strings"
	"context"
//	"time"

	LELib "acme/LEAccount/LELib"
    util "github.com/prr123/utility/utilLib"
)


func main() {

	numarg := len(os.Args)

    flags:=[]string{"dbg", "acnt", "type"}

	useStr := "/acnt=name /type=[prod|test] [/dbg]"
	helpStr := "help: program that creates a new account with Let's Encrypt CA\nThe account information is stored in yaml file und $LEDir!\n"

	if numarg > len(flags) + 1 {
		fmt.Printf("usage: %s %s\n",os.Args[0], useStr)
		fmt.Println("too many arguments in cl!")
		os.Exit(-1)
	}

	log.Printf("processing command line!\n")
	if numarg > 1 && os.Args[1] == "help" {
		fmt.Printf("%s", helpStr)
		fmt.Printf("usage: %s\n",useStr)
		os.Exit(1)
	}

    flagMap, err := util.ParseFlags(os.Args, flags)
	if err != nil {log.Fatalf("error -- util.ParseFlags: %v\n", err)}

	dbg := false
	_, ok := flagMap["dbg"]
	if ok {dbg = true}

	acntval, ok := flagMap["acnt"]
	if !ok { log.Fatalf("error -- acnt flag is required!\n")}
	if acntval.(string) == "none" {log. Fatalf("error -- acnt flag needs account name!\n")}
	if idx := strings.Index(acntval.(string), "."); idx > -1 {
		log.Fatalf("error -- acnt value includes an extension!\n")
	}
	acntNam := acntval.(string)

	tval, ok := flagMap["type"]
	if !ok { log.Fatalf("error -- type flag is required!\n")}
	if tval.(string) == "none" {log. Fatalf("error -- type flag needs a value!\n")}

	prod := false
	switch tval.(string) {
		case "prod": 
			prod = true
		case "test": 
			prod = false
		default:
			log.Fatalf("error -- invalid type flag value: %s!\n", tval.(string))
	}

	log.Printf("info -- account name: %s\n", acntNam)
	log.Printf("info -- prod:  %t\n", prod)
	log.Printf("info -- debug: %t\n", dbg)


    LEObj, err := LELib.InitLELib(acntNam, prod)
    if err != nil {log.Fatalf("error -- InitLELib: %v\n", err)}
	ctx := context.Background()
	LEObj.Ctx = ctx
	LEObj.Dbg = dbg
	err = LEObj.CreateLEAccount()
    if err != nil {log.Fatalf("error -- CreateLEAccount: %v\n", err)}
	fmt.Printf("created LE Account successfully!")

	acnt, err := LEObj.GetLEAccount()
    if err != nil {log.Fatalf("error -- InitLELib: %v\n", err)}
	LELib.PrintAcmeAccount(acnt)

	log.Printf("success testing account\n")
}

