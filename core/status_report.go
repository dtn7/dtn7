package core

import (
	"fmt"
	"io"
	"strings"

	"github.com/dtn7/cboring"
	"github.com/dtn7/dtn7-go/bundle"
)

// BundleStatusItem represents the a bundle status item, as used as an element
// in the bundle status information array of each Bundle Status Report.
type BundleStatusItem struct {
	Asserted        bool
	Time            bundle.DtnTime
	StatusRequested bool
}

func (bsi *BundleStatusItem) MarshalCbor(w io.Writer) error {
	var arrLen uint64 = 1
	if bsi.Asserted && bsi.StatusRequested {
		arrLen = 2
	}

	if err := cboring.WriteArrayLength(arrLen, w); err != nil {
		return err
	}

	if err := cboring.WriteBoolean(bsi.Asserted, w); err != nil {
		return err
	}

	if arrLen == 2 {
		if err := cboring.WriteUInt(uint64(bsi.Time), w); err != nil {
			return err
		}
	}

	return nil
}

func (bsi *BundleStatusItem) UnmarshalCbor(r io.Reader) error {
	var arrLen uint64
	if n, err := cboring.ReadArrayLength(r); err != nil {
		return err
	} else if n != 1 && n != 2 {
		return fmt.Errorf("BundleStatusItem: Array's length is %d, not 1 or 2", n)
	} else {
		arrLen = n
	}

	if b, err := cboring.ReadBoolean(r); err != nil {
		return err
	} else {
		bsi.Asserted = b
	}

	if arrLen == 2 {
		if n, err := cboring.ReadUInt(r); err != nil {
			return err
		} else {
			bsi.Time = bundle.DtnTime(n)
		}

		bsi.StatusRequested = true
	} else {
		bsi.StatusRequested = false
	}

	return nil
}

func (bsi BundleStatusItem) String() string {
	if !bsi.Asserted {
		return fmt.Sprintf("BundleStatusItem(%t)", bsi.Asserted)
	} else {
		return fmt.Sprintf("BundleStatusItem(%t, %v)", bsi.Asserted, bsi.Time)
	}
}

// NewBundleStatusItem returns a new BundleStatusItem, indicating an optional
// assertion - givenas asserted -, but no status time request.
func NewBundleStatusItem(asserted bool) BundleStatusItem {
	return BundleStatusItem{
		Asserted:        asserted,
		Time:            bundle.DtnTimeEpoch,
		StatusRequested: false,
	}
}

// NewTimeReportingBundleStatusItem returns a new BundleStatusItem, indicating
// both a positive assertion and a requested status time report.
func NewTimeReportingBundleStatusItem(time bundle.DtnTime) BundleStatusItem {
	return BundleStatusItem{
		Asserted:        true,
		Time:            time,
		StatusRequested: true,
	}
}

// StatusReportReason is the bundle status report reason code, which is used as
// the second element of the bundle status report array.
type StatusReportReason uint64

const (
	// NoInformation is the "No additional information" bundle status report
	// reason code.
	NoInformation StatusReportReason = 0

	// LifetimeExpired is the "Lifetime expired" bundle status report reason code.
	LifetimeExpired StatusReportReason = 1

	// ForwardUnidirectionalLink is the "Forwarded over unidirectional link"
	// bundle status report reason code.
	ForwardUnidirectionalLink StatusReportReason = 2

	// TransmissionCanceled is the "Transmission canceled" bundle status report
	// reason code.
	TransmissionCanceled StatusReportReason = 3

	// DepletedStorage is the "Depleted storage" bundle status report reason code.
	DepletedStorage StatusReportReason = 4

	// DestEndpointUnintelligible is the "Destination endpoint ID unintelligible"
	// bundle status report reason code.
	DestEndpointUnintelligible StatusReportReason = 5

	// NoRouteToDestination is the "No known route to destination from here"
	// bundle status report reason code.
	NoRouteToDestination StatusReportReason = 6

	// NoNextNodeContact is the "No timely contact with next node on route" bundle
	// status report reason code.
	NoNextNodeContact StatusReportReason = 7

	// BlockUnintelligible is the "Block unintelligible" bundle status report
	// reason code.
	BlockUnintelligible StatusReportReason = 8

	// HopLimitExceeded is the "Hop limit exceeded" bundle status report reason
	// code.
	HopLimitExceeded StatusReportReason = 9
)

func (srr StatusReportReason) String() string {
	switch srr {
	case NoInformation:
		return "No additional information"

	case LifetimeExpired:
		return "Lifetime expired"

	case ForwardUnidirectionalLink:
		return "Forward over unidirectional link"

	case TransmissionCanceled:
		return "Transmission canceled"

	case DepletedStorage:
		return "Depleted storage"

	case DestEndpointUnintelligible:
		return "Destination endpoint ID unintelligible"

	case NoRouteToDestination:
		return "No known route to destination from here"

	case NoNextNodeContact:
		return "No timely contact with next node on route"

	case BlockUnintelligible:
		return "Block unintelligible"

	case HopLimitExceeded:
		return "Hop limit exceeded"

	default:
		return "unknown"
	}
}

// StatusInformationPos describes the different bundle status information
// entries. Each bundle status report must contain at least the following
// bundle status items.
type StatusInformationPos int

const (
	// maxStatusInformationPos is the amount of different StatusInformationPos.
	maxStatusInformationPos int = 4

	// ReceivedBundle is the first bundle status information entry, indicating
	// the reporting node received this bundle.
	ReceivedBundle StatusInformationPos = 0

	// ForwardedBundle is the second bundle status information entry, indicating
	// the reporting node forwarded this bundle.
	ForwardedBundle StatusInformationPos = 1

	// DeliveredBundle is the third bundle status information entry, indicating
	// the reporting node delivered this bundle.
	DeliveredBundle StatusInformationPos = 2

	// DeletedBundle is the fourth bundle status information entry, indicating
	// the reporting node deleted this bundle.
	DeletedBundle StatusInformationPos = 3
)

func (sip StatusInformationPos) String() string {
	switch sip {
	case ReceivedBundle:
		return "received bundle"

	case ForwardedBundle:
		return "forwarded bundle"

	case DeliveredBundle:
		return "delivered bundle"

	case DeletedBundle:
		return "deleted bundle"

	default:
		return "unknown"
	}
}

// StatusReport is the bundle status report, used in an administrative record.
type StatusReport struct {
	StatusInformation []BundleStatusItem
	ReportReason      StatusReportReason
	SourceNode        bundle.EndpointID
	Timestamp         bundle.CreationTimestamp
}

// NewStatusReport creates a bundle status report for the given bundle and
// StatusInformationPos, which creates the right bundle status item. The
// bundle status report reason code will be used and the bundle status item
// gets the given timestamp.
func NewStatusReport(bndl bundle.Bundle, statusItem StatusInformationPos,
	reason StatusReportReason, time bundle.DtnTime) StatusReport {
	var sr = StatusReport{
		StatusInformation: make([]BundleStatusItem, maxStatusInformationPos),
		ReportReason:      reason,
		SourceNode:        bndl.PrimaryBlock.SourceNode,
		Timestamp:         bndl.PrimaryBlock.CreationTimestamp,
	}

	for i := 0; i < maxStatusInformationPos; i++ {
		sip := StatusInformationPos(i)

		switch {
		case sip == statusItem && bndl.PrimaryBlock.BundleControlFlags.Has(bundle.RequestStatusTime):
			sr.StatusInformation[i] = NewTimeReportingBundleStatusItem(time)

		case sip == statusItem:
			sr.StatusInformation[i] = NewBundleStatusItem(true)

		default:
			sr.StatusInformation[i] = NewBundleStatusItem(false)
		}
	}

	return sr
}

// StatusInformations returns an array of available StatusInformationPos.
func (sr StatusReport) StatusInformations() (sips []StatusInformationPos) {
	for i := 0; i < len(sr.StatusInformation); i++ {
		si := sr.StatusInformation[i]
		sip := StatusInformationPos(i)

		if si.Asserted {
			sips = append(sips, sip)
		}
	}

	return
}

func (sr *StatusReport) MarshalCbor(w io.Writer) error {
	// TODO: support fragmentation
	if err := cboring.WriteArrayLength(4, w); err != nil {
		return err
	}

	if err := cboring.WriteArrayLength(uint64(len(sr.StatusInformation)), w); err != nil {
		return err
	}
	for _, si := range sr.StatusInformation {
		if err := cboring.Marshal(&si, w); err != nil {
			return fmt.Errorf("Marshalling BundleStatusItem failed: %v", err)
		}
	}

	if err := cboring.WriteUInt(uint64(sr.ReportReason), w); err != nil {
		return err
	}

	if err := cboring.Marshal(&sr.SourceNode, w); err != nil {
		return fmt.Errorf("Marshalling EndpointID failed: %v", err)
	}

	if err := cboring.Marshal(&sr.Timestamp, w); err != nil {
		return fmt.Errorf("Marshalling CreationTimestamp failed: %v", err)
	}

	return nil
}

func (sr *StatusReport) UnmarshalCbor(r io.Reader) error {
	if n, err := cboring.ReadArrayLength(r); err != nil {
		return err
	} else if n != 4 {
		return fmt.Errorf("Expected array length 4, got %d", n)
	}

	if n, err := cboring.ReadArrayLength(r); err != nil {
		return err
	} else {
		sr.StatusInformation = make([]BundleStatusItem, int(n))
	}
	for i := 0; i < len(sr.StatusInformation); i++ {
		if err := cboring.Unmarshal(&sr.StatusInformation[i], r); err != nil {
			return fmt.Errorf("Unmarshalling BundleStatusItem failed: %v", err)
		}
	}

	if n, err := cboring.ReadUInt(r); err != nil {
		return err
	} else {
		sr.ReportReason = StatusReportReason(n)
	}

	if err := cboring.Unmarshal(&sr.SourceNode, r); err != nil {
		return fmt.Errorf("Unmarshalling EndpointID failed: %v", err)
	}

	if err := cboring.Unmarshal(&sr.Timestamp, r); err != nil {
		return fmt.Errorf("Unmarshalling CreationTimestamp failed: %v", err)
	}

	return nil
}

func (sr StatusReport) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "StatusReport([")

	for i := 0; i < len(sr.StatusInformation); i++ {
		si := sr.StatusInformation[i]
		sip := StatusInformationPos(i)

		if !si.Asserted {
			continue
		}

		if si.Time == bundle.DtnTimeEpoch {
			fmt.Fprintf(&b, "%v,", sip)
		} else {
			fmt.Fprintf(&b, "%v %v,", sip, si.Time)
		}
	}
	fmt.Fprintf(&b, "], ")

	fmt.Fprintf(&b, "%v, %v, %v", sr.ReportReason, sr.SourceNode, sr.Timestamp)

	return b.String()
}
