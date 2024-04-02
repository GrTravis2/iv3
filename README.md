Package of commands that operate Keyence IV3 cameras using TCP/IP procedures written in GO!

## Initialize a camera:
```go
    myCamera := Camera {
        location: "myDesk",
        description: "takes pictures of myDesk",
        ipAddress: "11.222.33.4",
        port: "8500", //iv3 defaults to port 8500
        delimiter: "\r\n", //iv3 defaults to carriage return
    }
```
## Use the camera object and additional arguments to "do stuff"

### Changing programs
```go
    response := ProgramChange(programNumber, myCamera) //reads response and returns success/unsuccessful
```

Types: response string, programNumber int, myCamera camera

### Read current program
```go
    response := ReadProgramNumber(myCamera)
```

Types: response string, myCamera camera

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
cameraResult := TriggerStatusResult(myCamera)
```
Types: cameraResult CameraResult, myCamera Camera

## More commands incoming...
