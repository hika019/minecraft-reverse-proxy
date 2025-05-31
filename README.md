# Minecraft Reverse Proxy

Minecraft Reverse Proxy is a TCP reverse proxy that receives connection requests to Minecraft servers and forwards them to different backend servers based on the domain. It is useful when you want to publish multiple Minecraft servers on a single port.

[X.com: Video in action](https://x.com/hika019/status/1928714602784739624)


## Usage

### 1. Build or Run

```sh
go run main.go
# or
go build -o mc-reverse-proxy
./mc-reverse-proxy
```

### 2. Example Configuration File (config.yml)

```yaml
listen: ":25565"
domains:
    - domain: hogehoge.example.com
        ip: 192.168.1.10
        port: 25565
    - domain: fugafuga.example.com
        ip: 192.168.1.11
        port: 25565
```

- `listen`: Address and port the proxy listens on
- `domains`: List of destination domains, IPs, and ports

### 3. Start the Server

Place `config.yml` in the same directory and start with the above command.

## License

MIT License

---

# Minecraft Reverse Proxy

Minecraft Reverse Proxyは、Minecraftサーバーへの接続要求を受け取り、ドメインごとに異なるバックエンドサーバーへ転送するTCPリバースプロキシです。
複数のMinecraftサーバーを1つのポートで公開したい場合などに利用できます。

## 使い方

### 1. ビルドまたは実行

```sh
go run main.go
# または
go build -o mc-reverse-proxy
./mc-reverse-proxy
```

### 2. 設定ファイル(config.yml)の例

```yaml
listen: ":25565"
domains:
    - domain: hogehoge.example.com
        ip: 192.168.1.10
        port: 25565
    - domain: fugafuga.example.com
        ip: 192.168.1.11
        port: 25565
```

- `listen`: プロキシが待ち受けるアドレスとポート
- `domains`: 転送先のドメイン・IP・ポートのリスト

### 3. サーバーの起動

`config.yml` を同じディレクトリに置き、上記コマンドで起動してください。

## ライセンス

MIT License
