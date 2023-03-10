/**
 * Alarm
 *
 * @author Afaan Bilal
 * @link https://afaan.dev/alarm
 */

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type Time struct {
	H int
	M int
	S int
}

type Alarm struct {
	time  Time
	label string
}

func main() {
	var alarms = [5]Alarm{
		{time: Time{6, 40, 0}, label: "Fajr"},
		{time: Time{13, 20, 0}, label: "Zuhr"},
		{time: Time{16, 15, 0}, label: "Asr"},
		{time: Time{17, 40, 0}, label: "Maghrib"},
		{time: Time{19, 20, 0}, label: "Isha"},
	}

	var alarmHandles = [5]chan string{}
	for i := 0; i < 5; i++ {
		alarmHandles[i] = SetAlarm(alarms[i].time, func() {
			fmt.Println("Alarm received: ", alarms[i].label)
			PlayBell()
		})
	}

	// Block.
	for i := 0; i < 5; i++ {
		<-alarmHandles[i]
	}
}

func SetAlarm(alarmTime Time, callback func()) (endRecSignal chan string) {
	endRecSignal = make(chan string)

	go func() {
		timeParts := strings.Split(time.Now().Format("15:04:05"), ":")
		hh, _ := strconv.Atoi(timeParts[0])
		mm, _ := strconv.Atoi(timeParts[1])
		ss, _ := strconv.Atoi(timeParts[2])

		startAlarm := GetDiffSeconds(Time{hh, mm, ss}, alarmTime)

		time.AfterFunc(time.Duration(startAlarm)*time.Second, func() {
			callback()
			endRecSignal <- "done"
			close(endRecSignal)
		})
	}()

	return
}

func GetDiffSeconds(fromTime, toTime Time) int {
	diff := GetSeconds(toTime) - GetSeconds(fromTime)

	if diff < 0 {
		return diff + 24*60*60
	} else {
		return diff
	}
}

func GetSeconds(time Time) int {
	return time.H*60*60 + time.M*60 + time.S
}

func PlayBell() {
	f, err := os.Open("bell.wav")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
	select {}
}
