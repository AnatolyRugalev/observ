//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.8 -recursive=false -src=./logq -build-tag=lgq     -type=Filter
//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.8 -recursive=false -src=./logt -build-tag=lqt     -type=Records,Group
//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.8 -recursive=false -src=./logt -build-tag=require -type=Require,RequireGroup

//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.8 -recursive=false -src=./metrq -build-tag=mcq     -type=Filter
//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.8 -recursive=false -src=./metrt -build-tag=mct     -type=Metrics,Group
//go:generate go run github.com/AnatolyRugalev/chaingen@v0.1.8 -recursive=false -src=./metrt -build-tag=require -type=Require,RequireGroup
package observ
