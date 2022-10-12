package utils

//import (
//	"fmt"
//	"github.com/hashicorp/consul/api"
//	"github.com/prometheus/common/log"
//	"strconv"
//	"strings"
//)
//
//const (
//	HTTPAddrEnvName      = "CONSUL_HTTP_ADDR"
//	HTTPTokenFileEnvName = "CONSUL_HTTP_TOKEN_FILE"
//	HTTPTokenEnvName     = "CONSUL_HTTP_TOKEN"
//	HTTPAuthEnvName      = "CONSUL_HTTP_AUTH"
//	HTTPSSLEnvName       = "CONSUL_HTTP_SSL"
//	HTTPCAFile           = "CONSUL_CACERT"
//	HTTPCAPath           = "CONSUL_CAPATH"
//	HTTPClientCert       = "CONSUL_CLIENT_CERT"
//	HTTPClientKey        = "CONSUL_CLIENT_KEY"
//	HTTPTLSServerName    = "CONSUL_TLS_SERVER_NAME"
//	HTTPSSLVerifyEnvName = "CONSUL_HTTP_SSL_VERIFY"
//	HTTPNamespaceEnvName = "CONSUL_NAMESPACE"
//	HTTPPartitionEnvName = "CONSUL_PARTITION"
//)
//
//// consul 配置项
//type consulOpt func(config *api.Config)
//
//type consulConfig func(string) consulOpt
//
//var keymap = map[string]string{
//	strings.ToLower("Address"):    HTTPAddrEnvName,
//	strings.ToLower("Token"):      HTTPTokenFileEnvName,
//	strings.ToLower("TokenFile"):  HTTPTokenEnvName,
//	strings.ToLower("Auth"):       HTTPAuthEnvName,
//	strings.ToLower("SSL"):        HTTPSSLEnvName,
//	strings.ToLower("TLSServer"):  HTTPTLSServerName,
//	strings.ToLower("CAFile"):     HTTPCAFile,
//	strings.ToLower("CAPath"):     HTTPCAPath,
//	strings.ToLower("ClientCert"): HTTPClientCert,
//	strings.ToLower("ClientKey"):  HTTPClientKey,
//	strings.ToLower("SSLVerify"):  HTTPSSLVerifyEnvName,
//	strings.ToLower("Namespace"):  HTTPNamespaceEnvName,
//	strings.ToLower("Partition"):  HTTPPartitionEnvName,
//}
//
//var configMap = map[string]consulConfig{
//	HTTPAddrEnvName:      consulHTTPAddrEnvName,
//	HTTPTokenFileEnvName: consulHTTPTokenFileEnvName,
//	HTTPTokenEnvName:     consulHTTPTokenEnvName,
//	HTTPAuthEnvName:      consulHTTPAuthEnvName,
//	HTTPSSLEnvName:       consulHTTPSSLEnvName,
//	HTTPTLSServerName:    consulHTTPTLSServerName,
//	HTTPCAFile:           consulHTTPCAFile,
//	HTTPCAPath:           consulHTTPCAPath,
//	HTTPClientCert:       consulHTTPClientCert,
//	HTTPClientKey:        consulHTTPClientKey,
//	HTTPSSLVerifyEnvName: consulHTTPSSLVerifyEnvName,
//	HTTPNamespaceEnvName: consulHTTPNamespaceEnvName,
//	HTTPPartitionEnvName: consulHTTPPartitionEnvName,
//}
//
//func consulHTTPAddrEnvName(address string) consulOpt {
//	return func(config *api.Config) {
//		config.Address = address
//	}
//}
//
//func consulHTTPTokenFileEnvName(tokenFile string) consulOpt {
//	return func(config *api.Config) {
//		config.TokenFile = tokenFile
//	}
//}
//
//func consulHTTPTokenEnvName(token string) consulOpt {
//	return func(config *api.Config) {
//		config.Token = token
//	}
//}
//
//func consulHTTPAuthEnvName(auth string) consulOpt {
//	return func(config *api.Config) {
//		var username, password string
//		if strings.Contains(auth, ":") {
//			split := strings.SplitN(auth, ":", 2)
//			username = split[0]
//			password = split[1]
//		} else {
//			username = auth
//		}
//
//		config.HttpAuth = &api.HttpBasicAuth{
//			Username: username,
//			Password: password,
//		}
//	}
//}
//
//func consulHTTPSSLEnvName(ssl string) consulOpt {
//	return func(config *api.Config) {
//		enabled, err := strconv.ParseBool(ssl)
//		if err != nil {
//			log.Warn(fmt.Sprintf("could not parse %s", HTTPSSLEnvName), "error", err)
//		}
//
//		if enabled {
//			config.Scheme = "https"
//		}
//	}
//}
//
//func consulHTTPTLSServerName(v string) consulOpt {
//	return func(config *api.Config) {
//		config.TLSConfig.Address = v
//	}
//}
//
//func consulHTTPCAFile(v string) consulOpt {
//	return func(config *api.Config) {
//		config.TLSConfig.CAFile = v
//	}
//}
//
//func consulHTTPCAPath(v string) consulOpt {
//	return func(config *api.Config) {
//		config.TLSConfig.CAPath = v
//	}
//}
//
//func consulHTTPClientCert(v string) consulOpt {
//	return func(config *api.Config) {
//		config.TLSConfig.CertFile = v
//	}
//}
//
//func consulHTTPClientKey(v string) consulOpt {
//	return func(config *api.Config) {
//		config.TLSConfig.KeyFile = v
//	}
//}
//
//func consulHTTPSSLVerifyEnvName(v string) consulOpt {
//	return func(config *api.Config) {
//		doVerify, err := strconv.ParseBool(v)
//		if err != nil {
//			log.Warn(fmt.Sprintf("could not parse %s", HTTPSSLVerifyEnvName), "error", err)
//		}
//		if !doVerify {
//			config.TLSConfig.InsecureSkipVerify = true
//		}
//	}
//}
//
//func consulHTTPNamespaceEnvName(v string) consulOpt {
//	return func(config *api.Config) {
//		config.Namespace = v
//	}
//}
//
//func consulHTTPPartitionEnvName(v string) consulOpt {
//	return func(config *api.Config) {
//		config.Partition = v
//	}
//}
