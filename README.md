# Minecraft Reverse Proxy

Minecraft Reverse Proxy is a TCP reverse proxy that receives connection requests to Minecraft servers and forwards them to different backend servers based on the domain. It is useful when you want to publish multiple Minecraft servers on a single port.

[X.com: Video in action](https://x.com/hika019/status/1928714602784739624)


## Usage

### 1. Build or Run

```sh
go run .
# or
go build -o mc-reverse-proxy
./mc-reverse-proxy
```

### 2. Example Configuration File (config.yml)

```yaml
listen: ":25565"
domains:
  - domain: "sample0.test.com"
    ip: "192.168.1.10"
    port: 25565
    allowed_ips:   # Allowed IPs for this domain (optional)
      - "127.0.0.1"
      - "192.168.1.100"
  - domain: "sample1.test.com"
    ip: "192.168.1.11"
    port: 25566
    # If allowed_ips is omitted, all IPs are allowed

allowed_ips: # Allowed IPs for the entire proxy (optional)
  - "127.0.0.1"
  - "192.168.1.200"
```

- `listen`: Address and port the proxy listens on
- `domains`: List of destination domains, IPs, ports, and (optionally) allowed IPs for each domain
- `allowed_ips`: List of allowed IPs for the entire proxy (optional, if omitted or empty, all IPs are allowed)

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
go run .
# または
go build -o mc-reverse-proxy
./mc-reverse-proxy
```

### 2. 設定ファイル(config.yml)の例

```yaml
listen: ":25565"
domains:
  - domain: "sample0.test.com"
    ip: "192.168.1.10"
    port: 25565
    allowed_ips:   # このドメインへのアクセス許可IP（省略可）
      - "127.0.0.1"
      - "192.168.1.100"
  - domain: "sample1.test.com"
    ip: "192.168.1.11"
    port: 25566
    # allowed_ips: を省略すると全許可

allowed_ips: # 全体のアクセス許可IP（省略可）
  - "127.0.0.1"
  - "192.168.1.200"
```

- `listen`: プロキシが待ち受けるアドレスとポート
- `domains`: 転送先のドメイン・IP・ポート・（任意で）アクセス許可IPリスト
- `allowed_ips`: プロキシ全体へのアクセス許可IPリスト（省略可、空なら全許可）

### 3. サーバーの起動

`config.yml` を同じディレクトリに置き、上記コマンドで起動してください。

## ライセンス

MIT License
