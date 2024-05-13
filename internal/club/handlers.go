package club

import "fmt"

// handleArrival обрабатывает приход клиента в клуб
func (cc *ComputerClub) handleArrival(e *Event) {
	fmt.Printf("%02d:%02d %d %s\n", e.eventTime.Hour(), e.eventTime.Minute(), e.eventID, e.client)

	if e.eventTime.Before(cc.openTime) || e.eventTime.After(cc.closeTime) {
		cc.logError(e, ErrNotOpenYet)
		return
	}

	if _, ok := cc.clients[e.client]; ok {
		cc.logError(e, ErrYouShallNotPass)
		return
	}
	cc.clients[e.client] = &Client{isInClub: true, arrivalTime: e.eventTime}
}

// handleSeating обрабатывает событие, когда клиент садится за стол
func (cc *ComputerClub) handleSeating(e *Event) {
	fmt.Printf("%02d:%02d %d %s %d\n", e.eventTime.Hour(), e.eventTime.Minute(), e.eventID, e.client, e.tableID)

	client, inClub := cc.clients[e.client]
	if !inClub {
		cc.logError(e, ErrClientUnknown)
		return
	}
	if client.currentTable != 0 {
		if client.currentTable == e.tableID {
			cc.logError(e, ErrPlaceIsBusy)
			return
		}
		cc.tables[client.currentTable-1].isOccupied = false
	}
	table := cc.tables[e.tableID-1]
	if table.isOccupied {
		cc.logError(e, ErrPlaceIsBusy)
		return
	}
	table.isOccupied = true
	table.clientName = e.client
	table.occupiedTime = e.eventTime
	client.currentTable = e.tableID
}

// handleWaiting обрабатывает событие ожидания клиента
func (cc *ComputerClub) handleWaiting(e *Event) {
	fmt.Printf("%02d:%02d %d %s\n", e.eventTime.Hour(), e.eventTime.Minute(), e.eventID, e.client)

	_, inClub := cc.clients[e.client]
	if !inClub {
		cc.logError(e, ErrClientUnknown)
		return
	}
	for _, table := range cc.tables {
		if !table.isOccupied {
			cc.logError(e, ErrICanWaitNoLonger)
			return
		}
	}
	if len(cc.queue) == cap(cc.queue) {
		cc.logDeparture(e.eventTime, e.client)
	}
	cc.queue <- e.client
}

// handleDeparture обрабатывает уход клиента
func (cc *ComputerClub) handleDeparture(e *Event) {
	defer delete(cc.clients, e.client)

	fmt.Printf("%02d:%02d %d %s\n", e.eventTime.Hour(), e.eventTime.Minute(), e.eventID, e.client)

	client, ok := cc.clients[e.client]
	if !ok {
		cc.logError(e, ErrClientUnknown)
		return
	}
	table := cc.tables[client.currentTable-1]
	table.occupiedDuration += e.eventTime.Sub(table.occupiedTime)
	if len(cc.queue) == 0 {
		table.isOccupied = false
		return
	}
	waitingClient := <-cc.queue
	table.clientName = waitingClient
	table.occupiedTime = e.eventTime
	seatedClient, _ := cc.clients[waitingClient]
	seatedClient.currentTable = client.currentTable
	cc.logSeating(e.eventTime, waitingClient, client.currentTable)
}
