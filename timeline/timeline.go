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
	Time() Time                        //defimpl:"read time"
	Dancer() dancer.Dancer             //defimpl:"read dancer"
	Position() geometry.Position       //defimpl:"read position"
	Direction() geometry.Direction     //defimpl:"read direction"
}

func (ds *DancerSnapshotImpl) DancerSnapshot() {}


// Timeline is used to record the positions and directions of a group
// of Dancers over time.
type Timeline interface {
	Timeline()
	Dancers() dancer.Dancers                    //defimpl:"read dancers"
	// The most recent time recorded in the Timeline. 
	MostRecent() Time                          //defimpl:"read mostRecent"
	DoSnapshots(func(DancerSnapshot) bool)     //defimpl:"iterate snapshots"
    	FindSnapshot(dancer.Dancer, Time) DancerSnapshot
	FindSnapshots(dancer dancer.Dancer, start, end Time) []DancerSnapshot
	MakeSnapshot(Time)
	// Bounds returns the most extreme coordinates of all dancers
	// throughout the Timeline.
	Bounds() (leftmost, rightmost geometry.Left, downmost, upmost geometry.Down)
}

func (tl *TimelineImpl) Timeline() {}

// FindSnapshot looks through the recorded DancerSnapshots for this
// Timeline for one that matches the specified Dancer and Time.
func (tl *TimelineImpl) FindSnapshot(d dancer.Dancer, time Time) DancerSnapshot {
	for _, s := range tl.snapshots {
		if d == s.Dancer() && time == s.Time() {
			return s
		}
	}
	return nil
}

func (tl *TimelineImpl) FindSnapshots(d dancer.Dancer, start, end Time) []DancerSnapshot {
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
func (tl *TimelineImpl) MakeSnapshot(time Time) {
	for _, d := range tl.Dancers() {
		s := &DancerSnapshotImpl {
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
	return &TimelineImpl{
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


func (timeline *TimelineImpl) Bounds() (leftmost, rightmost geometry.Left, downmost, upmost geometry.Down) {
	p := timeline.snapshots[0].Position()
	leftmost = p.Left
	rightmost = p.Left
	downmost = p.Down
	upmost = p.Down
	timeline.DoSnapshots(func (ds DancerSnapshot) bool {
		p = ds.Position()
		if p.Left > leftmost {
			leftmost = p.Left
		}
		if p.Left < rightmost {
			rightmost = p.Left
		}
		if p.Down > downmost {
			downmost = p.Down
		}
		if p.Down > upmost {
			upmost = p.Down
		}
		return true
	})
	return leftmost, rightmost, downmost, upmost
}

