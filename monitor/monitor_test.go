package monitor

import "testing"

func clear() {
	chans = make(map[key]interface{})
}

func TestAddNamed(t *testing.T) {

	clear()

	c := make(chan chan bool, 10)

	err := AddNamed("foo", "", c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	state := Get("foo", "")

	if state == nil {
		t.Fatal("state was nil")
	}

	if state.Len != 0 {
		t.Errorf("expected len 0 got %d", state.Len)
	}

	if state.Cap != 10 {
		t.Errorf("expected cap 10 got %d", state.Cap)
	}
}

func TestInstance(t *testing.T) {

	clear()

	c := make(chan chan bool, 10)
	d := make(chan bool, 20)

	AddNamed("c", "foo-1", c)
	AddNamed("d", "foo-2", d)

	states := GetAll()

	if len(states) != 2 {
		t.Errorf("expected 2 states but got %d", len(states))
	}

}

func TestGet(t *testing.T) {

	clear()

	c := make(chan struct {
		Bar int
		Meh string
	}, 10)

	if err := AddNamed("foo", "", c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	state := Get("foo", "")

	if state == nil {
		t.Fatal("state was nil")
	}

	if state.Len != 0 {
		t.Errorf("expected len 0 got %d", state.Len)
	}

	if state.Cap != 10 {
		t.Errorf("expected cap 10 got %d", state.Cap)
	}

	c <- struct {
		Bar int
		Meh string
	}{
		1, "meh",
	}

	state = Get("foo", "")

	if state == nil {
		t.Fatal("state was nil")
	}

	if state.Len != 1 {
		t.Errorf("expected len 1 got %d", state.Len)
	}

	if state.Cap != 10 {
		t.Errorf("expected cap 10 got %d", state.Cap)
	}

	<-c
	state = Get("foo", "")

	if state == nil {
		t.Fatal("state was nil")
	}

	if state.Len != 0 {
		t.Errorf("expected len 0 got %d", state.Len)
	}

	if state.Cap != 10 {
		t.Errorf("expected cap 10 got %d", state.Cap)
	}
}

func TestGetAll(t *testing.T) {

	clear()

	c := make(chan struct {
		Bar int
		Meh string
	}, 10)

	d := make(chan int, 4)

	foo := "c"
	if err := AddNamed(foo, "", c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bar := "d"
	if err := AddNamed(bar, "", d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c <- struct {
		Bar int
		Meh string
	}{
		1, "meh",
	}

	d <- 1
	d <- 2

	states := GetAll()
	if len(states) != 2 {
		t.Errorf("expected 2 states but got %d", len(states))
	}

	for k, state := range states {
		switch k {
		case foo:

			if state.Len != 1 {
				t.Errorf("expected %s state len 1 but got %d", k, state.Len)
			}
			if state.Cap != 10 {
				t.Errorf("expected %s state cap 10 but got %d", k, state.Cap)
			}
		case bar:

			if state.Len != 2 {
				t.Errorf("expected %s state len 2 but got %d", k, state.Len)
			}
			if state.Cap != 4 {
				t.Errorf("expected %s state cap 4 but got %d", k, state.Cap)
			}
		default:
			t.Errorf("invalid line reference: %s", k)
		}
	}
}
