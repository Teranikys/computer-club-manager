package club

import (
	"bufio"
	"math"
	"os"
	"sort"
	"time"
)

type ComputerClub struct {
	numTables  int
	openTime   time.Time
	closeTime  time.Time
	hourlyRate int
	events     []*Event
	tables     []*Table
	clients    map[string]*Client
	queue      chan string
}

type Event struct {
	eventTime time.Time
	eventID   int
	client    string
	tableID   int
}

type Table struct {
	isOccupied       bool
	clientName       string
	occupiedTime     time.Time
	occupiedDuration time.Duration
}

func (t *Table) revenue(hourlyRate int) int {
	return hourlyRate * int(math.Ceil(t.occupiedDuration.Hours()))
}

type Client struct {
	isInClub     bool
	currentTable int
	arrivalTime  time.Time
}

func newTables(numTables int) []*Table {
	tables := make([]*Table, 0, numTables)
	for i := 0; i < numTables; i++ {
		tables = append(tables, &Table{})
	}
	return tables
}

func NewComputerClub(file *os.File) (*ComputerClub, error) {
	cc := &ComputerClub{}

	scanner := bufio.NewScanner(file)

	err := cc.parseEvents(scanner)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

// ProcessEvents обрабатывает все события
func (cc *ComputerClub) ProcessEvents() {
	for _, e := range cc.events {
		switch e.eventID {
		case 1:
			cc.handleArrival(e)
		case 2:
			cc.handleSeating(e)
		case 3:
			cc.handleWaiting(e)
		case 4:
			cc.handleDeparture(e)
		}
	}
	cc.departDelayedClients()

	cc.logStats()
}

func (cc *ComputerClub) departDelayedClients() {
	departureClients := make([]string, 0)
	for clientName, client := range cc.clients {
		if client.isInClub {
			table := cc.tables[client.currentTable-1]
			table.isOccupied = false
			table.occupiedDuration += cc.closeTime.Sub(table.occupiedTime)
			departureClients = append(departureClients, clientName)
		}
	}
	sort.Slice(departureClients, func(i, j int) bool {
		return departureClients[i] < departureClients[j]
	})
	for _, client := range departureClients {
		cc.logDeparture(cc.closeTime, client)
	}
}
