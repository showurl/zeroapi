## proto file

### common.proto

```protobuf
syntax = "proto3";

package pb;
option go_package = "./pb";

message CommonReq {
  string id = 1;
}

```

### hello.proto

```protobuf
syntax = "proto3";

package pb;
option go_package = "./pb";
import "common.proto";

message StreamReq {
  string name = 1;
  CommonReq common = 2;
}

message StreamResp {
  string greet = 1;
}

service StreamGreeter {
  rpc greet(StreamReq) returns (StreamResp);
}
```

### ProtoSet file
```go
// protoset.go
package pb

import _ "embed"

//go:embed common.pb
var ProtoSetCommon []byte

//go:embed hello.pb
var ProtoSetHello []byte
```

