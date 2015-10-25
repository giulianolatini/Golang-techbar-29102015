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
export GOROOT="/usr/local/go"
export GOPATH="/home/wolf/Sviluppo/go"
export PATH=$GOPATH/bin:$GOROOT/bin:/usr/local/bin:$PATH
```

Riavviato il pc, così che all'avvio della sessione grafica le variabili saranno leggibili avviando Atom dall'interfaccia grafica e non solo dalla cli (ereditando così l'ambiente della shell di avvio).

[go-rename](https://atom.io/packages/go-rename)

[file-icons](https://atom.io/packages/file-icons) Configurazione `styles.less` nella cartella di configurazione di Atom

```less
 // style the background color of the tree view
 .tree-view {
     font-family: "Source Code Pro for Powerline";
     font-size: 12px;
 }

 .list-group li:not(.list-nested-item),
 .list-tree li:not(.list-nested-item),
 .list-group li.list-nested-item > .list-item,
 .list-tree li.list-nested-item > .list-item {
     line-height:18px;
 }

 .list-group .selected:before,
 .list-tree .selected:before {
     height:18px;
 }

 .list-tree.has-collapsable-children .list-nested-item > .list-tree > li,
 .list-tree.has-collapsable-children .list-nested-item > .list-group > li {
     padding-left:12px;
 }
```

[Zeal](https://zealdocs.org/download.html) - browser offline per documentazione, software opensource ispirato da Dash, disponibile per Linux e Windows. Installabile usando i seguenti comandi di UbuntuLinux e derivati:

```bash
sudo add-apt-repository ppa:zeal-developers/ppa
sudo apt-get update
sudo apt-get install zeal
```

[Atom plugin-Zeal](https://atom.io/packages/atom-zeal)

```bash
apm install atom-zeal
```

[Dash](https://kapeli.com/dash) - browser offline per documentazione disponibile solo per OSX

[Atom plugin-Dash](https://atom.io/packages/dash)

```bash
apm install dash
```

[highlight-line](https://atom.io/packages/highlight-line) This package allows customization of the line selection styles. In my case, I have added a dashed yellow border to the bottom and top my selection. I like the way it looks and helps me determine the range of selection specially at the last line where it could be a partial selection.

[gotags](https://github.com/jstemmer/gotags) There is an awesome tool called `gotags` that is ctags compatible generator for Go Language. It utilizes the power of AST and Parsing classes in the Go standard library to understand and capture all the structure, interfaces, variables and methods names. It generates a much better ctags list than the standard ctags standard tools.
