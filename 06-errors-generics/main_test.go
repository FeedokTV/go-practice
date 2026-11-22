package errorsgenerics

import (
	"errors"
	"fmt"
	"testing"
)

func fakeDBLayer(id string) error {
	return NotFoundError{ID: id}
}

func fakeServiceLayer(id string) error {
	return fmt.Errorf("db error: %w", fakeDBLayer(id))
}

func getObjectById(id string) error {
	return fmt.Errorf("app error: %w", fakeServiceLayer(id))
}

func getObjectByIdNonWrap(id string) error {
	return fmt.Errorf("app error: %v", fakeServiceLayer(id))
}

func TestWrapping(t *testing.T) {
	var nfError NotFoundError

	err := getObjectById("a1b2c3-dddd4-fffff")

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if !errors.As(err, &nfError) {
		t.Errorf("received unexpected type of error")
	}

	if nfError.ID != "a1b2c3-dddd4-fffff" {
		t.Errorf("expected id from error: a1b2c3-dddd4-fffff, got %s", nfError.ID)
	}

	var nfAnother NotFoundError
	if errors.As(getObjectByIdNonWrap("a1b2c3-dddd4-fffff"), &nfAnother) {
		t.Errorf("expected wrapped error, got straight")
	}

}

type TestEntityString struct {
	Name string
	ID   int
}

type TestEntityInt struct {
	Val int
}

func TestMaxByStr(t *testing.T) {

	tests := []struct {
		Name          string
		Entities      []TestEntityString
		GetterFunc    func(e TestEntityString) string
		ExpectedValue string
		ExpectedID    int
		ExpectOk      bool
	}{
		{
			"return biggest name",
			[]TestEntityString{{Name: "Book1", ID: 1}, {Name: "The Big Story Of Big Neverendless Night", ID: 2}},
			func(e TestEntityString) string { return e.Name },
			"The Big Story Of Big Neverendless Night",
			2,
			true,
		},
		{
			"return from three elems",
			[]TestEntityString{{Name: "Book1", ID: 1}, {Name: "Who Am I?", ID: 2}, {Name: "Zeo's Big Story Of Big Neverendless Night", ID: 3}},
			func(e TestEntityString) string { return e.Name },
			"Zeo's Big Story Of Big Neverendless Night",
			3,
			true,
		},
		{
			"from empty",
			[]TestEntityString{},
			func(e TestEntityString) string { return e.Name },
			"",
			0,
			false,
		},
		{
			"same keys",
			[]TestEntityString{{Name: "Book1", ID: 1}, {Name: "Book1", ID: 2}},
			func(e TestEntityString) string { return e.Name },
			"Book1",
			1,
			true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.Name, func(t *testing.T) {
			res, ok := MaxBy(tt.Entities, tt.GetterFunc)

			if tt.ExpectOk != ok {
				t.Errorf("expected status: %t, got %t", tt.ExpectOk, ok)
			}

			if res.Name != tt.ExpectedValue {
				t.Errorf("expected: %s, got: %s", tt.ExpectedValue, res.Name)
			}

			if res.ID != tt.ExpectedID {
				t.Errorf("expected obj num: %d, got: %d", tt.ExpectedID, res.ID)
			}
		})
	}
}

func TestMaxByInt(t *testing.T) {

	tests := []struct {
		Name          string
		Entities      []TestEntityInt
		GetterFunc    func(e TestEntityInt) int
		ExpectedValue int
		ExpectOk      bool
	}{
		{
			"return biggest Val",
			[]TestEntityInt{{Val: 100}, {Val: 1000}},
			func(e TestEntityInt) int { return e.Val },
			1000,
			true,
		},
		{
			"return from three elems",
			[]TestEntityInt{{Val: 100}, {Val: 1000}, {Val: 9999}},
			func(e TestEntityInt) int { return e.Val },
			9999,
			true,
		},
		{
			"negatives",
			[]TestEntityInt{{Val: -100}, {Val: -10}, {Val: -1}},
			func(e TestEntityInt) int { return e.Val },
			-1,
			true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.Name, func(t *testing.T) {
			res, ok := MaxBy(tt.Entities, tt.GetterFunc)

			if tt.ExpectOk != ok {
				t.Errorf("expected status: %t, got %t", tt.ExpectOk, ok)
			}

			if res.Val != tt.ExpectedValue {
				t.Errorf("expected: %d, got: %d", tt.ExpectedValue, res)
			}
		})
	}
}
