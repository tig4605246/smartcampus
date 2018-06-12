//Last edited time: 20180611
//Author: Kevin Xu Xi Ping
//Description: Agent for meter and chiller
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	//"math/rand"
	linuxproc "github.com/c9s/goprocinfo/linux"
	scmeter "smartcampus/meter"
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
	chiller := flag.Bool("chiller", false, "a bool")

	flag.Parse()

	d1 := []byte(strconv.Itoa(os.Getpid()))
	err := ioutil.WriteFile("/tmp/smartcampus_PID", d1, 0644)
	if err != nil {
		fmt.Println("Failed to write pid to /tmp/smartcampus")
	}

	//defer f.Close()

	if *help {
		fmt.Println("smartcampus Ver.", SC_VERSION)
		fmt.Println("For specifying gateway serial number, use -gwserial")
		fmt.Println("For specifying url, use -cpmUrl, -aemUrlm, -chillerUrl\n")
		fmt.Println("For using mac mapping file, toggle -macfile")
		fmt.Println("More info, please contact Kevin Xu, Email: tig4605246@gmail.com\n")
		os.Exit(0)
	} else if *test {
		sList := MapSerial(macFile)
		stats, _ := GetGwStat(cpuPath, diskPath)
		FunctionTest(gwSerial, cpmUrl, aemUrl, chillerUrl, sList, stats)
		os.Exit(0)
	} else if *meter {
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
			go scmeter.GetCpm70Data(gwSerial, cpmUrl, sList, stats, cpmLog)
			//fmt.Println("cpm70 result:", msg, " ", ret)
			go scmeter.GetAemdraData(gwSerial, aemUrl, sList, stats, aemLog)
			//fmt.Println("aemdra result:", msg, " ", ret)
			time.Sleep(30 * time.Second)
			CheckFile(cpmLog, aemLog)
		}

	} else if *chiller {
		TryChillerData()
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
