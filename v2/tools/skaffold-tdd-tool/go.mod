module github.com/redhat-marketplace/redhat-marketplace-operator/v2/tools/skaffold-tdd-tool

go 1.19

require (
	emperror.dev/errors v0.8.0
	github.com/GoogleContainerTools/skaffold v0.0.0-00010101000000-000000000000
	// github.com/GoogleContainerTools/skaffold v1.16.0
	github.com/caarlos0/env/v6 v6.4.0
	github.com/gdamore/tcell/v2 v2.0.1-0.20201017141208-acf90d56d591
	github.com/redhat-marketplace/redhat-marketplace-operator/v2 v2.0.0-20210729050326-8246afe36a7e // indirect
	github.com/rivo/tview v0.0.0-20201118063654-f007e9ad3893
	google.golang.org/grpc v1.52.0
	google.golang.org/protobuf v1.28.1
)

require github.com/redhat-marketplace/redhat-marketplace-operator/tests/v2 v2.0.0-00010101000000-000000000000

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/OneOfOne/xxhash v1.2.6 // indirect
	github.com/banzaicloud/k8s-objectmatcher v1.8.0 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/goph/emperror v0.17.2 // indirect
	github.com/gotidy/ptr v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/onsi/ginkgo/v2 v2.8.0 // indirect
	github.com/onsi/gomega v1.26.0 // indirect
	github.com/operator-framework/api v0.3.25 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.60.1 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.60.1 // indirect
	github.com/prometheus/client_golang v1.13.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/oauth2 v0.0.0-20221014153046-6fdb5e3db783 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/term v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221227171554-f9683d7f8bef // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.25.2 // indirect
	k8s.io/apiextensions-apiserver v0.25.0 // indirect
	k8s.io/apimachinery v0.25.2 // indirect
	k8s.io/client-go v12.0.0+incompatible // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220803164354-a70c9af30aea // indirect
	k8s.io/utils v0.0.0-20220922133306-665eaaec4324 // indirect
	sigs.k8s.io/controller-runtime v0.13.0 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible // Required by OLM
	github.com/GoogleContainerTools/skaffold => github.com/GoogleContainerTools/skaffold v1.16.0
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.6.0
	github.com/banzaicloud/k8s-objectmatcher => github.com/banzaicloud/k8s-objectmatcher v1.6.1
	github.com/containerd/containerd v1.4.0-0 => github.com/containerd/containerd v1.4.0
	github.com/coreos/prometheus-operator => github.com/prometheus-operator/prometheus-operator v0.41.0
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.4.0
	github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20190319215453-e7b5f7dbe98c+incompatible
	github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c => github.com/docker/docker v17.12.0-ce-rc1.0.20190319215453-e7b5f7dbe98c+incompatible
	github.com/operator-framework/operator-marketplace => github.com/operator-framework/operator-marketplace v0.0.0-20201110032404-0e3bd3db36a6
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.8.2-0.20200609102542-5d7e3e970602
	github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2 => ../../../airgap/v2
	github.com/redhat-marketplace/redhat-marketplace-operator/reporter/v2 => ../../../reporter/v2
	github.com/redhat-marketplace/redhat-marketplace-operator/tests/v2 => ../../../tests/v2
	github.com/redhat-marketplace/redhat-marketplace-operator/v2 => ../../
	k8s.io/api => k8s.io/api v0.24.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.24.7
	k8s.io/client-go => k8s.io/client-go v0.24.7
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.12.3
)
