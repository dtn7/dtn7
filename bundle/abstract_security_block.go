package bundle

// IDValueTuble implements the Tuple used in
type IDValueTuple struct {
	ID    uint64
	value []uint64 //TODO Pointer to Value Type?
}

// TargetSecResultsSet implements the security results array described in BPSEC 3.6.
type TargetSecResultsSet struct {
	resultsSet []IDValueTuple
}

// AbstractSecBlock implements the Abstract Security Block (ABS) data structure described in BPSEC 3.6.
type AbstractSecBlock struct {
	secTargets           []*block
	secContextID         uint64 //TODO Pointer to Security Context Type?
	secContextFlags      byte
	secSource            EndpointID
	secContextParameters []IDValueTuple
	secResults           []TargetSecResultsSet
}
