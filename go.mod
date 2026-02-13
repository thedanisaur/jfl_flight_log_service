module flight_log_service

go 1.24.0

toolchain go1.24.4

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/gofiber/fiber/v2 v2.42.0
	github.com/google/uuid v1.3.0
	github.com/thedanisaur/jfl_platform v0.0.0
)

require (
	cel.dev/expr v0.25.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/cel-go v0.27.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20220530130905-52f3993e8d6d // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.44.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20240823005443-9b4947da3948 // indirect
	golang.org/x/sys v0.31.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260128011058-8636f8732409 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260122232226-8e98ce8d340d // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/thedanisaur/jfl_platform => ../jfl_platform
