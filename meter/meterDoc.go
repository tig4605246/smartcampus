package meter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

//FuncConf : Function inputs reorganize as a struct
//Make it more flexible to meet "UNSTABLE" changes by "SOME PEOPLE"
type FuncConf struct {
	GwSerial  string
	CpmURL    []string
	AemURL    []string
	Stats     []float64
	CpmLog    *os.File
	AemLog    *os.File
	WoodHouse *bool
	//IM use different format
	ImCpmURL []string
	ImAemURL []string
}

//MacList : We use this to parse mac data from server
type MacList struct {
	MacDatas map[string]MacData
	RawMap   interface{}
}

//GetRawMap : Get MAC mapping table from bimo's server
//If we can't fetch the new one, use the old one
func (s *MacList) GetRawMap(macURL string) error {
	res, err := http.Get(macURL)
	if err != nil {
		return errors.New("MapMac: " + err.Error())
	}
	defer res.Body.Close()

	//Check status code
	if res.StatusCode != 200 {
		return errors.New("Return code is not 200")
	}
	err = json.NewDecoder(res.Body).Decode(&s.RawMap)
	//Check Parsing status
	if err != nil {
		return err
	}
	return nil
}

//SetField : Map meter info to struct
func (s *MacList) SetField() error {
	if s.RawMap == nil {
		return errors.New("RawMap is empty")
	}
	m := s.RawMap.([]interface{})
	if len(m) < 10 {
		return errors.New("RawMap is not OK")
	}
	if len(m) == 0 {
		return errors.New("Input tmpList is empty")
	}
	for i := 0; i < len(m); i++ {
		new := MacData{}
		mSub := m[i].(map[string]interface{})
		if mSubField, ok := mSub["MAC_Address"].(string); ok {
			new.MacAddress = mSubField
		}
		if mSubField, ok := mSub["DevID"].(int); ok {
			new.DevID = mSubField
		}
		if mSubField, ok := mSub["FLOOR"].(string); ok {
			new.Floor = mSubField
		}
		if mSubField, ok := mSub["GWID"].(string); ok {
			new.GwID = mSubField
		}
		if mSubField, ok := mSub["M_GWID"].(string); ok {
			new.MGwID = mSubField
		}
		if mSubField, ok := mSub["M_MAC"].(string); ok {
			new.MMac = mSubField
		}
		if mSubField, ok := mSub["NUM"].(string); ok {
			new.Num = mSubField
		}
		if mSubField, ok := mSub["PLACE"].(string); ok {
			new.Place = mSubField
		}
		if mSubField, ok := mSub["TERRITORY"].(string); ok {
			new.Territory = mSubField
		}
		if mSubField, ok := mSub["TYPE"].(string); ok {
			new.Type = mSubField
		}
		if mSubField, ok := mSub["meter_place"].(string); ok {
			new.MeterPlace = mSubField
		}
		if mSubField, ok := mSub["node_place"].(string); ok {
			new.NodePlace = mSubField
		}
		if mSubField, ok := mSub["Group"].(string); ok {
			new.Group = mSubField
		}
		if len(new.GwID) > 12 && len(new.MMac) > 7 {
			name := new.MMac[:8] + new.GwID[12:]
			s.MacDatas[name] = new
		}

	}

	return nil
}

//ShowDatas : Print contents
func (s *MacList) ShowDatas() {
	fmt.Println("how many: ", len(s.MacDatas))
	for k, v := range s.MacDatas {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("-----")
	}
}

//ShowRaws : Print contents
func (s *MacList) ShowRaws() {
	fmt.Println("Raw: ")
	fmt.Println(s.RawMap)
}

//MacData : mac data block, used by MeterList
type MacData struct {
	MacAddress string `json:"MAC_Address"`
	DevID      int    `json:"DevID"`
	Floor      string `json:"FLOOR"`
	GwID       string `json:"GWID"`
	MGwID      string `json:"M_GWID"`
	MMac       string `json:"M_MAC"`
	Num        string `json:"NUM"`
	Place      string `json:"PLACE"`
	Territory  string `json:"TERRITORY"`
	Type       string `json:"TYPE"`
	MeterPlace string `json:"meter_place"`
	NodePlace  string `json:"node_place"`
	Group      string `json:"Group"`
}

//DataForm : Dedicated for RT's server
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
	Get130        float64 `json:"GET_1_30,omitempty"`
	Get131        float64 `json:"GET_1_31,omitempty"`
	Get132        float64 `json:"GET_1_32,omitempty"`
	Get133        float64 `json:"GET_1_33,omitempty"`
	Get134        float64 `json:"GET_1_34,omitempty"`
	Get135        float64 `json:"GET_1_35,omitempty"`
	Get136        float64 `json:"GET_1_36,omitempty"`
	Get137        float64 `json:"GET_1_37,omitempty"`
	Get138        float64 `json:"GET_1_38,omitempty"`
}

func insertCpm(gwID string, stats []float64, timeUnix int64, postMac string, timeString string, value [32]float64, totalGen float64) DataForm {
	new := DataForm{
		Timestamp:     timeString,
		TimestampUnix: timeUnix,
		MacAddress:    postMac,
		GwId:          gwID,
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
	return new
}

func insertAem(gwID string, stats []float64, timeUnix int64, postMac string, timeString string, value [45]float64, totalGen float64) DataForm {
	new := DataForm{
		Timestamp:     timeString,
		TimestampUnix: timeUnix,
		MacAddress:    postMac,
		GwId:          gwID,
		CpuRate:       stats[0],
		StorageRate:   stats[1],
		Get11:         totalGen,
		Get12:         value[3],
		Get13:         value[4],
		Get14:         value[5],
		Get15:         value[6],
		Get16:         value[7],
		Get17:         value[8],
		Get18:         value[9],
		Get19:         value[10],
		Get110:        value[11],
		Get111:        value[12],
		Get112:        value[13],
		Get113:        value[14],
		Get114:        value[15],
		Get115:        value[16],
		Get116:        value[17],
		Get117:        value[18],
		Get118:        value[19],
		Get119:        value[20],
		Get120:        value[21],
		Get121:        value[22],
		Get122:        value[23],
		Get123:        value[24],
		Get124:        value[25],
		Get125:        value[26],
		Get126:        value[27],
		Get127:        value[28],
		Get128:        value[29],
		Get129:        value[30],
		Get130:        value[31],
		Get131:        value[32],
		Get132:        value[33],
		Get133:        value[34],
		Get134:        value[35],
		Get135:        value[36],
		Get136:        value[37],
		Get137:        value[39],
		Get138:        value[40],
	}
	return new
}
