module chaitin.cn/patronus/safeline-2/management/tcontrollerd

go 1.21

toolchain go1.21.3

require (
	chaitin.cn/dev/go/errors v0.0.0-20210324055134-dc5247602af6
	chaitin.cn/dev/go/log v0.0.0-20221220104336-05125760b10c
	chaitin.cn/dev/go/settings v0.0.0-20221220104336-05125760b10c
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.1
)

require (
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	chaitin.cn/dev/go/errors => ../internal/thirdparty/errors
	chaitin.cn/dev/go/log => ../internal/thirdparty/log
	chaitin.cn/dev/go/settings => ../internal/thirdparty/settings
)
