module github.com/redhat-marketplace/redhat-marketplace-operator/v2/tools/connect

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/chromedp/cdproto v0.0.0-20210429002609-5ec2b0624aec
	github.com/chromedp/chromedp v0.7.1
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.16.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/thediveo/enumflag v0.10.1
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a // indirect
	golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.4.0
