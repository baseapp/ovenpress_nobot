# httpdocs Minimal Server

A minimal static file server serving `index.html` on port `8080`.

## Run the Server

From the `httpdocs` directory:

```bash
cd httpdocs
go run server.go
```

To use a custom port (default is `8080`):

```bash
go run server.go -port 9000
```

Then open:

```text
http://localhost:8080
```

## Optional: Test Behind `go-away`

From the repository root, you can proxy port `8080` through `go-away` using the `minimal` challenge template:

```bash
./nobot.exe --bind :8081 --backend *=http://127.0.0.1:8080 --policy examples/generic.yml --policy-snippets examples/snippets/ --challenge-template minimal
```
