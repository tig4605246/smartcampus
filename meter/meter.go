package meter

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//GetCpm70Data : The param im is normally related to Industrial Management. Meter function will send data to multiple servers.
//Each server has its own list of post address, the function will post to them one by one in for loop
func GetCpm70Data(conf FuncConf, macTable MacList) (string, int) {
	cmd := exec.Command("/home/aaeon/API/cpm70-agent-tx", "--get-dev-status")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		//log.Fatal(err)
		return "cpm70-agent-tx Not found", -1
	}
	result := strings.Split(out.String(), "\n")
	line := 0
	if len(result) <= 2 {
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
				meterMac := subString[0][6:14] + meterSerialString
				if val, ok := macTable.MacDatas[meterMac]; ok {
					//fmt.Println("Find addr: ", meterMac)
					gwID = val.GwID
					postMac = val.MacAddress
				} else {
					postMac = "aa:bb:02" + ":" + conf.GwSerial + ":" + "99" + ":" + meterSerialString
					gwID = "meter_" + conf.GwSerial + "_" + "99" + "_" + meterSerialString
				}
			}

			//Format time
			if len(subString[1]) > 17 {
				subString[1] = subString[1][:10] + " " + subString[1][11:13] + ":" + subString[1][14:16] + ":" + subString[1][17:]
			}
			catchTime, _ := time.Parse("2006-01-02 15:04:05", subString[1])
			timeString := catchTime.Format("2006-01-02 15:04:05")
			timeUnix := catchTime.Unix() - 28800
			totalGen, _ := strconv.ParseFloat(subString[28], 64)
			var value [32]float64
			for i := 2; i < len(subString); i++ {
				value[i], _ = strconv.ParseFloat(subString[i], 64)
			}

			//Post to Bimo & Peter servers
			new := insertCpm(gwID, conf.Stats, timeUnix, postMac, timeString, value, totalGen)
			jsonVal, err := json.Marshal(new)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonVal, "", "\t") // For debug print
			postToServer(jsonVal, conf.CpmURL, conf.CpmLog)
			//fmt.Println(prettyJSON.String())
			//Post to IM server
			postToServer(jsonVal, conf.ImCpmURL, conf.CpmLog)

		}
	}
	return "GetCpm70Data run Successfully", 0
}

//GetAemdraData : The param im is normally related to Industrial Management. Meter function will send data to multiple servers.
//Each server has its own list of post address, the function will post to them one by one in for loop
func GetAemdraData(conf FuncConf, macTable MacList) (string, int) {
	cmd := exec.Command("/home/aaeon/API/aemdra-agent-tx", "--get-dev-status")
	var out bytes.Buffer
	deviceList := make(map[string]bool)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "aemdra-agent-tx Not found", -1
	}
	result := strings.Split(out.String(), "\n")
	line := 0
	if len(result) <= 2 {
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
				meterMac := subString[0][6:14] + meterSerialString
				if val, ok := macTable.MacDatas[meterMac]; ok {
					gwID = val.GwID
					postMac = val.MacAddress
				} else {
					postMac = "aa:bb:02" + ":" + conf.GwSerial + ":" + "99" + ":" + meterSerialString
					gwID = "meter_" + conf.GwSerial + "_" + "99" + "_" + meterSerialString
				}
			}

			//Format time
			if len(subString[1]) > 17 {
				subString[1] = subString[1][:10] + " " + subString[1][11:13] + ":" + subString[1][14:16] + ":" + subString[1][17:]
			}
			catchTime, _ := time.Parse("2006-01-02 15:04:05", subString[1])
			timeString := catchTime.Format("2006-01-02 15:04:05")
			timeUnix := catchTime.Unix() - 28800
			totalGen, _ := strconv.ParseFloat(subString[36], 64)
			var value [45]float64
			for i := 2; i < len(subString); i++ {
				value[i], _ = strconv.ParseFloat(subString[i], 64)
			}

			//Post to Bimo's server
			new := insertAem(gwID, conf.Stats, timeUnix, postMac, timeString, value, totalGen)
			jsonVal, err := json.Marshal(new)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonVal, "", "\t")
			postToServer(jsonVal, conf.AemURL, conf.AemLog)

			//Post to IM server
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
