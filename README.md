# GraphQLのお試し
## 概要
- [zennの記事](https://zenn.dev/hsaki/books/golang-graphql) を参考

## 動作確認
```sh
docker compose up -d
docker compose exec app bash

# サーバー起動
go run ./server.go
```
- サーバー起動後、`localhost:8080` で左画面にクエリを入力して実行するとレスポンス確認可

## メモ
- [gqlgen](https://github.com/99designs/gqlgen) によってコードを生成する
- 取得(`query`)、追加・更新・削除(`mutation`)
- `gqlgen.yml` と `schema.graphqls` を用意して `gqlgen generate` によりリゾルバやモデルのコードを生成する
  - リゾルバとはデータ操作を行うもののことで、実態は特定のフィールドのデータを返す関数のこと
- `sqlboiler` 使ってコード自動生成
```sh
# sqlboiler.tomlを作成して以下を実行
sqlboiler mysql
```
- ページネーションの[before,after,first,last部分について](https://tekrog.com/graphql-relay-cursor)
