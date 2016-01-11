package monitor

import (
	"fmt"
	"reflect"
	"sync"
)

var chans = make(map[string]interface{})
var chmu sync.RWMutex

// AddNamed adds a channel to be monitor and associates the channel with this name.
func AddNamed(name string, channel interface{}) error {

	//reflect on the input to get the correct channel type.
	if reflect.TypeOf(channel).Kind() != reflect.Chan {
		return fmt.Errorf("invalid input type %v for input param channel, must be of type chan", channel)
	}

	chmu.Lock()
	defer chmu.Unlock()
	if _, found := chans[name]; found {
		return fmt.Errorf("channel with name: %s already being monitored.", name)
	}
	chans[name] = channel

	return nil
}

// ChanState struct holding Length and Capacity.
type ChanState struct {
	Len int `json:"length"`
	Cap int `json:"capacity"`
}

// Get returns the channel state for a give channel name.
func Get(name string) *ChanState {

	chmu.RLock()
	defer chmu.RUnlock()
	ch, found := chans[name]
	if !found {
		return nil
	}

	return &ChanState{
		Len: reflect.ValueOf(ch).Len(),
		Cap: reflect.ValueOf(ch).Cap(),
	}

}

// Get the channel states map[string]*ChanState of all the monitored channels. Keyed by channel name.
func GetAll() map[string]*ChanState {

	results := make(map[string]*ChanState)

	chmu.RLock()
	defer chmu.RUnlock()
	for name, _ := range chans {
		results[name] = Get(name)
	}

	return results

}
