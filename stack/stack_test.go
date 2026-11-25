package stack

import (
	"reflect"
	"testing"
)

func TestSubjectExample(t *testing.T) {
	// Exemple du sujet : 2 1 3 6 8 5
	s := New([]int{2, 1, 3, 6, 8, 5})
	for _, op := range []string{"sa", "pb", "pb", "pb", "rb", "rrr", "pa", "pa", "pa"} {
		if err := s.Exec(op); err != nil {
			t.Fatal(err)
		}
	}
	want := []int{1, 2, 3, 5, 6, 8}
	if !reflect.DeepEqual(s.A, want) || len(s.B) != 0 {
		t.Fatalf("a=%v b=%v, want a=%v b=[]", s.A, s.B, want)
	}
}

func TestOperations(t *testing.T) {
	tests := []struct {
		op   string
		a, b []int
	}{
		{"sa", []int{2, 1, 3}, []int{4, 5, 6}},
		{"sb", []int{1, 2, 3}, []int{5, 4, 6}},
		{"ss", []int{2, 1, 3}, []int{5, 4, 6}},
		{"pa", []int{4, 1, 2, 3}, []int{5, 6}},
		{"pb", []int{2, 3}, []int{1, 4, 5, 6}},
		{"ra", []int{2, 3, 1}, []int{4, 5, 6}},
		{"rb", []int{1, 2, 3}, []int{5, 6, 4}},
		{"rr", []int{2, 3, 1}, []int{5, 6, 4}},
		{"rra", []int{3, 1, 2}, []int{4, 5, 6}},
		{"rrb", []int{1, 2, 3}, []int{6, 4, 5}},
		{"rrr", []int{3, 1, 2}, []int{6, 4, 5}},
	}
	for _, tt := range tests {
		s := &Stacks{A: []int{1, 2, 3}, B: []int{4, 5, 6}}
		if err := s.Exec(tt.op); err != nil {
			t.Fatalf("%s: %v", tt.op, err)
		}
		if !reflect.DeepEqual(s.A, tt.a) || !reflect.DeepEqual(s.B, tt.b) {
			t.Errorf("%s: a=%v b=%v, want a=%v b=%v", tt.op, s.A, s.B, tt.a, tt.b)
		}
	}
}

func TestEmptyStacksDoNotPanic(t *testing.T) {
	s := New(nil)
	for _, op := range []string{"sa", "sb", "ss", "pa", "pb", "ra", "rb", "rr", "rra", "rrb", "rrr"} {
		if err := s.Exec(op); err != nil {
			t.Fatal(err)
		}
	}
}

func TestUnknownInstruction(t *testing.T) {
	s := New([]int{1, 2})
	for _, op := range []string{"", "SA", "sa ", "rrrr", "p"} {
		if s.Exec(op) == nil {
			t.Errorf("%q devrait être refusée", op)
		}
	}
}

func TestParseArgs(t *testing.T) {
	ok := map[string][]int{
		"2 1 3 6 5 8": {2, 1, 3, 6, 5, 8},
		"  -5   +3 0": {-5, 3, 0},
		"":            {},
		"2147483647":  {2147483647},
		"-2147483648": {-2147483648},
	}
	for in, want := range ok {
		got, err := ParseArgs([]string{in})
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("ParseArgs(%q) = %v, %v ; want %v", in, got, err, want)
		}
	}

	bad := []string{"0 one 2 3", "1 2 1", "2147483648", "1.5", "3 - 4", "--1"}
	for _, in := range bad {
		if _, err := ParseArgs([]string{in}); err == nil {
			t.Errorf("ParseArgs(%q) devrait échouer", in)
		}
	}
}

func TestParseArgsMultipleArgs(t *testing.T) {
	got, err := ParseArgs([]string{"3", "2 1", "0"})
	if err != nil || !reflect.DeepEqual(got, []int{3, 2, 1, 0}) {
		t.Fatalf("got %v, %v", got, err)
	}
}
