module github.com/obbaa-477/pppoe-relay-vnf

go 1.18

require (
	github.com/confluentinc/confluent-kafka-go v1.9.2
	github.com/google/gopacket v1.1.19
	github.com/mdlayher/raw v0.1.0
	github.com/obbaa-477/common v1.0.0
	github.com/onsi/ginkgo/v2 v2.3.1
	github.com/onsi/gomega v1.22.1
	github.com/stretchr/testify v1.8.2
	go.mongodb.org/mongo-driver v1.11.2
	golang.org/x/exp v0.0.0-20230515195305-f3d0a9c9a5cc
	google.golang.org/grpc v1.52.3
	google.golang.org/protobuf v1.28.1
)

replace github.com/obbaa-477/common => ../common

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sync v0.2.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230131230820-1c016267d619 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// https://github.com/mdlayher/raw/pull/64
replace github.com/mdlayher/raw => github.com/pyther/raw v0.0.0-20200508193324-eb26248ef18b

// https://github.com/google/gopacket/pull/781
replace github.com/google/gopacket => github.com/pyther/gopacket v1.1.18-0.20200502044149-9afa69325031
