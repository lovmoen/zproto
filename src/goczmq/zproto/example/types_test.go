package example

import (
	"testing"

	"github.com/zeromq/goczmq"
)

// Yay! Test function.
func TestTypes(t *testing.T) {

	// Create pair of sockets we can send through

	// Output socket
	output := goczmq.NewSock(goczmq.DEALER)
	defer output.Destroy()

	routingId := "Shout"
	output.SetIdentity(routingId)
	_, err := output.Bind("inproc://selftest-types")
	if err != nil {
		t.Fatal(err)
	}
	defer output.Unbind("inproc://selftest-types")

	// Input socket
	input := goczmq.NewSock(goczmq.ROUTER)
	defer input.Destroy()

	err = input.Connect("inproc://selftest-types")
	if err != nil {
		t.Fatal(err)
	}
	defer input.Disconnect("inproc://selftest-types")

	// Create a Types message and send it through the wire
	types := NewTypes()

	types.sequence = 123

	types.ClientForename = "Life is short but Now lasts for ever"

	types.ClientSurname = "Life is short but Now lasts for ever"

	types.ClientMobile = "Life is short but Now lasts for ever"

	types.ClientEmail = "Life is short but Now lasts for ever"

	types.SupplierForename = "Life is short but Now lasts for ever"

	types.SupplierSurname = "Life is short but Now lasts for ever"

	types.SupplierMobile = "Life is short but Now lasts for ever"

	types.SupplierEmail = "Life is short but Now lasts for ever"

	err = types.Send(output)
	if err != nil {
		t.Fatal(err)
	}
	transit, err := Recv(input)
	if err != nil {
		t.Fatal(err)
	}

	tr := transit.(*Types)

	if tr.sequence != 123 {
		t.Fatalf("expected %d, got %d", 123, tr.sequence)
	}

	if tr.ClientForename != "Life is short but Now lasts for ever" {
		t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.ClientForename)
	}

	if tr.ClientSurname != "Life is short but Now lasts for ever" {
		t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.ClientSurname)
	}

	if tr.ClientMobile != "Life is short but Now lasts for ever" {
		t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.ClientMobile)
	}

	if tr.ClientEmail != "Life is short but Now lasts for ever" {
		t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.ClientEmail)
	}

	if tr.SupplierForename != "Life is short but Now lasts for ever" {
		t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.SupplierForename)
	}

	if tr.SupplierSurname != "Life is short but Now lasts for ever" {
		t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.SupplierSurname)
	}

	if tr.SupplierMobile != "Life is short but Now lasts for ever" {
		t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.SupplierMobile)
	}

	if tr.SupplierEmail != "Life is short but Now lasts for ever" {
		t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.SupplierEmail)
	}

	err = tr.Send(input)
	if err != nil {
		t.Fatal(err)
	}

	transit, err = Recv(output)
	if err != nil {
		t.Fatal(err)
	}

	if routingId != string(tr.RoutingId()) {
		t.Fatalf("expected %s, got %s", routingId, string(tr.RoutingId()))
	}
}
