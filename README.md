# assignment-1

Setup:
1. Create a directory 
```
mddir -p $HOME/go/src/github.com/vineet-garg
cd $HOME/go/src/github.com/vineet-garg 
```

2. Clone the repo 
``` 
git clone https://github.com/vineet-garg/assignment-1.git
```
3. step into the directory 
```
cd assignment-1
```
4. build 
```
go build
```
5. start the server 
```
./assignment-1
```
6. sample APIs:
    ```
    curl -X POST http://localhost:8080/hash -d "password=angryMonkey"
    ```
    ```
    1
    ```
    ```
    curl -XGET http://localhost:8080/hash/1
    ```
    ```
    ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==
    ```
    ```
    curl -X GET http://localhost:8080/stats
    ```
    ```
    {"average":19,"total":1}
    ```
    ```
    curl -X POST http://localhost:8080/shutdown
    ```
7. stop the server using HTTPS request as above  or ```~c```
