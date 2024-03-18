# gostac-api
### a STAC api written in go with fiber, gorm, elasticsearch and postgres  
#### https://documenter.getpostman.com/view/12888943/VVBXwQnu   
-------

### RUN POSTGRES API LOCALLY (localhost:6002):   
```$ docker compose up database```  
```$ cd pg-api```   
```$ go build```  
```$ go run app.go```  

### RUN ELASTICSEARCH API LOCALLY (localhost:6003):   
```$ docker compose up elasticsearch```  
```$ cd es-api```   
```$ go build```  
```$ go run app.go```  
    
### TEST LOCALLY:       
```$ make test```
   
### PSQL:
```$ docker exec -it stac-db bash```
```$ psql```

### RUN IN DOCKER (localhost:6002):  
```$ make database```  
```$ make api```  

---- 
### Developer notes:    
#### Identify failing tests:
```go test -v github.com/jonhealy1/goapi-stac/tests 2>&1 | grep "FAIL\|---"```    
#### Run specific test:   
```go test github.com/jonhealy1/goapi-stac/tests -run TestEsGetCollection```

