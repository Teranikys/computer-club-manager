package club

import (
	"fmt"
	"time"
)

// logError выводит ошибку в консоль
func (cc *ComputerClub) logError(e *Event, errMessage string) {
	fmt.Printf("%02d:%02d 13 %s\n", e.eventTime.Hour(), e.eventTime.Minute(), errMessage)
}

// logSeating выводит событие о посадке клиента за стол
func (cc *ComputerClub) logSeating(time time.Time, client string, tableID int) {
	fmt.Printf("%02d:%02d 12 %s %d\n", time.Hour(), time.Minute(), client, tableID)
}

// logDeparture выводит событие об уходе клиента
func (cc *ComputerClub) logDeparture(time time.Time, clientName string) {
	fmt.Printf("%02d:%02d 11 %s\n", time.Hour(), time.Minute(), clientName)
}

func (cc *ComputerClub) logStats() {
	fmt.Printf("%02d:%02d\n", cc.closeTime.Hour(), cc.closeTime.Minute())
	for tableID, table := range cc.tables {
		fmt.Printf("%d %d %02d:%02d\n", tableID+1, table.revenue(cc.hourlyRate), int(table.occupiedDuration.Hours()), int(table.occupiedDuration.Minutes())%60)
	}
}
