package channel

import "testing"

func TestAdd(t *testing.T) {

	m := New()

	c := make(chan chan bool, 10)

	name, err := m.Add(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	state := m.Get(name)

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

	if err := m.AddNamed("foo", c); err != nil {
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
		t.Errorf("expected len 1 got %d", state.Len)
	}

	if state.Cap != 10 {
		t.Errorf("expected cap 10 got %d", state.Cap)
	}

	<-c
	state = m.Get("foo")

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

	m := New()

	c := make(chan struct {
		Bar int
		Meh string
	}, 10)

	d := make(chan int, 4)

	foo, err := m.Add(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bar, err := m.Add(d)
	if err != nil {
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
