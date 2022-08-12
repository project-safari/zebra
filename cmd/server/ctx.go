package main

// set constant values for ctx.
type CtxKey string

const (
	ResourcesCtxKey = CtxKey("resources")
	AuthCtxKey      = CtxKey("authKey")
	ClaimsCtxKey    = CtxKey("claims")
)
