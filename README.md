# goapi-stac  
### a STAC api written in go with fiber, gorm, elasticsearch, postgres and kafka
#### https://documenter.getpostman.com/view/12888943/VVBXwQnu   
-------

### RUN LOCALLY (localhost:6002):   
```$ make database```  
```$ go build```  
```$ go run app.go```  
    
### TEST LOCALLY:       
```$ make test```
   
### PSQL:
```$ docker exec -it stac-db bash```   
```$ psql```

### RUN IN DOCKER (localhost:6002):  
```$ make database```   
```$ make msg```   
```$ make api```  

---- 
### Developer notes:    
#### Identify failing tests:
```go test -v github.com/jonhealy1/goapi-stac/tests 2>&1 | grep "FAIL\|---"```    
#### Run specific test:   
```go test github.com/jonhealy1/goapi-stac/tests -run TestEsGetCollection```

