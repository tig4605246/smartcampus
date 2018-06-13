package meter

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	Get116        float64 `json:"GET_1_16"`
	Get117        float64 `json:"GET_1_17"`
	Get118        float64 `json:"GET_1_18"`
	Get119        float64 `json:"GET_1_19"`
	Get120        float64 `json:"GET_1_20"`
	Get121        float64 `json:"GET_1_21"`
	Get122        float64 `json:"GET_1_22"`
	Get123        float64 `json:"GET_1_23"`
	Get124        float64 `json:"GET_1_24"`
	Get125        float64 `json:"GET_1_25"`
	Get126        float64 `json:"GET_1_26"`
	Get127        float64 `json:"GET_1_27"`
	Get128        float64 `json:"GET_1_28"`
	Get129        float64 `json:"GET_1_29"`
	//Get12         float64 `json:"GET_1_2"`
}

func GetCpm70Data(gwSerial string, cpmUrl string, sList map[string]string, stats []float64, logFile *os.File) (string, int) {
	cmd := exec.Command("/home/aaeon/API/cpm70-agent-tx", "--get-dev-status")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return "cpm70-agent-tx Not found", -1
	}
	//fmt.Printf("Result:\n %s", )
	result := strings.Split(out.String(), "\n")
	line := 0
	if len(result) <= 2 {
		logFile.WriteString("agent's return value is not valid, raw message:\n" + out.String())
		return "not valid", 0
	}
	for _, m := range result {
		var postMac string
		var gwId string
		// fmt.Println("Line ", line, ":\n", m)
		line++
		subString := strings.Split(m, ";")
		// fmt.Println("Len of subString is", len(subString))
		if len(subString) >= 31 {
			// fmt.Println("get first ", subString[0], "\n split it ")

			//Format MAC and GWID

			meterSerialNum := subString[0][14:16]
			meterMac := subString[0][6:14]
			if val, ok := sList[meterMac]; ok {
				//postMac = "aa:bb:02" + ":" + subString[0][4:6] + ":" + val + ":" + meterSerialNum
				//gwId = "meter_" + subString[0][4:6] + "_" + val + "_" + meterSerialNum
				gwId = "meter_" + gwSerial + "_" + val + "_" + meterSerialNum
				postMac = "aa:bb:02" + ":" + gwSerial + ":" + val + ":" + meterSerialNum
			} else {
				//postMac = "aa:bb:02" + ":" + subString[0][4:6] + ":" + "99" + ":" + meterSerialNum
				//gwId = "meter_" + subString[0][4:6] + "_" + "99" + "_" + meterSerialNum
				postMac = "aa:bb:02" + ":" + gwSerial + ":" + "99" + ":" + meterSerialNum
				gwId = "meter_" + gwSerial + "_" + "99" + "_" + meterSerialNum
			}

			// fmt.Println("meter serial: ", meterSerialNum)
			// fmt.Println("meter Mac: ", meterMac)
			// fmt.Println("Post Mac: ", postMac)
			// fmt.Println("GW ID: ", gwId)

			//Format time

			subString[1] = subString[1][:10] + " " + subString[1][11:13] + ":" + subString[1][14:16] + ":" + subString[1][17:]
			catchTime, _ := time.Parse("2006-01-02 15:04:05", subString[1])
			timeString := catchTime.Format("2006-01-02 15:04:05")
			timeUnix := catchTime.Unix()
			// fmt.Println("time: ", timeString)
			// fmt.Println("Unix: ", timeUnix)
			totalGen, _ := strconv.ParseFloat(subString[28], 64)
			// if val, ok := cpmLastTotal[meterMac]; ok {
			// 	totalGen = totalGen - val
			// } else {
			// 	cpmLastTotal[meterMac] = totalGen
			// 	totalGen = 0
			// }
			var value [32]float64
			for i := 2; i < len(subString); i++ {
				value[i], _ = strconv.ParseFloat(subString[i], 64)
			}

			//Form JSON
			new := DataForm{
				Timestamp:     timeString,
				TimestampUnix: timeUnix,
				MacAddress:    postMac,
				GwId:          gwId,
				CpuRate:       stats[0],
				StorageRate:   stats[1],
				Get11:         totalGen,
				Get12:         value[2],
				Get13:         value[3],
				Get14:         value[4],
				Get15:         value[5],
				Get16:         value[6],
				Get17:         value[7],
				Get18:         value[8],
				Get19:         value[9],
				Get110:        value[10],
				Get111:        value[11],
				Get112:        value[12],
				Get113:        value[13],
				Get114:        value[14],
				Get115:        value[15],
				Get116:        value[16],
				Get117:        value[17],
				Get118:        value[18],
				Get119:        value[19],
				Get120:        value[20],
				Get121:        value[21],
				Get122:        value[22],
				Get123:        value[23],
				Get124:        value[24],
				Get125:        value[25],
				Get126:        value[26],
				Get127:        value[27],
				Get128:        value[29],
				Get129:        value[30],
			}
			jsonVal, err := json.Marshal(new)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			logFile.WriteString("json:\n" + string(prettyJSON.Bytes()) + "\n")
			logFile.WriteString(cpmUrl + "\n")
			res, err := http.Post(cpmUrl, "application/json", bytes.NewBuffer(jsonVal))
			if err != nil {
				logFile.WriteString("Post failed")
				logFile.WriteString(err.Error() + "\n")
				return "fail", -1
			}
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			logFile.WriteString("Post return:\n" + string(body) + "\n" + res.Status)

		}
	}
	return "Success", 0
}

func GetAemdraData(gwSerial string, cpmUrl string, sList map[string]string, stats []float64, logFile *os.File) (string, int) {
	cmd := exec.Command("/home/aaeon/API/aemdra-agent-tx", "--get-dev-status")
	var out bytes.Buffer
	var deviceList map[string]bool
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return "aemdra-agent-tx Not found", -1
	}
	//fmt.Printf("Result:\n %s", )
	result := strings.Split(out.String(), "\n")
	line := 0
	if len(result) <= 2 {
		logFile.WriteString("agent's return value is not valid, raw message:\n" + out.String())
		return "not valid", 0
	}
	for _, m := range result {
		var postMac string
		var gwId string
		// fmt.Println("Line ", line, ":\n", m)
		line++
		subString := strings.Split(m, ";")
		// fmt.Println("Len of subString is", len(subString))
		if len(subString) >= 31 {
			// fmt.Println("get first ", subString[0], "\n split it ")
			if _, ok := deviceList[subString[0]]; ok {
				continue
			} else {
				deviceList[subString[0]] = true
			}

			//Format MAC and GWID

			meterSerialNum := subString[0][14:16]
			meterMac := subString[0][6:14]
			if val, ok := sList[meterMac]; ok {
				//do something here
				//postMac = "aa:bb:02" + ":" + subString[0][4:6] + ":" + val + ":" + meterSerialNum
				//gwId = "meter_" + subString[0][4:6] + "_" + val + "_" + meterSerialNum
				gwId = "meter_" + gwSerial + "_" + val + "_" + meterSerialNum
				postMac = "aa:bb:02" + ":" + gwSerial + ":" + val + ":" + meterSerialNum
			} else {
				//postMac = "aa:bb:02" + ":" + subString[0][4:6] + ":" + "99" + ":" + meterSerialNum
				//gwId = "meter_" + subString[0][4:6] + "_" + "99" + "_" + meterSerialNum
				postMac = "aa:bb:02" + ":" + gwSerial + ":" + "99" + ":" + meterSerialNum
				gwId = "meter_" + gwSerial + "_" + "99" + "_" + meterSerialNum
			}

			// fmt.Println("meter serial: ", meterSerialNum)
			// fmt.Println("meter Mac: ", meterMac)
			// fmt.Println("Post Mac: ", postMac)
			// fmt.Println("GW ID: ", gwId)

			//Format time

			subString[1] = subString[1][:10] + " " + subString[1][11:13] + ":" + subString[1][14:16] + ":" + subString[1][17:]
			catchTime, _ := time.Parse("2006-01-02 15:04:05", subString[1])
			timeString := catchTime.Format("2006-01-02 15:04:05")
			timeUnix := catchTime.Unix()
			// fmt.Println("time: ", timeString)
			// fmt.Println("Unix: ", timeUnix)
			totalGen, _ := strconv.ParseFloat(subString[36], 64)
			// if aemLastTotal == 0 {
			// 	aemLastTotal = totalGen
			// 	totalGen = 0
			// } else {
			// 	totalGen = totalGen - aemLastTotal
			// }
			//Form JSON
			new := DataForm{
				Timestamp:     timeString,
				TimestampUnix: timeUnix,
				MacAddress:    postMac,
				GwId:          gwId,
				CpuRate:       stats[0],
				StorageRate:   stats[1],
				Get11:         totalGen,
				//Get12:	       r2,
			}
			jsonVal, err := json.Marshal(new)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			logFile.WriteString("json:\n" + string(prettyJSON.Bytes()) + "\n")
			logFile.WriteString(cpmUrl + "\n")
			res, err := http.Post(cpmUrl, "application/json", bytes.NewBuffer(jsonVal))
			if err != nil {
				logFile.WriteString("\nPost failed\n")
				logFile.WriteString(err.Error() + "\n")
				return "fail", -1
			}
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			logFile.WriteString("Post return:\n" + string(body) + "\n" + res.Status)

		}

	}
	return "Success", 0
}
