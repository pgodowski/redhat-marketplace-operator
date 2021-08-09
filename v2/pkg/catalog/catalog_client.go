package catalog

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	emperror "emperror.dev/errors"
	"github.com/go-logr/logr"
	marketplacev1beta1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/config"
	prom "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/prometheus"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils"
	. "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	ListForVersionEndpoint = "list-for-version"
	GetSystemMeterdefinitionTemplatesEndpoint = "get-system-meterdefs"
)

type CatalogClientBuilder struct {
	Url      string
	Insecure bool
}

type CatalogClient struct {
	endpoint   *url.URL
	httpClient http.Client
}

func NewCatalogClientBuilder(cfg *config.OperatorConfig) *CatalogClientBuilder {
	builder := &CatalogClientBuilder{}

	if cfg.URL != "" {
		builder.Url = cfg.MeterdefinitionCatalog.FileServerValues.URL
	}

	return builder
}

func (b *CatalogClientBuilder) NewCatalogServerClient(client client.Client, deployedNamespace string, kubeInterface kubernetes.Interface, reqLogger logr.Logger) (*CatalogClient, error) {
	service, err := getCatalogServerService(deployedNamespace, client, reqLogger)
	if err != nil {
		return nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	cert, err := getCertFromConfigMap(client, deployedNamespace, reqLogger)
	if err != nil {
		return nil, err
	}

	saClient := prom.NewServiceAccountClient(deployedNamespace, kubeInterface)
	authToken, err := saClient.NewServiceAccountToken(utils.OPERATOR_SERVICE_ACCOUNT, utils.FileServerAudience, 3600, reqLogger)
	if err != nil {
		return nil, err
	}

	if service != nil && len(cert) != 0 && authToken != "" {
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}

		ok := caCertPool.AppendCertsFromPEM(cert)
		if !ok {
			err = emperror.New("failed to append cert to cert pool")
			reqLogger.Error(err, "cert pool error")
			return nil, err
		}

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}

		var transport http.RoundTripper = &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
		}

		transport = WithBearerAuth(transport, authToken)

		catalogServerClient := http.Client{
			Transport: transport,
			Timeout:   1 * time.Second,
		}

		url,err := url.Parse(b.Url)
		if err != nil {
			return nil,err
		}

		reqLogger.Info("Catalog Server client created successfully")
		return &CatalogClient{
			httpClient: catalogServerClient,
			endpoint: url,
		}, nil
	}

	return nil, &ExecResult{
		ReconcileResult: reconcile.Result{},
		Err:             emperror.New("catalog server client prerequisites not ready"),
	}
}

func(c *CatalogClient) ListMeterdefintionsFromFileServer(csvName string, version string, namespace string,reqLogger logr.Logger) ([]string, []marketplacev1beta1.MeterDefinition, *ExecResult) {
	reqLogger.Info("retrieving meterdefinitions", "csvName", csvName, "csvVersion", version)

	// url := fmt.Sprintf("https://rhm-meterdefinition-file-server.openshift-redhat-marketplace.svc.cluster.local:8200/list-for-version/%s/%s", csvName, version)
	url,err := concatPaths(c.endpoint.String(),ListForVersionEndpoint,csvName,version)
	if err != nil {
		return nil,nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	response, err := c.httpClient.Get(url.String())
	return ReturnMeterdefs(csvName,namespace,*response,err,reqLogger)
}

func (c *CatalogClient) GetSystemMeterdefs(csvName string, version string, namespace string, reqLogger logr.Logger) ([]string, []marketplacev1beta1.MeterDefinition, *ExecResult) {

	reqLogger.Info("retrieving system meterdefinitions", "csvName", csvName, "csvVersion", version)

	// url := fmt.Sprintf("https://rhm-meterdefinition-file-server.openshift-redhat-marketplace.svc.cluster.local:8200/get-system-meterdefs/%s/%s/%s", csvName, version, namespace)
	url,err := concatPaths(c.endpoint.String(),GetSystemMeterdefinitionTemplatesEndpoint,csvName,version,namespace)
	if err != nil {
		return nil,nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	response, err := c.httpClient.Get(url.String())
	return ReturnMeterdefs(csvName,namespace,*response,err,reqLogger)

}

func  ReturnMeterdefs (csvName string, namespace string,response http.Response,err error,reqLogger logr.Logger) ([]string, []marketplacev1beta1.MeterDefinition, *ExecResult){
	if err != nil {
		reqLogger.Error(err, "Error on GET to Catalog Server")
		if err == io.EOF {
			reqLogger.Error(err, "system meterdefintion not found")
			return nil, nil, &ExecResult{
				ReconcileResult: reconcile.Result{},
				Err:             emperror.New("empty response"),
			}
		}

		reqLogger.Error(err, "Error querying file server")
		return nil, nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	meterDefNames := []string{}
	mdefSlice := []marketplacev1beta1.MeterDefinition{}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		reqLogger.Error(err, "error reading body")
		return nil, nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	if response.StatusCode == 404 {
		reqLogger.Info(response.Status)
		err = emperror.New(response.Status)
		return nil, nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	reqLogger.Info("response data", "data", string(data))

	meterDefsData := strings.Replace(string(data), "<<NAMESPACE-PLACEHOLDER>>", namespace, -1)
	err = yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(meterDefsData)), 100).Decode(&mdefSlice)
	if err != nil {
		reqLogger.Error(err, "error decoding response from fetchGlobalMeterdefinitions()")
		return nil, nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	for _, meterDefItem := range mdefSlice {
		meterDefNames = append(meterDefNames, meterDefItem.ObjectMeta.Name)
	}

	reqLogger.Info("meterdefintions returned from file server", csvName, meterDefNames)

	return meterDefNames, mdefSlice, &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}

func concatPaths(basePath string, paths ...string) (*url.URL, error) {

    u, err := url.Parse(basePath)

    if err != nil {
        return nil, err
    }

    temp := append([]string{u.Path}, paths...)

    concatenatedPaths := path.Join(temp...)

    u.Path = concatenatedPaths

    return u, nil
}


func getCatalogServerService(deployedNamespace string, client client.Client, reqLogger logr.InfoLogger) (*corev1.Service, error) {
	service := &corev1.Service{}

	err := client.Get(context.TODO(), types.NamespacedName{Namespace: deployedNamespace, Name: utils.CATALOG_SERVER_SERVICE_NAME}, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func getCertFromConfigMap(client client.Client, deployedNamespace string, reqLogger logr.Logger) ([]byte, error) {
	cm := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Namespace: deployedNamespace, Name: "serving-certs-ca-bundle"}, cm)
	if err != nil {
		return nil, err
	}

	reqLogger.Info("extracting cert from config map")

	out, ok := cm.Data["service-ca.crt"]

	if !ok {
		err = emperror.New("Error retrieving cert from config map")
		return nil, err
	}

	cert := []byte(out)
	return cert, nil

}

func WithBearerAuth(rt http.RoundTripper, token string) http.RoundTripper {
	addHead := WithHeader(rt)
	addHead.Header.Set("Authorization", "Bearer "+token)
	return addHead
}

type withHeader struct {
	http.Header
	rt http.RoundTripper
}

func WithHeader(rt http.RoundTripper) withHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return withHeader{Header: make(http.Header), rt: rt}
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}
