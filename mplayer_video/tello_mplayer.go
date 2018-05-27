package main

import (
	"fmt"
	"os/exec"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8890")

	work := func() {
		mplayer := exec.Command("mplayer", "-fps", "60", "-xy", "500", "-")
		mplayerIn, _ := mplayer.StdinPipe()
		if err := mplayer.Start(); err != nil {
			fmt.Println(err)
			return
		}

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			drone.StartVideo()
			drone.SetVideoEncoderRate(4)
			gobot.Every(1*time.Second, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := mplayerIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
