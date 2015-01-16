package main

import (
	"fmt"
	"math"
	"syscall"
)

type TimeWindow struct {
	timeTotal int64
	timeMin   float64
	timeMax   float64
	timeAvg   int64
	count     int64
	ec        map[uint16]int
}

type FullResult struct {
	status uint16
	str    string
}

type ResultSet struct {
	Options   *Options
	FullStats map[string]interface{}
	TimeStats []TimeWindow
	Stats     map[uint16]int
	remaining uint64
	opStart   int64

	curTFrame  int64
	curWinTime int64
	winBegin   int64
}

func (fr *FullResult) SetStatus(value string, status uint16) {
	if value == "" {
		fr.status = 1
	}

}

func (rs *ResultSet) setResCode(rc uint16, key string, value string, expectedValue string) {

	rs.remaining--
	rs.Stats[rc]++

	if rs.Options.Full == true {
		rs.FullStats[key] = value
	}

	if rs.Options.TimeRes == 0 {
		return
	}

	var tv syscall.Timeval
	if err := syscall.Gettimeofday(&tv); err != nil {
		fmt.Printf("Error on gettimeofday %v \n", err)
	}

	curTimeInMSecs := (int64(tv.Sec) * 1000) + (int64(tv.Usec) / 1000)
	opsDuration := curTimeInMSecs - rs.opStart
	rs.curTFrame = tv.Sec - (tv.Sec % rs.Options.TimeRes)

	var win TimeWindow

	if rs.curWinTime == 0 {
		rs.curWinTime = rs.curTFrame
		rs.winBegin = rs.curTFrame
		_ = append(rs.TimeStats, win)
	} else if rs.curWinTime < rs.curTFrame {

	}

	lastWin := &rs.TimeStats[len(rs.TimeStats)-1]
	lastWin.count++
	lastWin.timeTotal += opsDuration
	lastWin.timeMin = math.Min(float64(lastWin.timeMin), float64(opsDuration))
	lastWin.timeMax = math.Max(float64(lastWin.timeMax), float64(opsDuration))
	lastWin.ec[rc]++
}

func (rs *ResultSet) ResultsJson(res *ResultResponse) {
	for rc, count := range rs.Stats {
		res.Summary[string(rc)] = count
	}

	if rs.Options.TimeRes == 0 {
		return
	}

	res.Timings.Base = rs.winBegin
	res.Timings.Step = rs.Options.TimeRes

	for _, winstat := range rs.TimeStats {
		win := Window{}
		win.Count = winstat.count
		win.Min = int64(winstat.timeMin)
		win.Max = int64(winstat.timeMax)
		if winstat.count == 0 {
			win.Avg = winstat.timeTotal / winstat.count
		} else {
			win.Avg = 0
		}
		for rc, count := range winstat.ec {
			win.Errors[string(rc)] = count
		}
		_ = append(res.Timings.Windows, win)
	}

}
