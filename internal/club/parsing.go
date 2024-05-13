package club

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseTime(input string) (time.Time, error) {
	t, err := time.Parse(`15:04`, input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// parseEvents загружает события из файла из файла
func (cc *ComputerClub) parseEvents(scanner *bufio.Scanner) error {
	var err error

	// Parse number of tables
	if !scanner.Scan() {
		return fmt.Errorf("failed to read number of tables")
	}
	cc.numTables, err = strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println(scanner.Text())
		return fmt.Errorf("invalid number of tables format")
	}

	cc.tables = newTables(cc.numTables)
	cc.queue = make(chan string, cc.numTables)

	// Parse open and close time
	if !scanner.Scan() {
		return fmt.Errorf("failed to read opening hours")
	}
	times := strings.Split(scanner.Text(), " ")
	if len(times) != 2 {
		fmt.Println(scanner.Text())
		return fmt.Errorf("invalid opening hours format")
	}
	cc.openTime, err = parseTime(times[0])
	if err != nil {
		fmt.Println(scanner.Text())
		return err
	}
	cc.closeTime, err = parseTime(times[1])
	if err != nil {
		fmt.Println(scanner.Text())
		return err
	}

	// Parse hourly rate
	if !scanner.Scan() {
		return fmt.Errorf("failed to read hourly rate")
	}
	cc.hourlyRate, err = strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println(scanner.Text())
		return fmt.Errorf("invalid hourly rate format")
	}

	cc.clients = make(map[string]*Client)
	// Load events
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) < 3 {
			fmt.Println(scanner.Text())
			return fmt.Errorf("invalid event format in line: %s", line)
		}

		eventTime, err := parseTime(parts[0])
		if err != nil {
			fmt.Println(scanner.Text())
			return fmt.Errorf("invalid time format in line: %s", line)
		}

		eventID, err := strconv.Atoi(parts[1])
		if err != nil || (eventID < 1 || eventID > 4) {
			fmt.Println(scanner.Text())
			return fmt.Errorf("invalid event ID in line: %s", line)
		}

		clientName := parts[2]
		var tableID int
		if eventID == 2 {
			if len(parts) != 4 {
				fmt.Println(scanner.Text())
				return fmt.Errorf("invalid event format in line: %s", line)
			}
			tableID, err = strconv.Atoi(parts[3])
			if err != nil || tableID < 1 || tableID > cc.numTables {
				fmt.Println(scanner.Text())
				return fmt.Errorf("invalid table ID in line: %s", line)
			}
		}

		cc.events = append(cc.events, &Event{
			eventTime: eventTime,
			eventID:   eventID,
			client:    clientName,
			tableID:   tableID,
		})
	}

	return nil
}
