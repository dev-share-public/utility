package main

import (
	"log"
	"utility/setup"
)

func main() {
	AppSetUp, err := setup.FirstSetup()
	if err != nil {
		// log.Fatal(err)
		panic("Error SetUp :: " + err.Error())
	}
	// fmt.Println("RUN http://127.0.0.1:" + AppSetUp.FiberV2RunPort)
	if err := AppSetUp.FiberV2.Listen(":" + AppSetUp.FiberV2RunPort); err != nil {
		log.Fatalf("Server Faoled to start : %v", err)
	}
}
