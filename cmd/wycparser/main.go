package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/huntresslabs/win-service-updater/updater"
)

func main() {
	iuc, err := updater.ParseWYC(os.Args[1])
	if nil != err {
		log.Fatal(err)
	}
	spew.Dump("%+v", iuc)

	// v := reflect.ValueOf(iuc)
	// for i := 0; i < v.NumField(); i++ {
	// 	tlv, ok := v.Field(i).Interface().(updater.TLV)
	// 	if !ok {
	// 		// log.Fatal("could not covert to TLV")
	// 		tlvArr, ok := v.Field(i).Interface().([]updater.TLV)
	// 		if !ok {
	// 			// log.Fatal("could not covert to TLV")
	// 		}
	// 		for _, tlv := range tlvArr {
	// 			updater.DisplayTLV(&tlv)
	// 		}
	// 	}
	// 	updater.DisplayTLV(&tlv)
	// }
}
