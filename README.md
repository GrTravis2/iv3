Package of commands that operate Keyence IV3 cameras using TCP/IP procedures written in GO!

## Initialize a camera:
```go
myCamera := Camera {
    location    string,
    description string,
    ipAddress   string,
    port        string, //iv3 defaults to port "8500"
    delimiter   string, //iv3 defaults to carriage return "\r\n"
    }
```
## Use the camera object and additional arguments to "do stuff"

### Changing programs
```go
var programNumber int
var myCamera Camera
var response string

response = ProgramChange(programNumber, myCamera)
//reads response and returns success/unsuccessful
```

### Read current program
```go
var myCamera Camera
var response string

response = ReadProgramNumber(myCamera)
```



### Command camera to trigger and return program results, **this one is more complicated...**

two more structs are used in the response for this command:
```go
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
```
**Results are returned in the form of a camera result struct made of summary values and individual tool results**
```go
var myCamera Camera
var cameraResult CameraResult

cameraResult = TriggerStatusResult(myCamera)
```
### Command to read current operating status of the camera 
```go
var operating bool

operating = OperatingStatus(myCamera)
//if camera in "run" mode operating == true
//else operating == false or ("program" mode / not operating)

```


## List of commands are still a WIP...
