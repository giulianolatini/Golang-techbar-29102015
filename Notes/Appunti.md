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

```css
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

> [highlight-line](https://atom.io/packages/highlight-line) This package allows customization of the line selection styles. In my case, I have added a dashed yellow border to the bottom and top my selection. I like the way it looks and helps me determine the range of selection specially at the last line where it could be a partial selection.
>
> [gotags](https://github.com/jstemmer/gotags) There is an awesome tool called `gotags` that is ctags compatible generator for Go Language. It utilizes the power of AST and Parsing classes in the Go standard library to understand and capture all the structure, interfaces, variables and methods names. It generates a much better ctags list than the standard ctags standard tools.

#### A cosa serve _ nell'import???

importa il package richiesto anche se non viene usato esplicitamente nel codice prodotto!!! (Ad esempio il package sqlite3 che include in database/sql il supporto per lo sqlite3)

> [Import for side effect](https://golang.org/doc/effective_go.html#blank_unused)
>
> An unused import like fmt or io in the previous example should eventually be used or removed: blank assignments identify code as a work in progress. But sometimes it is useful to import a package only for its side effects, without any explicit use. For example, during its init function, the net/http/pprof package registers HTTP handlers that provide debugging information. It has an exported API, but most clients need only the handler registration and access the data through a web page. To import the package only for its side effects, rename the package to the blank identifier:
>
> import _ "net/http/pprof" This form of import makes clear that the package is being imported for its side effects, because there is no other possible use of the package: in this file, it doesn't have a name. (If it did, and we didn't use that name, the compiler would reject the program.)
