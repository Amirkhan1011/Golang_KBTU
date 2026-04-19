package main

import "testing"

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		t.Errorf("Add(2, 3) = %d; want %d", got, want)
	}
}

func TestAddTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 2, 3, 5},
		{"positive + zero", 5, 0, 5},
		{"negative + positive", -1, 4, 3},
		{"both negative", -2, -3, -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
	}{
		{"normal division", 10, 2, 5, false},
		{"negative dividend", -10, 2, -5, false},
		{"zero dividend", 0, 5, 0, false},
		{"division by zero", 5, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Divide(%d, %d) error = %v; wantErr %v", tt.a, tt.b, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Divide(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 5, 3, 2},
		{"positive minus zero", 5, 0, 5},
		{"negative minus positive", -3, 4, -7},
		{"both negative", -2, -3, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
