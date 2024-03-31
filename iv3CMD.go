package main

import (
	"fmt"
	"net"
	"os"
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

func (cameraName Camera) connString() string {
	return fmt.Sprintf("%v:%v", cameraName.ipAddress, cameraName.port)
}

func (cameraName Camera) info() {
	fmt.Printf("Camera location: %v\n", cameraName.location)
	fmt.Printf("Brief description: %v\n", cameraName.description)
	fmt.Printf("Connect to %v at ip:port : %v\n", cameraName, cameraName.connString())
}

func Iv3CmdTemplate(prefix string, arg string, cameraName Camera) string {
	//put together input string
	inputString := fmt.Sprintf("%v,%v%v", prefix, arg, cameraName.delimiter)

	conn, err := net.Dial("tcp", cameraName.connString())
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

	for {
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
		os.Exit(0)
	} else {
		fmt.Printf("Program change unsuccessful, please try again\n")
		os.Exit(1)
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
		os.Exit(0)
	} else {
		fmt.Printf("Program change unsuccessful, please try again\n")
		os.Exit(1)
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

func main() {}
