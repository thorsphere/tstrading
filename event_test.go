// Copyright (c) 2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public Licence v3.0
// that can be found in the LICENSE file.
package tstrading_test

// Import standard library packages, tseventserver, tsfio and tserrs
import (
	"strings" // strings for building string output in tests
	"testing" // testing for writing test cases
	"time"    // time for working with time and dates

	"github.com/thorsphere/tserr"     // tserr for custom error handling
	"github.com/thorsphere/tsfio"     // tsfio for file input/output operations, including handling golden files
	"github.com/thorsphere/tstrading" // tstrading for testing
)

var (
	// Define some sample events for testing purposes
	evNfp *tstrading.Event = &tstrading.Event{
		Name:     "Non-Farm Payrolls",
		Time:     time.Date(2024, 7, 5, 8, 30, 0, 0, time.UTC),
		Country:  "US",
		Actual:   new(200.0),
		Estimate: new(180.0),
		Previous: new(150.0),
		Unit:     "K",
		Impact:   tstrading.ImpactHigh,
		Source:   "Bureau of Labor Statistics",
	}
	evGdp24 *tstrading.Event = &tstrading.Event{
		Name:     "GDP Growth Rate",
		Time:     time.Date(2024, 7, 10, 8, 30, 0, 0, time.UTC),
		Country:  "US",
		Actual:   new(3.5),
		Estimate: new(3.0),
		Previous: new(2.8),
		Unit:     "%",
		Impact:   tstrading.ImpactMedium,
		Source:   "Bureau of Economic Analysis",
	}
	evGdp30 *tstrading.Event = &tstrading.Event{
		Name:     "GDP Growth Rate",
		Time:     time.Date(2030, 7, 10, 8, 30, 0, 0, time.UTC),
		Country:  "US",
		Actual:   nil,
		Estimate: nil,
		Previous: nil,
		Unit:     "%",
		Impact:   tstrading.ImpactLow,
		Source:   "Bureau of Economic Analysis",
	}
	// Define a slice of events for testing purposes
	evs []*tstrading.Event = []*tstrading.Event{
		evNfp,
		evGdp24,
		evGdp30,
	}
)

// TestEvents tests the String method of the Event struct by comparing the output to a golden file.
func TestEvents(t *testing.T) {
	// Create a formatted string representation of the sample events using the String method of the Event struct
	var out strings.Builder
	// Iterate over each event in the sample events slice and append its string representation to the output string
	for _, ev := range evs {
		out.WriteString(ev.String())
		out.WriteString("\n")
	}
	// Compare the output to a golden file using the EvalGoldenFile function from the tsfio package,
	// and if there is an error, fail the test with the error message
	if e := tsfio.EvalGoldenFile(&tsfio.Testcase{Name: "events", Data: out.String()}); e != nil {
		t.Fatal(e)
	}
}

// TestWrongImpact tests the String method of the ImpactLevel type with an invalid impact level value
// and expects the output to be "unknown".
func TestWrongImpact(t *testing.T) {
	var i tstrading.ImpactLevel = 99 // Invalid impact level
	// The expected output for an invalid impact level should be "unknown"
	expected := "unknown"
	// Get the actual string representation of the impact level using the String method
	actual := i.String()
	// If the actual output does not match the expected output, fail the test with an error message
	// indicating the mismatch
	if actual != expected {
		t.Fatal(tserr.EqualStr(&tserr.EqualStrArgs{Var: "ImpactLevel", Actual: actual, Want: expected}))
	}
}

// TestStringNil1 tests the String method of the Event struct with a nil event.
func TestStringNil1(t *testing.T) {
	var ev *tstrading.Event = nil
	// Get the string representation of the nil event
	actual := ev.String()
	// The expected string representation of a nil event should be "<nil>"
	expected := "<nil>"
	// Check if the string representation of a nil event is "<nil>"
	if actual != expected {
		// If the string representation is not "<nil>", fail the test with an error message
		t.Fatal(tserr.EqualStr(&tserr.EqualStrArgs{Var: "String", Actual: actual, Want: expected}))
	}
}

// TestNearEqual tests the NearEqual method of the Event struct by comparing two events that are nearly equal.
func TestNearEqual(t *testing.T) {
	// Define a sample event for testing purposes
	ev1 := evNfp
	ev2 := &tstrading.Event{
		Name:     ev1.Name,
		Time:     ev1.Time,
		Country:  ev1.Country,
		Actual:   new(*(ev1.Actual)),
		Estimate: new(*(ev1.Estimate)),
		Previous: new(*(ev1.Previous)),
		Unit:     ev1.Unit,
		Impact:   ev1.Impact,
		Source:   ev1.Source,
	}
	// Check if the events are nearly equal
	if !ev1.NearEqual(ev2) {
		t.Fatal(tserr.Equal(&tserr.EqualArgs{X: ev1.String(), Y: ev2.String()}))
	}
	// Modify the actual value of the second event
	act := *(ev2.Actual) + 0.1
	ev2.Actual = &act
	// Check if the events are not nearly equal
	if ev1.NearEqual(ev2) {
		t.Fatal(tserr.NotEqual(&tserr.NotEqualArgs{X: ev1.String(), Y: ev2.String()}))
	}
}

// TestNearEqualNil tests the NearEqual method of the Event struct with nil events.
func TestNearEqualNil1(t *testing.T) {
	// Define a sample event for testing purposes
	ev1 := evNfp
	// Define a nil event for testing purposes
	var ev2 *tstrading.Event = nil
	// Define another nil event for testing purposes
	var ev3 *tstrading.Event = nil
	// Check if the events are nearly equal
	if ev1.NearEqual(ev2) {
		// If the events are nearly equal, fail the test with an error message
		t.Fatal(tserr.NotEqual(&tserr.NotEqualArgs{X: ev1.String(), Y: "nil"}))
	}
	// Check if the events are nearly equal
	if ev2.NearEqual(ev1) {
		// If the events are nearly equal, fail the test with an error message
		t.Fatal(tserr.NotEqual(&tserr.NotEqualArgs{X: "nil", Y: ev1.String()}))
	}
	// Check if the events are nearly equal
	if ev2.NearEqual(ev3) {
		// If the events are nearly equal, fail the test with an error message
		t.Fatal(tserr.NotEqual(&tserr.NotEqualArgs{X: "nil", Y: "nil"}))
	}
}

// TestGenerateDocID tests the GenerateDocID method of the Event struct by comparing the document IDs
// of the sample events to a golden file.
func TestGenerateDocID(t *testing.T) {
	// Create a strings.Builder to build the document IDs for the sample events
	b := strings.Builder{}
	// Iterate over each event in the sample events slice and append its document ID to the builder
	for _, ev := range evs {
		docID, err := ev.GenerateDocID()
		if err != nil {
			t.Fatal(tserr.Op(&tserr.OpArgs{Op: "GenerateDocID", Fn: ev.Name, Err: err}))
		}
		b.WriteString(docID)
		b.WriteRune(',')
	}
	// Compare the document IDs to a golden file using the EvalGoldenFile function from the tsfio package,
	// and if there is an error, fail the test with the error message
	if e := tsfio.EvalGoldenFile(&tsfio.Testcase{Name: "docid", Data: b.String()}); e != nil {
		t.Fatal(e)
	}
}

// TestGenerateDocIDNil tests the GenerateDocID method of the Event struct with a nil event.
func TestGenerateDocIDNil(t *testing.T) {
	// Create a nil event
	var ev *tstrading.Event = nil
	// Get the document ID of the nil event
	docID, err := ev.GenerateDocID()
	// The expected document ID for a nil event should be an empty string
	expected := ""
	// Check if the document ID of a nil event is an empty string
	if docID != expected {
		// If the document ID is not an empty string, fail the test with an error message
		t.Fatal(tserr.EqualStr(&tserr.EqualStrArgs{Var: "DocID", Actual: docID, Want: expected}))
	}
	// Check if the error is nil
	if err == nil {
		// If the error is not nil, fail the test with an error message
		t.Fatal(tserr.NilFailed("GenerateDocID"))
	}
}

func TestNewPeriodForDate(t *testing.T) {
	// Create a new period for the current date
	p := tstrading.NewPeriodForDate(time.Now())
	// Check if the period has the correct start and end times
	if p.From.Year() != time.Now().Year() ||
		p.From.Month() != time.Now().Month() ||
		p.From.Day() != time.Now().Day() ||
		p.From.Hour() != 0 ||
		p.From.Minute() != 0 ||
		p.From.Second() != 0 ||
		p.From.Nanosecond() != 0 ||
		p.To.Year() != time.Now().Year() ||
		p.To.Month() != time.Now().Month() ||
		p.To.Day() != time.Now().Day() ||
		p.To.Hour() != 23 ||
		p.To.Minute() != 59 ||
		p.To.Second() != 59 ||
		p.To.Nanosecond() != 999999999 {
		t.Fatal(tserr.Equal(&tserr.EqualArgs{X: "TODO", Y: "Period{From:2026-07-05 00:00:00 +0000 UTC, To:2026-07-06 23:59:59.999999999 +0000 UTC}"}))
	}
}
