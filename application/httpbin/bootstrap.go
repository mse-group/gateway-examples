package main

type Args struct {
	ServerOptions ServerOptions
	Version string
}

type ServerOptions struct {
	Addr string
	TLS  TLSOptions
}

type TLSOptions struct {
	Enable bool
	VerifyClient bool
}
