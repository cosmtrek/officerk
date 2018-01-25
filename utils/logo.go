package utils

import (
	"fmt"
)

var logo = `
  ______    ______    ______   __                                __    __ 
 /      \  /      \  /      \ /  |                              /  |  /  |
/$$$$$$  |/$$$$$$  |/$$$$$$  |$$/   _______   ______    ______  $$ | /$$/ 
$$ |  $$ |$$ |_ $$/ $$ |_ $$/ /  | /       | /      \  /      \ $$ |/$$/  
$$ |  $$ |$$   |    $$   |    $$ |/$$$$$$$/ /$$$$$$  |/$$$$$$  |$$  $$<   
$$ |  $$ |$$$$/     $$$$/     $$ |$$ |      $$    $$ |$$ |  $$/ $$$$$  \  
$$ \__$$ |$$ |      $$ |      $$ |$$ \_____ $$$$$$$$/ $$ |      $$ |$$  \ 
$$    $$/ $$ |      $$ |      $$ |$$       |$$       |$$ |      $$ | $$  |
 $$$$$$/  $$/       $$/       $$/  $$$$$$$/  $$$$$$$/ $$/       $$/   $$/

              A Distributed DAG Based CronJob System

- build timestamp: %s
- commit: %s
-------------------------------------------------------------------------`

// MasterLogo ...
func MasterLogo(timestamp string, commit string) string {
	l := fmt.Sprintf(logo, timestamp, commit)
	master := `
                   MASTER - I'm Your Father!
=========================================================================`
	return fmt.Sprintf("%s%s\n", l, master)
}

// NodeLogo ...
func NodeLogo(timestamp string, commit string) string {
	l := fmt.Sprintf(logo, timestamp, commit)
	node := `
                   NODE - Where's My Father?
=========================================================================`
	return fmt.Sprintf("%s%s\n", l, node)
}
