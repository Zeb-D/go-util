package common

import (
	"time"
)

//length is 13,time
func TimeToMillisecond(t time.Time) int64 {
	return TimeToDuration(t, time.Millisecond)
}

func TimeToDuration(t time.Time, duration time.Duration) int64 {
	if ISSupportDuration(duration) {
		return t.UnixNano() / int64(duration)
	} else {
		//not support other args
		panic("duration args not time.Duration const")
	}
}

func MillisecondToTime(ms int64) time.Time {
	return DurationToTime(ms, time.Millisecond)
}

func DurationToTime(ms int64, duration time.Duration) time.Time {
	if ISSupportDuration(duration) {
		return time.Unix(0, ms*int64(duration))
	} else {
		//not support other args
		panic("duration args not time.Duration const")
	}
}

func ISSupportDuration(duration time.Duration) bool {
	switch duration {
	case time.Nanosecond, time.Microsecond,
		time.Millisecond, time.Second, time.Minute, time.Hour:
		return true
	default:
		//not support other args
		return false
	}
}
