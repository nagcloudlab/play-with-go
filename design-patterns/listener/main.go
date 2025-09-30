package main

import "fmt"

/*

	- door
		- open()
		- close()


rules for 'listener/observer	patterns

-> subject & observer must be loosely coupled
-> subject should not know about observer
-> observer should not know about subject
-> observer should be able to register & unregister itself
-> subject should be able to notify all observers when there is a change in its state
-> observer should be able to update itself when there is a change in subject's state

*/

//-------------------------------------------------------------
// DoorListener Interface
// -------------------------------------------------------------

type DoorEvent struct {
	Number int
	Floor  int
}

type DoorListener interface {
	On(event DoorEvent)
	Off(event DoorEvent)
}

// -------------------------------------------------------------
// Light
// -------------------------------------------------------------

type Light struct {
}

func (l *Light) On(event DoorEvent) {
	number := event.Number
	floor := event.Floor
	fmt.Printf("Light On - Number: %d, Floor: %d\n", number, floor)
}

func (l *Light) Off(event DoorEvent) {
	number := event.Number
	floor := event.Floor
	fmt.Printf("Light Off - Number: %d, Floor: %d\n", number, floor)
}

// -------------------------------------------------------------
// Fan
// -------------------------------------------------------------

type Fan struct {
}

func (f *Fan) On(event DoorEvent) {
	number := event.Number
	floor := event.Floor
	fmt.Printf("Fan Started - Number: %d, Floor: %d\n", number, floor)
}

func (f *Fan) Off(event DoorEvent) {
	number := event.Number
	floor := event.Floor
	fmt.Printf("Fan Stopped - Number: %d, Floor: %d\n", number, floor)
}

//-------------------------------------------------------------
// Monitor
// -------------------------------------------------------------

type Monitor struct {
}

func (m *Monitor) On(event DoorEvent) {
	number := event.Number
	floor := event.Floor
	fmt.Printf("Monitor On - Number: %d, Floor: %d\n", number, floor)
}

func (m *Monitor) Off(event DoorEvent) {
	number := event.Number
	floor := event.Floor
	fmt.Printf("Monitor Off - Number: %d, Floor: %d\n", number, floor)
}

// -------------------------------------------------------------
// Door
// -------------------------------------------------------------
type Door struct {
	listeners []DoorListener
}

func (d *Door) RegisterListener(listener DoorListener) {
	d.listeners = append(d.listeners, listener)
}

func (d *Door) UnregisterListener(listener DoorListener) {
	for i, l := range d.listeners {
		if l == listener {
			d.listeners = append(d.listeners[:i], d.listeners[i+1:]...)
			break
		}
	}
}

func (d *Door) Open() {
	fmt.Println("Door Opened")
	event := DoorEvent{Number: 101, Floor: 1}
	for _, listener := range d.listeners {
		listener.On(event)
	}
}

func (d *Door) Close() {
	fmt.Println("Door Closed")
	event := DoorEvent{Number: 101, Floor: 1}
	for _, listener := range d.listeners {
		listener.Off(event)
	}
}

// -------------------------------------------------------------
// Main
// -------------------------------------------------------------

func main() {
	d := &Door{}

	light := &Light{}
	fan := &Fan{}

	d.RegisterListener(light)
	d.RegisterListener(fan)

	// simulate door open and close
	d.Open()
	fmt.Println("......")
	d.Close()

	fmt.Println("......")
	monitor := &Monitor{}
	d.RegisterListener(monitor)
	d.UnregisterListener(fan)

	d.Open()
	fmt.Println("......")
	d.Close()
}
