package pb

import _ "embed"

//go:embed common.pb
var ProtoSetCommon []byte

//go:embed hello.pb
var ProtoSetHello []byte
