/*
Copyright 2021 The KodeRover Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	"net/http"
	"time"

	"github.com/koderover/zadig/pkg/microservice/aslan/core"
	"github.com/koderover/zadig/pkg/microservice/aslan/server/rest"
	"github.com/koderover/zadig/pkg/tool/kube/client"
	"github.com/koderover/zadig/pkg/tool/log"
)

func Serve(ctx context.Context) error {
	go func() {
		if err := client.Start(ctx); err != nil {
			panic(err)
		}
	}()

	core.Setup()
	defer core.TearDown()

	log.Infof("App Aslan Started at %s", time.Now())

	engine := rest.NewEngine()
	server := &http.Server{Addr: ":25000", Handler: engine}

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Errorf("Failed to stop server, error: %s", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Errorf("Failed to start http server, error: %s", err)
		return err
	}

	return nil
}
