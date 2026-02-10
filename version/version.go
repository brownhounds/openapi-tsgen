package version

import _ "embed"

//go:generate cp ../VERSION VERSION
//go:embed VERSION
var Ver string
