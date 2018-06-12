package chiller

import (
	"github.com/goburrow/modbus"
)

func TryChillerData() {
	// Modbus RTU/ASCII
	handler := modbus.NewRTUClientHandler("/dev/ttyS1")
	handler.BaudRate = 19200
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(15, 2)

}
