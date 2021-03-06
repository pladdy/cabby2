package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	cabby "github.com/pladdy/cabby2"
)

func TestGetToken(t *testing.T) {
	tests := []struct {
		url      string
		index    int
		expected string
	}{
		{"/api_root/collections/collection_id/objects/stix_id", 0, ""},
		{"/api_root/collections/collection_id/objects/stix_id", 1, "api_root"},
		{"/api_root/collections/collection_id/objects/stix_id", 3, "collection_id"},
		{"/api_root/collections/collection_id/objects/stix_id", 5, "stix_id"},
		{"/api_root/collections/collection_id/objects/stix_id", 7, ""},
	}

	for _, test := range tests {
		result := getToken(test.url, test.index)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestLastURLPathToken(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/collections/", "collections"},
		{"/collections/someId", "someId"},
	}

	for _, test := range tests {
		result := lastURLPathToken(test.path)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestTakeAddedAfter(t *testing.T) {
	tests := []struct {
		request    *http.Request
		addedAfter string
	}{
		{httptest.NewRequest("GET", "/foo/bar/baz", nil), ""},
		{httptest.NewRequest("GET", "/foo/bar/baz?added_after=invalid", nil), ""},
		{httptest.NewRequest("GET", "/foo/bar/baz?added_after=2016-02-21T05:01:01.000Z", nil), "2016-02-21T05:01:01Z"},
		{httptest.NewRequest("GET", "/foo/bar/baz?added_after=2016-02-21T05:01:01.123Z", nil), "2016-02-21T05:01:01.123Z"},
	}

	for _, test := range tests {
		result := takeAddedAfter(test.request)
		if result != test.addedAfter {
			t.Error("Got:", result, "Expected:", test.addedAfter)
		}
	}
}

func TestTakeCollectionID(t *testing.T) {
	tests := []struct {
		request *http.Request
		id      string
	}{
		{httptest.NewRequest("GET", "/api_root_path/collections/collectionID", nil), "collectionID"},
		{httptest.NewRequest("GET", "/api_root_path/collections/collectionID/objects", nil), "collectionID"},
		{httptest.NewRequest("GET", "/api_root_path/collections/", nil), ""},
	}

	for _, test := range tests {
		result := takeCollectionID(test.request)
		if result != test.id {
			t.Error("Got:", result, "Expected:", test.id)
		}
	}
}

func TestTakeMatchIDs(t *testing.T) {
	tests := []struct {
		request *http.Request
		matchID string
	}{
		{httptest.NewRequest("GET", "/foo/bar/baz", nil), ""},
		{httptest.NewRequest("GET", "/foo/bar/baz?match[id]=some-id", nil), "some-id"},
		{httptest.NewRequest("GET", "/foo/bar/baz?match[id]=id1,id2,id3", nil), "id1,id2,id3"},
		{httptest.NewRequest("GET", "/foo/bar/baz?match[id]=id1&match[id]=id2&match[id]=id3", nil), "id1,id2,id3"},
	}

	for _, test := range tests {
		result := takeMatchIDs(test.request)
		if result != test.matchID {
			t.Error("Got:", result, "Expected:", test.matchID)
		}
	}
}

func TestTakeMatchTypes(t *testing.T) {
	tests := []struct {
		request   *http.Request
		matchType string
	}{
		{httptest.NewRequest("GET", "/foo/bar/baz", nil), ""},
		{httptest.NewRequest("GET", "/foo/bar/baz?match[type]=some-type", nil), "some-type"},
		{httptest.NewRequest("GET", "/foo/bar/baz?match[type]=type1,type2,type3", nil), "type1,type2,type3"},
		{httptest.NewRequest("GET", "/foo/bar/baz?match[type]=type1&match[type]=type2&match[type]=type3", nil),
			"type1,type2,type3"},
	}

	for _, test := range tests {
		result := takeMatchTypes(test.request)
		if result != test.matchType {
			t.Error("Got:", result, "Expected:", test.matchType)
		}
	}
}

func TestTakeMatchVersions(t *testing.T) {
	tests := []struct {
		request      *http.Request
		matchVersion string
	}{
		{httptest.NewRequest("GET", "/foo/bar/baz", nil), ""},
		{httptest.NewRequest("GET", "/foo/bar/baz?match[version]=some-version", nil), "some-version"},
		{httptest.NewRequest("GET", "/foo/bar/baz?match[version]=version1,version2,version3", nil),
			"version1,version2,version3"},
		{httptest.NewRequest("GET", "/foo/bar/baz?match[version]=version1&match[version]=version2&match[version]=version3", nil),
			"version1,version2,version3"},
	}

	for _, test := range tests {
		result := takeMatchVersions(test.request)
		if result != test.matchVersion {
			t.Error("Got:", result, "Expected:", test.matchVersion)
		}
	}
}

func TestTakeCollectionAccessInvalidCollection(t *testing.T) {
	// create a request with a valid context BUT a path with an invalid collection in it
	req := httptest.NewRequest("GET", "/foo/bar/baz", nil)

	ca := takeCollectionAccess(req)
	empty := cabby.CollectionAccess{}

	if ca != empty {
		t.Error("Got:", ca, "Expected:", empty)
	}
}
