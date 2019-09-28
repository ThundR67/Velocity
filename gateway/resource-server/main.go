package main

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	micro "github.com/micro/go-micro"
	"go.uber.org/zap"

	"github.com/SonicRoshan/Velocity/gateway/resource-server/resources/main/users"
	"github.com/SonicRoshan/Velocity/global/clients"
	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/logger"
	"github.com/SonicRoshan/Velocity/global/utils"
)

var srv micro.Service
var log = logger.GetLogger("resource_server.logs")

//middleware is the middleware for graphql handler
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		jwtToken := r.Header.Get(config.ResourceSrvConfigJWTHeader)

		log.Debug("MiddleWare Got JWT Token", zap.String("Token", jwtToken))

		jwtSrv := clients.NewJWTClient(srv)
		valid, userID, scopes, err := jwtSrv.ValidateToken(jwtToken, config.TokenTypeAccess)

		log.Debug("Got Request", zap.String("UserID", userID), zap.Strings("Scopes", scopes))

		if !valid || err != nil || scopes == nil {
			utils.GatewayRespond(w, nil, config.InvalidTokenMsg, nil, log)
			return
		}

		ctx := context.WithValue(r.Context(), config.ResourceSrvConfigUserIDKey, userID)
		ctx = context.WithValue(ctx, config.ResourceSrvConfigScopesKey, scopes)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})

}

func main() {

	defer utils.HandlePanic(log)

	srv = utils.CreateService(config.ResourceServerService)
	fields, _ := users.GetResource(srv)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, _ := graphql.NewSchema(schemaConfig)

	graphQLHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	http.Handle("/graphql", middleware(graphQLHandler))
	http.ListenAndServe(":8080", nil)
}
