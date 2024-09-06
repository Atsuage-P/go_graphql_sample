package main

import (
	"bytes"
	"context"
	"database/sql"
	"log"
	"mygql/env"
	"mygql/graph"
	"mygql/graph/services"
	"mygql/internal"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	boil.DebugMode = true

	cnf := env.LoadEnv()
	db := ConnectDB(&cnf.DB)
	service := services.New(db)

	srv := handler.NewDefaultServer(internal.NewExecutableSchema(internal.Config{
		Resolvers: &graph.Resolver{
			Srv:     service,
			Loaders: graph.NewLoaders(service),
		},
		// クエリ複雑度の個別制限
		Complexity: graph.ComplexityConfig(),
	}))
	// クエリ複雑度の一括制限
	srv.Use(extension.FixedComplexityLimit(10))
	// リクエストを受け取った時に呼ばれるミドルウェア
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		log.Println("before OperationHandler")
		res := next(ctx)
		defer log.Println("after OperationHandler")
		return res
	})
	// レスポンスを作成する段階で呼ばれるミドルウェア
	srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		log.Println("before ResponseHandler")
		res := next(ctx)
		defer log.Println("after ResponseHandler")
		return res
	})
	// ルートリゾルバの実行前後用のミドルウェア
	srv.AroundRootFields(func(ctx context.Context, next graphql.RootResolver) graphql.Marshaler {
		log.Println("before RootResolver")
		res := next(ctx)
		defer func() {
			var b bytes.Buffer
			res.MarshalGQL(&b)
			log.Println("after RootResolver", b.String())
		}()
		return res
	})
	// レスポンスに含めるjsonフィールドを1つ作る処理の前後用のミドルウェア
	srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res any, err error) {
		log.Println("before Resolver")
		res, err = next(ctx)
		defer log.Println("after Resolver", res)
		return
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func ConnectDB(cnf *env.DBConfig) *sql.DB {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("DB connect Error: %v", err)
	}
	c := mysql.Config{
		DBName:    cnf.Name,
		User:      cnf.User,
		Passwd:    cnf.Password,
		Addr:      cnf.Host + ":" + cnf.Port,
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}
	db, err := sql.Open(cnf.MS, c.FormatDSN())
	if err != nil {
		log.Fatalf("DB connect Error: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("DB Ping Error: %v", err)
	}
	db.SetConnMaxLifetime(time.Duration(cnf.MaxLifeTimeMin))
	db.SetMaxOpenConns(cnf.MaxOpenConns)
	db.SetMaxIdleConns(cnf.MaxOpenIdleConns)

	return db
}
