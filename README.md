# tstrading

[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/thorsphere/tstrading)](https://pkg.go.dev/mod/github.com/thorsphere/tstrading)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thorsphere/tstrading)

[![Go Report Card](https://goreportcard.com/badge/github.com/thorsphere/tstrading)](https://goreportcard.com/report/github.com/thorsphere/tstrading)
[![CodeFactor](https://www.codefactor.io/repository/github/thorsphere/tstrading/badge)](https://www.codefactor.io/repository/github/thorsphere/tstrading)
![OSS Lifecycle](https://img.shields.io/osslifecycle/thorsphere/tstrading)
![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/thorsphere/tstrading)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/thorsphere/tstrading)
![GitHub last commit](https://img.shields.io/github/last-commit/thorsphere/tstrading)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/thorsphere/tstrading)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/thorsphere/tstrading)
![GitHub Top Language](https://img.shields.io/github/languages/top/thorsphere/tstrading)
![GitHub](https://img.shields.io/github/license/thorsphere/tstrading)

A Go package providing domain types and data models for trading. It focuses on economic calendar events, their impact levels, and time periods, with support for formatting, comparison, and NoSQL document‑ID generation.

## Types

### `Event`

Represents a single economic calendar event.

```go
type Event struct {
    ID          int64       `json:"id"`
    Name        string      `json:"name"`
    Time        time.Time   `json:"time"`
    Country     string      `json:"country"`
    Currency    *string     `json:"currency"`
    Actual      *float64    `json:"actual"`
    Estimate    *float64    `json:"estimate"`
    Previous    *float64    `json:"previous"`
    Unit        *string     `json:"unit"`
    Precision   int         `json:"precision"`
    Change      *float64    `json:"change"`
    ChangePct   *float64    `json:"change_pct"`
    Surprise    *float64    `json:"surprise"`
    SurprisePct *float64    `json:"surprise_pct"`
    Impact      ImpactLevel `json:"impact"`
    Source      string      `json:"source"`
}
```

**Key methods:**
- `String() string` – returns a formatted table string; handles nil receiver gracefully. Rows with nil values are omitted.
- `NearEqual(other *Event) bool` – compares two events for near‑equality (excluding the `ID` field). Float comparisons use a precision‑based tolerance (half of the smallest displayable unit). Pointer fields (`Unit`, `Currency`, float pointers) are compared nil‑safely.
- `GenerateDocID() (string, error)` – produces a deterministic, URL‑safe document ID from `Time`, `Country`, and `Name` using SHA‑256.

### `Period`

A time span defined by `From` and `To`.

```go
type Period struct {
    From time.Time
    To   time.Time
}
```

**Constructor:**
- `NewPeriodForDate(date time.Time) *Period` – creates a `Period` covering the full day in UTC (00:00:00 – 23:59:59.999999999).

### `ImpactLevel`

An integer‑based enum representing the expected market impact of an event.

```go
type ImpactLevel int

const (
    ImpactLow    ImpactLevel = iota
    ImpactMedium
    ImpactHigh
)
```

`ImpactLevel` implements the `fmt.Stringer` interface, returning `"low"`, `"medium"`, `"high"`, or `"unknown"` for any invalid value.

## Dependencies

- [`tstable`](https://github.com/thorsphere/tstable) – table formatting used by `Event.String()`
- [`tserr`](https://github.com/thorsphere/tserr) – error and test helpers
- [`lpstats`](https://github.com/thorsphere/lpstats) – float pointer comparison, formatting, and string pointer utilities

## License

GNU Affero General Public License v3.0 (see the LICENSE file).
