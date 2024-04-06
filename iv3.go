package iv3

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Camera struct {
	location    string
	description string
	ipAddress   string
	port        string
	delimiter   string
	//programList []program //Leaving note here to eventually add configurable progams to help with mapping and error handling
}

/* Potential new struct to describe iv3 program
type iv3Program struct {
	programNumber int
	description string
	alias string
	toolList []tool
}

//Potential new struct to describe individual tools within an iv3 program
type programTool struct {
	toolNumber int
	toolType int //may need to create lookup table to map tool types to int. Likely in order from user manual - see page 4-28.
				 //This will be useful when it comes time to convert matching rate field into something practical for use across many tools...
	useUpperLimit bool
	threshold int
}

func (programTool programTool) fitMatchRate()? {
	return fmt.Sprintf("%v:%v", cameraName.ipAddress, cameraName.port)
}
*/

type CameraResult struct {
	resultNumber    int
	totalPassResult bool
	toolResult      []ToolResult
}

type ToolResult struct {
	toolNumber     int
	toolPassResult bool
	matchingRate   int
}

func mapJudgement(judgement string) bool {
	var output bool
	switch judgement {
	case "OK":
		output = true
	case "NG":
		output = false
	}
	return output
}

func parseCameraResult(cameraResult []string) CameraResult {
	//fmt.Println("Camera result string: ", cameraResult)
	judgement := mapJudgement(cameraResult[2]) //map input string into bool to fit result struct

	var toolResults []ToolResult //attempt to loop over left over tool results and add to tool result array
	toolResultOnly := cameraResult[3:]
	//fmt.Println("Tool result string: ", toolResultOnly)
	toolCount := (len(toolResultOnly) / 3)
	for i := range toolCount {
		//Need to call strconv func to turn given strings to proper types for struct
		//fmt.Printf("i: %v\n", i)
		toolNumber, err := strconv.Atoi(toolResultOnly[3*i])
		if err != nil {
			fmt.Printf("Unable to interpret tool number response. Try again or escalate to engineering.\n Reported tool number: %v\n", toolResultOnly[3*i])
			panic(err)
		}
		matchingRate, err := strconv.Atoi(toolResultOnly[(3*i)+2])
		if err != nil {
			fmt.Printf("Unable to matching rate response. Try again or escalate to engineering.\n Reported matching rate: %v\n", toolResultOnly[(3*i)+2])
			panic(err)
		}
		mapResults := ToolResult{
			toolNumber:     toolNumber,
			toolPassResult: mapJudgement(toolResultOnly[(3*i)+1]),
			matchingRate:   matchingRate,
		}
		toolResults = append(toolResults, mapResults)
	}
	resultNumber, err := strconv.Atoi(cameraResult[1])
	if err != nil {
		fmt.Printf("Unable to interpret result number response. Try again or escalate to engineering.\n Reported result number: %v\n", cameraResult[1])
		panic(err)
	}
	outputResult := CameraResult{
		resultNumber:    resultNumber,
		totalPassResult: judgement,
		toolResult:      toolResults,
	}
	return outputResult
}

func (cameraName Camera) ConnString() string {
	return fmt.Sprintf("%v:%v", cameraName.ipAddress, cameraName.port)
}

func (cameraName Camera) Info() {
	fmt.Printf("Camera location: %v\n", cameraName.location)
	fmt.Printf("Brief description: %v\n", cameraName.description)
	fmt.Printf("Connect to %v at ip:port : %v\n", cameraName, cameraName.ConnString())
}

func Iv3CmdTemplate(prefix string, arg string, cameraName Camera) string {
	//put together input string
	inputString := fmt.Sprintf("%v,%v%v", prefix, arg, cameraName.delimiter)

	conn, err := net.Dial("tcp", cameraName.ConnString())
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}
	defer conn.Close()

	//send data
	data := []byte(inputString)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}
	//read response
	buffer := make([]byte, 1024)

	//read data from client
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}
	//response := string(buffer[:n])
	var response = string(buffer[:n])
	conn.Close()
	responseTrim := strings.Trim(response, "\r")
	return responseTrim
}

func ProgramChange(programNumber int, cameraName Camera) {
	//configure prefix and input of command based on template starting point
	prefix := "PW"
	arg := fmt.Sprint(programNumber) //convert program # to str
	response := Iv3CmdTemplate(prefix, arg, cameraName)
	//handle expected response and act accordingly
	responseSplit := strings.Split(response, ",") //responseSplit[0] always returns prefix. Rest of msg will depend on CMD - see manual!
	//handle expected response and act accordingly
	if responseSplit[0] == "PW" {
		fmt.Printf("Program Change Successful.\n")
	} else {
		fmt.Printf("Program change unsuccessful, please try again\n")
	}
}

func ReadProgramNumber(cameraName Camera) {
	//configure prefix and input of command based on template starting point
	prefix := "PR"
	arg := "" //cmd takes no arg, leave blank
	response := Iv3CmdTemplate(prefix, arg, cameraName)
	//handle expected response and act accordingly
	responseSplit := strings.Split(response, ",") //responseSplit[0] always returns prefix. Rest of msg will depend on CMD - see manual!
	if responseSplit[0] == "PR" {
		fmt.Printf("Program Read Successful.\n")
		fmt.Printf("Current Program Number: %v\n", responseSplit[1])
	} else {
		fmt.Printf("Program change unsuccessful, please try again\n")
	}
}

func TriggerStatusResult(cameraName Camera) CameraResult {
	//configure prefix and input of command based on template starting point
	prefix := "T2"
	arg := "" //cmd takes no arg, leave blank
	response := Iv3CmdTemplate(prefix, arg, cameraName)
	//handle expected response and act accordingly
	responseSplit := strings.Split(response, ",") //responseSplit[0] always returns prefix. Rest of msg will depend on CMD - see manual!
	/*
		Camerea judgement response is complex, summary below but see manual for full details:
		"RT,aa,bb,cc,dd,eeeeee,cc,dd,eeeeee,cc,dd,eeeeee..."
		responseSplit[0] - prefix response "RT"
		responseSplit[1] - result number from 0 to 32767, in theory should be same as trigger # - "aaaaa"
		responseSplit[2] - overall camera reseult, either "OK" or "NG" for program pass or fail - "bb"
		Start of individual tool results - each tool has 3 comma seperated values ("cc", "dd", "eeeeee")
		The tool result pattern will repeat matching the number of tools loaded in the current camera program up to a max 64 tools.
		responseSplit[3] - tool number from 01 to 64 - "cc"
		responseSplit[4] - individual tool result either "OK" or "NG" - "dd"
		responseSplit[5] - tool matching rate, format depends on type of tool returning results. Too many possibilites to explain here - "eeeeee"
	*/
	if responseSplit[0] == "RT" {
		fmt.Printf("Camera trigger successful.\n")
		fmt.Printf("Current result number: %v\n", responseSplit[1])
	} else {
		fmt.Printf("camera trigger unsuccessful, please try again\n")
	}
	cameraResult := parseCameraResult(responseSplit)
	return cameraResult
}

func OperatingStatus(cameraName Camera) bool { // is the camera currently in run mode or program mode?
	//configure prefix and input of command based on template starting point
	prefix := "RM"
	arg := "" //cmd takes no arg, leave blank
	response := Iv3CmdTemplate(prefix, arg, cameraName)
	//handle expected response and act accordingly
	responseSplit := strings.Split(response, ",")

	operating, err := strconv.ParseBool(responseSplit[1])
	if responseSplit[0] == "RM" && err == nil {
		fmt.Printf("Status reading successful. Camera in operating mode: %v", operating)

	} else {
		fmt.Printf("Status reading unsuccessful. Try again.")
	}
	return operating
}

type SensorState struct {
	busy bool
	//second bit reserved for system - will always return 0 - discard it!
	imageCapture bool
	sdCard       bool
	sdCardFull   bool
	warning      bool
	err          bool
}

func SensorStatus(cameraName Camera) SensorState {
	//configure prefix and input of command based on template starting point
	prefix := "SR"
	arg := "" //cmd takes no arg, leave blank
	response := Iv3CmdTemplate(prefix, arg, cameraName)
	//handle expected response and act accordingly
	responseSplit := strings.Split(response, ",")
	var sensorState SensorState
	if responseSplit[0] == "SR" {
		fmt.Printf("Status check successul.")
		if responseSplit[1] == "1" {
			sensorState.busy = true
		}
		//skip responseSplit[2] - bit reserved by camera system
		if responseSplit[3] == "1" {
			sensorState.imageCapture = true
		}
		if responseSplit[4] == "1" {
			sensorState.sdCard = true
		}
		if responseSplit[5] == "1" {
			sensorState.sdCardFull = true
		}
		if responseSplit[6] == "1" {
			sensorState.warning = true
		}
		if responseSplit[7] == "1" {
			sensorState.err = true
		}
	}
	return sensorState
}

type ProgramStats struct {
	maxTime      int
	minTime      int
	avgTime      int
	trigCount    int
	okCount      int
	ngCount      int
	trigErrCount int
	toolStats    []ToolStat
}

type ToolStat struct {
	toolNumber      int
	maxMatchingRate int
	minMatchingRate int
}

func StatInfoReading(cameraName Camera) ProgramStats {
	//configure prefix and input of command based on template starting point
	prefix := "STR"
	arg := "" //cmd takes no arg, leave blank
	response := Iv3CmdTemplate(prefix, arg, cameraName)
	//handle expected response and act accordingly
	responseSplit := strings.Split(response, ",")
	if responseSplit[0] == "STR" {
		fmt.Printf("Camera statistics reading successful.\n")
	} else {
		fmt.Printf("camera statistics reading unsuccessful, please try again\n")
	}
	maxTime, err := strconv.Atoi(responseSplit[1])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported max processing time: %v\n", responseSplit[1])
		panic(err)
	}
	minTime, err := strconv.Atoi(responseSplit[2])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported min processing time: %v\n", responseSplit[2])
		panic(err)
	}
	avgTime, err := strconv.Atoi(responseSplit[3])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported average processing time: %v\n", responseSplit[3])
		panic(err)
	}
	trigCount, err := strconv.Atoi(responseSplit[4])
	if err != nil {
		fmt.Printf("Unable to interpretresponse. Try again or escalate to engineering.\n Reported triger count: %v\n", responseSplit[4])
		panic(err)
	}
	okCount, err := strconv.Atoi(responseSplit[5])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported ok count: %v\n", responseSplit[5])
		panic(err)
	}
	ngCount, err := strconv.Atoi(responseSplit[6])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported not good count: %v\n", responseSplit[6])
		panic(err)
	}
	trigErrCount, err := strconv.Atoi(responseSplit[7])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported tool number: %v\n", responseSplit[7])
		panic(err)
	}
	var toolResult []ToolStat
	toolResultOnly := responseSplit[7:]
	toolCount := (len(toolResultOnly) / 3)
	for i := range toolCount {
		//Need to call strconv func to turn given strings to proper types for struct
		//fmt.Printf("i: %v\n", i)
		toolNumber, err := strconv.Atoi(toolResultOnly[3*i])
		if err != nil {
			fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported tool number: %v\n", toolResultOnly[3*i])
			panic(err)
		}
		maxMatchingRate, err := strconv.Atoi(toolResultOnly[3*i])
		if err != nil {
			fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported max matching rate: %v\n", toolResultOnly[(3*i)+1])
			panic(err)
		}
		minMatchingRate, err := strconv.Atoi(toolResultOnly[(3*i)+2])
		if err != nil {
			fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported min matching rate: %v\n", toolResultOnly[(3*i)+2])
			panic(err)
		}
		mapResults := ToolStat{
			toolNumber:      toolNumber,
			maxMatchingRate: maxMatchingRate,
			minMatchingRate: minMatchingRate,
		}
		toolResult = append(toolResult, mapResults)
	}
	programStats := ProgramStats{
		maxTime:      maxTime,
		minTime:      minTime,
		avgTime:      avgTime,
		trigCount:    trigCount,
		okCount:      okCount,
		ngCount:      ngCount,
		trigErrCount: trigErrCount,
		toolStats:    toolResult,
	}
	return programStats
}

func ThresholdWrite(cameraName Camera, toolNumber int, thresholdValue int) int {
	//configure prefix and input of command based on template starting point
	prefix := "DW"
	arg := fmt.Sprintf("%v,0,%v", toolNumber, thresholdValue) //prefix, args delimiter
	response := Iv3CmdTemplate(prefix, arg, cameraName)
	//handle expected response and act accordingly
	responseSplit := strings.Split(response, ",")
	if responseSplit[0] == "DW" {
		fmt.Printf("Program threshold change successful.\n")
	} else {
		fmt.Printf("Program threshold change unsuccessful, please try again\n")
	}
	toolNumber, err := strconv.Atoi(responseSplit[2])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported min matching rate: %v\n", responseSplit[2])
		panic(err)
	}
	return toolNumber
}

type ToolThreshold struct {
	toolNumber     int
	upperLimit     bool
	thresholdValue int
}

func ThresholdRead(cameraName Camera, toolNumber int) ToolThreshold {
	//configure prefix and input of command based on template starting point
	prefix := "DR"
	arg := fmt.Sprintf("%v,0", toolNumber) //prefix, args delimiter
	response := Iv3CmdTemplate(prefix, arg, cameraName)
	//handle expected response and act accordingly
	responseSplit := strings.Split(response, ",")
	if responseSplit[0] == "DR" {
		fmt.Printf("Program threshold change successful.\n")
	} else {
		fmt.Printf("Program threshold change unsuccessful, please try again\n")
	}
	toolNumber, err := strconv.Atoi(responseSplit[1])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported tool number time: %v\n", responseSplit[1])
		panic(err)
	}
	upperLimit, err := strconv.ParseBool(responseSplit[2])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported limit type: %v\n", responseSplit[2])
		panic(err)
	}
	thresholdValue, err := strconv.Atoi(responseSplit[3])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported threshold value: %v\n", responseSplit[3])
		panic(err)
	}
	ToolThreshold := ToolThreshold{
		toolNumber:     toolNumber,
		upperLimit:     upperLimit,
		thresholdValue: thresholdValue,
	}
	return ToolThreshold
}

func StatReset(cameraName Camera) {
	//configure prefix and input of command based on template starting point
	prefix := "STC"
	arg := "" //no args
	response := Iv3CmdTemplate(prefix, arg, cameraName)
	//handle expected response and act accordingly
	responseSplit := strings.Split(response, ",")
	if responseSplit[0] == "STC" {
		fmt.Printf("Statistics reset successful.\n")
	} else {
		fmt.Printf("Statistics reset unsuccessful, please try again\n")
	}
}

func ErrorRead(cameraName Camera) int {
	//configure prefix and input of command based on template starting point
	prefix := "RER"
	arg := "" //no args
	response := Iv3CmdTemplate(prefix, arg, cameraName)
	//handle expected response and act accordingly
	responseSplit := strings.Split(response, ",")
	if responseSplit[0] == "RER" {
		fmt.Printf("Statistics reset successful.\n")
	} else {
		fmt.Printf("Statistics reset unsuccessful, please try again\n")
	}
	errNumber, err := strconv.Atoi(responseSplit[1])
	if err != nil {
		fmt.Printf("Unable to interpret response. Try again or escalate to engineering.\n Reported error number: %v\n", responseSplit[1])
		panic(err)
	}
	return errNumber
}
