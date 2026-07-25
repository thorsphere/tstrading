# tstrading

[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/thorsphere/tstrading)](https://pkg.go.dev/mod/github.com/thorsphere/tstrading)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thorsphere/tstrading)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/thorsphere/tstrading)
![GitHub Top Language](https://img.shields.io/github/languages/top/thorsphere/tstrading)
[![CodeFactor](https://www.codefactor.io/repository/github/thorsphere/tstrading/badge)](https://www.codefactor.io/repository/github/thorsphere/tstrading)
![OSS Lifecycle](https://img.shields.io/osslifecycle/thorsphere/tstrading)

A Go package providing domain types and data models for trading. It focuses on economic calendar events, their impact levels, and time periods, with support for formatting, comparison, and NoSQL document‑ID generation.

## Types

### `Event`

Represents a single economic calendar event.

```go
type Event struct {
    ID          int64       `json:"id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
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

An integer‑based enum representing the expected market impact of an economic event.

```go
type ImpactLevel int

const (
    ImpactUnknown ImpactLevel = iota // default zero value – used when impact is not set or unrecognised
    ImpactLow                        // 1
    ImpactMedium                     // 2
    ImpactHigh                       // 3
)
```

`ImpactLevel` implements the `fmt.Stringer` interface:

- `ImpactLow` → `"low"`
- `ImpactMedium` → `"medium"`
- `ImpactHigh` → `"high"`
- `ImpactUnknown` (or any invalid value) → `"unknown"`

## Dependencies

- [`tstable`](https://github.com/thorsphere/tstable) – table formatting used by `Event.String()`
- [`tserr`](https://github.com/thorsphere/tserr) – error and test helpers
- [`lpstats`](https://github.com/thorsphere/lpstats) – float pointer comparison, formatting, and string pointer utilities

## Documentation & Resources

- [Go Package Documentation](https://pkg.go.dev/github.com/thorsphere/tstrading) — Complete API reference
- [Open Source Insights](https://deps.dev/go/github.com%2Fthorsphere%2Ftstrading) — Dependency analysis

## ⚖️ License & Commercial Usage

Copyright (c) 2026 thorsphere. All rights reserved.

This project is licensed under the **Functional Source License v1.1 (FSL-1.1-ALv2)**. 

* The use, modification, and redistribution of this Go package is completely free for private, educational, non-commercial, and internal purposes. 
* If you are a company or institution looking to use this package in a commercial product, service, or business environment, you must secure a commercial license.
* Each version of this software automatically converts to the fully open-source Apache License, Version 2.0 on the second anniversary of its release.

For full details, please see the [LICENSE](LICENSE) file.

### 💼 Commercial Licensing & Inquiries

To purchase a commercial license or discuss support options, please reach out directly:

* 📩 **Contact:** business at thorsphere dot com
* 💬 **Response Time:** Usually within a couple of business days.

*Please include your company name and a brief overview of your use case so I can provide the right licensing details.*
