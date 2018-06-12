package chiller

import (
	"fmt"
	"github.com/goburrow/modbus"
	"log"
	"time"
)

func TryChillerData() {
	// Modbus RTU/ASCII
	handler := modbus.NewRTUClientHandler("/dev/ttyS1")
	handler.BaudRate = 19200
	handler.DataBits = 8
	handler.Parity = "E"
	handler.StopBits = 2
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadInputRegisters(0, 9)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("0 to 20")
	vNum := 1
	for i := 0; i < len(results); i = i + 2 {
		fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/10.0, " Celsius")
		vNum = vNum + 1
	}
	fmt.Println(results)

	results, err = client.ReadHoldingRegisters(44, 26)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("44 to 70")
	vNum = 44
	for i := 0; i < len(results); i = i + 2 {

		if i == 6 || i == 10 {
			fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/100.0, " Bar")
		} else if i == 30 {
			fmt.Println("value ", vNum, " ", (float64(results[i])*256 + float64(results[i+1])), " Hour")
		} else if i == 26 {
			fmt.Println("value ", vNum, " ", (float64(results[i])*256 + float64(results[i+1])), " times")
		} else if i == 18 {
			fmt.Println("value ", vNum, " ", (-65535.0+(float64(results[i])*256+float64(results[i+1])))/10.0, " Celsius")
		} else {
			fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/10.0)
		}
		vNum = vNum + 1
	}

	fmt.Println(results)

}
