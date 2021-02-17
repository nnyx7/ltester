package main

import (
	"flag"
	"fmt"
)

// Params hold parameters needed to initialize Ltester
type Params struct {
	url    string
	method string
	// headers
	// input
	numRequests int
	duration    int
	warmUp      int
	change      int
	period      int
	respFile    string
}

func main() {
	var params Params
	setFlags(&params)
	flag.Parse()

	lt, err := ltesterFromParams(&params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(lt.execute())
	}

	XValues, YValues := fromFile(lt.respFile)
	draw(XValues, YValues, "result.png")
}
