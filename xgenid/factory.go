package xgenid

import (
	"github.com/duchiporexia/goutils/xerr"
	"github.com/duchiporexia/goutils/xlog"
	"github.com/duchiporexia/goutils/xmsg"
	"runtime"
	"sync"
	"time"
)

var (
	nanosInMilli = time.Millisecond.Nanoseconds()
	epoch        = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano() / nanosInMilli
)

const (
	invalidServerId           = 0
	defaultServerBits         = 10
	serverIDMask        int   = 1<<defaultServerBits - 1
	defaultSequenceBits       = 12
	sequenceMax         int64 = 1 << defaultSequenceBits
	defaultTimeBits           = defaultServerBits + defaultSequenceBits
)

var (
	serverId = invalidServerId
	sequence int64
	lastTime int64
	mux      = &sync.Mutex{}
)

func setServerId(id int) {
	xlog.InfoF("[******SERVER ID******]: %v", id)
	serverId = id & serverIDMask
}

func resetServerId() {
	xlog.InfoF("resetServerId: %v to %v", serverId, invalidServerId)
	serverId = invalidServerId
}

// IdN requests the next n ids
func IdN(n int) ([]int64, error) {
	mux.Lock()
	defer mux.Unlock()
	if serverId == invalidServerId {
		return nil, xerr.ErrInvalidServerId
	}

	t := waitNext(lastTime)
	ids := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		id := id(t)
		ids = append(ids, id)
	}
	return ids, nil
}

func Id() (int64, error) {
	mux.Lock()
	defer mux.Unlock()
	if serverId == invalidServerId {
		return 0, xerr.ErrInvalidServerId
	}

	t := waitNext(lastTime)
	return id(t), nil
}

func LongID() (xmsg.LongID, error) {
	id, err := Id()
	if err != nil {
		return 0, err
	}
	return xmsg.LongID(id), nil
}

func IdStr() (string, error) {
	id, err := Id()
	if err != nil {
		return "", err
	}
	return xmsg.LongID(id).String(), nil
}

func id(timeNow int64) int64 {
	// less than is just for idN
	if timeNow <= lastTime {
		sequence = sequence + 1
		if sequence == sequenceMax {
			sequence = timeNow & 1
			lastTime++
		}
	} else {
		sequence = timeNow & 1
		lastTime = timeNow
	}
	return (lastTime << defaultTimeBits) | (sequence << defaultServerBits) | int64(serverId)
}

func waitNext(lastTime int64) int64 {
	var timeNow = getNowInMs()
	var diff = lastTime - timeNow
	if diff <= 0 {
		return timeNow
	}
	// Case 1: There is no slot for the current ms.
	// --- So, adjust timeNow to be same as lastTime
	// Case 2: WARNING, clock is moving backwards.
	if diff <= 5 {
		runtime.Gosched()
	} else {
		time.Sleep(time.Millisecond * time.Duration(diff))
	}
	return waitNext(lastTime)
}

func getNowInMs() int64 {
	return time.Now().UnixNano()/nanosInMilli - epoch
}
