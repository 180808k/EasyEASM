package amass

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

func RunAmass(seedDomain string, results chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Running Amass on %s\n", seedDomain)

	// First command to run Amass
	cmd := exec.Command("amass", "enum", "--passive", "-nocolor", "-d", seedDomain)
	var amassOut bytes.Buffer
	cmd.Stdout = &amassOut
	cmd.Stderr = &amassOut // Capture stderr as well

	// Run the command and capture any error
	err := cmd.Run()
	if err != nil {
		// Log the error and the output from Amass
		log.Printf("Error running Amass on %s: %v\nOutput: %s\n", seedDomain, err, amassOut.String())
		return // Exit the function without crashing the application
	}

	// Now run the oam_subs command
	cmd = exec.Command("oam_subs", "-names", "-d", seedDomain)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		// Log the error and the output from oam_subs
		log.Printf("Error running oam_subs on %s: %v\nOutput: %s\n", seedDomain, err, out.String())
		return // Exit the function without crashing the application
	}

	// Process the output
	for _, domain := range strings.Split(out.String(), "\n") {
		if strings.Contains(domain, seedDomain) && len(domain) != 0 {
			results <- domain
		}
	}
	fmt.Printf("Amass Run completed for %s\n", seedDomain)
}
