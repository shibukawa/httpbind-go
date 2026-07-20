# 型付きテンプレート言語 要件レビュー

添付の「Compact Typed Template Language Specification」を、初期実装の判断単位に分解したレビュー用文書です。正本は `.knowledge` の英語概念群であり、本書は人間向けの要約です。

## 結論

HTML と SQL は、型宣言・式・コンポーネント・構造的な `if` / `for` を共通化し、本文の解析・安全規則・Go コード生成を出力形式ごとに分離します。テンプレートは実行時に解釈せず、リフレクションなしの Go コードへ事前生成します。

推奨パッケージ構成は次のとおりです。

- `templates/htmlbind`: HTML の解析、検証、文脈別エスケープ、Go 生成
- `templates/sqlbind`: SQL の解析、構造化句、パラメータ化、結果契約、Go 生成
- `templates/internal`: 共通の型、式、シンボル、制御構文 AST、診断

既存のルート `sqlbind` は生成コードが利用する行スキャン用ランタイムとして維持します。`templates/sqlbind` はテンプレートのコンパイラ層なので責務は重なりません。

## 共通言語コア

- プリミティブ、レコード、配列、optional、基本 enum を扱う。
- コンポーネントの公開性は名前の大文字・小文字ではなく `export` で決める。
- 変数、フィールド、配列添字、リテラル、型付き関数呼び出し、比較、論理演算、基本算術、null 判定、三項演算を扱う。
- 標準関数と型付き外部関数を認め、汎用のユーザー定義値関数は初期対象外とする。
- `if` / `else` / `else if` と HTML 向け `for` を、各フォーマットが許可した構造位置だけで認識する。
- 宣言した出力型が本文パーサー、挿入規則、生成 API、SQL の件数契約を決める。

### 宣言キーワード

HTMLとSQLの出力宣言を総称してtemplate declarationと呼び、形式別のlowercaseキーワードを必須にします。

```text
export component UserCard(user: User): html {
  ...
}

export statement FindUser(id: int): sql.optional<UserRow> {
  ...
}
```

- `component`: HTML宣言。出力型は`html`のみ。
- `statement`: SQL宣言。出力型は`sql.*`のみ。
- キーワードと出力型が一致しなければコンパイルエラー。
- `export`がなければモジュール内private。
- 共通ASTは`TemplateDecl`、形式別ASTは`HTMLComponentDecl`と`SQLStatementDecl`。

### 名前と大小文字

大小文字とword formを字句クラスの契約として検証します。

| 分類 | 規則 | 例 |
| --- | --- | --- |
| SQLキーワード | UPPERCASE | `SELECT`, `LEFT JOIN`, `IS NULL` |
| DSLキーワード | lowercase | `export`, `component`, `statement`, `if`, `subquery` |
| ユーザー定義シンボル | PascalCase | `UserRow`, `UserCard`, `FindUser` |
| DSL引数・フィールド・変数 | lowerCamelCase | `tenantID`, `minimumAge` |
| SQL schema・table・column・alias | lower_snake_case | `user_accounts`, `created_at` |
| HTML組み込み名 | lowercase / kebab-case | `div`, `aria-label` |
| SQL組み込み関数・型名 | lowercase | `count`, `coalesce`, `integer` |
| 組み込み出力型 | lowercase | `html`, `sql.exec`, `sql.relation` |

誤ったcaseのSQLキーワードは識別子として再解釈せず、正しいUPPERCASEを診断します。ユーザー定義名や識別子を暗黙変換せず、ユーザー定義シンボルはcase-sensitiveで解決します。PostgreSQL初期版はlowercaseのunquoted identifierだけを扱い、mixed-case quoted identifierは延期します。

## HTML 要件

- 静的な要素名・属性名、コンポーネント呼び出し、子要素、子位置の `if` / `for` を扱う。
- テキスト、通常属性、URL 属性、boolean 属性を区別して検証・出力する。
- optional 属性は値がない場合に属性全体を省略する。
- 通常文字列は文脈別に必ずエスケープし、生 HTML は暗黙変換できない `trusted_html` 型だけに許可する。
- 動的タグ名、動的属性名、属性 spread、属性値内のブロック制御は初期対象外とする。
- 生成関数は `io.Writer` へ直接書き込み、完全な DOM を構築しない。

## 明示的なエスケープ制御

次の4つをコンパイラ組み込みの intrinsic として初期機能に含めます。

- `RawHTML(value)`: `trusted_html` を返し、HTML の子要素位置だけで無加工出力する。
- `JsonForScript(value)`: `script_json` を返し、`<script>` 内のデータ位置へ安全な JSON を挿入する。
- `RawCSS(value)`: `trusted_css` を返し、`<style>` 内容だけで無加工出力する。
- `RawJavaScript(value)`: `trusted_javascript` を返し、`<script>` 内容だけで無加工出力する。

`Raw*` は sanitizer ではなく、呼び出し側が内容を信頼済みと表明する危険な境界です。通常文字列からの暗黙変換や、4つの専用型どうしの変換は認めません。専用型を誤った位置へ挿入した場合はコンパイルエラーにします。

`JsonForScript` は `RawJavaScript` と異なり、安全なデータ変換です。静的に JSON 化できる型だけを受け付け、通常の JSON エスケープに加えて `<`、`>`、`&`、U+2028、U+2029 を安全な表現へ変換し、値による `</script>` 終端を防ぎます。入力を JavaScript コードとして評価することはありません。

対象外とするもの:

- 未信頼 HTML・CSS・JavaScript の sanitizer
- CSP nonce / hash の管理
- 任意の動的 CSS 値を安全に構築する API

## SQL 要件

- 通常の式挿入は SQL 文字列化せず、生成設定に応じたプレースホルダーと bind 引数へ変換する。
- `sql.exec`、`sql.one<T>`、`sql.optional<T>`、`sql.many<T>`、`sql.predicate` を公開契約として扱う。
- `where`、predicate group、join、set、order by、insert、returning を構造化し、区切り・空句・パラメータ番号を生成側で管理する。
- SELECT 列と RETURNING 列の形は静的に保つ。解析で件数を証明できなくても、宣言済み API を暗黙変更しない。
- `one` と `optional` の件数を静的に証明できない場合は実行時検査する。
- UPDATE / DELETE の動的 WHERE が空なら実行を拒否する。空の動的 SET も実行前に拒否する。
- 任意の識別子挿入、一般的な SQL ループ、動的な結果列、bulk insert は初期対象外とする。

### 型付きサブクエリ

privateな`statement`は`sql.relation<T>`を返し、FROM/JOINに構造的に埋め込めます。

```text
statement ActiveUsers(
  tenantID: int
): sql.relation<ActiveUserRow> {
  SELECT id, name
  FROM users
  WHERE tenant_id = {tenantID}
    AND active = TRUE
}

export statement ListActiveUsers(
  tenantID: int
): sql.many<ActiveUserRow> {
  SELECT u.id, u.name
  FROM subquery ActiveUsers(tenantID) AS u
  ORDER BY u.id
}
```

- `sql.relation<T>`は初期版ではprivateで、単独の実行APIを生成しない。
- `FROM subquery`と`JOIN subquery`で利用する。
- lower_snake_caseのaliasを必須にする。
- 選択列と`T`を検証し、外側のalias参照も`T`に対して型検査する。
- relationのSQL文字列同士は連結しない。内側ASTを外側ASTへ展開する。
- AST展開後にdialect loweringとplaceholder/Args生成を一度だけ実行する。
- 引数は明示的に渡し、外側aliasの暗黙参照、再帰呼び出し、動的結果列は禁止する。
- scalar subquery、CTE、correlated/LATERAL subquery、recursive CTEは延期する。

### プレースホルダー

テンプレートでは `where id = {id}` のように値式だけを記述します。`$1`や`?`などのbind placeholderをテンプレート作者が直接管理することはありません。文字列リテラルとコメント以外に手書きplaceholderがあればコンパイルエラーにします。

値を評価した時点で、`Args`への追加とplaceholderの出力を一つの操作として行います。動的な`where`で引数数が実行時に変わる場合も、生成コードが採用した方式で番号を管理します。

初期方式:

- `dollar_numbered`: `$1`, `$2`, ...
- `question`: `?`, `?`, ...

placeholder方式はdialectとは別のコード生成オプションですが、デフォルトは選択dialectに従います。

### 低レベル・高レベルAPI

生成SQLは二層のAPIを持ちます。

低レベルAPI:

```go
type Statement struct {
    SQL  string
    Args []any
}

func BuildFindUser(id int) (Statement, error)
```

DBへ接続せず、SQL文字列と`QueryContext` / `ExecContext`へ渡すargsを返します。空の動的`SET`や危険な空`WHERE`もここで検出します。

高レベルAPI:

```go
func FindUser(ctx context.Context, db Queryer, id int) (*User, error)
func UpdateUser(ctx context.Context, db Execer, id int, name string) (sql.Result, error)
```

低レベルbuilderを呼び、`sql.DB`、`sql.Conn`、`sql.Tx`互換の最小interfaceで実行・scan・件数検査を行います。

- `sql.exec`: `ExecContext`
- `sql.one<T>`: `QueryContext`で0行・1行・複数行を検査
- `sql.optional<T>`: `QueryContext`で0行・1行・複数行を検査
- `sql.many<T>`: `QueryContext`で全行をscan

`QueryRowContext`だけでは複数行を検出できないため、静的に最大1行と証明できる場合以外は使用しません。

### Dialect方針

初期dialectはPostgreSQLとします。厳格な型とスキーマ情報を静的テンプレート型・結果型検証の基準にしつつ、初期ASTと構文はportable subsetに保ちます。

初期portable subset:

- SELECT、INSERT、UPDATE、DELETE
- JOIN、WHERE、ORDER BY、LIMIT、OFFSET
- 基本的なRETURNING
- bind値
- 個別placeholderへ展開するIN

SQLiteはPostgreSQL固有機能を広く追加する前の第2dialectとします。追加時には、動的型affinity、STRICT table、日付・時刻・decimal・booleanの保存形式、パラメータ上限、RETURNING制限をdialect loweringで扱います。PostgreSQLの配列、`ANY`、JSONBなどは後続のdialect固有最適化です。

### Dialectの決定時期

Dialectとplaceholder方式はコードジェネレーター実行時に固定し、生成アプリケーションの実行時APIには渡しません。

```text
typed SQL IR
  -> dialect capability検証
  -> dialect lowering
  -> placeholder appender焼き込み
  -> Statement builderとdatabase/sql wrapper生成
```

生成後のAPIはdialect引数、placeholder引数、driverからの自動判定を持ちません。同じアプリケーションで複数dialectを使う場合は、dialectごとに別パッケージまたは別artifactを生成します。

## 生成と実行時要件

- 実行時テンプレート解析、reflection、動的型検索、文字列評価、virtual DOM を使わない。
- HTML は静的部分を直接ストリーム出力し、書き込みエラーを保持する。
- SQL は `context.Context`、DBTX 互換 executor、型付き引数を受ける API を生成する。
- SQL の問い合わせ・scan・件数検査のエラーを保持する。
- 既存の reflection-free、TinyGo / WASM、ランタイム依存分離の方針に従う。

## 推奨実装順

1. `component` / `statement`宣言、大小文字規約、型、式、署名
2. HTML 構造、文脈別エスケープ、コンポーネント、`if` / `for`
3. `RawHTML`、`JsonForScript`、`RawCSS`、`RawJavaScript` と script / style 文脈検証
4. PostgreSQL向けSQL IR、結果契約、静的なSQL文
5. `sql.relation<T>`のFROM/JOIN展開とplaceholder一括生成
6. `Statement` builderと`database/sql`実行wrapper
7. 構造化SQL句と更新・削除の安全ガード
8. SQLite dialect lowering

## レビューが必要な点

- [ ] パッケージ名を `templates/htmlbind` と `templates/sqlbind` で確定する。
- [ ] 初期マイルストーンで indexed `for` を含めるか決める。
- [ ] enum の明示値、SQL フィールド名 annotation、匿名 SQL row 型を初期対象に含めるか決める。
- [ ] Go における optional 値と `sql.optional<T>` の表現を決める。
- [ ] `Raw*` の利用時に警告を出すか、明示許可だけで十分とするか決める。
- [ ] `JsonForScript` が利用する生成 JSON codec と対応可能な型を確定する。
