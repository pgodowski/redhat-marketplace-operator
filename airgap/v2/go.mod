module github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/canonical/go-dqlite v1.8.0
	github.com/go-co-op/gocron v1.5.0
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.16.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.9.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.0
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83
	google.golang.org/grpc v1.41.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.5
	sigs.k8s.io/controller-runtime v0.10.2
)

replace (
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.4.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
)
