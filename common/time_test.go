package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeToMillisecond(t *testing.T) {
	now := time.Now()
	ms := TimeToMillisecond(now)
	ss := TimeToDuration(time.Now(), time.Second)
	fmt.Printf("ms:%d,ss:%d \n", ms, ss)
	tt := MillisecondToTime(ms)
	tt2 := DurationToTime(ss, time.Second)
	fmt.Printf("time:%s,time2:%s \n", tt, tt2)
	assert.NotEqual(t, tt, now)
}
