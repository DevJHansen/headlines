package utils

import "testing"

func TestGetImgUrlFromAtr(t *testing.T) {
	got := GetImgUrlFromStyleAtr("background: url('https://i0.wp.com/informante.web.na/wp-content/uploads/2024/08/Two-die-in-car-accident-on-Gobabis-road-Informante-Image-workspace_Facebook.jpg?fit=1344%2C888&ssl=1') no-repeat center;")
	expect := "https://i0.wp.com/informante.web.na/wp-content/uploads/2024/08/Two-die-in-car-accident-on-Gobabis-road-Informante-Image-workspace_Facebook.jpg?fit=1344%2C888&ssl=1"

	if got != expect {
		t.Errorf("Got %q, expected: %q", got, expect)
	}
}
