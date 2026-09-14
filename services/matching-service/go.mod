module jobradar/matching-service

go 1.27.0

require (
	github.com/segmentio/kafka-go v0.4.51
	github.com/stretchr/testify v1.8.0
	google.golang.org/grpc v1.83.2
	jobradar/events v0.0.0
	jobradar/proto v0.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace jobradar/proto => ../../proto

replace jobradar/events => ../../events
