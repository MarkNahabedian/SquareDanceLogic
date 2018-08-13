// Package timeline records Dancer position and direction over time.
package timeline

import "fmt"
import "squaredance/geometry"
import "squaredance/dancer"

// Time represents a time for the purpose of a Timeline.  Applications
// might just increment it for successive snapshots or might use it to
// correspond to beats in the music that the dancers are dancing to.
type Time int

// DancerSnapshot encapsulates the Position and Direction of a
// specific Dancer at a specific Time.
type DancerSnapshot interface {
	DancerSnapshot()
	Time() Time
	Dancer() dancer.Dancer
	Position() geometry.Position
	Direction() geometry.Direction
}

type implDancerSnapshot struct {
	time Time
	dancer dancer.Dancer
	position geometry.Position
	direction geometry.Direction
}

func (ds *implDancerSnapshot) DancerSnapshot() {}
func (ds *implDancerSnapshot) Time() Time { return ds.time }
func (ds *implDancerSnapshot) Dancer() dancer.Dancer { return ds.dancer }
func (ds *implDancerSnapshot) Position() geometry.Position { return ds.position }
func (ds *implDancerSnapshot) Direction() geometry.Direction { return ds.direction }


// Timeline is used to record the positions and directions of a group
// of Dancers over time.
type Timeline interface {
	Timeline()
	Dancers() dancer.Dancers
	MostRecent() Time    // The most recent time recorded in the Timeline. 
	FindSnapshot(dancer.Dancer, Time) DancerSnapshot
	FindSnapshots(dancer dancer.Dancer, start, end Time) []DancerSnapshot
	MakeSnapshot(Time)
}

type implTimeline struct {
	dancers dancer.Dancers
	mostRecent Time
	snapshots []DancerSnapshot
}

func (tl *implTimeline) Timeline() {}
func (tl *implTimeline) Dancers() dancer.Dancers { return tl.dancers }
func (tl *implTimeline) MostRecent() Time { return tl.mostRecent }

// FindSnapshot looks through the recorded DancerSnapshots for this
// Timeline for one that matches the specified Dancer and Time.
func (tl *implTimeline) FindSnapshot(d dancer.Dancer, time Time) DancerSnapshot {
	for _, s := range tl.snapshots {
		if d == s.Dancer() && time == s.Time() {
			return s
		}
	}
	return nil
}

func (tl *implTimeline) FindSnapshots(d dancer.Dancer, start, end Time) []DancerSnapshot {
	found := []DancerSnapshot{}
	for _, s := range tl.snapshots {
		if d == s.Dancer() && start <= s.Time() && s.Time() < end {
			found = append(found, s)
		}
	}
	return found
}

// MakeSnapshot records DancerSnapshots for the tracked dancers
// labeled with the specified Time.
func (tl *implTimeline) MakeSnapshot(time Time) {
	for _, d := range tl.Dancers() {
		s := &implDancerSnapshot {
			time: time,
			dancer: d,
			position: d.Position(),
			direction: d.Direction(),
		}
		tl.snapshots = append(tl.snapshots, s)
	}
	if time > tl.mostRecent {
		tl.mostRecent = time
	}
}

// NewTimeline returns a new, empty Timeline for tracking the
// specified dancers.
func NewTimeline(dancers dancer.Dancers) Timeline {
	return &implTimeline{
		dancers: dancers,
		mostRecent: Time(0),
		snapshots: []DancerSnapshot{},
	}
}


// ShowHistory writes the position and directiion of each Dancer
// over time to standard output.
func ShowHistory(tl Timeline) {
	for _, d := range tl.Dancers() {
		fmt.Printf("\nDancer %s\n", d)
		for _, s := range tl.FindSnapshots(d, -1, tl.MostRecent() + 1) {
			fmt.Printf("    %3d  %s  %s\n", s.Time(), s.Position(), s.Direction())
		}
	}
}

