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

//GetCpm70Data : The param im is normally related to Industrial Management. Meter function will send data to multiple servers.
//Each server has its own list of post address, the function will post to them one by one in for loop
func GetCpm70Data(conf FuncConf) (string, int) {
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
		conf.CpmLog.WriteString("agent's return value is not valid, raw message:\n" + out.String())
		return "not valid", 0
	}
	for _, m := range result {
		var postMac string
		var gwID string
		line++
		subString := strings.Split(m, ";")
		if len(subString) >= 31 {
			//Format MAC and GWID
			tmpMeterSerial, err := strconv.ParseInt(subString[0][14:16], 16, 32)
			if err != nil {
				conf.CpmLog.WriteString("parse " + subString[0][14:16] + " to int failed\n")
			}
			var meterSerialString string
			if tmpMeterSerial < 10 {
				meterSerialString = "0" + strconv.FormatInt(tmpMeterSerial, 10)
			} else {
				meterSerialString = strconv.FormatInt(tmpMeterSerial, 10)
			}

			//Wood House has special config (if true, use wood house config)
			if *(conf.WoodHouse) {
				gwID = "space_02"
				postMac = "aa:bb:05:01:01:02"
			} else {
				meterMac := subString[0][6:14]
				if val, ok := conf.SList[meterMac]; ok {
					gwID = "meter_" + conf.GwSerial + "_" + val + "_" + meterSerialString
					postMac = "aa:bb:03" + ":" + conf.GwSerial + ":" + val + ":" + meterSerialString
				} else {
					postMac = "aa:bb:03" + ":" + conf.GwSerial + ":" + "99" + ":" + meterSerialString
					gwID = "meter_" + conf.GwSerial + "_" + "99" + "_" + meterSerialString
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

			//Post to Bimo & Peter servers
			new := insertCpm(gwID, conf.Stats, timeUnix, postMac, timeString, value, totalGen)
			jsonVal, err := json.Marshal(new)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			conf.CpmLog.WriteString("json:\n" + string(prettyJSON.Bytes()) + "\n")
			postToServer(jsonVal, conf.CpmURL, conf.CpmLog)

			//Post to IM server
			imData := ImWrap1{}
			cpmDataIm := insertCpmIm(conf.GwSerial, timeString, value, subString[0])
			imData.CpmRow = append(imData.CpmRow, cpmDataIm)
			jsonVal, err = json.Marshal(imData)
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			conf.CpmLog.WriteString("json:\n" + string(prettyJSON.Bytes()) + "\n")
			postToServer(jsonVal, conf.ImCpmURL, conf.CpmLog)

		}
	}
	return "GetCpm70Data run Successfully", 0
}

//GetAemdraData : The param im is normally related to Industrial Management. Meter function will send data to multiple servers.
//Each server has its own list of post address, the function will post to them one by one in for loop
func GetAemdraData(conf FuncConf) (string, int) {
	cmd := exec.Command("/home/aaeon/API/aemdra-agent-tx", "--get-dev-status")
	var out bytes.Buffer
	deviceList := make(map[string]bool)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return "aemdra-agent-tx Not found", -1
	}
	result := strings.Split(out.String(), "\n")
	line := 0
	if len(result) <= 2 {
		conf.AemLog.WriteString("agent's return value is not valid, raw message:\n" + out.String())
		return "not valid", 0
	}
	for _, m := range result {
		var postMac string
		var gwID string
		line++
		subString := strings.Split(m, ";")
		if len(subString) >= 31 {
			if _, ok := deviceList[subString[0]]; ok {
				continue
			} else {
				deviceList[subString[0]] = true
			}

			//Format MAC and GWID
			tmpMeterSerial, err := strconv.ParseInt(subString[0][14:16], 16, 32)
			if err != nil {
				conf.AemLog.WriteString("parse " + subString[0][14:16] + " to int failed\n")
				continue
			}
			var meterSerialString string
			if tmpMeterSerial < 10 {
				meterSerialString = "0" + strconv.FormatInt(tmpMeterSerial, 10)
			} else {
				meterSerialString = strconv.FormatInt(tmpMeterSerial, 10)
			}

			//Wood House has special config (if true, use wood house config)
			if *(conf.WoodHouse) {
				gwID = "space_02"
				postMac = "aa:bb:05:01:01:02"
			} else {
				meterMac := subString[0][6:14]
				if val, ok := conf.SList[meterMac]; ok {
					gwID = "meter_" + conf.GwSerial + "_" + val + "_" + meterSerialString
					postMac = "aa:bb:03" + ":" + conf.GwSerial + ":" + val + ":" + meterSerialString
				} else {
					postMac = "aa:bb:03" + ":" + conf.GwSerial + ":" + "99" + ":" + meterSerialString
					gwID = "meter_" + conf.GwSerial + "_" + "99" + "_" + meterSerialString
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

			//Post to Bimo & Peter servers
			new := insertAem(gwID, conf.Stats, timeUnix, postMac, timeString, value, totalGen)
			jsonVal, err := json.Marshal(new)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			conf.AemLog.WriteString("json:\n" + string(prettyJSON.Bytes()) + "\n")
			postToServer(jsonVal, conf.AemURL, conf.AemLog)

			//Post to IM server
			imData := ImWrap2{}
			aemData := insertAemIm(conf.GwSerial, timeString, value, subString[0])
			imData.AemRow = append(imData.AemRow, aemData)
			jsonVal, err = json.Marshal(imData)
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			conf.AemLog.WriteString("json:\n" + string(prettyJSON.Bytes()) + "\n")
			postToServer(jsonVal, conf.ImAemURL, conf.AemLog)

		}

	}
	return "GetAemdraData run Successfully", 0
}

func postToServer(jsonVal []byte, URL []string, logFile *os.File) {
	for i := 0; i < len(URL); i++ {
		res, err := http.Post(URL[i], "application/json", bytes.NewBuffer(jsonVal))
		if err != nil {
			logFile.WriteString("Post failed")
			logFile.WriteString(err.Error() + "\n")
		} else {
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			logFile.WriteString(URL[i] + " Post return:\n" + string(body) + "\n" + res.Status + "\n")
		}

	}
	return

}

//ImTest : Testing function. To help IM test their server
func ImTest(imURL string) {
	//Post to IM server
	imData := ImWrap1{}
	cpmData := CpmForm{
		Timestamp: "2006-01-02 15:04:06",
		GwId:      "IIC3NTUST-0005",
		DevID:     "33000509b53b300e",
		Get11:     1,
		Get12:     2,
		Get13:     3,
		Get14:     4,
		Get15:     5,
		Get16:     6,
		Get17:     1,
		Get18:     1,
		Get19:     1,
		Get110:    1,
		Get111:    1,
		Get112:    1,
		Get113:    1,
		Get114:    1,
		Get115:    1,
		Get116:    1,
		Get117:    1,
		Get118:    1,
		Get119:    1,
		Get120:    1,
		Get121:    1,
		Get122:    1,
		Get123:    1,
		Get124:    1,
		Get125:    1,
		Get126:    1,
		Get127:    1,
		Get128:    1,
		Get129:    1,
	}
	imData.CpmRow = append(imData.CpmRow, cpmData)
	jsonVal, err := json.Marshal(imData)
	fmt.Println(string(jsonVal))
	fmt.Println(imData.CpmRow)
	if err != nil {

		fmt.Println(err.Error())
	}
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonVal, "", "\t")
	if err != nil {

		fmt.Println(err.Error())
	}
	fmt.Println("json:\n" + string(prettyJSON.Bytes()) + "\n")
	fmt.Println(imURL + "\n")
	res, err := http.Post(imURL, "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {

		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("IM Post return:\n" + string(body) + "\n" + res.Status)
}
