// Copyright 2021-now by lifenjoiner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lifenjoiner/pd/dispatcher"
	"github.com/lifenjoiner/pd/hoststat"
	"github.com/lifenjoiner/pd/proxypool"
	"github.com/lifenjoiner/pd/server/tcp"
	"github.com/lifenjoiner/pd/statichost"
)

var (
	name = "pd"
	version = "0.1.0"
	description = "A local proxy dispatcher."
)

// ServeFromConfig starts the serving.
func ServeFromConfig(config *Config) {
	svrConf := &config.SvrConf
	dispatcher.GlobalStaticHosts = statichost.MapStaticFiles(config.Blocked, config.Direct)
	dispatcher.GlobalHostStats = hoststat.MapStatsFile(config.StatFile)
	dispatcher.StartProbeDirect(config.NetProbeUrl, svrConf.UpstreamTimeout)
	go func(){
		dispatcher.GlobalProxyPool = proxypool.InitProxyPool(svrConf.Proxies, svrConf.ProxyProbeUrl, svrConf.UpstreamTimeout)
	}()

	for _, listen := range config.Listens {
		s := &tcp.TCPServer{listen, svrConf}
		go s.ListenAndServe()
	}
}

func main() {
	cfg := parseConfig()
	log.Printf("%v v%v %v", name, version, description)
	ServeFromConfig(cfg)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
