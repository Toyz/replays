package replays

type ReplayType int32
const (
	Highlight       ReplayType = 0
	PlayOfTheGame              = 2
	ManualHighlight            = 8
)

func (rt ReplayType) String() string {
	switch rt {
	case Highlight:
		return "Highlight"
	case PlayOfTheGame:
		return "Play of the game"
	case ManualHighlight:
		return "Manual Highlight"
	default:
		return "Unknown Type"
	}
}

type ReplayQuality int32
const (
	Low    ReplayQuality = 30
	Medium               = 50
	High                 = 80
	Ultra                = 100
)

func (rq ReplayQuality) String() string {
	switch rq {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Ultra:
		return "Ultra"
	default:
		return "Unknown Type"
	}
}
