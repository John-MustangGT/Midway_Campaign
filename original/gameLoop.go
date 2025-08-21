package main

import (
   "fmt"
   "strings"
   "strconv"
   "time"
)


func gameLoop() {
	for !gameOver {
		clearScreen()
		displayMap()
		displayStatus() 
		displayContacts()
		fmt.Printf("            CAP - ON DECK - -- BELOW --\n")
		fmt.Printf("            F4F SBD TBD F4F SBD TBD\n")
		displayCarriers()
		
		if gameOver {
			break
		}
		
		// Process player command
		fmt.Print("COMMAND: ")
		if !scanner.Scan() {
			break
		}
		command := strings.TrimSpace(strings.ToUpper(scanner.Text()))
		
		if command == "" {
			advanceTime(0)
			continue
		}
		
		if len(command) == 1 && command >= "0" && command <= "9" {
			hours, _ := strconv.Atoi(command)
			advanceTime(hours * 60)
			continue
		}
		
		// Handle commands (existing code)
		switch {
		case strings.HasPrefix(command, "T"):
			changeCourse()
		case strings.HasPrefix(command, "A"):
			armStrike()  
		case strings.HasPrefix(command, "L"):
			launchStrike()
		case strings.HasPrefix(command, "CA"):
			setCap()
		case strings.HasPrefix(command, "CL"):
			clearDeck()
		default:
			fmt.Println("COMMANDS ARE:")
			fmt.Println("T-CHANGE TF COURSE  CA-SET CAP")
			fmt.Println("A-ARM STRIKE        CL-CLEAR DECK") 
			fmt.Println("L-LAUNCH STRIKE      #-WAIT # HOURS")
			time.Sleep(3 * time.Second)
			continue
		}
		
		// CRITICAL: Process AI turn after player command (BASIC lines 880-1500)
		processAITurn()
	}
	
	endGame()
}
