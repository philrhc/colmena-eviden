/*
Copyright © 2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This work has been implemented within the context of COLMENA project.
*/
package main

import (
	"colmena/sla-management-svc/app/assessment"
	"colmena/sla-management-svc/app/assessment/monitor"
	"colmena/sla-management-svc/app/assessment/monitor/genericadapter"
	"colmena/sla-management-svc/app/assessment/monitor/prometheus"
	"colmena/sla-management-svc/app/assessment/monitor/testadapter"
	"colmena/sla-management-svc/app/assessment/notifier"
	grpc_notifier "colmena/sla-management-svc/app/assessment/notifier/grpc"
	"colmena/sla-management-svc/app/assessment/notifier/lognotifier"
	"colmena/sla-management-svc/app/assessment/notifier/rest"
	"colmena/sla-management-svc/app/common/cfg"
	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"
	"colmena/sla-management-svc/app/repositories/memrepository"
	grpc "colmena/sla-management-svc/protobuf/server"
	restAPI "colmena/sla-management-svc/rest-api"
	"sync"

	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// path used in logs
const pathLOG string = "SLA > "

/*
Main function. Environment variables used by the SLA & QoS Manager:
  - PROMETHEUS_ADDRESS (e.g., "http://localhost:9090")
  - MONITORING_ADAPTER (e.g., "prometheus")
  - NOTIFIER_ADAPTER (e.g., "rest_endpoint", "rpc")
  - NOTIFICATION_ENDPOINT (e.g., "http://localhost:10090")
  - GRPC_PORT (e.g., "8099")
*/
func main() {
	// tests: environment variables
	// os.Setenv("MONITORING_ADAPTER", "prometheus")
	// os.Setenv("PROMETHEUS_ADDRESS", "http://192.168.137.47:9090") //"http://192.168.137.25:9091") // http://localhost:9090
	//os.Setenv("NOTIFIER_ADAPTER", "rest_endpoint")                // "grpc"
	//os.Setenv("NOTIFICATION_ENDPOINT", "http://localhost:10090")  // "localhost:8099"
	//os.Setenv("GRPC_PORT", "8099")

	logs.GetLogger().Info("Starting SLA & QoS Manager [version: v0.1.0; date: 2024.10.18; id: 1] ...")

	// main configuration
	// variables are set through environment variables (i.e. using Kubernetes or Docker deployment files)
	config := createMainConfig()
	logMainConfig(config)

	checkPeriod := asSeconds(config, cfg.CheckPeriodPropertyName)
	trasientTime := asSeconds(config, cfg.TransientTimePropertyName)

	// REPOSITORY (DB)
	logs.GetLogger().Info(pathLOG + "Setting Database Adapter ...")
	repo := buildRepositoryAdapter(config)

	// NOTIFIER
	logs.GetLogger().Info(pathLOG + "Setting Notifier / Subscriber adapter ...")
	notifier := buildNotifierAdapter(config)

	// MONITORING ADAPTER
	logs.GetLogger().Info(pathLOG + "Setting Monitoring Adapter ...")
	adapter := buildMonitoringAdapter(config)

	// VALIDATOR
	validater := model.NewDefaultValidator(false, true)

	// start application - thread
	logs.GetLogger().Info(pathLOG + "Initializing QAA assessment process [THREAD] ...")
	aCfg := assessment.Config{
		Repo:      repo,
		Adapter:   adapter,
		Notifier:  notifier,
		Transient: trasientTime,
	}
	go createValidationThread(checkPeriod, aCfg) // assessment thread
	time.Sleep(2 * time.Second)

	// servers:
	var wg sync.WaitGroup

	// REST API server - thread
	a, _ := restAPI.New(aCfg, repo, validater, adapter)
	logs.GetLogger().Info(pathLOG + "Initializing SLA REST API server [THREAD] ...")
	wg.Add(1)
	go a.InitializeRESTAPI() // rest api thread
	time.Sleep(2 * time.Second)

	// gRPC server - thread
	grpc_port := os.Getenv("GRPC_PORT")

	if grpc_port != "" {
		logs.GetLogger().Info(pathLOG + "Initializing gRPC server [THREAD] ...")
		wg.Add(1)

		//go server.CreateServer(&wg, "8099")
		go grpc.InitializegRPCServer(&wg, repo)
		time.Sleep(2 * time.Second)
	}

	logs.GetLogger().Info("\t-----------------------------------------------------------------")

	wg.Wait()
}

///////////////////////////////////////////////////////////////////////////////

// buildRepositoryAdapter
func buildRepositoryAdapter(config *viper.Viper) model.IRepository {
	logs.GetLogger().Warn(pathLOG + "[Repository Adapter] Using default Database Adapter [memory repository] ... ")
	repo, errRepo := memrepository.New()
	if errRepo != nil {
		logs.GetLogger().Fatal(pathLOG+"[Repository Adapter] Error creating repository: ", errRepo.Error())
	}
	return repo
}

// buildNotifierAdapter
func buildNotifierAdapter(config *viper.Viper) notifier.ViolationNotifier {
	aType := config.GetString(cfg.NotifierAdapterPropertyName)
	switch aType {
	case "rest_endpoint":
		logs.GetLogger().Info(pathLOG + "[Notifier Adapter] Using REST-ENDPOINT notifier adapter ...")
		return rest.New(config)

	case "grpc":
		logs.GetLogger().Info(pathLOG + "[Notifier Adapter] Using gRPC notifier adapter ...")
		return grpc_notifier.New(config)

	default:
		logs.GetLogger().Warn(pathLOG + "[Notifier Adapter] Using Default Notifier (no subscriber) ...")
		return lognotifier.LogNotifier{}
	}
}

// buildMonitoringAdapter
func buildMonitoringAdapter(config *viper.Viper) monitor.MonitoringAdapter {
	aType := config.GetString(cfg.MonitoringAdapterPropertyName)
	if os.Getenv(cfg.MonitoringAdapterPropertyName) == prometheus.Name {
		aType = prometheus.Name
	} else if os.Getenv(cfg.MonitoringAdapterPropertyName) == testadapter.Name {
		aType = testadapter.Name
	}

	switch aType {
	case prometheus.Name:
		logs.GetLogger().Info(pathLOG + "[Monitoring Adapter] Using Prometheus adapter ...")
		promadapter := prometheus.New(config)
		adapter := genericadapter.New(
			"prometheus",
			promadapter.Retrieve(),
			genericadapter.Identity)
		return adapter
	default:
		logs.GetLogger().Info(pathLOG + "[Monitoring Adapter] Using Test adapter ...")
		adapter := genericadapter.New(
			"default",
			testadapter.New(config).Retrieve(),
			genericadapter.Identity)
		return adapter
	}
}

/*
Creates the main Viper configuration
*/
func createMainConfig() *viper.Viper {
	logs.GetLogger().Info(pathLOG + "[Configuration] Generating QA Agent configuration values (based on Viper) ...")

	logs.GetLogger().Debug(pathLOG + "[Configuration] Defined QA Agent environment variables:")
	logs.GetLogger().Debug(pathLOG + "    Repository Adapter ..... " + os.Getenv(cfg.RepositoryAdapterPropertyName))
	logs.GetLogger().Debug(pathLOG + "    Notifier Adapter ....... " + os.Getenv(cfg.NotifierAdapterPropertyName))
	logs.GetLogger().Debug(pathLOG + "    Monitoring Adapter ..... " + os.Getenv(cfg.MonitoringAdapterPropertyName))
	logs.GetLogger().Debug(pathLOG + "    Cluster Key ............ " + os.Getenv(cfg.ClusterKeyPropertyName))
	logs.GetLogger().Debug(pathLOG + "    Check Period Time ...... " + os.Getenv(cfg.CheckPeriodPropertyName))
	logs.GetLogger().Debug(pathLOG + "    Transient Time ......... " + os.Getenv(cfg.TransientTimePropertyName))

	// new viper.Viper - CONFIGURATION OBJECT
	config := viper.New()
	config.SetEnvPrefix(cfg.ConfigPrefix) // Env vars start with 'qaa_'
	config.AutomaticEnv()

	// QAA CONFIGURATION:
	logs.GetLogger().Info(pathLOG + "[Configuration] Setting configuration values ...")
	// CheckPeriod
	if os.Getenv(cfg.CheckPeriodPropertyName) == "" {
		config.SetDefault(cfg.CheckPeriodPropertyName, cfg.DefaultCheckPeriod)
	} else {
		config.SetDefault(cfg.CheckPeriodPropertyName, os.Getenv(cfg.CheckPeriodPropertyName))
	}
	// TransientTime
	config.SetDefault(cfg.TransientTimePropertyName, cfg.DefaultTransientTime)
	// ADAPTERS
	// Repository
	if os.Getenv(cfg.RepositoryAdapterPropertyName) == "" {
		config.SetDefault(cfg.RepositoryAdapterPropertyName, cfg.DefaultRepositoryType)
	} else {
		config.SetDefault(cfg.RepositoryAdapterPropertyName, os.Getenv(cfg.RepositoryAdapterPropertyName))
	}
	// Notifier
	if os.Getenv(cfg.NotifierAdapterPropertyName) == "" {
		config.SetDefault(cfg.NotifierAdapterPropertyName, cfg.DefaultNotifierType)
	} else {
		config.SetDefault(cfg.NotifierAdapterPropertyName, os.Getenv(cfg.NotifierAdapterPropertyName))
	}

	if os.Getenv(cfg.NotificationURLPropertyName) == "" {
		config.SetDefault(cfg.NotificationURLPropertyName, cfg.DefaultNotificationURL)
	} else {
		config.SetDefault(cfg.NotificationURLPropertyName, os.Getenv(cfg.NotificationURLPropertyName))
	}

	// Monitoring
	if os.Getenv(cfg.MonitoringAdapterPropertyName) == "" {
		config.SetDefault(cfg.MonitoringAdapterPropertyName, cfg.DefaultMonitoringAdapterType)
	} else {
		config.SetDefault(cfg.MonitoringAdapterPropertyName, os.Getenv(cfg.MonitoringAdapterPropertyName))
	}
	// ClusterKey
	if os.Getenv(cfg.ClusterKeyPropertyName) == "" {
		config.SetDefault(cfg.ClusterKeyPropertyName, cfg.DefaultClusterKey)
	} else {
		config.SetDefault(cfg.ClusterKeyPropertyName, os.Getenv(cfg.ClusterKeyPropertyName))
	}

	logs.GetLogger().Debug(pathLOG + "[Configuration] Returning configuration object ...")
	return config
}

// asSeconds
func asSeconds(config *viper.Viper, field string) time.Duration {
	raw := config.GetString(field)
	// if it is already a valid duration, return directly
	if _, err := time.ParseDuration(raw); err == nil {
		return config.GetDuration(field)
	}

	// if not, assume it is (decimal) number of seconds; read as ms and convert to seconds.
	ms := config.GetFloat64(field)
	return time.Duration(ms*1000) * time.Millisecond
}

// shortDur
func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

func logMainConfig(config *viper.Viper) {
	logs.GetLogger().Info(pathLOG + "[Main Configuration] Loading initial configuration values ... ")

	checkPeriod := asSeconds(config, cfg.CheckPeriodPropertyName)
	repoType := config.GetString(cfg.RepositoryAdapterPropertyName)
	notifierType := config.GetString(cfg.NotifierAdapterPropertyName)
	adapterType := config.GetString(cfg.MonitoringAdapterPropertyName)
	//transientTime := asSeconds(config, cfg.TransientTimePropertyName)

	logs.GetLogger().Info(pathLOG + "[Main Configuration] QoS Agent initialization values:\n" +
		"\t-----------------------------------------------------------------\n" +
		//"\tConfigfile:              " + config.ConfigFileUsed() + "\n" +
		"\tRepository type (DB):    " + repoType + "\n" +
		"\tMonitoring Adapter type: " + adapterType + "\n" +
		"\tNotifier type:           " + notifierType + "\n" +
		//"\tTransient time:          " + shortDur(transientTime) + "\n" +
		"\tCheck period:            " + shortDur(checkPeriod) + "\n" +
		"\t-----------------------------------------------------------------")
}

// createValidationThread
func createValidationThread(checkPeriod time.Duration, cfg assessment.Config) {
	ticker := time.NewTicker(checkPeriod)

	for {
		<-ticker.C
		cfg.Now = time.Now()
		assessment.AssessActiveQoSDefinitions(cfg)
	}
}
