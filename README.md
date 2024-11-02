# go-graphql-practice
Go GraphQLによるTODOリストの構築(gqlgen)

## Resolverのテストについて
Resolver単体は最初はすごく薄いと思うが(Service層の内容をそのまま返しているだけなので)、Service層のテストができていれば、Resolverのテストはメリットが薄いという考えもあるが、開発が進んでくると、同じFieldを複数のResolverから取得するようになってくると思うので、リクエストレベルできちんと意図した型のデータが返ってくるか確かめたいケースが出てくると考える

→ フロントエンドも含めてE2Eで捕捉する方針が一案

CIでの課金の心配がない状態であれば(Playwrightとか使えば)、バックエンドの修正だけでもE2E走らせてあげるのもアリかな？と思う一方、そうなるとCIの実行時間が長くなりすぎて、リリース速度に影響しそう

→ バックエンドの不具合はバックエンドのテストで拾えるのが理想かと考える

Resolverで指定したFieldが取得できるのかをテストしたい(純粋なhttptestを使えばできる)

## 参考
- https://qiita.com/zigenin/items/df9359bf2f209d08f117
- https://github.com/99designs/gqlgen?tab=readme-ov-file
- https://qiita.com/hiroyky/items/4d7764172e73ff54f18b
- autobind
	- https://qiita.com/ryota-yamamoto/items/3f15f476f17db047ef5d
	- https://tech.layerx.co.jp/entry/2021/10/22/171242
- カスタムエラーを返す
	- https://www.wheatandcat.me/entry/2022/03/22/084127
- gqlgenとhttptestを使用してテストを行う
	- https://zenn.dev/ygsn/articles/b5d843774ed524
	- https://budougumi0617.github.io/2020/05/29/go-testing-httptest/

## やってみた所管
- sqlboilerと組み合わせて行うと、go-playground/validatorが使いづらい...
	- sqlboilerによる自動生成モデルとgqlgenによるinputの自動生成型ファイルにvalidateタグを仕込むわけにはいかん...
	- となると、validate用に別途構造体の作成・マッピングがいる
	- コードベースでバリデーションルールを書けるozzo-validationの方が使いやすいかも


go get github.com/99designs/gqlgen/codegen/config@v0.17.55
go get github.com/99designs/gqlgen/internal/imports@v0.17.55
go get github.com/99designs/gqlgen@v0.17.55
go run github.com/99designs/gqlgen
