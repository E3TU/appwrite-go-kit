package functions

import (
	"testing"

	"github.com/appwrite/sdk-for-go/v7/client"
)

func TestNewReturnsSDKService(t *testing.T) {
	s := New(client.New())
	if s.Functions() == nil {
		t.Fatal("Functions() = nil")
	}
}
