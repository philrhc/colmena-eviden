/*
  COLMENA-DESCRIPTION-SERVICE
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
package cfg

import "time"

const (
	// ConfigPrefix is the prefix of env vars that configure the QoS-Aerts-Agent
	ConfigPrefix string = "qaa"

	// Repository
	// RepositoryAdapterPropertyName is the name of the property repository adater type
	RepositoryAdapterPropertyName string = "repository_adapter"
	// DefaultRepositoryType is the name of the default repository
	DefaultRepositoryType string = "memory"

	// Notifier
	// NotifierAdapterPropertyName is the name of the property notifier adater type
	NotifierAdapterPropertyName string = "NOTIFIER_ADAPTER"
	// DefaultNotifierType is the name of the default notifier
	DefaultNotifierType string = "default"
	// RestNotifierType is the name of the REST notifier
	RestNotifierType string = "rest_endpoint"
	// gRPCNotifierType is the name of the gRPC notifier
	GRPCNotifierType string = "grpc"
	// RabbitMQNotifierType is the name of the RabbitMQ notifier
	RabbitMQNotifierType string = "rabbitmq"
	// NotificationURLPropertyName is the name of the property notificationUrl
	NotificationURLPropertyName string = "NOTIFICATION_ENDPOINT"
	// DefaultNotificationURL is the name of the default notifier
	DefaultNotificationURL string = "http://localhost:10090"

	// Context Zenoh Endpoint
	ContextZenohEndpointPropertyName string = "CONTEXT_ZENOH_ENDPOINT"
	DefaultContextZenohEndpoint      string = "http://localhost:8000"
	ContextZenohContextsPropertyName string = "CONTEXT_ZENOH_CONTEXTS"
	DefaultContextZenohContexts      string = "colmena/contexts"

	// Monitoring
	// MonitoringAdapterPropertyName is the name of the property monitoring adater type
	MonitoringAdapterPropertyName string = "MONITORING_ADAPTER"
	// DefaultMonitoringAdapterType is the name of the default adapter
	DefaultMonitoringAdapterType string = "dummy"

	// COMPOSE_PROJECT_NAME
	ComposeProjectPropertyName string = "COMPOSE_PROJECT_NAME"
	AgentIdPropertyName        string = "AGENT_ID"

	// Assessment
	// CheckPeriodPropertyName is the name of the property CheckPeriod
	CheckPeriodPropertyName string = "checkPeriod"
	// DefaultCheckPeriod is the default number of seconds of the periodic assessment execution
	DefaultCheckPeriod time.Duration = 30 * time.Second

	// TransientTimePropertyName is the name of the property that holds the number of
	// seconds to wait until a new violation for a guarantee term is raised
	TransientTimePropertyName string = "transientTime"
	// DefaultTransientTime is the default number of seconds after a violation
	// to raise a violation for the same guarantee term
	DefaultTransientTime time.Duration = 0

	// Assessment
	ASSESSMENT_X string = "ASSESSMENT_X"
	ASSESSMENT_Y string = "ASSESSMENT_Y"
	ASSESSMENT_Z string = "ASSESSMENT_Z"
)
