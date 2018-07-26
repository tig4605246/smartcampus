package meter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetCpm70Data(gwSerial string, cpmUrl string, sList map[string]string, stats []float64, logFile *os.File, woodHouse *bool, imUrl string) (string, int) {
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
			tmpMeterSerial, err := strconv.ParseInt(subString[0][14:16], 16, 32)
			if err != nil {
				logFile.WriteString("parse " + subString[0][14:16] + " to int failed\n")
			}
			var meterSerialString string
			if tmpMeterSerial < 10 {
				meterSerialString = "0" + strconv.FormatInt(tmpMeterSerial, 10)
			} else {
				meterSerialString = strconv.FormatInt(tmpMeterSerial, 10)
			}
			//logFile.WriteString("serial string: " + meterSerialString + "\n")
			if *woodHouse {
				//Wood House
				gwId = "space_02"
				postMac = "aa:bb:05:01:01:02"
			} else {
				meterMac := subString[0][6:14]
				if val, ok := sList[meterMac]; ok {
					gwId = "meter_" + gwSerial + "_" + val + "_" + meterSerialString
					postMac = "aa:bb:03" + ":" + gwSerial + ":" + val + ":" + meterSerialString
				} else {
					postMac = "aa:bb:03" + ":" + gwSerial + ":" + "99" + ":" + meterSerialString
					gwId = "meter_" + gwSerial + "_" + "99" + "_" + meterSerialString
				}
			}

			//Format time
			if len(subString[1]) > 17 {
				subString[1] = subString[1][:10] + " " + subString[1][11:13] + ":" + subString[1][14:16] + ":" + subString[1][17:]
			}
			catchTime, _ := time.Parse("2006-01-02 15:04:05", subString[1])
			timeString := catchTime.Format("2006-01-02 15:04:05")
			timeUnix := catchTime.Unix()
			totalGen, _ := strconv.ParseFloat(subString[28], 64)
			var value [32]float64
			for i := 2; i < len(subString); i++ {
				value[i], _ = strconv.ParseFloat(subString[i], 64)
			}

			//Form JSON
			new := InsertCpm(gwId, stats, timeUnix, postMac, timeString, value, totalGen)
			jsonVal, err := json.Marshal(new)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			logFile.WriteString("json:\n" + string(prettyJSON.Bytes()) + "\n")
			logFile.WriteString(cpmUrl + "\n")
			res, err := http.Post(cpmUrl, "application/json", bytes.NewBuffer(jsonVal))
			if err != nil {
				logFile.WriteString("Post failed")
				logFile.WriteString(err.Error() + "\n")
				//return "fail", -1
			}
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			logFile.WriteString("IoT Post return:\n" + string(body) + "\n" + res.Status)

			//Post to IM server
			imData := ImWrap1{}
			cpmData := InsertCpmIm(gwSerial, timeString, value, subString[0])
			imData.CpmRow = append(imData.CpmRow, cpmData)
			jsonVal, err = json.Marshal(imData)
			//var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			logFile.WriteString("json:\n" + string(prettyJSON.Bytes()) + "\n")
			logFile.WriteString(imUrl + "\n")
			res, err = http.Post(imUrl, "application/json", bytes.NewBuffer(jsonVal))
			if err != nil {
				logFile.WriteString("Post failed")
				logFile.WriteString(err.Error() + "\n")
				//return "fail", -1
			}
			defer res.Body.Close()
			body, _ = ioutil.ReadAll(res.Body)
			logFile.WriteString("IM Post return:\n" + string(body) + "\n" + res.Status)

		}
	}
	return "Success", 0
}

func GetAemdraData(gwSerial string, cpmUrl string, sList map[string]string, stats []float64, logFile *os.File, woodHouse *bool, imUrl string) (string, int) {
	cmd := exec.Command("/home/aaeon/API/aemdra-agent-tx", "--get-dev-status")
	var out bytes.Buffer
	deviceList := make(map[string]bool)
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

			tmpMeterSerial, err := strconv.ParseInt(subString[0][14:16], 16, 32)
			if err != nil {
				logFile.WriteString("parse " + subString[0][14:16] + " to int failed\n")
				continue
			}
			var meterSerialString string
			if tmpMeterSerial < 10 {
				meterSerialString = "0" + strconv.FormatInt(tmpMeterSerial, 10)
			} else {
				meterSerialString = strconv.FormatInt(tmpMeterSerial, 10)
			}
			//logFile.WriteString("serial string: " + meterSerialString + "\n")

			if *woodHouse {
				//Wood House
				gwId = "space_02"
				postMac = "aa:bb:05:01:01:02"
			} else {
				meterMac := subString[0][6:14]
				if val, ok := sList[meterMac]; ok {
					gwId = "meter_" + gwSerial + "_" + val + "_" + meterSerialString
					postMac = "aa:bb:03" + ":" + gwSerial + ":" + val + ":" + meterSerialString
				} else {
					postMac = "aa:bb:03" + ":" + gwSerial + ":" + "99" + ":" + meterSerialString
					gwId = "meter_" + gwSerial + "_" + "99" + "_" + meterSerialString
				}
			}

			//Format time

			if len(subString[1]) > 17 {
				subString[1] = subString[1][:10] + " " + subString[1][11:13] + ":" + subString[1][14:16] + ":" + subString[1][17:]
			}
			catchTime, _ := time.Parse("2006-01-02 15:04:05", subString[1])
			timeString := catchTime.Format("2006-01-02 15:04:05")
			timeUnix := catchTime.Unix()
			totalGen, _ := strconv.ParseFloat(subString[36], 64)
			var value [45]float64
			for i := 2; i < len(subString); i++ {
				value[i], _ = strconv.ParseFloat(subString[i], 64)
			}
			//Form JSON
			new := InsertAem(gwId, stats, timeUnix, postMac, timeString, value, totalGen)
			jsonVal, err := json.Marshal(new)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			logFile.WriteString("json:\n" + string(prettyJSON.Bytes()) + "\n")
			logFile.WriteString(cpmUrl + "\n")
			res, err := http.Post(cpmUrl, "application/json", bytes.NewBuffer(jsonVal))
			if err != nil {
				logFile.WriteString("\nPost failed\n")
				logFile.WriteString(err.Error() + "\n")
				//return "fail", -1
			}
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			logFile.WriteString("IoT Post return:\n" + string(body) + "\n" + res.Status)

			//Post to IM server
			imData := ImWrap2{}
			aemData := InsertAemIm(gwSerial, timeString, value, subString[0])
			imData.AemRow = append(imData.AemRow, aemData)
			jsonVal, err = json.Marshal(imData)
			//var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			logFile.WriteString("json:\n" + string(prettyJSON.Bytes()) + "\n")
			logFile.WriteString(imUrl + "\n")
			res, err = http.Post(imUrl, "application/json", bytes.NewBuffer(jsonVal))
			if err != nil {
				logFile.WriteString("Post failed")
				logFile.WriteString(err.Error() + "\n")
				//return "fail", -1
			}
			defer res.Body.Close()
			body, _ = ioutil.ReadAll(res.Body)
			logFile.WriteString("IM Post return:\n" + string(body) + "\n" + res.Status)

		}

	}
	return "Success", 0
}
