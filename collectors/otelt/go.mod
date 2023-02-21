module github.com/AnatolyRugalev/observ/collectors/logrust

go 1.19

replace github.com/AnatolyRugalev/observ v0.0.1 => ./../..

require (
	github.com/AnatolyRugalev/observ v0.0.1
	go.opentelemetry.io/otel v1.13.0
	go.opentelemetry.io/otel/sdk/metric v0.36.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/samber/lo v1.37.0 // indirect
	go.opentelemetry.io/otel/metric v0.36.0 // indirect
	go.opentelemetry.io/otel/sdk v1.13.0 // indirect
	go.opentelemetry.io/otel/trace v1.13.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/sys v0.0.0-20220919091848-fb04ddd9f9c8 // indirect
)
