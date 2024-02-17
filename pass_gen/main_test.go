package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	g := NewGenerator(true, true, true, true)
	password := g.GeneratePassword(12)

	if len(password) != 12 {
		t.Errorf("Generated password length is incorrect, got: %d, want: %d", len(password), 12)
	}

	alphabet := g.alphabet.GetAlphabet()
	for _, char := range password {
		if !strings.ContainsRune(alphabet, char) {
			t.Errorf("Generated password contains invalid character: %c", char)
		}
	}
}

func TestPasswordStrength(t *testing.T) {
	tests := []struct {
		password string
	}{
		{"a5JA!?~GCZ^t0)qj^dZRE1L!#s&l"}, //example strong password
		{"dYkX&\\0"},                     // example good password
		{"qy_p/uq~m"},                    // example medium password
		{"password"},                     // example weak password
	}

	strong := NewPassword(tests[0].password)
	strongStrength := strong.PasswordStrength()
	assert.GreaterOrEqual(t, strongStrength, 6)

	good := NewPassword(tests[1].password)
	goodStrength := good.PasswordStrength()
	assert.GreaterOrEqual(t, goodStrength, 4)

	medium := NewPassword(tests[2].password)
	mediumStrength := medium.PasswordStrength()
	assert.GreaterOrEqual(t, mediumStrength, 3)

	weak := NewPassword(tests[3].password)
	weakStrength := weak.PasswordStrength()
	assert.GreaterOrEqual(t, weakStrength, 2)

	// for _, test := range tests {
	// 	p := NewPassword(test.password)
	// 	strength := p.PasswordStrength()
	// 	if strength != test.expected {
	// 		t.Errorf("Password strength calculation failed for %s. Got: %d, Expected: %d", test.password, strength, test.expected)
	// 	}
	// }
}
func TestCalculateScore(t *testing.T) {
	tests := []struct {
		password string
		expected string
	}{
		{"a5JA!?~GCZ^t0)qj^dZRE1L!#s&l", "This is a very good password!"},
		{"dY6!kX&0", "Good password, but you can still do better!"},
		{"qy_p/uq~m", "Medium password. Try making it better!"},
		{"password", "This is a weak password. Generate a new one!"},
	}

	for _, test := range tests {
		p := NewPassword(test.password)
		score := p.CalculateScore()
		if score != test.expected {
			t.Errorf("Score calculation failed for %s. Got: %s, Expected: %s", test.password, score, test.expected)
		}
	}
}
