// +build windows

package main

import (
	"dogfooter-control/control/cmd/service"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"log"
	"os"
	"time"
)

var windowSvcName string
func main() {
	windowSvcName = "DogfooterControlService"
	//eLog, err := eventlog.Open(windowSvcName)
	//if err != nil {
	//	return
	//}
	//defer eLog.Close()

	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Panicf("failed to determine if we are running in an interactive session: %v", err)
	}

	run := debug.Run
	if !isIntSess { // windows service
		run = svc.Run
	} else {
		log.Println("Service run in interactive session")
		windowSvcName = ""
	}

	err = run(windowSvcName, &WindowService{})
	if err != nil {
		return
	}
}

type WindowService struct{}

func (m *WindowService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	go service.Run(windowSvcName)
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				fmt.Fprintf(os.Stderr, "svc.Interrogate")
				changes <- c.CurrentStatus
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				fmt.Fprintf(os.Stderr, "svc.Stop or svc.Shutdown")
				break loop
			//case svc.Pause:
			//	fmt.Fprintf(os.Stderr, "svc.Pause")
			//	changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			//case svc.Continue:
			//	fmt.Fprintf(os.Stderr, "svc.Continue")
			//	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}
