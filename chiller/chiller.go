package chiller

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"github.com/goburrow/modbus"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type DataForm struct {
	Timestamp     string  `json:"Timestamp"`
	TimestampUnix int64   `json:"Timestamp_Unix"`
	MacAddress    string  `json:"MAC_Address"`
	GwId          string  `json:"GW_ID"`
	CpuRate       float64 `json:"CPU_rate"`
	StorageRate   float64 `json:"Storage_rate"`
	Get11         float64 `json:"GET_1_1"`
	Get12         float64 `json:"GET_1_2"`
	Get13         float64 `json:"GET_1_3"`
	Get14         float64 `json:"GET_1_4"`
	Get15         float64 `json:"GET_1_5"`
	Get16         float64 `json:"GET_1_6"`
	Get17         float64 `json:"GET_1_7"`
	Get18         float64 `json:"GET_1_8"`
	Get19         float64 `json:"GET_1_9"`
	Get110        float64 `json:"GET_1_10"`
	Get111        float64 `json:"GET_1_11"`
	Get112        float64 `json:"GET_1_12"`
	Get113        float64 `json:"GET_1_13"`
	Get114        float64 `json:"GET_1_14"`
	Get115        float64 `json:"GET_1_15"`
}

func GetChillerData(gwId string, postMac string, chillerUrl string, stats []float64, logFile *os.File) {
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
		logFile.WriteString(err.Error())
		return
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadInputRegisters(0, 8)
	if err != nil {
		logFile.WriteString(err.Error())
		return
	}
	//logFile.WriteString("0 to 7\n")
	var value [20]float64
	getCount := 0
	vNum := 1
	for i := 0; i < len(results); i = i + 2 {
		//fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/10.0, " Celsius")
		value[getCount] = (float64(results[i])*256 + float64(results[i+1])) / 10.0
		getCount = getCount + 1
		vNum = vNum + 1
	}
	//fmt.Println(results)

	results, err = client.ReadHoldingRegisters(44, 26)
	if err != nil {
		log.Fatal(err)
		return
	}
	//fmt.Println("44 to 69")
	vNum = 44
	for i := 0; i < len(results); i = i + 2 {

		if i == 0 {
			//機電流
			//fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/10.0, " A")
			value[getCount] = (float64(results[i])*256 + float64(results[i+1])) / 10.0
			getCount = getCount + 1
		} else if i == 6 {
			//高壓 顯示
			//fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/100.0, " Bar")
			value[getCount] = (float64(results[i])*256 + float64(results[i+1])) / 100.0
			getCount = getCount + 1
		} else if i == 10 {
			//低壓 顯示
			//fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/100.0, " Bar")
		} else if i == 14 {
			//飽和蒸發
			//fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/10.0, " Celsius")
			value[getCount] = (float64(results[i])*256 + float64(results[i+1])) / 10.0
			getCount = getCount + 1
		} else if i == 16 {
			//膨脹開度
			//fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/10.0, " %")
			value[getCount] = (float64(results[i])*256 + float64(results[i+1])) / 10.0
			getCount = getCount + 1
		} else if i == 18 {
			//過熱度
			if results[i] == 255 {
				//fmt.Println("value ", vNum, " ", (-65535.0+(float64(results[i])*256+float64(results[i+1])))/10.0, " Celsius")
				value[getCount] = (-65535.0 + (float64(results[i])*256 + float64(results[i+1]))) / 10.0
			} else {
				//fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/10.0, " Celsius")
				value[getCount] = (float64(results[i])*256 + float64(results[i+1])) / 10.0
			}
			getCount = getCount + 1
		} else if i == 26 {
			//啟動次數
			//fmt.Println("value ", vNum, " ", (float64(results[i])*256 + float64(results[i+1])), " times")
			value[getCount] = (float64(results[i])*256 + float64(results[i+1]))
			getCount = getCount + 1
		} else if i == 30 {
			//運轉時數
			//fmt.Println("value ", vNum, " ", (float64(results[i])*256 + float64(results[i+1])), " Hour")
			value[getCount] = (float64(results[i])*256 + float64(results[i+1]))
			getCount = getCount + 1
		} else if i == 50 {
			//冷凝溫度
			//fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/10.0, " Celsius")
			//value[getCount] = (float64(results[i])*256 + float64(results[i+1])) / 10.0
			//getCount = getCount + 1
		} else {
			//fmt.Println("value ", vNum, " ", (float64(results[i])*256+float64(results[i+1]))/10.0)
		}
		vNum = vNum + 1
	}
	//fmt.Println(results)

	//Post them

	//Get time
	catchTime := time.Now()
	timeString := catchTime.Format("2006-01-02 15:04:05")
	timeUnix := catchTime.Unix()
	//Form JSON
	new := DataForm{
		Timestamp:     timeString,
		TimestampUnix: timeUnix,
		MacAddress:    postMac,
		GwId:          gwId,
		CpuRate:       stats[0],
		StorageRate:   stats[1],
		Get11:         value[0],
		Get12:         value[1],
		Get13:         value[2],
		Get14:         value[3],
		Get15:         value[4],
		Get16:         value[5],
		Get17:         value[6],
		Get18:         value[7],
		Get19:         value[8],
		Get110:        value[9],
		Get111:        value[10],
		Get112:        value[11],
		Get113:        value[12],
		Get114:        value[13],
		Get115:        value[14],
	}
	jsonVal, err := json.Marshal(new)
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonVal, "", "\t")
	logFile.WriteString(string(prettyJSON.Bytes()) + "\n")
	res, err := http.Post(chillerUrl, "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {
		logFile.WriteString("Post failed" + "\n")
		logFile.WriteString(err.Error() + "\n")
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	logFile.WriteString("Post return:\n" + string(body) + "\n" + res.Status + "\n")

}
