// Package dahuaevents contains structs for Dahua events.
package dahuaevents

import "github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"

const (
	CodeCrossLineDetection    = "CrossLineDetection"
	CodeCrossRegionDetection  = "CrossRegionDetection"
	CodeFaceDetection         = "FaceDetection"
	CodeIntelliFrame          = "IntelliFrame"
	CodeInterVideoAccess      = "InterVideoAccess"
	CodeLeFunctionStatusSync  = "LeFunctionStatusSync"
	CodeLeftDetection         = "LeftDetection"
	CodeLoginFailure          = "LoginFailure"
	CodeManNumDetection       = "ManNumDetection"
	CodeNTPAdjustTime         = "NTPAdjustTime"
	CodeNewFile               = "NewFile"
	CodeNumberStat            = "NumberStat"
	CodeQueueNumDetection     = "QueueNumDetection"
	CodeReboot                = "Reboot"
	CodeRtspSessionDisconnect = "RtspSessionDisconnect"
	CodeSceneChange           = "SceneChange"
	CodeSmartMotionHuman      = "SmartMotionHuman"
	CodeStorageChange         = "StorageChange"
	CodeTakenAwayDetection    = "TakenAwayDetection"
	CodeTimeChange            = "TimeChange"
	CodeVideoMotion           = "VideoMotion"
	CodeVideoMotionInfo       = "VideoMotionInfo"
	CodeWanderDetection       = "WanderDetection"
	CodeSystemState           = "SystemState"
)

type Action string

const (
	ActionStart = "Start"
	ActionStop  = "Stop"
	ActionPulse = "Pulse"
	ActionState = "State"
)

type InterVideoAccess struct {
	// Type can be "WebLogin", "WebAllLogout".
	Type string `json:"Type"`
}

type TimeChange struct {
	BeforeModifyTime dahuarpc.Timestamp `json:"BeforeModifyTime"`
	ModifiedTime     dahuarpc.Timestamp `json:"ModifiedTime"`
}

type NTPAdjustTime struct {
	Address string             `json:"Address"`
	Before  dahuarpc.Timestamp `json:"Before"`
	Result  bool               `json:"result"`
}

type VideoMotion struct {
	LocaleTime        dahuarpc.Timestamp `json:"LocaleTime"`
	Utc               int                `json:"UTC"`
	Name              string             `json:"Name"`
	SmartMotionEnable bool               `json:"SmartMotionEnable"`
}

type SmartMotionHuman struct {
	LocaleTime dahuarpc.Timestamp `json:"LocaleTime"`
	Utc        int                `json:"UTC"`
	Name       string             `json:"Name"`
}

type LeFunctionStatusSync struct {
	// Function can be "WightLight".
	Function string `json:"Function"`
	Status   bool   `json:"Status"`
}

// CrossLineDetection is "Tripwire".
type CrossLineDetection = Detection

// CrossRegionDetection is "Intrusion".
type CrossRegionDetection = Detection

// WanderDetection is "Loitering Detection".
type WanderDetection = Detection

// LeftDetection is "Abandoned Object".
type LeftDetection = Detection

// TakenAwayDetection is "Missing Object".
type TakenAwayDetection = Detection

// ManNumDetection is "In Area No.".
type ManNumDetection = Detection

// QueueNumDetection is "Queuing".
type QueueNumDetection = Detection

// NumberStat is "People Counting".
type NumberStat = Detection

type Detection struct {
	AreaID    int `json:"AreaID"`
	CfgRuleID int `json:"CfgRuleId"`
	// Action can be "Cross", "".
	Action string `json:"Action"`
	// Class can be "Normal", "FaceDetection", "NumberStat".
	Class         string    `json:"Class"`
	EnteredNumber int       `json:"EnteredNumber"`
	CountInGroup  int       `json:"CountInGroup"`
	DetectRegion  [4][2]int `json:"DetectRegion"`
	// Direction can be "Enter", "Leave", "LeftToRight", "RightToLeft".
	Direction    string `json:"Direction"`
	EventID      int    `json:"EventID"`
	EventSeq     int    `json:"EventSeq"`
	ExitedNumber int    `json:"ExitedNumber"`
	// Faces is used when Class is "FaceDetection".
	Faces         []DetectionFace `json:"Faces"`
	FrameSequence int             `json:"FrameSequence"`
	GroupID       int             `json:"GroupID"`
	// ManList is used when Class is "NumberStat".
	ManList  []DetectionManList `json:"ManList"`
	Mark     int                `json:"Mark"`
	Name     string             `json:"Name"`
	Number   int                `json:"Number"`
	Object   DetectionObject    `json:"Object"`
	Objects  []DetectionObject  `json:"Objects"`
	Pts      float64            `json:"PTS"`
	PresetID int                `json:"PresetID"`
	Priority int                `json:"Priority"`
	RuleID   int                `json:"RuleID"`
	RuleID0  int                `json:"RuleId"`
	Source   float64            `json:"Source"`
	// Type can be "ExitOver", "".
	Type  string `json:"Type"`
	Utc   int    `json:"UTC"`
	Utcms int    `json:"UTCMS"`
}

type DetectionManList struct {
	BoundingBox [4]int `json:"BoundingBox"`
}

type DetectionFace struct {
	BoundingBox [4]int `json:"BoundingBox"`
	Center      [2]int `json:"Center"`
	ObjectID    int    `json:"ObjectID"`
	ObjectType  string `json:"ObjectType"`
	RelativeID  int    `json:"RelativeID"`
}

type DetectionObject struct {
	// Action can be "Appear".
	Action       string                   `json:"Action"`
	BoundingBox  [4]int                   `json:"BoundingBox"`
	BrandYear    int                      `json:"BrandYear"`
	CarLogoIndex int                      `json:"CarLogoIndex"`
	CarWindow    DetectionObjectCarWindow `json:"CarWindow"`
	// Category can be "Unknown".
	Category       string `json:"Category"`
	Center         [2]int `json:"Center"`
	Confidence     int    `json:"Confidence"`
	FrameSequence  int    `json:"FrameSequence"`
	LowerBodyColor [4]int `json:"LowerBodyColor"`
	MainColor      [4]int `json:"MainColor"`
	ObjectID       int    `json:"ObjectID"`
	// ObjectType can be "Human", "Vehicle", "HumanFace".
	ObjectType        string  `json:"ObjectType"`
	RelativeID        int     `json:"RelativeID"`
	SerialUUID        string  `json:"SerialUUID"`
	Source            float64 `json:"Source"`
	Speed             int     `json:"Speed"`
	SpeedTypeInternal int     `json:"SpeedTypeInternal"`
	SubBrand          int     `json:"SubBrand"`
	// Text can be "Unknown".
	Text       string `json:"Text"`
	CarLenMode int    `json:"carLenMode"`
	CarLength  int    `json:"carLength"`
}

type DetectionObjectCarWindow struct {
	BoundingBox [4]int `json:"BoundingBox"`
}

type IntelliFrame struct {
	// Action can be "Start", "Stop".
	Action string `json:"Action"`
}

type NewFile struct {
	CountInGroup int `json:"CountInGroup"`
	// Event can be "CrossRegionDetection", "FaceDetection".
	Event       string `json:"Event"`
	File        string `json:"File"`
	GroupID     int    `json:"GroupID"`
	Index       int    `json:"Index"`
	MailTimeout int    `json:"MailTimeout"`
	Size        int    `json:"Size"`
	// StoragePoint can be "Temporary", "NULL".
	StoragePoint string `json:"StoragePoint"`
}

type RtspSessionDisconnect struct {
	// Device is an IP address such as "192.168.60.8".
	Device     string `json:"Device"`
	StreamType string `json:"StreamType"`
	UserAgent  string `json:"UserAgent"`
}

type SceneChange struct {
	EventID int     `json:"EventID"`
	Pts     float64 `json:"PTS"`
	Utc     float64 `json:"UTC"`
}

type LoginFailure struct {
	Address string `json:"Address"`
	Name    string `json:"Name"`
	Type    string `json:"Type"`
	Utc     int    `json:"UTC"`
}

type SystemState struct {
	// State can be "Active".
	State string `json:"State"`
}
