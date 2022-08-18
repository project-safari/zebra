package main

// Set constant values for ctx.
type CtxKey string

// CtxKey constants for resources, key, claims.
const (
	ResourcesCtxKey = CtxKey("resources")
	AuthCtxKey      = CtxKey("authKey")
	ClaimsCtxKey    = CtxKey("claims")
)
