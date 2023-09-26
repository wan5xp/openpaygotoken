module github.com/wan5xp/openpaygotoken

go 1.20

require github.com/wan5xp/openpaygotoken/pkg/openpaygotoken v0.0.0-20190108105601-1b9a9b2b2f2f

require github.com/wan5xp/openpaygotoken/pkg/simulators v0.0.0-20190108105601-1b9a9b2b2f2f

require (
	github.com/aead/siphash v1.0.1 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9
)

replace github.com/wan5xp/openpaygotoken/pkg/openpaygotoken => ./pkg/openpaygotoken

replace github.com/wan5xp/openpaygotoken/pkg/simulators => ./pkg/simulators
