Light command line app for manipulating wav audio file

### Feature
* lufs  
Detect loudness of audio content, see [loudness](https://en.wikipedia.org/wiki/EBU_R_128)
* normalize  
set loudness to target typical lufs (-16) see [https://tech.ebu.ch/publications/r128/](https://) 
* split  
split multichannel (input) file to multiple single channel file

### QuickStart
#### run or build

``` sh
git clone htpps://github.com/odit-bit/soundfreak
cd soundfreak
```
```go
go build -o ./app ./cli 
```
or just use `go run`

``` go
go run ./cli lufs ./myAwesomeMaster.wav //get lufs
//or
go run./cli --help //for available command
```



