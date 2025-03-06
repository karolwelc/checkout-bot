package securedtouch

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var (
	mouseEvents    = []string{"click", "dblclick", "mousedown", "mousemove", "mouseout", "mouseover", "mouseup", "wheel"}
	indirectEvents = []string{"change", "fullscreenchange", "invalid", "languagechange", "orientationchange", "seeked", "seeking", "selectstart", "selectionchange", "submit", "volumechange", "reset"}
)

type InputData struct {
	AppSessId    string
	LocationHref string
	DeviceId     string
	CheckSum     string
	StToken      string
}

func GenSecuredTouchData(inputData InputData) []byte {
	rand.Seed(time.Now().Unix())
	var motionData motionData

	var appSessId = inputData.AppSessId
	var locationHref = inputData.LocationHref
	var deviceId = inputData.DeviceId
	var checksum = inputData.CheckSum
	var StToken = inputData.StToken

	windowId := uuid.New().String()
	motionData.UsernameTs = time.Now().Unix() - 1

	motionData.ApplicationID = "asos"
	motionData.DeviceID = deviceId
	motionData.DeviceType = "Chrome(106.0.0.0)-Windows(10)"
	motionData.AppSessionID = appSessId
	motionData.StToken = StToken

	//motionData.KeyboardInteractionPayloads = append(motionData.KeyboardInteractionPayloads)

	for i := 0; i < chooseNumber(2, 5); i++ {
		var mouseInteractionPayload mouseInteractionPayload
		var events []mouseEvent
		var mouseEvent mouseEvent
		for x := 0; x < chooseNumber(3, 6); x++ {
			mouseEvent = createMouseEvent(mouseEvent)
			events = append(events, mouseEvent)
			if mouseEvent.Type == "mouseout" {
				mouseInteractionPayload.AdditionalData.Mouseout = 1
			} else if mouseEvent.Type == "mouseover" {
				mouseInteractionPayload.AdditionalData.Mouseover = 1
			}
		}
		mouseInteractionPayload.Events = events
		mouseInteractionPayload.Identified = false
		mouseInteractionPayload.AdditionalData.WindowID = windowId
		mouseInteractionPayload.AdditionalData.LocationHref = locationHref
		mouseInteractionPayload.AdditionalData.Checksum = checksum
		mouseInteractionPayload.AdditionalData.InnerWidth = 1920
		mouseInteractionPayload.AdditionalData.InnerHeight = 926
		mouseInteractionPayload.AdditionalData.OuterWidth = 1920
		mouseInteractionPayload.AdditionalData.OuterWidth = 1040
		mouseInteractionPayload.AdditionalData.SnapshotsReduceFactor = 0
		mouseInteractionPayload.AdditionalData.EventsWereReduced = true
		motionData.MouseInteractionPayloads = append(motionData.MouseInteractionPayloads, mouseInteractionPayload)

	}

	motionData.IndirectEventsCounters = struct{}{}
	motionData.Gestures = []interface{}{}
	motionData.MetricsData = struct{}{}
	motionData.AccelerometerData = []interface{}{}
	motionData.GyroscopeData = []interface{}{}
	motionData.LinearAccelerometerData = []interface{}{}
	motionData.RotationData = []interface{}{}
	motionData.Index = 2
	motionData.PayloadID = deviceId
	motionData.Tags = []interface{}{}
	motionData.Environment.Ops = 0
	motionData.Environment.WebGl = ""
	motionData.Environment.DevicePixelRatio = 1
	motionData.Environment.ScreenWidth = 1920
	motionData.Environment.ScreenHeight = 1080
	motionData.IsMobile = false
	motionData.Username = deviceId

	b, err := json.Marshal(motionData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
	return Encrypt(string(b))
}

func createMouseEvent(previousMe mouseEvent) mouseEvent {
	var event mouseEvent

	var (
		timepassed int = chooseNumber(985, 3587)
	)

	event.Type = mouseEvents[rand.Intn(len(mouseEvents))]
	if previousMe.EventTs == 0 {
		event.EventTs = chooseNumber(3051, 6712)
	} else {
		event.EventTs = previousMe.EventTs + timepassed
	}
	if previousMe.EpochTs == 0 {
		event.EpochTs = time.Now().UnixMilli()
	} else {
		event.EpochTs = previousMe.EpochTs + int64(timepassed)
	}
	event.Button = 0
	if event.Type == "mousedown" {
		event.Buttons = 1
	} else {
		event.Buttons = 0
	}
	X := chooseNumber(800, 900)
	Y := chooseNumber(350, 450)
	event.ClientX = X
	event.ClientY = Y
	event.MovementX = chooseNumber(-1, 3)
	event.MovementY = chooseNumber(-1, 3)
	event.OffsetX = X
	event.OffsetY = Y
	event.PageX = X
	event.PageY = Y
	event.ScreenX = -1 * chooseNumber(800, 1300)
	event.ScreenY = chooseNumber(430, 650)
	if event.Type == "mousedown" {
		event.Which = 1
	} else {
		event.Which = 0
	}
	event.ModifierKeys = []string{"NumLock"}
	if event.Type == "mousedown" || event.Type == "mouseup" || event.Type == "click" || event.Type == "wheel" || event.Type == "dblclick" {
		Xf := float64(chooseNumber(40, 100)) + truncate(rand.Float64(), 5)
		Yf := float64(chooseNumber(40, 100)) + truncate(rand.Float64(), 5)
		event.TargetBottom = float64(chooseNumber(400, 500)) + truncate(rand.Float64(), 5)
		event.TargetHeight = float64(chooseNumber(40, 100)) + truncate(rand.Float64(), 5)
		event.TargetLeft = Xf
		event.TargetRight = float64(chooseNumber(40, 100)) + truncate(rand.Float64(), 5)
		event.TargetTop = Yf
		event.TargetWidth = float64(chooseNumber(40, 100)) + truncate(rand.Float64(), 5)
		event.TargetX = Xf
		event.TargetY = Yf
	}
	return event
}

type keyboardInteraction struct {
	StID      string `json:"stId"`
	ElementID string `json:"elementId"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Events    []struct {
		Type           string        `json:"type"`
		EventTs        float64       `json:"eventTs"`
		EpochTs        int64         `json:"epochTs"`
		ModifierKeys   []interface{} `json:"modifierKeys"`
		SelectionStart interface{}   `json:"selectionStart"`
		SelectionEnd   interface{}   `json:"selectionEnd"`
		Key            interface{}   `json:"key"`
		KeyCode        interface{}   `json:"keyCode"`
		KeystrokeID    interface{}   `json:"keystrokeId"`
		CurrentLength  int           `json:"currentLength"`
	} `json:"events"`
	Identified     bool `json:"identified"`
	Counter        int  `json:"counter"`
	AdditionalData struct {
		WindowID              string `json:"windowId"`
		LocationHref          string `json:"locationHref"`
		Checksum              string `json:"checksum"`
		InnerWidth            int    `json:"innerWidth"`
		InnerHeight           int    `json:"innerHeight"`
		OuterWidth            int    `json:"outerWidth"`
		OuterHeight           int    `json:"outerHeight"`
		SnapshotsReduceFactor int    `json:"snapshotsReduceFactor"`
		EventsWereReduced     bool   `json:"eventsWereReduced"`
	} `json:"additionalData"`
}

type mouseInteractionPayload struct {
	Events         []mouseEvent `json:"events"`
	Identified     bool         `json:"identified"`
	Counter        int          `json:"counter"`
	AdditionalData struct {
		WindowID              string `json:"windowId"`
		LocationHref          string `json:"locationHref"`
		Checksum              string `json:"checksum"`
		Mouseout              int    `json:"mouseout,omitempty"`
		Mouseover             int    `json:"mouseover,omitempty"`
		InnerWidth            int    `json:"innerWidth"`
		InnerHeight           int    `json:"innerHeight"`
		OuterWidth            int    `json:"outerWidth"`
		OuterHeight           int    `json:"outerHeight"`
		SnapshotsReduceFactor int    `json:"snapshotsReduceFactor"`
		EventsWereReduced     bool   `json:"eventsWereReduced"`
	} `json:"additionalData,omitempty"`
}

type motionData struct {
	ApplicationID               string                    `json:"applicationId"`
	DeviceID                    string                    `json:"deviceId"`
	DeviceType                  string                    `json:"deviceType"`
	AppSessionID                string                    `json:"appSessionId"`
	StToken                     string                    `json:"stToken"`
	KeyboardInteractionPayloads []keyboardInteraction     `json:"keyboardInteractionPayloads,omitempty"`
	MouseInteractionPayloads    []mouseInteractionPayload `json:"mouseInteractionPayloads"`
	IndirectEventsPayload       []struct {
		Category       string `json:"category"`
		Type           string `json:"type"`
		EventTs        int    `json:"eventTs"`
		EpochTs        int64  `json:"epochTs"`
		AdditionalData struct {
			WindowID     string `json:"windowId"`
			LocationHref string `json:"locationHref"`
			Checksum     string `json:"checksum"`
		} `json:"additionalData,omitempty"`
	} `json:"indirectEventsPayload,omitempty"`
	IndirectEventsCounters struct {
	} `json:"indirectEventsCounters,omitempty"`
	Gestures    []interface{} `json:"gestures"`
	MetricsData struct {
	} `json:"metricsData"`
	AccelerometerData       []interface{} `json:"accelerometerData"`
	GyroscopeData           []interface{} `json:"gyroscopeData"`
	LinearAccelerometerData []interface{} `json:"linearAccelerometerData"`
	RotationData            []interface{} `json:"rotationData"`
	Index                   int           `json:"index"`
	PayloadID               string        `json:"payloadId"`
	Tags                    []interface{} `json:"tags"`
	Environment             struct {
		Ops              int    `json:"ops"`
		WebGl            string `json:"webGl"`
		DevicePixelRatio int    `json:"devicePixelRatio"`
		ScreenWidth      int    `json:"screenWidth"`
		ScreenHeight     int    `json:"screenHeight"`
	} `json:"environment"`
	IsMobile   bool   `json:"isMobile"`
	UsernameTs int64  `json:"usernameTs"`
	Username   string `json:"username"`
}

type mouseEvent struct {
	Type         string   `json:"type"`
	EventTs      int      `json:"eventTs"`
	EpochTs      int64    `json:"epochTs"`
	Button       int      `json:"button"`
	Buttons      int      `json:"buttons"`
	ClientX      int      `json:"clientX"`
	ClientY      int      `json:"clientY"`
	MovementX    int      `json:"movementX"`
	MovementY    int      `json:"movementY"`
	OffsetX      int      `json:"offsetX"`
	OffsetY      int      `json:"offsetY"`
	PageX        int      `json:"pageX"`
	PageY        int      `json:"pageY"`
	ScreenX      int      `json:"screenX"`
	ScreenY      int      `json:"screenY"`
	Which        int      `json:"which"`
	ModifierKeys []string `json:"modifierKeys"`
	TargetBottom float64  `json:"targetBottom,omitempty"`
	TargetHeight float64  `json:"targetHeight,omitempty"`
	TargetLeft   float64  `json:"targetLeft,omitempty"`
	TargetRight  float64  `json:"targetRight,omitempty"`
	TargetTop    float64  `json:"targetTop,omitempty"`
	TargetWidth  float64  `json:"targetWidth,omitempty"`
	TargetX      float64  `json:"targetX,omitempty"`
	TargetY      float64  `json:"targetY,omitempty"`
}

func chooseNumber(min int, max int) int {
	return rand.Intn(max-min) + min
}

func truncate(x float64, n int) float64 {
	return math.Trunc(x*math.Pow(10, float64(n))) * math.Pow(10, -float64(n))
}
