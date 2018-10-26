# scan

Port scanner

# usage

(needs go and docker installed)

```
$ go get -u github.com/voutasaurus/scan
$ cd ${GOPATH:-~/go}/src/github.com/voutasaurus/scan
$ ./dockerbuild.sh
++ GOOS=linux
++ go build .
++ docker build -t scan .
Sending build context to Docker daemon  3.017MB
Step 1/3 : FROM scratch
 ---> 
Step 2/3 : COPY scan /
 ---> 73e921262191
Step 3/3 : ENTRYPOINT ["/scan"]
 ---> Running in 731ab35525f2
Removing intermediate container 731ab35525f2
 ---> af1035c2ae62
Successfully built af1035c2ae62
Successfully tagged scan:latest
++ rm scan
++ docker run scan
```

```
Usage of /scan:
  -ip string
    	set IP address to scan (default "127.0.0.1")
  -limit int
    	max bytes to read from each port (default 4096)
  -max int
    	maximum port (default 1024)
  -min int
    	minimum port (default 1)
  -timeout int
    	seconds to scan for (default 1)
  -v	show debug log
```
