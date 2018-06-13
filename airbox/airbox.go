package airbox

import (
	"fmt"
	"io/ioutil"
	//    "log"
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	//"strings"
	"time"
)

type Airbox struct {
	//Timestamp     string  `json:"Timestamp"`
	//TimestampUnix int64   `json:"Timestamp_Unix"`
	MacAddress  string  `json:"MAC_Address"`
	GwId        string  `json:"GW_ID"`
	CpuRate     float64 `json:"CPU_rate"`
	StorageRate int     `json:"Storage_rate"`
	Get11       float64 `json:"GET_1_1"`
	Get12       float64 `json:"GET_1_2"`
	//Get21         float64 `json:"GET_2_1"`
	//Get22         float64 `json:"GET_2_2"`
	//Get23         float64 `json:"GET_2_3"`
	//Set11         float64 `json:"SET_1_1"`
	//Set12         float64 `json:"SET_1_2"`
}

func AirBox() {
	for {
		for i := 0; i < 4; i++ {
			go AirPost(i)
		}
		time.Sleep(10 * time.Second)
	}
}

func AirPost(gwid int) {
	mac := []string{"aa:bb:cc:11:11:18", "aa:bb:cc:11:11:17", "aa:bb:cc:11:11:14", "aa:bb:cc:11:11:12"}
	id := []string{"airbox_18", "airbox_17", "airbox_14", "airbox_12"}

	//nowTime := time.Now()
	//timeString := nowTime.Format("2006-01-02 15:04:05")
	//timeUnix := nowTime.Unix()
	r1 := rand.Float64() * 50.0
	r2 := rand.Float64() * 50.0
	//r3 := rand.Float64() * 50.0
	//r4 := rand.Float64() * 50.0
	//r5 := rand.Float64() * 50.0
	//r6 := 0.0
	new := Airbox{
		//Timestamp:     timeString,
		//TimestampUnix: timeUnix,
		MacAddress:  mac[gwid],
		GwId:        id[gwid],
		CpuRate:     1.0,
		StorageRate: 1,
		Get11:       r1,
		Get12:       r2,
		//Get22:         r3,
		//Get23:         r4,
		//Set11:         r5,
		//Set12:         r6,
	}
	//fmt.Println(new)
	jsonVal, err := json.Marshal(new)
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonVal, "", "\t")
	fmt.Println("json:\n", string(prettyJSON.Bytes()))
	res, err := http.Post("https://beta2-api.dforcepro.com/gateway/v2/rawdata", "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {
		fmt.Println("Post failed")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("Post return:\n", string(body))

}
