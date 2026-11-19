package discordlog

import (
	"os"
	"testing"
)

func TestShareQRLink(t *testing.T) {
	if os.Getenv("WESH_TEST_DISCORD") != "1" {
		t.Skip()
	}
	err := ShareQRLink("go test -v", QRCodeRoom, "Lorem Ipsum", "https://tone.onl", "https://tone.onl")
	if err != nil {
		t.Fatal(err)
	}
}
