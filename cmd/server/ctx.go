package main

type CtxKey string

const (
	ResourcesCtxKey = CtxKey("resources")
	AuthCtxKey      = CtxKey("authKey")
	ClaimsCtxKey    = CtxKey("claims")
)
