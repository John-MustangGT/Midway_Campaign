package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Game state structures
type Fleet struct {
	x, y, speed, course, fleetType, damage, aa float64
}

type Carrier struct {
	fleet, f4f, sbd, tbd, deckF4f, deckSbd, deckTbd, cap, damage float64
}

type Strike struct {
	f4f, val, sbd, kate, tbd, target, arrivalTime, returnTime, escort, bomberType float64
	launched int
}

// Global game variables
var (
	fleets    [6]Fleet
	carriers  [9]Carrier
	strikes   [10]Strike
	weights   [6]float64
	fx, fy, fz [6]int
	c1        [4]int
	gameDay   = 3
	gameTime  = 720.0
	pi        = 0.017453293
	gameOver  = false
	skip      = false
	victory0  = 0.0
	victory1  = 0.0
	c5        = 0
	c6        = 1
	c7        = 0.0
	scanner   = bufio.NewScanner(os.Stdin)
)

func main() {
	showIntro()
	initializeGame()
	gameLoop()
}

func showIntro() {
	fmt.Println("MICROCOMPUTER")
	fmt.Println("GAMES,INC.")
	fmt.Println("A DIVISION OF....")
	time.Sleep(2 * time.Second)
	
	fmt.Println()
	fmt.Println("THE AVALON HILL GAME COMPANY")
	time.Sleep(2 * time.Second)
	
	fmt.Println()
	fmt.Println("COPYRIGHT (C)")
	fmt.Println("AVALON HILL GAME CO. 1982.")
	fmt.Println("ALL RIGHTS RESERVED.")
	fmt.Println()
	fmt.Println("COMPUTER PROGRAM AND")
	fmt.Println("AUDIO VISUAL DISPLAY")
	fmt.Println("COPYRIGHTED.")
	fmt.Println()
	fmt.Println("PRESENTS ....")
	time.Sleep(2 * time.Second)
}

func initializeGame() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
	
	// Initialize fleet data
	fleetData := [][]float64{
		{0, 1, 0, 25, 0.1, 0.02},
		{0, 1, 0, 18, 0.2, 0.01},
		{0, 1, 0, 25, 0.1, 0.01},
		{0, 3, 0, 25, 0.1, 0.06},
		{0, 4, 0, 25, 0.1, 0.04},
		{2, 5, 0, 0, 0.25, 0.04},
	}
	
	positions := [][]float64{
		{270, 90, 525},
		{230, 60, 560},
		{230, 60, 560},
		{25, 20, 380},
		{25, 20, 380},
		{0, 0, 0},
	}
	
	for i := 0; i < 6; i++ {
		fleets[i].fleetType = fleetData[i][0]
		fleets[i].damage = fleetData[i][1]
		fleets[i].speed = fleetData[i][3]
		fleets[i].aa = fleetData[i][4]
		
		// Set positions
		l := positions[i][2] + 175*rand.Float64() - 200*rand.Float64()
		if i < 3 {
			l -= 200 * rand.Float64()
		}
		j := (positions[i][0] + positions[i][1]*rand.Float64()) * pi
		
		fleets[i].x = 850 - l*math.Sin(j)
		fleets[i].y = 450 - l*math.Cos(j)
		if i < 3 {
			fleets[i].x = 850 - l*math.Sin(j)
			fleets[i].y = 450 - l*math.Cos(j)
		}
		
		if i >= 3 {
			if fleets[i].x > 1124 {
				fleets[i].x = 1124
			}
			if fleets[i].y > 1149 {
				fleets[i].y = 1149
			}
		}
		
		j = j + 180*pi
		if j > 180*pi {
			j += 360 * pi
		}
		
		if i < 3 {
			fleets[i].course = j
		} else if i != 5 {
			fleets[i].course = 205 * pi
		} else {
			fleets[i].course = 0
		}
	}
	
	// Initialize carrier data
	carrierData := [][]float64{
		{0, 21, 21, 21},
		{0, 30, 23, 30},
		{0, 21, 21, 21},
		{0, 21, 21, 21},
		{3, 27, 38, 14},
		{3, 27, 35, 15},
		{4, 25, 37, 13},
		{5, 14, 14, 10},
		{1, 15, 0, 15},
	}
	
	for i := 0; i < 9; i++ {
		carriers[i].fleet = carrierData[i][0]
		carriers[i].f4f = carrierData[i][1]
		carriers[i].sbd = carrierData[i][2]
		carriers[i].tbd = carrierData[i][3]
		carriers[i].damage = 0
	}
	
	// Set weights
	weights = [6]float64{1.5, 1.4, 1.3, 1.3, 1.2, 1}
	
	// Initialize strikes
	for i := 0; i < 10; i++ {
		strikes[i].launched = -1
	}
	
	// Set carrier 8's values
	carriers[8].cap = carriers[8].f4f
	carriers[8].f4f = 0
	
   // In initializeGame(), replace the deck initialization section:
   // Initialize deck positions - BASIC line 100 shows deck = hangar initially
   for i := 4; i <= 7; i++ {
   	// Copy hangar to deck initially (BASIC line 100)
   	carriers[i].deckF4f = carriers[i].f4f
   	carriers[i].deckSbd = carriers[i].sbd  
   	carriers[i].deckTbd = carriers[i].tbd
   	// Then clear hangar (BASIC line 100 shows J-3 gets cleared)
   	carriers[i].f4f = 0
   	carriers[i].sbd = 0
   	carriers[i].tbd = 0
   	carriers[i].cap = 0
   }
	
	// Set initial courses for TF-16 and TF-17
	for i := 3; i <= 4; i++ {
		angle := getAngle(i, 4)
		fleets[i].course = angle
	}
}

func getAngle(from, to int) float64 {
	dx := fleets[to].x - fleets[from].x
	dy := fleets[to].y - fleets[from].y
	if dy == 0 {
		if dx < 0 {
			return (90 - 180) * pi
		}
		return 90 * pi
	}
	angle := math.Atan(dx / dy)
	if dy > 0 {
		if angle < 0 {
			angle -= 360 * pi
		}
		return angle
	}
	return angle + 180*pi
}

func changeCourse() {
	fmt.Print("WHICH TASK FORCE: ")
	if !scanner.Scan() {
		return
	}
	tf := strings.TrimSpace(scanner.Text())
	if tf == "" {
		return
	}
	
	var fleetNum int
	if strings.Contains(tf, "16") || tf == "3" {
		fleetNum = 3
	} else if strings.Contains(tf, "17") || tf == "4" {
		fleetNum = 4
	} else {
		return
	}
	
	fmt.Printf("NEW COURSE FOR TF-%d: ", fleetNum+13)
	if !scanner.Scan() {
		return
	}
	courseStr := strings.TrimSpace(scanner.Text())
	if courseStr == "" {
		return
	}
	
	course, err := strconv.ParseFloat(courseStr, 64)
	if err != nil || course < 0 || course > 360 {
		return
	}
	
	fleets[fleetNum].course = course * pi
	displayContacts()
}

func setCap() {
	carrier := selectCarrier()
	if carrier == -1 {
		return
	}
	
	fmt.Printf("F4F's FOR %s CAP: ", getCarrierName(carrier))
	if !scanner.Scan() {
		return
	}
	
	capStr := strings.TrimSpace(scanner.Text())
	cap, err := strconv.Atoi(capStr)
	if err != nil {
		return
	}
	
	// Return current CAP to hangar
	carriers[carrier].f4f += carriers[carrier].cap
	carriers[carrier].cap = 0
	
	if float64(cap) <= carriers[carrier].f4f {
		carriers[carrier].cap = float64(cap)
		carriers[carrier].f4f -= float64(cap)
	} else {
		carriers[carrier].cap = carriers[carrier].f4f
		carriers[carrier].f4f = 0
		remaining := float64(cap) - carriers[carrier].cap
		
		if remaining < carriers[carrier].deckF4f {
			carriers[carrier].cap += remaining
			carriers[carrier].deckF4f -= remaining
		} else {
			carriers[carrier].cap += carriers[carrier].deckF4f
			carriers[carrier].deckF4f = 0
		}
	}
	
	displayCarriers()
}

func selectCarrier() int {
	fmt.Print("WHICH CARRIER: ")
	if !scanner.Scan() {
		return -1
	}
	
	carrierStr := strings.TrimSpace(strings.ToUpper(scanner.Text()))
	if carrierStr == "" {
		return -1
	}
	
	var carrier int = -1
	switch carrierStr {
	case "E", "ENTERPRISE":
		carrier = 4
	case "H", "HORNET":
		carrier = 5
	case "Y", "YORKTOWN":
		carrier = 6
	case "M", "MIDWAY":
		carrier = 7
	}
	
	if carrier == -1 {
		return -1
	}
	
	if carriers[carrier].damage >= 60 {
		fmt.Printf("%s IS NOT OPERATIONAL.\n", getCarrierName(carrier))
		time.Sleep(1 * time.Second)
		return -1
	}
	
	return carrier
}

func advanceTime(minutes int) {
	if minutes == 0 {
		minutes = 30 + int(30*rand.Float64())
	}
	
	gameTime += float64(minutes)
	if gameTime >= 1440 {
		gameDay--
		gameTime -= 1440
		if gameDay <= 0 {
			gameOver = true
		}
	}
	
	// Move fleets
	for i := 0; i < 5; i++ {
		fleets[i].x += float64(minutes) * fleets[i].speed * math.Sin(fleets[i].course) / 60
		fleets[i].y += float64(minutes) * fleets[i].speed * math.Cos(fleets[i].course) / 60
	}
	
	// Process strikes and other game logic
	processStrikes()
	processReconnaissance()
	processAI()
}

func processStrikes() {
	// Process strikes arriving and returning
	for i := 0; i < 10; i++ {
		if strikes[i].launched == -1 {
			continue
		}
		
		if gameTime >= strikes[i].arrivalTime && strikes[i].escort != -1 {
			// Strike arrives at target
			fmt.Printf("AIR STRIKE ATTACKING TARGET!\n")
			// Combat resolution would go here
			strikes[i].escort = -1
		}
		
		if gameTime >= strikes[i].returnTime {
			// Strike returns
			if strikes[i].launched < 4 {
				// Return to US carriers
				carriers[strikes[i].launched].f4f += strikes[i].f4f
				carriers[strikes[i].launched].sbd += strikes[i].sbd
				carriers[strikes[i].launched].tbd += strikes[i].tbd
			}
			strikes[i].launched = -1
		}
	}
}

func processReconnaissance() {
	// Simplified reconnaissance logic
	if gameTime > 240 && gameTime < 1140 {
		for i := 0; i < 3; i++ {
			if rand.Float64() < 0.043 {
				fleets[i].damage = 1 // Spotted
				fmt.Printf("PBY SPOTS JAPANESE %s.\n", getFleetName(i))
				gameOver = false // Force display update
			}
		}
	}
}

func processAI() {
	// Simplified AI logic for Japanese fleet movement
	for i := 0; i < 3; i++ {
		if i == 0 { // Carrier group
			// Head toward Midway if not spotted
			if fleets[i].damage < 2 {
				angle := getAngle(i, 5)
				fleets[i].course = angle
			}
		}
	}
}

// Utility functions
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func displayMap() {
	// Clear previous positions
	for i := 0; i < 6; i++ {
		if fz[i] == 1 {
			// Clear old position
			x := fx[i]
			y := fy[i]
			if x >= 0 && x <= 22 && y >= 0 && y <= 11 {
				// This would clear the old position in a real terminal
			}
		}
		fz[i] = 0
		fx[i] = int(fleets[i].x*0.02 + 0.5)
		fy[i] = int(fleets[i].y*0.01 + 0.5)
	}
	
	// Create map grid
	mapGrid := make([][]string, 12)
	for i := range mapGrid {
		mapGrid[i] = make([]string, 23)
		for j := range mapGrid[i] {
			mapGrid[i][j] = "."
		}
	}
	
	// Place fleets on map
	contactNum := 0
	for i := 0; i < 6; i++ {
		x := fx[i]
		y := fy[i]
		
		if x >= 0 && x <= 22 && y >= 0 && y <= 11 {
			fz[i] = 1
			var symbol string
			
			if i <= 2 {
				// Japanese fleets - only show if spotted
				if fleets[i].damage == 0 {
					continue // Not spotted, don't show
				}
				contactNum++
				symbol = fmt.Sprintf("%d", contactNum)
			} else {
				// US fleets and Midway - always show
				switch i {
				case 3:
					symbol = "6" // TF-16
				case 4:
					symbol = "7" // TF-17
				case 5:
					symbol = "M" // Midway
				}
			}
			
			mapGrid[11-y][x] = symbol // Invert Y for proper display
		}
	}
	
	// Display the map
	fmt.Printf("%d JUNE 1942  %02d:%02d\n", gameDay, int(gameTime)/60, int(gameTime)%60)
	fmt.Println()
	
	for i := 0; i < 12; i++ {
		for j := 0; j < 23; j++ {
			fmt.Print(mapGrid[i][j])
			if j < 22 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func displayStatus() {
	fmt.Printf("                    TF-16                 TF-17\n")
	fmt.Println()
}

func displayContacts() {
	contactNum := 0
	fmt.Println()
	for i := 0; i < 3; i++ {
		if fleets[i].damage > 0 {
			contactNum++
			fmt.Printf("CONTACT  %d ", contactNum)
			if fleets[i].damage >= 2 {
				switch i {
				case 0:
					fmt.Print("CV")
				case 1:
					fmt.Print("TT")
				case 2:
					fmt.Print("CA")
				}
			} else {
				fmt.Print("??")
			}
			
			// Calculate bearing from Midway (fleet 5) to contact
			bearing := getBearingFromMidway(i)
			range_nm := getRangeFromMidway(i)
			
			fmt.Printf("   %03.0f° %4.0f\n", bearing, range_nm)
		}
	}
	
	// Clear remaining contact lines
	for i := contactNum + 1; i <= 3; i++ {
		fmt.Printf("                             \n")
	}
}

func displayTime() {
	hours := int(gameTime) / 60
	minutes := int(gameTime) % 60
	fmt.Printf("%d JUNE 1942  %02d:%02d\n", gameDay, hours, minutes)
}

func getDistance(fleet1, fleet2 int) float64 {
	dx := fleets[fleet1].x - fleets[fleet2].x
	dy := fleets[fleet1].y - fleets[fleet2].y
	return math.Sqrt(dx*dx + dy*dy)
}

func getBearingFromMidway(targetFleet int) float64 {
	// Calculate bearing from Midway (fleet 5) to target fleet
	dx := fleets[targetFleet].x - fleets[5].x
	dy := fleets[targetFleet].y - fleets[5].y
	
	if dy == 0 {
		if dx < 0 {
			return 270.0
		}
		return 90.0
	}
	
	// Calculate angle in radians
	angle := math.Atan(dx / dy)
	
	// Convert to degrees and adjust for proper bearing
	bearing := angle / pi
	
	if dy > 0 {
		// Target is north of Midway
		bearing = bearing + 0.0
	} else {
		// Target is south of Midway  
		bearing = bearing + 180.0
	}
	
	// Ensure bearing is 0-360
	for bearing < 0 {
		bearing += 360.0
	}
	for bearing >= 360 {
		bearing -= 360.0
	}
	
	return bearing
}

func getRangeFromMidway(targetFleet int) float64 {
	// Calculate range in nautical miles from Midway to target
	dx := fleets[targetFleet].x - fleets[5].x
	dy := fleets[targetFleet].y - fleets[5].y
	return math.Sqrt(dx*dx + dy*dy)
}

func getCarrierName(carrier int) string {
	names := []string{
		"AKAGI", "KAGA", "SORYU", "HIRYU",
		"ENTERPRISE", "HORNET", "YORKTOWN", "MIDWAY", "ZUIHO",
	}
	if carrier >= 0 && carrier < len(names) {
		return names[carrier]
	}
	return "UNKNOWN"
}

func getFleetName(fleet int) string {
	names := []string{
		"CARRIER GROUP", "TRANSPORT GROUP", "CRUISER GROUP",
		"TASK FORCE 16", "TASK FORCE 17", "MIDWAY ISLAND",
	}
	if fleet >= 0 && fleet < len(names) {
		return names[fleet]
	}
	return "UNKNOWN"
}

func endGame() {
	fmt.Println("\nTHE BATTLE IS OVER. REPORT:")
	fmt.Println("CARRIER    DAMAGE")
	fmt.Println("__________ ______")
	
	// Calculate victory points and display results
	usVictory := 0.0
	japVictory := 0.0
	
	for i := 0; i <= 3; i++ {
		fmt.Printf("%-10s ", getCarrierName(i))
		if carriers[i].damage >= 100 {
			fmt.Println("SUNK")
			usVictory += 700
		} else if carriers[i].damage >= 60 {
			fmt.Println("HEAVY")
			usVictory += 200
		} else if carriers[i].damage > 0 {
			fmt.Println("LIGHT")
			usVictory += 100
		} else {
			fmt.Println("NONE")
		}
	}
	
	fmt.Println("__________ ______")
	
	for i := 4; i <= 7; i++ {
		fmt.Printf("%-10s ", getCarrierName(i))
		if carriers[i].damage >= 100 {
			if i == 7 {
				fmt.Println("DESTROYED")
			} else {
				fmt.Println("SUNK")
			}
			japVictory += 700
		} else if carriers[i].damage >= 60 {
			fmt.Println("HEAVY")
			japVictory += 200
		} else if carriers[i].damage > 0 {
			fmt.Println("LIGHT")
			japVictory += 100
		} else {
			fmt.Println("NONE")
		}
	}
	
	victory := usVictory - japVictory
	var victor, victoryType string
	
	if victory >= 0 {
		victor = "UNITED STATES"
	} else {
		victor = "JAPANESE"
		victory = -victory
	}
	
	if victory >= 2000 {
		victoryType = "STRATEGIC"
	} else if victory >= 1000 {
		victoryType = "TACTICAL"
	} else {
		victoryType = "MARGINAL"
	}
	
	fmt.Printf("\n%s %s VICTORY\n", victor, victoryType)
	
	fmt.Print("\nPLAY AGAIN (Y/N)? ")
	if scanner.Scan() {
		response := strings.ToUpper(strings.TrimSpace(scanner.Text()))
		if response == "Y" {
			main()
		}
	}
}
