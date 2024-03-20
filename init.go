package einit

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func reapZombies() error {
	//static void sigreap(int signo) {
	//	while (waitpid(-1, NULL, WNOHANG) > 0)
	//	;
	//}
	wstatus := syscall.WaitStatus(0)
	for {
		childPid, err := syscall.Wait4(-1, &wstatus, syscall.WNOHANG, nil)
		if err != nil {
			return err
		}
		if childPid > 0 {
			fmt.Printf("pid is %d \n", childPid)
		} else if childPid == 0 {
			return nil
		}
	}
}

func waitSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGCHLD)
	go func() {
		for {
			sig := <-sigChan
			if sig == syscall.SIGCHLD {
				err := reapZombies()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
}

func ReapZombiePid() {
	go func() {
		waitSignal()
	}()
}
