package sqlite

import (
	"context"
	"testing"

	"github.com/pladdy/cabby2/tester"
)

func TestStatusServiceCreateStatus(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ds.StatusService()

	test := tester.Status

	err := s.CreateStatus(context.Background(), test)
	if err != nil {
		t.Error("Got:", err)
	}

	result, err := s.Status(context.Background(), test.ID.String())
	if err != nil {
		t.Error("Got:", err)
	}

	passed := tester.CompareStatus(result, test)
	if !passed {
		t.Error("Comparison failed")
	}
}

func TestStatusServiceCreateStatusFail(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ds.StatusService()

	_, err := ds.DB.Exec("drop table taxii_status")
	if err != nil {
		t.Fatal(err)
	}

	test := tester.Status

	err = s.CreateStatus(context.Background(), test)
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestStatusServiceStatus(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ds.StatusService()

	// create a status
	expected := tester.Status
	err := s.CreateStatus(context.Background(), expected)
	if err != nil {
		t.Error("Got:", err)
	}

	result, err := s.Status(context.Background(), expected.ID.String())
	if err != nil {
		t.Error("Got:", err, "Expected no error")
	}

	passed := tester.CompareStatus(result, expected)
	if !passed {
		t.Error("Comparison failed")
	}
}

func TestStatusServiceStatusQueryErr(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ds.StatusService()

	_, err := ds.DB.Exec("drop table taxii_status")
	if err != nil {
		t.Fatal(err)
	}

	expected := tester.Status

	_, err = s.Status(context.Background(), expected.ID.String())
	if err == nil {
		t.Error("Got:", err, "Expected an error")
	}
}

func TestStatusServiceUpdateStatus(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ds.StatusService()

	// create a status
	expected := tester.Status
	err := s.CreateStatus(context.Background(), expected)
	if err != nil {
		t.Error("Got:", err)
	}

	// update the status
	expected.TotalCount = 3
	expected.SuccessCount = 0
	expected.FailureCount = 1

	err = s.UpdateStatus(context.Background(), expected)
	if err != nil {
		t.Error("Got:", err)
	}

	// verify it's updated
	expected.PendingCount = 2
	result, err := s.Status(context.Background(), expected.ID.String())
	if err != nil {
		t.Error("Got:", err)
	}

	passed := tester.CompareStatus(result, expected)
	if !passed {
		t.Error("Comparison failed")
	}

	// complete the status and check
	expected.SuccessCount = 2
	err = s.UpdateStatus(context.Background(), expected)
	if err != nil {
		t.Error("Got:", err)
	}

	// verify it's updated
	expected.PendingCount = 0
	expected.Status = "complete"

	result, err = s.Status(context.Background(), expected.ID.String())
	if err != nil {
		t.Error("Got:", err)
	}

	passed = tester.CompareStatus(result, expected)
	if !passed {
		t.Error("Comparison failed")
	}
}
