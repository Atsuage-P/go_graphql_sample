package main

import (
	"database/sql"
	"log"
	"mygql/env"
	"mygql/graph"
	"mygql/graph/services"
	"mygql/internal"
	"net/http"
	"os"
	"time"

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
		Complexity: graph.ComplexityConfig(),
	}))
	srv.Use(extension.FixedComplexityLimit(10))

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
