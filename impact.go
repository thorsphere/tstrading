// Copyright (c) 2026 thorsphere.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tstrading

// ImpactLevel represents the expected market impact of an economic event.
type ImpactLevel int

// Define constants for the different impact levels.
const (
	ImpactUnknown ImpactLevel = iota // iota starts at 0, so ImpactUnknown = 0
	ImpactLow                        // ImpactLow = 1
	ImpactMedium                     // ImpactMedium = 2
	ImpactHigh                       // ImpactHigh = 3
)

// String returns a string representation of the ImpactLevel.
func (i ImpactLevel) String() string {
	switch i {
	case ImpactLow:
		return "low" // Return "low" for ImpactLow
	case ImpactMedium:
		return "medium" // Return "medium" for ImpactMedium
	case ImpactHigh:
		return "high" // Return "high" for ImpactHigh
	default:
		return "unknown" // Return "unknown" for any undefined impact levels
	}
}
