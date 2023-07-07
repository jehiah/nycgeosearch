package nycgeosearch

import (
	"context"
	"testing"
)

func TestReverseGeocode(t *testing.T) {
	got, err := PlanningLabs.ReverseGeocode(context.Background(), Location{
		Lat: 40.7484, Lng: -73.9857,
	}, Options{Size: 1})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", got)
	if len(got.Features) != 1 {
		t.Fatalf("expected at 1 response got %d", len(got.Features))
	}
	if s := got.Features[0].PropertyMustString("name"); s != "340 FIFTH AVENUE" {
		t.Fatalf("got %q", s)
	}
}
