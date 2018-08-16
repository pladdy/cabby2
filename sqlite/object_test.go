package sqlite

import (
	"testing"

	"github.com/pladdy/cabby2/tester"
)

func TestObjectServiceObject(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ObjectService{DB: ds.DB}

	expected := tester.Object

	result, err := s.Object(expected.CollectionID.String(), expected.ID.String())
	if err != nil {
		t.Error("Got:", err, "Expected no error")
	}

	tester.CompareObject(result, expected, t)
}

func TestObjectServiceObjectQueryErr(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ObjectService{DB: ds.DB}

	_, err := s.DB.Exec("drop table stix_objects")
	if err != nil {
		t.Fatal(err)
	}

	expected := tester.Object
	_, err = s.Object(expected.CollectionID.String(), expected.ID.String())
	if err == nil {
		t.Error("Got:", err, "Expected an error")
	}
}

func TestObjectServiceObjectInvalidRawID(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ObjectService{DB: ds.DB}

	_, err := s.DB.Exec(`insert into stix_objects (id, type, created, modified, object, collection_id)
	                     values ('fail', 'fail', 'fail', 'fail', '{"fail": true}', 'fail')`)
	if err != nil {
		t.Fatal(err)
	}

	expected := tester.Object
	_, err = s.Object(expected.CollectionID.String(), "fail")
	if err == nil {
		t.Error("Got:", err, "Expected an error")
	}
}

func TestObjectsServiceObjects(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ObjectService{DB: ds.DB}

	expected := tester.Object

	results, err := s.Objects(expected.CollectionID.String())
	if err != nil {
		t.Error("Got:", err, "Expected no error")
	}

	if len(results) <= 0 {
		t.Error("Got:", len(results), "Expected: > 0")
	}

	result := results[0]
	tester.CompareObject(result, expected, t)
}

func TestObjectsServiceObjectsQueryErr(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ObjectService{DB: ds.DB}

	_, err := s.DB.Exec("drop table stix_objects")
	if err != nil {
		t.Fatal(err)
	}

	expected := tester.Object

	_, err = s.Objects(expected.CollectionID.String())
	if err == nil {
		t.Error("Got:", err, "Expected an error")
	}
}

func TestObjectServiceObjectsInvalidRawID(t *testing.T) {
	setupSQLite()
	ds := testDataStore()
	s := ObjectService{DB: ds.DB}

	_, err := s.DB.Exec(`insert into stix_objects (id, type, created, modified, object, collection_id)
	                     values ('fail', 'fail', 'fail', 'fail', '{"fail": true}', 'fail')`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Objects("fail")
	if err == nil {
		t.Error("Got:", err, "Expected an error")
	}
}