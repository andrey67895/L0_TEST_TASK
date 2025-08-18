package domain

import (
	"testing"
)

func TestOrder_IsValid(t *testing.T) {
	cases := []struct {
		name  string
		order Order
		want  bool
	}{
		{
			name:  "Валидный UID заказа",
			order: Order{OrderUID: "12345"},
			want:  true,
		},
		{
			name:  "Невалидный UID заказа",
			order: Order{OrderUID: ""},
			want:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.order.IsValid()
			if got != tc.want {
				t.Errorf("IsValid() = %v, want %v", got, tc.want)
			}
		})
	}
}
