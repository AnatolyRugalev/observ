//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.10 -recursive=false -src=./logq -build-tag=level1 -type=Filter,Group,Records
//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.10 -recursive=false -src=./logt -build-tag=level1 -type=Records,Group
//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.10 -recursive=false -src=./logt -build-tag=level2 -type=Require,RequireGroup

//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.10 -recursive=false -src=./metrq -build-tag=metrq   -type=Filter,Group
//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.10 -recursive=false -src=./metrt -build-tag=metrt   -type=Metrics,Group
//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.10 -recursive=false -src=./metrt -build-tag=require -type=Require,RequireGroup
package observ
