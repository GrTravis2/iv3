# To operate a keyence iv3 camera over a network...

## Initialize a camera with NewCamera
```Go
func NewCamera(ipAddress []int) *Camera

type Camera struct {
	ip        []int
	port      int
	delimiter rune
}

// init iv3 camera with default values, change thru setters
myCamera := NewCamera([]int{127, 0, 0, 1})
```

## Pass camera to a messenger object
```Go
func NewMessenger(name string, c *Camera.Camera) *Messenger

// Messenger is a collection of cameras, and handles all network logic.
// All you need to know is the name the camera was assigned!
type Messenger struct {
	Cameras map[string]*Camera.Camera
}

myMessenger := NewMessenger("exampleCamera", myCamera)
```

## Send the camera a command to make it "do stuff"
```Go
func (m *Messenger) Send(name string, msg Message) (Response, error)

// tell the "exampleCamera" to take a picture, evaluate, and send back results
response, err := myMessenger.Send("exampleCamera", Operate.Trig())
```

## See packages for available commands and iv3 manual for usage
- **Operate -> commands for operating camera and reading/writing data**
- **~~Status~~ TODO! -> commands for reading and clearing camera warnings/errors**
- **~~Config~~ TODO! -> commands for editting camera settings**

