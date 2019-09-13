package main

import (
	"log"
	"os"
	"reflect"

	"github.com/huntresslabs/win-service-updater/updater"
)

func main() {
	iuc, err := updater.ParseWYC(os.Args[1])
	if nil != err {
		log.Fatal(err)
	}
	// fmt.Printf("%+v", iuc)

	v := reflect.ValueOf(iuc)

	// values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		tlv, ok := v.Field(i).Interface().(updater.TLV)
		if !ok {
			// log.Fatal("could not covert to TLV")
			tlvArr, ok := v.Field(i).Interface().([]updater.TLV)
			if !ok {
				// log.Fatal("could not covert to TLV")
			}
			for _, tlv := range tlvArr {
				updater.DisplayTLV(&tlv)
			}
		}
		updater.DisplayTLV(&tlv)
	}
}
