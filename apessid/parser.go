package apessid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yankiwi/wlc_exporter/rpc"

	log "github.com/sirupsen/logrus"
)

// Parse parses cli output and tries to find interfaces with related stats
func (c *apessidCollector) Parse(ostype string, output string) (map[string]Apess, error) {
	log.Debugf("OS: %s\n", ostype)
	switch ostype {
	case rpc.ArubaController:
		return c.ParseArubaController(output)
	// case rpc.ArubaInstant:
	// 	return c.ParseArubaInstant(output)
	default:
		return nil, errors.New("'show ap essid' is not implemented for " + ostype)
	}
}

func findWordDistancesSequential(header string) []int {
	words := strings.Fields(header) // Split the header string into words
	var indices []int               // Slice to store the starting index of each word

	lastIndex := 0 // Initialize last index to start search from beginning
	for _, word := range words {
		// Find the index of the next occurrence by starting the search after the last found index
		nextIndex := strings.Index(header[lastIndex:], word) + lastIndex
		if nextIndex < lastIndex {
			break // If the word is not found, break the loop
		}
		indices = append(indices, nextIndex)
		// Update lastIndex to the end of the current word to ensure sequential search
		lastIndex = nextIndex + len(word)
	}

	// Calculate the distances between the starting indices of consecutive words
	var distances []int
	for i := 0; i < len(indices)-1; i++ {
		distances = append(distances, indices[i+1]-indices[i])
	}
	return distances
}

func parseData(lines []string, distances []int) []map[string]string {

	headers := strings.Fields("ESSID APs MBSSID_Tx_BSS Clients Vlans Encryption") // Define headers based on your headerString2

	distances = append(distances, 20)
	var results []map[string]string // Slice of maps to hold the parsed data

	for _, line := range lines {
		lineMap := make(map[string]string) // Map to store the data for each line
		start := 0                         // Starting index for slicing

		for i, dist := range distances { // Iterate over the distances
			// Check if it's the last header and handle it differently
			if i == len(headers)-1 {
				lineMap[headers[i]] = strings.TrimSpace(line[start:])
				break
			}

			end := start + dist // Calculate end index for slicing
			if end > len(line) {
				// If the end index is out of range, break the loop to avoid a panic
				break
			}
			lineMap[headers[i]] = strings.TrimSpace(line[start:end]) // Slice the line and trim spaces
			start = end                                              // Update start index for next iteration
		}
		results = append(results, lineMap) // Append the map to the results slice
	}

	return results
}

// Code to remove data between paranthesis
// func removeBetweenParentheses(s string) string {
// 	start := strings.Index(s, "(")
// 	end := strings.Index(s, ")")

// 	// If there is no "(" or ")" in the string, or they are in the wrong order, return the original string
// 	if start == -1 || end == -1 || start > end {
// 		return s
// 	}

// 	// Remove the substring from start to end
// 	return s[:start] + s[end+1:]
// }

// Parse parses ArubaInstant cli output and tries to find interfaces with related stats
func (c *apessidCollector) ParseArubaController(output string) (map[string]Apess, error) {
	aps := make(map[string]Apess)

	headers := ""
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if strings.Contains(line, "Name") && strings.Contains(line, "IP Address") {
			var sb strings.Builder // Use a StringBuilder to efficiently build the new header
			for j := 0; j < len(line); j++ {
				if line[j] == ' ' && !(j+1 < len(line) && line[j+1] == ' ') && !(j-1 >= 0 && line[j-1] == ' ') {
					sb.WriteByte('_') // Replace a single space with an underscore
				} else {
					sb.WriteByte(line[j]) // Keep the original character
				}
			}
			lines[i] = sb.String() // Set the modified header back in the lines
			headers = sb.String()
			break
		}
	}

	// fmt.Println(headers)

	var startIndex, endIndex int
	for i, line := range lines {
		if strings.Contains(line, "----") {
			startIndex = i + 1
		}
		endIndex = i
	}

	apes := Apess{}
	name := ""

	// Check if we have valid start and end indexes
	if startIndex > 0 && endIndex > 0 && endIndex > startIndex {
		// Extract the data lines
		dataLines := lines[startIndex:endIndex]

		lengths := findWordDistancesSequential(headers)

		parsedData := parseData(dataLines, lengths)

		for _, lineData := range parsedData {
			fmt.Println(lineData)
			apes.essid = lineData["ESSID"]
			apes.aps = lineData["APs"]
			apes.mbssidtxbss = lineData["MBSSID_Tx_BSS"]
			apes.clients = lineData["Clients"]
			apes.vlans = lineData["VLAN(s)"]
			apes.encryption = lineData["Encryption"]
			aps[name] = apes
		}
		return aps, nil
	} else {
		fmt.Println("The data format is not as expected.")
	}
	return aps, nil
}
