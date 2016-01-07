package channel

import (
	"fmt"
	"reflect"
	"runtime"
)

// Monitor holds channels of interest and can return capacity and length at request.
type Monitor struct {
	chans map[string]interface{}
}

func New() *Monitor {
	return &Monitor{
		chans: make(map[string]interface{}),
	}
}

// AddNamed adds a channel to be monitor and associates the channel with this name.
func (this *Monitor) AddNamed(name string, channel interface{}) error {

	//reflect on the input to get the correct channel type.
	if reflect.TypeOf(channel).Kind() != reflect.Chan {
		return fmt.Errorf("invalid input type %v for input param channel, must be of type chan", channel)
	}

	if _, found := this.chans[name]; found {
		return fmt.Errorf("channel with name: %s already being monitored.", name)
	}

	this.chans[name] = channel

	return nil
}

// Add a channel
func (this *Monitor) Add(channel interface{}) (string, error) {

	//name the channel using the callers file and line.
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	name := fmt.Sprintf("%s:%d %s\n", file, line, f.Name())

	return name, this.AddNamed(name, channel)
}

// ChanState struct holding Length and Capacity.
type ChanState struct {
	Len int `json:"length"`
	Cap int `json:"capacity"`
}

// Get returns the channel state for a give channel name.
func (this *Monitor) Get(name string) *ChanState {

	ch, found := this.chans[name]
	if !found {
		return nil
	}

	return &ChanState{
		Len: reflect.ValueOf(ch).Len(),
		Cap: reflect.ValueOf(ch).Cap(),
	}

}

// Get the channel states map[string]*ChanState of all the monitored channels. Keyed by channel name.
func (this *Monitor) GetAll() map[string]*ChanState {

	results := make(map[string]*ChanState)

	for name, _ := range this.chans {
		results[name] = this.Get(name)
	}

	return results

}
