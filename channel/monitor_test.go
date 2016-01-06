package channel

import "testing"

func TestAdd(t *testing.T) {

	m := New()

	c := make(chan chan bool, 10)

	if err := m.Add("foo", c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	state := m.Get("foo")

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

func TestGet(t *testing.T) {

	m := New()

	c := make(chan struct {
		Bar int
		Meh string
	}, 10)

	if err := m.Add("foo", c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	state := m.Get("foo")

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

	state = m.Get("foo")

	if state == nil {
		t.Fatal("state was nil")
	}

	if state.Len != 1 {
		t.Errorf("expected len 0 got %d", state.Len)
	}

	if state.Cap != 10 {
		t.Errorf("expected cap 10 got %d", state.Cap)
	}
}

func TestGetAll(t *testing.T) {

	m := New()

	c := make(chan struct {
		Bar int
		Meh string
	}, 10)

	d := make(chan int, 4)

	if err := m.Add("foo", c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := m.Add("bar", d); err != nil {
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

	states := m.GetAll()
	if len(states) != 2 {
		t.Errorf("expected 2 states but got %d", len(states))
	}

	if states["bar"].Len != 2 {
		t.Errorf("expected bar state len 2 but got %d", states["bar"].Len)
	}
	if states["bar"].Cap != 4 {
		t.Errorf("expected bar state cap 4 but got %d", states["bar"].Cap)
	}

	if states["foo"].Len != 1 {
		t.Errorf("expected foo state len 1 but got %d", states["foo"].Len)
	}
	if states["foo"].Cap != 10 {
		t.Errorf("expected foo state cap 10 but got %d", states["foo"].Cap)
	}

}
