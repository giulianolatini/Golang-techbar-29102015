Notes
=====

Appunti di sviluppo
-------------------

### 1) Attivare go-plus in Atom

dopo aver installato [go-plus](https://atom.io/packages/go-plus) con il comando:

```bash
apm install go-plus
```

vanno installati i [go-tools](http://marcio.io/2015/07/supercharging-atom-editor-for-go-development/)

```bash
go get -u golang.org/x/tools/cmd/...
go get -u github.com/golang/lint/golint
```

Per rendere accessibili al plug-in i tools di go vanno istanziate le variabili oppurtune nel `/etc/profile`:

```bash
# Setting Golang enviroment globally for users
export GOROOT="/usr/local/go"                                                                                                              export GOPATH="/home/wolf/Sviluppo/go"
export PATH=$GOPATH/bin:$GOROOT/bin:/usr/local/bin:$PATH
```
