package main

import "errors"

var (
	// ErrSolutionFail error cannot compute a solution
	ErrSolutionFail = errors.New("exceeded 2^20 iterations failed to find solution")

	// ErrResourceEmpty error empty hashcash resource
	ErrResourceEmpty = errors.New("empty hashcash resource")
)
