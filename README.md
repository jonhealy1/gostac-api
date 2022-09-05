# go-stac-api-postgres  
### a STAC api written in go with fiber, gorm and postgres   
-------
### PROGRESS:  
- Sep.4/2022 - only collection CRUD routes are working at this time   
- Sep.5/2022 - item CRUD routes added 
- added item collection route (still needs better formatting) 

### TODO: 
- Search route   
- import ENV variables   
- add swagger docs  
- add tests!  
  
### RUN LOCALLY (localhost:6002):  
- Public postman collection available here: https://www.getpostman.com/collections/a16d074dcd961569278b 

```$ make database```  
```$ go build```  
```$ go run app.go```  
   