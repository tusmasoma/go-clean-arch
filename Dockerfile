# ベースイメージ
FROM golang:1.21.3

# MySQLクライアントをインストール
RUN apt-get update && apt-get install -y default-mysql-client

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係をコピー
COPY go.mod ./
COPY go.sum ./

# 依存関係のインストール
WORKDIR /app
RUN go mod download

# Air をインストール
RUN go install github.com/cosmtrek/air@v1.29.0

WORKDIR /app

# ソースコードをコピー
COPY . .

# エントリポイントスクリプトをコピー
COPY entrypoint.sh /usr/local/bin/

# エントリポイントスクリプトを実行可能にする
RUN chmod +x /usr/local/bin/entrypoint.sh

# コンテナが起動するときに実行されるコマンド (バイナリにしたgolangのファイルを実行)
ENTRYPOINT ["entrypoint.sh"]

# Air を使用してアプリケーションを起動する
WORKDIR /app
CMD ["air", "-c", ".air.toml"]