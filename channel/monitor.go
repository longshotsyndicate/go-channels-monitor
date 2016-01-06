package channel

import (
	"fmt"
	"reflect"
)

type Monitor struct {
	chans map[string]interface{}
}

func New() *Monitor {
	return &Monitor{
		chans: make(map[string]interface{}),
	}
}

func (this *Monitor) Add(name string, channel interface{}) error {

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

type ChanState struct {
	Len int
	Cap int
}

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

func (this *Monitor) GetAll() map[string]*ChanState {

	results := make(map[string]*ChanState)

	for name, _ := range this.chans {
		results[name] = this.Get(name)
	}

	return results

}
