/*
 * Copyright (c) 2023 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"io/ioutil"
	"os"

	"github.com/deepflowio/deepflow/server/ingester/config"

	logging "github.com/op/go-logging"
	yaml "gopkg.in/yaml.v2"
)

var log = logging.MustGetLogger("flow_log.config")

const (
	DefaultThrottle          = 50000
	DefaultDecoderQueueCount = 2
	DefaultDecoderQueueSize  = 1 << 14
	DefaultBrokerQueueSize   = 1 << 14
	DefaultFlowLogTTL        = 72 // hour
)

var DefaultOtlpExportDatas = []string{"cbpf-net-span", "ebpf-sys-span"}
var DefaultOtlpExportDataTypes = []string{"service_info", "tracing_info", "network_layer", "flow_info", "transport_layer", "application_layer", "metrics"}

type FlowLogTTL struct {
	L4FlowLog int `yaml:"l4-flow-log"`
	L7FlowLog int `yaml:"l7-flow-log"`
	L4Packet  int `yaml:"l4-packet"`
}

type ExporterConfig struct {
	Enabled                     bool              `yaml:"enabled"`
	Addr                        string            `yaml:"addr"`
	QueueCount                  int               `yaml:"queue-count"`
	QueueSize                   int               `yaml:"queue-size"`
	ExportDatas                 []string          `yaml:"export-datas"`
	ExportDataTypes             []string          `yaml:"export-data-types"`
	ExportCustomK8sLabelsRegexp string            `yaml:"export-custom-k8s-labels-regexp"`
	ExportOnlyWithTraceID       bool              `yaml:"export-only-with-traceid"`
	GrpcHeaders                 map[string]string `yaml:"grpc-headers"`
}

type Config struct {
	Base              *config.Config
	CKWriterConfig    config.CKWriterConfig `yaml:"flowlog-ck-writer"`
	Throttle          int                   `yaml:"throttle"`
	L4Throttle        int                   `yaml:"l4-throttle"`
	L7Throttle        int                   `yaml:"l7-throttle"`
	FlowLogTTL        FlowLogTTL            `yaml:"flow-log-ttl-hour"`
	DecoderQueueCount int                   `yaml:"flow-log-decoder-queue-count"`
	DecoderQueueSize  int                   `yaml:"flow-log-decoder-queue-size"`
	Exporter          ExporterConfig        `yaml:"otlp-exporter"`
}
type FlowLogConfig struct {
	FlowLog Config `yaml:"ingester"`
}

func (c *Config) Validate() error {
	if c.DecoderQueueCount == 0 {
		c.DecoderQueueCount = DefaultDecoderQueueCount
	}

	if c.FlowLogTTL.L4FlowLog == 0 {
		c.FlowLogTTL.L4FlowLog = DefaultFlowLogTTL
	}

	if c.FlowLogTTL.L7FlowLog == 0 {
		c.FlowLogTTL.L7FlowLog = DefaultFlowLogTTL
	}

	if c.FlowLogTTL.L4Packet == 0 {
		c.FlowLogTTL.L4Packet = DefaultFlowLogTTL
	}

	if len(c.Exporter.ExportDatas) == 0 {
		c.Exporter.ExportDatas = DefaultOtlpExportDatas
	}

	if len(c.Exporter.ExportDataTypes) == 0 {
		c.Exporter.ExportDataTypes = DefaultOtlpExportDataTypes
	}

	return nil
}

func Load(base *config.Config, path string) *Config {
	config := &FlowLogConfig{
		FlowLog: Config{
			Base:              base,
			Throttle:          DefaultThrottle,
			DecoderQueueCount: DefaultDecoderQueueCount,
			DecoderQueueSize:  DefaultDecoderQueueSize,
			CKWriterConfig:    config.CKWriterConfig{QueueCount: 1, QueueSize: 1000000, BatchSize: 512000, FlushTimeout: 10},
			FlowLogTTL:        FlowLogTTL{DefaultFlowLogTTL, DefaultFlowLogTTL, DefaultFlowLogTTL},
			Exporter: ExporterConfig{
				false,
				"127.0.0.1:4317",
				2,
				100000,
				DefaultOtlpExportDatas,
				DefaultOtlpExportDataTypes,
				"",
				false,
				nil,
			},
		},
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Info("no config file, use defaults")
		return &config.FlowLog
	}
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Warning("Read config file error:", err)
		config.FlowLog.Validate()
		return &config.FlowLog
	}
	if err = yaml.Unmarshal(configBytes, &config); err != nil {
		log.Error("Unmarshal yaml error:", err)
		os.Exit(1)
	}

	if err = config.FlowLog.Validate(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return &config.FlowLog
}
