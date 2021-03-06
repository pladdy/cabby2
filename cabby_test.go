package cabby

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

func TestAPIRootValidate(t *testing.T) {
	tests := []struct {
		apiRoot     APIRoot
		expectError bool
	}{
		{APIRoot{Path: "foo"}, true},
		{APIRoot{}, true},
		{APIRoot{Path: "foo", Title: "title"}, true},
		{APIRoot{Path: "foo", Title: "title", Versions: []string{"taxii-2.1"}}, true},
		{APIRoot{Path: "foo", Title: "title", Versions: []string{TaxiiVersion}}, false},
	}

	for _, test := range tests {
		result := test.apiRoot.Validate()

		if test.expectError && result == nil {
			t.Error("Got:", result, "Expected:", test.expectError)
		}
	}
}

func TestAPIRootIncludesMinVersion(t *testing.T) {
	tests := []struct {
		apiRoot  APIRoot
		expected bool
	}{
		{APIRoot{Path: "foo"}, false},
		{APIRoot{}, false},
		{APIRoot{Path: "foo", Title: "title"}, false},
		{APIRoot{Versions: []string{TaxiiVersion}}, true},
		{APIRoot{Versions: []string{TaxiiVersion, TaxiiVersion}}, true},
		{APIRoot{Versions: []string{TaxiiVersion, "taxii-2.1"}}, true},
	}

	for _, test := range tests {
		result := test.apiRoot.IncludesMinVersion(test.apiRoot.Versions)

		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestConfigParse(t *testing.T) {
	c := Config{}.Parse("config/cabby.example.json")

	if c.Host != "localhost" {
		t.Error("Got:", "localhost", "Expected:", "localhost")
	}
	if c.Port != 1234 {
		t.Error("Got:", strconv.Itoa(1234), "Expected:", strconv.Itoa(1234))
	}
}

func TestParseConfigNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered", r)
		}
	}()

	_ = Config{}.Parse("foo/bar")
	t.Error("Failed to panic with an unknown resource")
}

func TestConfigParseInvalidJSON(t *testing.T) {
	invalidJSON := "invalid.json"

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered", r)
			os.Remove(invalidJSON)
		}
	}()

	ioutil.WriteFile(invalidJSON, []byte("invalid"), 0644)
	Config{}.Parse(invalidJSON)
	t.Error("Failed to panic with an unknown resource")
}

func TestDiscoveryValidate(t *testing.T) {
	tests := []struct {
		discovery   Discovery
		expectError bool
	}{
		{Discovery{Title: "foo"}, false},
		{Discovery{}, true},
	}

	for _, test := range tests {
		result := test.discovery.Validate()

		if test.expectError && result == nil {
			t.Error("Got:", result, "Expected:", test.expectError)
		}
	}
}

func TestNewCollection(t *testing.T) {
	tests := []struct {
		idString    string
		shouldError bool
	}{
		{"invalid", true},
		{uuid.Must(uuid.NewV4()).String(), false},
	}

	for _, test := range tests {
		c, err := NewCollection(test.idString)

		if test.shouldError && err == nil {
			t.Error("Test with id of", test.idString, "should produce an error!")
		}

		if err == nil && c.ID.String() != test.idString {
			t.Error("Got:", c.ID.String(), "Expected:", test.idString)
		}
	}

	// test if 'collections' is passed; return a uuid
	_, err := NewCollection("collections")
	if err != nil {
		t.Error("Got:", err, "Expected no error")
	}
}

func TestCollectionValidate(t *testing.T) {
	validID, _ := NewID()
	validTitle := "a title"

	tests := []struct {
		collection  Collection
		expectError bool
	}{
		{Collection{ID: validID, Title: validTitle}, false},
		{Collection{Title: validTitle}, true},
		{Collection{ID: validID}, true},
		{Collection{}, true},
	}

	for _, test := range tests {
		result := test.collection.Validate()

		if test.expectError && result == nil {
			t.Error("Got:", result, "Expected:", test.expectError)
		}
	}
}

func TestNewRange(t *testing.T) {
	invalidRange := Range{First: -1, Last: -1}

	tests := []struct {
		input       string
		resultRange Range
		isError     bool
	}{
		{"items 0-10", Range{First: 0, Last: 10}, false},
		{"items 0 10", invalidRange, true},
		{"items 10", invalidRange, true},
		{"", invalidRange, false},
	}

	for _, test := range tests {
		result, err := NewRange(test.input)
		if result != test.resultRange {
			t.Error("Got:", result, "Expected:", test.resultRange)
		}

		if err != nil && test.isError == false {
			t.Error("Got:", err, "Expected: no error")
		}
	}
}

func TestRangeString(t *testing.T) {
	tests := []struct {
		testRange Range
		expected  string
	}{
		{Range{First: 0, Last: 0}, "items 0-0"},
		{Range{First: 0, Last: 0, Total: 50}, "items 0-0/50"},
	}

	for _, test := range tests {
		result := test.testRange.String()
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestRangeValid(t *testing.T) {
	tests := []struct {
		testRange Range
		expected  bool
	}{
		{Range{First: 1, Last: 0}, false},
		{Range{First: 0, Last: 0}, true},
		{Range{First: 0, Last: -1}, false},
	}

	for _, test := range tests {
		result := test.testRange.Valid()
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestNewID(t *testing.T) {
	_, err := NewID()
	if err != nil {
		t.Error("Expected no error:", err)
	}
}

func TestIDFromString(t *testing.T) {
	uid := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	id, err := IDFromString(uid)
	if err != nil {
		t.Error("Expected no error:", err)
	}

	if id.String() != uid {
		t.Error("Got:", id.String(), "Expected:", uid)
	}
}

func TestIDFromStringBadInput(t *testing.T) {
	_, err := IDFromString("")
	if err == nil {
		t.Error("Expected an error")
	}

	_, err = IDFromString("fail")
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestIDUsingString(t *testing.T) {
	id, err := IDUsingString("")
	if err != nil {
		t.Fatal(err)
	}

	expected := "101afa45-b9bd-5b31-8734-0a59e5cc3db3"
	if id.String() != expected {
		t.Error("Got:", id.String(), "Expected:", expected)
	}
}

func TestIDIsEmpty(t *testing.T) {
	id, err := NewID()
	if err != nil {
		t.Fatal(err)
	}

	if id.IsEmpty() == true {
		t.Error("Expected to NOT be empty")
	}

	emptyID := ID{}
	if emptyID.IsEmpty() == false {
		t.Error("Expected ID to be empty")
	}
}

func TestNewStatus(t *testing.T) {
	_, err := NewStatus(1)
	if err != nil {
		t.Error("Got:", err, "Expected: no error")
	}
}

func TestNewStatusError(t *testing.T) {
	_, err := NewStatus(0)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestUserDefined(t *testing.T) {
	tests := []struct {
		user     User
		expected bool
	}{
		{user: User{Email: "foo"}, expected: true},
		{user: User{}, expected: false},
	}

	for _, test := range tests {
		result := test.user.Defined()
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestUserValidate(t *testing.T) {
	tests := []struct {
		user        User
		expectError bool
	}{
		{User{Email: "foo"}, true},
		{User{}, true},
		{User{Email: "no@no.no"}, false},
		{User{Email: "some-person@yaoo.co.uk"}, false},
	}

	for _, test := range tests {
		result := test.user.Validate()

		if test.expectError && result == nil {
			t.Error("Got:", result, "Expected:", test.expectError)
		}
	}
}
