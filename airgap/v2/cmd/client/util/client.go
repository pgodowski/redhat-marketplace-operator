// Copyright 2021 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"emperror.dev/errors"
	"github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/apis/fileretriever"
	"github.com/redhat-marketplace/redhat-marketplace-operator/airgap/v2/apis/filesender"
	"github.com/spf13/viper"
	"github.com/twitchtv/twirp"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

// InitClient Initializes gRPC client connection
func InitClient() (*grpc.ClientConn, error) {
	// Create connection
	insecure := viper.GetBool("insecure")
	var conn *grpc.ClientConn
	var err error

	// Fetch target address
	address := viper.GetString("address")
	if len(strings.TrimSpace(address)) == 0 {
		return nil, fmt.Errorf("target address is blank/empty")
	}

	options := []grpc.DialOption{}

	if insecure {
		options = append(options, grpc.WithInsecure())
	} else {
		cert := viper.GetString("certificate-path")
		caCertPool := x509.NewCertPool()

		caCert, err := ioutil.ReadFile(cert)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load cert file")
		}

		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("failed to append certs %s", cert)
		}

		creds := credentials.NewTLS(&tls.Config{
			RootCAs: caCertPool,
		})

		//creds.OverrideServerName("rhm-data-service.openshift-redhat-marketplace.svc.cluster.local")
		options = append(options, grpc.WithTransportCredentials(creds))

		tokenString := viper.GetString("token")
		if tokenString != "" {
			oauth2Token := &oauth2.Token{
				AccessToken: tokenString,
			}
			fmt.Println(tokenString)
			perRPC := oauth.NewOauthAccess(oauth2Token)
			options = append(options, grpc.WithPerRPCCredentials(perRPC))
		}
	}

	conn, err = grpc.Dial(address, options...)

	if err != nil {
		return nil, fmt.Errorf("connection error: %v", err)
	}
	return conn, nil
}

type Closeable func() error

func InitFileRetrieverProtobufClient() (context.Context, fileretriever.FileRetriever, error) {
	address := viper.GetString("address")
	token := viper.GetString("token")
	tlsConfig, err := getTlS()

	if err != nil {
		return nil, nil, err
	}

	client := fileretriever.NewFileRetrieverProtobufClient(address,
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
		twirp.WithClientPathPrefix(""))

	header := make(http.Header)
	header.Add("authorization", "bearer "+token)

	ctx := context.Background()
	ctx, err = twirp.WithHTTPRequestHeaders(ctx, header)

	if err != nil {
		return ctx, client, err
	}

	return ctx, client, nil
}

func InitFileSenderProtobufClient() (context.Context, filesender.FileSender, error) {
	address := viper.GetString("address")
	token := viper.GetString("token")
	tlsConfig, err := getTlS()

	if err != nil {
		return nil, nil, err
	}

	client := filesender.NewFileSenderProtobufClient(address,
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
		twirp.WithClientPathPrefix(""))

	header := make(http.Header)
	header.Add("authorization", "bearer "+token)

	ctx := context.Background()
	ctx, err = twirp.WithHTTPRequestHeaders(ctx, header)

	if err != nil {
		return ctx, client, err
	}

	return ctx, client, nil
}

func getTlS() (tlsConfig *tls.Config, err error) {
	tlsEnabled := viper.GetBool("tls")
	serverName := viper.GetString("serverName")
	insecure := viper.GetBool("insecure")

	if tlsEnabled {
		if insecure {
			tlsConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			var caCert []byte
			cert := viper.GetString("ca-path")
			caCertPool := x509.NewCertPool()

			caCert, err = ioutil.ReadFile(cert)
			if err != nil {
				return
			}

			ok := caCertPool.AppendCertsFromPEM(caCert)
			if !ok {
				return nil, fmt.Errorf("failed to load ca %s", cert)
			}

			tlsConfig = &tls.Config{
				RootCAs:    caCertPool,
				ServerName: serverName,
			}
		}
	}
	return
}
