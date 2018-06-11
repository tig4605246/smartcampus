//Last edited time: 20180611
//Author: Kevin Xu Xi Ping
//Description: Agent for meter and chiller
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	//"math/rand"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"net/http"
	"strconv"
	"time"
)

//Default values
const (
	SC_VERSION  = "0.1"
	CPM_URL     = "https://beta2-api.dforcepro.com/gateway/v1/rawdata"
	AEM_URL     = "https://beta2-api.dforcepro.com/gateway/v1/rawdata"
	CHILLER_URL = "https://beta2-api.dforcepro.com/gateway/v1/rawdata"
	GWSERIAL    = "0003"
	MAC_FILE    = "./macFile"
)

// var (
// 	cpmLastTotal map[string]float64
// 	aemLastTotal map[string]float64
// )

type DataForm struct {
	Timestamp     string  `json:"Timestamp"`
	TimestampUnix int64   `json:"Timestamp_Unix"`
	MacAddress    string  `json:"MAC_Address"`
	GwId          string  `json:"GW_ID"`
	CpuRate       float64 `json:"CPU_rate"`
	StorageRate   float64 `json:"Storage_rate"`
	Get11         float64 `json:"GET_1_1"`
	//Get12         float64 `json:"GET_1_2"`
}

func main() {
	var cpmUrl string
	var aemUrl string
	var chillerUrl string
	var cpuPath string
	var diskPath string
	var gwSerial string

	flag.StringVar(&cpmUrl, "cpmurl", CPM_URL, "a string var")
	flag.StringVar(&aemUrl, "aemurl", AEM_URL, "a string var")
	flag.StringVar(&chillerUrl, "chillerurl", CHILLER_URL, "a string var")
	flag.StringVar(&cpuPath, "cpupath", "/proc/stat", "a string var")
	flag.StringVar(&diskPath, "diskpath", "/dev/mmcblk0p1", "a string var")
	flag.StringVar(&gwSerial, "gwserial", GWSERIAL, "a string var")

	help := flag.Bool("help", false, "a bool")
	meter := flag.Bool("meter", false, "a bool")
	test := flag.Bool("test", false, "a bool")
	macFile := flag.Bool("macfile", false, "a bool")

	flag.Parse()

	d1 := []byte(strconv.Itoa(os.Getpid()))
	err := ioutil.WriteFile("/tmp/smartcampus_PID", d1, 0644)
	if err != nil {
		fmt.Println("Failed to write pid to /tmp/smartcampus")
	}
	//defer f.Close()

	if *help {
		fmt.Println("For specifying url, use -cpmUrl, -aemUrlm, -chillerUrl\n")
		fmt.Println("More info, please contact Kevin Xu, Email: tig4605246@gmail.com\n")
		os.Exit(0)
	}
	if *test {
		sList := MapSerial(macFile)
		stats, _ := GetGwStat(cpuPath, diskPath)
		FunctionTest(gwSerial, cpmUrl, aemUrl, chillerUrl, sList, stats)
		os.Exit(0)
	}
	if *meter {
		cpmLog, err := os.Create("./cpmLog")
		//_, err = cpmLog.WriteString("1234")
		if err != nil {
			// handle the error here
			fmt.Println("Can't create cpmLog")
			return
		}
		aemLog, err := os.Create("./aemLog")
		//_, err = aemLog.WriteString("1234")
		if err != nil {
			// handle the error here
			fmt.Println("Can't create aemLog")
			return
		}
		CheckFile(cpmLog, aemLog)
		defer cpmLog.Close()
		defer aemLog.Close()
		for {
			sList := MapSerial(macFile)
			stats, _ := GetGwStat(cpuPath, diskPath)
			// fmt.Println("stat: ", stats)
			go GetCpm70Data(gwSerial, cpmUrl, sList, stats, cpmLog)
			//fmt.Println("cpm70 result:", msg, " ", ret)
			go GetAemdraData(gwSerial, aemUrl, sList, stats, aemLog)
			//fmt.Println("aemdra result:", msg, " ", ret)
			time.Sleep(30 * time.Second)
			CheckFile(cpmLog, aemLog)
		}

	}
	fmt.Println("Usage:\nsmartermeter [-help] [-config] [-meter] [-cpmUrl=] [-aemUrl=] [-cpuPath] [-diskPath]")
	return
}

func FunctionTest(gwSerial string, cpmUrl string, aemUrl string, chillerUrl string, sList map[string]string, stats []float64) {
	fmt.Println("Gateway Serial:", gwSerial)
	fmt.Println("Cpu:", stats[0], " Disk:", stats[1])
	fmt.Println("url config:")
	fmt.Println("cpm url :\n", cpmUrl)
	fmt.Println("aem url :\n", aemUrl)
	fmt.Println("chiller url is:\n", chillerUrl)
	fmt.Println("List meter's matching table")
	for name, val := range sList {
		fmt.Println(name, " ", val)
	}
	return
}

func Version() {
	fmt.Println("SmartCampus Agent Version: ", SC_VERSION)
}

func MapSerial(macFile *bool) map[string]string {
	sList := make(map[string]string)
	if !*macFile {
		sList["09b52f35"] = "01"
		sList["09b52f13"] = "01"
		sList["09b52f21"] = "02"
		sList["09b53b05"] = "03"
		sList["09b53b79"] = "04"
		sList["09b53b49"] = "01" //AD
		sList["09b52f1e"] = "01"
		sList["09b53b18"] = "01"
		sList["09b53b1c"] = "01"
		sList["09b53b1c"] = "01"
		sList["09b53b21"] = "01"
		sList["09b52f5a"] = "01"
		sList["09b52f02"] = "01"
		sList["09b52f47"] = "01"
		sList["09b52f48"] = "01"
		sList["09b52f10"] = "01"
		sList["09b53b30"] = "01"
		sList["09b4decb"] = "01"
		sList["09b53b30"] = "98" //test
		//99 for unknown
	} else {
		file, err := ioutil.ReadFile(MAC_FILE)
		//_, err = cpmLog.WriteString("1234")
		if err != nil {
			// handle the error here
			log.Fatal("Can't Open macFile")
		}
		subString := strings.Split(string(file), "\n")
		for i := 0; i < len(subString); i++ {
			if len(subString[i]) > 0 {
				macMap := strings.Split(subString[i], ":")
				sList[macMap[0]] = macMap[1]
			}
			//fmt.Println(i, ":", subString[i])
		}
	}

	return sList
}

func GetGwStat(cpuPath string, diskPath string) ([]float64, int) {
	cStat, err := linuxproc.ReadStat(cpuPath)
	if err != nil {
		log.Fatal("cStat read fail")
		return []float64{0, 0}, -1
	}
	dStat, err := linuxproc.ReadDisk(diskPath)
	if err != nil {
		log.Fatal("dStat read fail")
		return []float64{0, 0}, -1
	}
	return []float64{float64((cStat.CPUStats[0].Nice + cStat.CPUStats[0].System)) / float64((cStat.CPUStats[0].Nice + cStat.CPUStats[0].System + cStat.CPUStats[0].Idle)), (float64(dStat.Used*100) / float64(dStat.All))}, 0
}

func CheckFile(cpmLog *os.File, aemLog *os.File) {
	// get the cpmLog size
	stat, err := cpmLog.Stat()
	if err != nil {
		return
	}
	if stat.Size() > 1000000 {
		cpmLog.Truncate(0)
	}
	stat, err = aemLog.Stat()
	if err != nil {
		return
	}
	if stat.Size() > 1000000 {
		aemLog.Truncate(0)
	}
	//fmt.Println("File size is ", stat.Size())
	return
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
			logFile.WriteString("json:\n" + string(prettyJSON.Bytes()))
			logFile.WriteString(cpmUrl)
			res, err := http.Post(cpmUrl, "application/json", bytes.NewBuffer(jsonVal))
			if err != nil {
				logFile.WriteString("Post failed")
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
			logFile.WriteString("json:\n" + string(prettyJSON.Bytes()))
			logFile.WriteString(cpmUrl)
			res, err := http.Post(cpmUrl, "application/json", bytes.NewBuffer(jsonVal))
			if err != nil {
				logFile.WriteString("Post failed")
			}
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			logFile.WriteString("Post return:\n" + string(body) + "\n" + res.Status)

		}

	}
	return "Success", 0
}
