package meter

//Dedicated for RT's server
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

//Dedicated for IM
type CpmForm struct {
	Timestamp string  `json:"lastReportTime"`
	GwId      string  `json:"GWID"`
	DevID     string  `json:"devID"`
	Get11     float64 `json:"wire"`
	Get12     float64 `json:"freq"`
	Get13     float64 `json:"ua"`
	Get14     float64 `json:"ub"`
	Get15     float64 `json:"uc"`
	Get16     float64 `json:"u_avg"`
	Get17     float64 `json:"uab"`
	Get18     float64 `json:"ubc"`
	Get19     float64 `json:"uca"`
	Get110    float64 `json:"uln_avg"`
	Get111    float64 `json:"ia"`
	Get112    float64 `json:"ib"`
	Get113    float64 `json:"ic"`
	Get114    float64 `json:"i_avg"`
	Get115    float64 `json:"pa"`
	Get116    float64 `json:"pb"`
	Get117    float64 `json:"pc"`
	Get118    float64 `json:"p_sum"`
	Get119    float64 `json:"sa"`
	Get120    float64 `json:"sb"`
	Get121    float64 `json:"sc"`
	Get122    float64 `json:"s_sum"`
	Get123    float64 `json:"pfa"`
	Get124    float64 `json:"pfb"`
	Get125    float64 `json:"pfc"`
	Get126    float64 `json:"pf_avg"`
	Get127    float64 `json:"ae_tot"`
	Get128    float64 `json:"uavg_thd"`
	Get129    float64 `json:"iavg_thd"`
}

//Dedicated for IM
type AemForm struct {
	Timestamp string  `json:"lastReportTime"`
	GwId      string  `json:"GWID"`
	DevID     string  `json:"devID"`
	Get11     float64 `json:"blockId"`
	Get12     float64 `json:"wire"`
	Get13     float64 `json:"freq"`
	Get14     float64 `json:"ua"`
	Get15     float64 `json:"ub"`
	Get16     float64 `json:"uc"`
	Get17     float64 `json:"u_avg"`
	Get18     float64 `json:"uab"`
	Get19     float64 `json:"ubc"`
	Get110    float64 `json:"uca"`
	Get111    float64 `json:"uln_avg"`
	Get112    float64 `json:"ia"`
	Get113    float64 `json:"ib"`
	Get114    float64 `json:"ic"`
	Get115    float64 `json:"i_avg"`
	Get116    float64 `json:"pa"`
	Get117    float64 `json:"pb"`
	Get118    float64 `json:"pc"`
	Get119    float64 `json:"p_sum"`
	Get120    float64 `json:"qa"`
	Get121    float64 `json:"qb"`
	Get122    float64 `json:"qc"`
	Get123    float64 `json:"q_sum"`
	Get124    float64 `json:"sa"`
	Get125    float64 `json:"sb"`
	Get126    float64 `json:"sc"`
	Get127    float64 `json:"s_sum"`
	Get128    float64 `json:"pfa"`
	Get129    float64 `json:"pfb"`
	Get130    float64 `json:"pfc"`
	Get131    float64 `json:"pf_avg"`
	Get132    float64 `json:"aea"`
	Get133    float64 `json:"aeb"`
	Get134    float64 `json:"aec"`
	Get135    float64 `json:"ae_tot"`
	Get136    float64 `json:"rea"`
	Get137    float64 `json:"reb"`
	Get138    float64 `json:"rec"`
	Get139    float64 `json:"re_tot"`
}
type ImWrap1 struct {
	CpmRow []CpmForm `json:"rows,omitempty"`
}

type ImWrap2 struct {
	AemRow []AemForm `json:"rows,omitempty"`
}

func InsertCpm(gwId string, stats []float64, timeUnix int64, postMac string, timeString string, value [32]float64, totalGen float64) DataForm {
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
	return new
}

func InsertCpmIm(gwSerial string, timeString string, value [32]float64, devId string) CpmForm {
	cpmData := CpmForm{
		Timestamp: timeString,
		GwId:      gwSerial,
		DevID:     devId,
		Get11:     value[2],
		Get12:     value[3],
		Get13:     value[4],
		Get14:     value[5],
		Get15:     value[6],
		Get16:     value[7],
		Get17:     value[8],
		Get18:     value[9],
		Get19:     value[10],
		Get110:    value[11],
		Get111:    value[12],
		Get112:    value[13],
		Get113:    value[14],
		Get114:    value[15],
		Get115:    value[16],
		Get116:    value[17],
		Get117:    value[18],
		Get118:    value[19],
		Get119:    value[20],
		Get120:    value[21],
		Get121:    value[22],
		Get122:    value[23],
		Get123:    value[24],
		Get124:    value[25],
		Get125:    value[26],
		Get126:    value[27],
		Get127:    value[28],
		Get128:    value[29],
		Get129:    value[30],
	}
	return cpmData
}

func InsertAem(gwId string, stats []float64, timeUnix int64, postMac string, timeString string, value [45]float64, totalGen float64) DataForm {
	new := DataForm{
		Timestamp:     timeString,
		TimestampUnix: timeUnix,
		MacAddress:    postMac,
		GwId:          gwId,
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

func InsertAemIm(gwSerial string, timeString string, value [45]float64, devId string) AemForm {
	aemData := AemForm{
		Timestamp: timeString,
		GwId:      gwSerial,
		DevID:     devId,
		Get11:     value[2],
		Get12:     value[3],
		Get13:     value[4],
		Get14:     value[5],
		Get15:     value[6],
		Get16:     value[7],
		Get17:     value[8],
		Get18:     value[9],
		Get19:     value[10],
		Get110:    value[11],
		Get111:    value[12],
		Get112:    value[13],
		Get113:    value[14],
		Get114:    value[15],
		Get115:    value[16],
		Get116:    value[17],
		Get117:    value[18],
		Get118:    value[19],
		Get119:    value[20],
		Get120:    value[21],
		Get121:    value[22],
		Get122:    value[23],
		Get123:    value[24],
		Get124:    value[25],
		Get125:    value[26],
		Get126:    value[27],
		Get127:    value[28],
		Get128:    value[29],
		Get129:    value[30],
		Get130:    value[31],
		Get131:    value[32],
		Get132:    value[33],
		Get133:    value[34],
		Get134:    value[35],
		Get135:    value[36],
		Get136:    value[37],
		Get137:    value[38],
		Get138:    value[39],
		Get139:    value[40],
	}
	return aemData
}

func ImTest(imUrl string) {
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
	fmt.Println(imUrl + "\n")
	res, err := http.Post(imUrl, "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {

		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("IM Post return:\n" + string(body) + "\n" + res.Status)
}
