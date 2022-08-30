package config

import (
	"goravel/app/providers"

	"github.com/goravel/framework/cache"
	"github.com/goravel/framework/console"
	"github.com/goravel/framework/contracts"
	"github.com/goravel/framework/database"
	foundation "github.com/goravel/framework/foundation/providers"
	"github.com/goravel/framework/http"
	"github.com/goravel/framework/log"
	"github.com/goravel/framework/queue"
	"github.com/goravel/framework/route"
	"github.com/goravel/framework/schedule"
	"github.com/goravel/framework/support/facades"
)

//Boot Start all init methods of the current folder to bootstrap all config.
func Boot() {}

func init() {
	config := facades.Config
	config.Add("app", map[string]interface{}{
		//Application Name
		//This value is the name of your application. This value is used when the
		//framework needs to place the application's name in a notification or
		//any other location as required by the application or its packages.
		"name": config.Env("APP_NAME", "nft"),

		//Application Environment
		//This value determines the "environment" your application is currently
		//running in. This may determine how you prefer to configure various
		//services the application utilizes. Set this in your ".env" file.
		"env": config.Env("APP_ENV", "production"),

		//Application Debug Mode
		"debug": config.Env("APP_DEBUG", false),

		//Encryption Key
		//32 character string, otherwise these encrypted strings
		//will not be safe. Please do this before deploying an application!
		"key": config.Env("APP_KEY", ""),

		//Application URL
		"url": config.Env("APP_URL", "http://localhost"),

		//Application host, http server listening address.
		"host": config.Env("APP_HOST", "127.0.0.1:3000"),

		"grpc_host": config.Env("GRPC_HOST", ""),

		//Autoloaded service providers
		//The service providers listed here will be automatically loaded on the
		//request to your application. Feel free to add your own services to
		//this array to grant expanded functionality to your applications.
		"providers": []contracts.ServiceProvider{
			&log.ServiceProvider{},
			&console.ServiceProvider{},
			&database.ServiceProvider{},
			&cache.ServiceProvider{},
			&http.ServiceProvider{},
			&foundation.ArtisanServiceProvider{},
			&route.ServiceProvider{},
			&schedule.ServiceProvider{},
			&queue.ServiceProvider{},
			&providers.AppServiceProvider{},
			&providers.RouteServiceProvider{},
			&providers.GrpcServiceProvider{},
			&providers.ConsoleServiceProvider{},
			&providers.QueueServiceProvider{},
		},
	})
}
