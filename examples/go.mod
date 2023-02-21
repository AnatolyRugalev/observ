module github.com/AnatolyRugalev/observ/examples

go 1.19

replace (
	github.com/AnatolyRugalev/observ v0.0.1 => ./..
	github.com/AnatolyRugalev/observ/collectors/logrust v0.0.1 => ./../collectors/logrust
	github.com/AnatolyRugalev/observ/collectors/otelt v0.0.1 => ./../collectors/otelt
	github.com/AnatolyRugalev/observ/collectors/prometheust v0.0.1 => ./../collectors/prometheust
)

require (
	github.com/AnatolyRugalev/observ v0.0.1
	github.com/AnatolyRugalev/observ/collectors/logrust v0.0.1
	github.com/AnatolyRugalev/observ/collectors/otelt v0.0.1
	github.com/AnatolyRugalev/observ/collectors/prometheust v0.0.1
	github.com/prometheus/client_golang v1.14.0
	github.com/samber/lo v1.37.0
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.1
	go.opentelemetry.io/otel v1.13.0
	go.opentelemetry.io/otel/metric v0.36.0
	go.opentelemetry.io/otel/sdk/metric v0.36.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	go.opentelemetry.io/otel/sdk v1.13.0 // indirect
	go.opentelemetry.io/otel/trace v1.13.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/sys v0.0.0-20220919091848-fb04ddd9f9c8 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
