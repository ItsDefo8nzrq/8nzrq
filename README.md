<div id="header" align="center">
<h1>💥GoBoom💥</h1>
  <p><b>A simple tool to "DDoS" a webserver via multiple threads each using a different proxy</b></p>
</div>



## Features 📁

- Auto harvest of multiple sources for free proxy
- Possibility to add file with your own proxy (ip:port format)
- Possibility to set your own threads limit
- Possibility to use without proxy


## Usage/Examples 📖
### Help
```shell
go run GoBoom.go -h

usage: GoBoom [-h|--help] -d|--domain "<value>" [-t|--threads "<value>"]
              [-p|--proxy-file "<value>" [-p|--proxy-file "<value>" ...]]
              [-x|--proxy-mult <integer>] [-m|--mode <integer>]

              Boom some website by proxy

Arguments:

  -h  --help        Print help information
  -d  --domain      Domain to boom
  -t  --threads     Number of core to use. Default: max
  -p  --proxy-file  Proxy file(s), separate with a ',' each files. Format of
                    file(s) must be ip:port. Default: []
  -x  --proxy-mult  You can multiply the working proxys detected with this
                    option. Default: 12
  -m  --mode        Mode of attack, 1 for pass all traffic trough proxy, 2
                    don't use proxy. Default: 1
```
### With Golang
```shell
git clone https://github.com/ugomeguerditchian/GoBoom
cd GoBoom
go run GoBoom.go -d example.com
```
### With binaries
Open it in a terminal and add your args
```shell
GoBoom.exe -d example.com 

```
### Add your own proxy file
You can add your own file containing proxy :
```shell
    go run GoBoom.go -d example.com -p C:\myfile1.txt,C:\myfile2.txt

```

### Use GoBoom without proxy
You have to specify the number of threads you want to use
```shell
    go run GoBoom.go -d example.com -t 100 -m 2

```
In this case GoBoom will use all core of your CPU, you have to adapt it with the number of threads you want to use


## Authors 🖋

- [@ugomeguerditchian](https://github.com/ugomeguerditchian)

## Contributors 🖊

- [@lisandro-git](https://github.com/lisandro-git)

