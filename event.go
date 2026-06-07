// Copyright (c) 2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tstrading

// Import standard library packages and lpstats for utility functions.
import (
	"crypto/sha256" // sha256 for hashing event details to generate a unique document ID for NoSQL databases
	"encoding/hex"  // hex for encoding the hash output as a hexadecimal string
	"fmt"           // fmt for string formatting
	"math"
	"strings"
	"time" // time for handling event timestamps and periods

	"github.com/thorsphere/lpstats" // lpstats for utility functions to compare float pointers and format them as strings
	"github.com/thorsphere/tserr"
	"github.com/thorsphere/tstable" // tstable for formatting tables
)

// Event represents a single calendar event with its details.
// Thought: Implement a definition for country codes, currency codes, and units of measurement and use it here for better type safety and validation.
type Event struct {
	ID          int64       `json:"id"`           // Unique identifier for the event, e.g., a database primary key or a UUID. It is expected to be set by the database.
	Name        string      `json:"name"`         // Name of the economic event, e.g., "Non-Farm Payrolls", "GDP Growth Rate"
	Time        time.Time   `json:"time"`         // Date and time of the event in UTC, when the data is released or expected to be released
	Country     string      `json:"country"`      // ISO 3166-1 alpha-2 two-letter country code
	Currency    *string     `json:"currency"`     // Currency of the values, e.g., "USD", "EUR", "JPY". Pointer because it can be nil if the value is not yet released
	Actual      *float64    `json:"actual"`       // Pointer, because it can be nil if the value is not yet released
	Estimate    *float64    `json:"estimate"`     // Pointer, because it can be nil if the value is not yet released
	Previous    *float64    `json:"previous"`     // Pointer, because it can be nil if the value is not yet released
	Unit        *string     `json:"unit"`         // Unit of measurement for the values, e.g., "%", "K", "M", "B". Pointer because it can be nil if the value is not yet released
	Precision   int         `json:"precision"`    // Number of decimal places to round the values to, e.g., 2 for USD, 0 for JPY
	Change      *float64    `json:"change"`       // Pointer, because it can be nil if the value is not yet released
	ChangePct   *float64    `json:"change_pct"`   // Pointer, because it can be nil if the value is not yet released
	Surprise    *float64    `json:"surprise"`     // Pointer, because it can be nil if the value is not yet released
	SurprisePct *float64    `json:"surprise_pct"` // Pointer, because it can be nil if the value is not yet released
	Impact      ImpactLevel `json:"impact"`       // Impact level of the event
	Source      string      `json:"source"`       // Source of the data, e.g., "Official Government Website"
}

// Period represents a time period with a start and end date.
type Period struct {
	From time.Time // Start date and time of the period
	To   time.Time // End date and time of the period
}

// NewPeriodForDate creates a Period that spans the entire given date (00:00:00 to 23:59:59) in UTC.
func NewPeriodForDate(date time.Time) *Period {
	// Normalize the input date to the start of the day in UTC
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	// The end of the day is calculated by adding one day to the start of the day and
	// then subtracting one nanosecond to get the last moment of the specified date.
	// Calculate the end of the day by adding 24 hours and subtracting a nanosecond to get 23:59:59.999999999
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)
	// Return a new Period struct with the calculated start and end times for the specified date
	return &Period{From: startOfDay, To: endOfDay}
}

// fmtVal formats a float pointer value with its unit and currency.
// If the value is nil, it returns an empty string.
// If the value is not nil, it formats it with the currency (if available), the value, and the unit (if available).
func (ev *Event) fmtVal(value *float64) string {
	// If the value is nil, return an empty string
	if value == nil {
		return ""
	}
	// Create a new string builder
	var b strings.Builder
	// If the currency is not nil, append it to the string builder
	if ev.Currency != nil && *ev.Currency != "" {
		// Append the currency to the string builder
		b.WriteString(*ev.Currency)
		// Append a space to the string builder
		b.WriteByte(' ')
	}
	// Append the value to the string builder
	b.WriteString(lpstats.FmtFloatPtr(value, ev.Precision))
	// If the unit is not nil, append it to the string builder
	if ev.Unit != nil && *ev.Unit != "" {
		// Append the unit to the string builder
		b.WriteString(*ev.Unit)
	}
	// Return the formatted string
	return b.String()
}

// fmtPct formats a float pointer value as a percentage.
// If the value is nil, it returns an empty string.
// If the value is not nil, it formats it as a percentage with a percent sign.
func (ev *Event) fmtPct(value *float64) string {
	// If the value is nil, return an empty string
	if value == nil {
		return ""
	}
	// Create a new string builder
	var b strings.Builder
	// Append the value to the string builder
	b.WriteString(lpstats.FmtFloatPtr(value, 1))
	// Append a percent sign to the string builder
	b.WriteByte('%')
	// Return the formatted string
	return b.String()
}

// String returns a formatted string representation of the Event.
func (ev *Event) String() string {
	if ev == nil {
		return "<nil>"
	}
	// Create a new table
	tbl, err := tstable.New([]string{"Event", ev.Name})
	// If there is an error, return an empty string
	if err != nil {
		return "<error>"
	}
	// Add the event details to the table
	tbl.AddRow([]string{"Country", ev.Country})
	tbl.AddRow([]string{"Time", ev.Time.Format(time.RFC3339)})
	tbl.AddRow([]string{"Impact", ev.Impact.String()})
	addRowIfVal(tbl, "Actual", ev.Actual, ev.fmtVal)
	addRowIfVal(tbl, "Estimate", ev.Estimate, ev.fmtVal)
	addRowIfVal(tbl, "Previous", ev.Previous, ev.fmtVal)
	addRowIfVal(tbl, "Surprise", ev.Surprise, ev.fmtVal)
	addRowIfVal(tbl, "Surprise in %", ev.SurprisePct, ev.fmtPct)
	addRowIfVal(tbl, "Change", ev.Change, ev.fmtVal)
	addRowIfVal(tbl, "Change in %", ev.ChangePct, ev.fmtPct)
	tbl.AddRow([]string{"Source", ev.Source})
	// Return the formatted table as a string
	return tbl.String()
}

// addRowIfVal adds a row to the table with the label and formatted value if the value is not nil.
func addRowIfVal(tbl *tstable.Table, label string, val *float64, f func(*float64) string) {
	// If the value is not nil, add a row to the table with the label and formatted value
	if val != nil {
		// Add a row to the table with the label and formatted value
		tbl.AddRow([]string{label, f(val)})
	}
}

// NearEqual compares two Event instances for near-equality, taking into account all fields including
// the pointer fields for Actual, Estimate, and Previous. It does not compare the ID field,
// as it is expected to be set by the database and may not be the same for two events that are otherwise identical.
// It returns true if all fields are equal, and false otherwise.
func (ev *Event) NearEqual(other *Event) bool {
	// Check if one of the events is nil
	if ev == nil || other == nil {
		// If one of the events is nil, return false
		return false
	}
	// Tolerance: half of the smallest displayable unit.
	tolerance := 0.5 * math.Pow10(-ev.Precision)
	// Check if all fields are equal
	return ev.Name == other.Name &&
		ev.Time.Equal(other.Time) &&
		ev.Country == other.Country &&
		lpstats.NearEqualFloatPtr(ev.Actual, other.Actual, tolerance) &&
		lpstats.NearEqualFloatPtr(ev.Estimate, other.Estimate, tolerance) &&
		lpstats.NearEqualFloatPtr(ev.Previous, other.Previous, tolerance) &&
		lpstats.NearEqualFloatPtr(ev.Surprise, other.Surprise, tolerance) &&
		lpstats.NearEqualFloatPtr(ev.SurprisePct, other.SurprisePct, tolerance) &&
		lpstats.NearEqualFloatPtr(ev.Change, other.Change, tolerance) &&
		lpstats.NearEqualFloatPtr(ev.ChangePct, other.ChangePct, tolerance) &&
		lpstats.EqualStrPtr(ev.Unit, other.Unit) &&
		lpstats.EqualStrPtr(ev.Currency, other.Currency) &&
		ev.Precision == other.Precision &&
		ev.Impact == other.Impact &&
		ev.Source == other.Source
}

// GenerateDocID creates a deterministic, safe, unique document ID for NoSQL databases (like Firestore)
// based on the event's Time, Country, and Name.
func (ev *Event) GenerateDocID() (string, error) {
	if ev == nil {
		return "", tserr.NilPtr()
	}
	// Concatenate the fields that make an event unique
	key := fmt.Sprintf("%s|%s|%s", ev.Time.UTC().Format(time.RFC3339), ev.Country, ev.Name)
	// Hash the string to ensure it is URL-safe and of predictable length
	hash := sha256.Sum256([]byte(key))
	// Encode the hash as a hexadecimal string to use as the document ID and nil to indicate success
	return hex.EncodeToString(hash[:]), nil
}
