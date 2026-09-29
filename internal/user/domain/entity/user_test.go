package entity

import (
	"errors"
	"testing"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

func ptr(s string) *string { return &s }

func TestApplyUpdate(t *testing.T) {
	u := &User{Name: "Alice", Email: "alice@example.com"}
	if err := u.ApplyUpdate(ptr(" Bob "), ptr("BOB@example.com")); err != nil {
		t.Fatal(err)
	}
	if u.Name != "Bob" || u.Email != "bob@example.com" {
		t.Errorf("got %+v", u)
	}

	tests := []struct {
		name  string
		in    [2]*string
		field string
	}{
		{"empty body", [2]*string{nil, nil}, "body"},
		{"bad email", [2]*string{nil, ptr("bad")}, "email"},
		{"blank name", [2]*string{ptr("  "), nil}, "name"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := u.ApplyUpdate(tc.in[0], tc.in[1])
			var e *errs.Error
			if !errors.As(err, &e) || e.Kind != errs.KindValidation || e.Fields[tc.field] == "" {
				t.Fatalf("want validation error on %q, got %v", tc.field, err)
			}
			if u.Name != "Bob" || u.Email != "bob@example.com" {
				t.Errorf("invalid update changed the user: %+v", u)
			}
		})
	}
}
