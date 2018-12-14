//Last edited time: 2018/12/11
//Author: NTUST, BMW Lab, Xu Xi Ping
//Description: Agent for meter and chiller
package main

import (
	"flag"
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"io/ioutil"
	"log"
	"os"
	"smartcampus/airbox"
	scchiller "smartcampus/chiller"
	scmeter "smartcampus/meter"
	"strconv"
	"strings"
	"time"
)

//Default values
const (
	SCVersion  = "1.9"
	GWSerial   = "03"
	CpmURL     = "https://smartcampus.et.ntust.edu.tw:10002/meter/cpm"
	AemURL     = "https://smartcampus.et.ntust.edu.tw:10002/meter/aemdra"
	ChillerURL = "https://smartcampus.et.ntust.edu.tw:10002/chiller/rawdata"
	MacURL     = "https://smartcampus.et.ntust.edu.tw:10002/meter/devices"
	IMCpmURL   = "http://140.118.101.97:4000/cpm72/gw/data"
	IMAemURL   = "http://140.118.101.97:4000/aemdra/gw/data"
)

func main() {
	var cpmURL string
	var aemURL string
	var chillerURL string
	var cpuPath string
	var diskPath string
	var GWSerial string
	var gwID string
	var postMac string
	var imAemURL string
	var imCpmURL string
	var macURL string

	macTable := scmeter.MacList{}
	macTable.MacDatas = make(map[string]scmeter.MacData)

	flag.StringVar(&cpmURL, "cpmurl", CpmURL, "a string var")
	flag.StringVar(&aemURL, "aemurl", AemURL, "a string var")
	flag.StringVar(&chillerURL, "chillerurl", ChillerURL, "a string var")
	flag.StringVar(&cpuPath, "cpupath", "/proc/stat", "a string var")
	flag.StringVar(&diskPath, "diskpath", "/dev/mmcblk0p1", "a string var")
	flag.StringVar(&GWSerial, "gwserial", GWSerial, "a string var")
	flag.StringVar(&gwID, "gwid", "chiller_01", "a string var")
	flag.StringVar(&postMac, "postmac", "aa:bb:03:01:01:01", "a string var")
	flag.StringVar(&imAemURL, "imaemurl", IMAemURL, "a string var")
	flag.StringVar(&imCpmURL, "imcpmurl", IMCpmURL, "a string var")
	flag.StringVar(&macURL, "macurl", MacURL, "a string var")

	help := flag.Bool("help", false, "a bool")
	meter := flag.Bool("meter", false, "a bool")
	test := flag.Bool("test", false, "a bool")
	chiller := flag.Bool("chiller", false, "a bool")
	versionFlag := flag.Bool("version", false, "a bool")
	airboxTest := flag.Bool("airbox", false, "a bool")
	woodHouse := flag.Bool("woodhouse", false, "a bool")

	flag.Parse()
	d1 := []byte(strconv.Itoa(os.Getpid()))
	err := ioutil.WriteFile("/tmp/smartcampus_PID", d1, 0644)
	if err != nil {
		fmt.Println("Failed to write pid to /tmp/smartcampus")
	}

	if *help {
		fmt.Println("smartcampus Ver.", SCVersion)
		fmt.Println("For specifying gateway serial number, use -GWSerial")
		fmt.Println("For specifying url, use -cpmURL, -aemUrlm, -chillerURL")
		fmt.Println("For change mac mapping url, toggle -macurl")
		fmt.Println("For more info, please Refer to https://github.com/tig4605246/smartcampus")
		return
	} else if *versionFlag {
		version()
		return
	} else if *test {
		stats, _ := getGwStat(cpuPath, diskPath)
		functionTest(GWSerial, cpmURL, aemURL, chillerURL, stats)
		return
	} else if *airboxTest {
		fmt.Println("Now posting Airbox fake data")
		airbox.AirBox()
		return
	} else if *meter {
		//Initialize the input struct
		scConfig, res := initConf(GWSerial, cpmURL, aemURL, cpuPath, diskPath, woodHouse, imCpmURL, imAemURL)
		if res != "success" {
			fmt.Println("Error while creating config struct: ", res)
			os.Exit(0)
		}
		checkFile(scConfig.CpmLog, scConfig.AemLog)
		defer scConfig.CpmLog.Close()
		defer scConfig.AemLog.Close()
		macTable.GetRawMap(MacURL)
		macTable.SetField()
		//Parse URL of cpm and aem respectively
		timeTick := 0
		for {
			scConfig.Stats, _ = getGwStat(cpuPath, diskPath)
			scmeter.GetCpm70Data(scConfig, macTable)
			scmeter.GetAemdraData(scConfig, macTable)
			time.Sleep(30 * time.Second)
			timeTick++
			if timeTick > 120 {
				checkFile(scConfig.CpmLog, scConfig.AemLog)
				macTable.GetRawMap(MacURL)
				macTable.SetField()
				timeTick = 0
			}

		}

	} else if *chiller {

		chillerLog, err := os.Create("./chillerLog")
		if err != nil {
			fmt.Println("Can't create chillerLog")
			return
		}
		aemLog, err := os.Create("./aemLog")
		if err != nil {
			fmt.Println("Can't create aemLog")
			return
		}
		checkFile(chillerLog, aemLog)
		defer chillerLog.Close()
		defer aemLog.Close()
		for {
			stats, _ := getGwStat(cpuPath, diskPath)
			go scchiller.GetChillerData(gwID, postMac, chillerURL, stats, chillerLog)
			time.Sleep(30 * time.Second)
			checkFile(chillerLog, aemLog)
		}

	}
	return
}

func functionTest(GWSerial string, cpmURL string, aemURL string, chillerURL string, stats []float64) {
	fmt.Println("Gateway Serial:", GWSerial)
	fmt.Println("Cpu:", stats[0], " Disk:", stats[1])
	fmt.Println("url config:")
	fmt.Println("cpm url :\n", cpmURL)
	fmt.Println("aem url :\n", aemURL)
	fmt.Println("chiller url is:\n", chillerURL)
	return
}

func version() {
	fmt.Println("SmartCampus Agent Version: ", SCVersion)
}

func getGwStat(cpuPath string, diskPath string) ([]float64, int) {
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

func checkFile(cpmLog *os.File, aemLog *os.File) {
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
	return
}

func initConf(GWSerial string, cpmURL string, aemURL string, cpuPath string, diskPath string, woodHouse *bool, imCpmURL string, imAemURL string) (scmeter.FuncConf, string) {
	var err error
	scConfig := scmeter.FuncConf{}
	scConfig.CpmURL = strings.Split(cpmURL, "^")
	scConfig.AemURL = strings.Split(aemURL, "^")
	scConfig.GwSerial = GWSerial
	scConfig.Stats, _ = getGwStat(cpuPath, diskPath)
	scConfig.WoodHouse = woodHouse
	scConfig.ImCpmURL = strings.Split(imCpmURL, "^")
	scConfig.ImAemURL = strings.Split(imAemURL, "^")
	scConfig.CpmLog = new(os.File)
	scConfig.CpmLog, err = os.Create("./cpmLog")
	if err != nil {
		fmt.Println("Can't create cpmLog")
		return scConfig, "fail to create log"
	}
	scConfig.AemLog, err = os.Create("./aemLog")
	if err != nil {
		fmt.Println("Can't create aemLog")
		return scConfig, "fail to create log"
	}
	return scConfig, "success"
}
