/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2023 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package data

//edgex 的采集数据库

import (
	"context"
	"os"

	"github.com/gorilla/mux"

	"github.com/edgexfoundry/go-mod-bootstrap/v3/bootstrap"
	"github.com/edgexfoundry/go-mod-bootstrap/v3/bootstrap/flags"
	"github.com/edgexfoundry/go-mod-bootstrap/v3/bootstrap/handlers"
	"github.com/edgexfoundry/go-mod-bootstrap/v3/bootstrap/interfaces"
	"github.com/edgexfoundry/go-mod-bootstrap/v3/bootstrap/startup"
	bootstrapConfig "github.com/edgexfoundry/go-mod-bootstrap/v3/config"
	"github.com/edgexfoundry/go-mod-bootstrap/v3/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/common"

	"github.com/edgexfoundry/edgex-go"
	"github.com/edgexfoundry/edgex-go/internal/core/data/application"
	"github.com/edgexfoundry/edgex-go/internal/core/data/config"
	"github.com/edgexfoundry/edgex-go/internal/core/data/container"
	pkgHandlers "github.com/edgexfoundry/edgex-go/internal/pkg/bootstrap/handlers"
)

func Main(ctx context.Context, cancel context.CancelFunc, router *mux.Router) {

	//服务启动计时
	startupTimer := startup.NewStartUpTimer(common.CoreDataServiceKey)

	// All common command-line flags have been moved to DefaultCommonFlags. Service specific flags can be add here,
	// by inserting service specific flag prior to call to commonFlags.Parse().
	// Example:
	// 		flags.FlagSet.StringVar(&myvar, "m", "", "Specify a ....")
	//      ....
	//      flags.Parse(os.Args[1:])
	//
	f := flags.New()
	//通用的服务参数， 解析参数
	f.Parse(os.Args[1:])

	//将 配置对象 添加到 依赖注入容器中
	configuration := &config.ConfigurationStruct{}
	dic := di.NewContainer(di.ServiceConstructorMap{
		container.ConfigurationName: func(get di.Get) interface{} {
			return configuration
		},
	})

	//启动 rest 服务器
	httpServer := handlers.NewHttpServer(router, true)

	//开始启动服务
	bootstrap.Run(
		ctx,
		cancel,
		f,
		common.CoreDataServiceKey,
		common.ConfigStemCore,
		configuration,
		startupTimer,
		dic,
		true,
		bootstrapConfig.ServiceTypeOther,
		[]interfaces.BootstrapHandler{
			// 新建数据库并启动数据库
			pkgHandlers.NewDatabase(httpServer, configuration, container.DBClientInterfaceName).BootstrapHandler, // add v2 db client bootstrap handler

			// 创建并初始化 Message Client
			handlers.MessagingBootstrapHandler,

			// 服务指标监测
			handlers.NewServiceMetrics(common.CoreDataServiceKey).BootstrapHandler, // Must be after Messaging

			// App 启动
			application.BootstrapHandler, // Must be after Service Metrics and before next handler

			//初始化 REST 路由 并且 订阅客户端事件
			NewBootstrap(router, common.CoreDataServiceKey).BootstrapHandler,

			// 启动 HTTP 服务器
			httpServer.BootstrapHandler,

			// 输出启动成功的消息
			handlers.NewStartMessage(common.CoreDataServiceKey, edgex.Version).BootstrapHandler,
		},
	)
}
