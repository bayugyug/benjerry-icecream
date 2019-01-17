## benjerry-icecream

* [x] Sample golang rest-api (ben&jerry icecream)


### Pre-Requisite
	
	- Please run this in your command line to ensure packages are in-place.
	  (normally these will be handled when compiling the api binary)
	
```sh


```

### Compile

```sh

     git clone https://github.com/bayugyug/benjerry-icecream.git && cd benjerry-icecream

     git pull && make clean && make

```

### Required Data Preparation

    - Create sample mysql db
	
	- Needs to create the api database and grant the necessary permissions
	
	- Refer the testdata/dump.sql & testdata/dump_test.sql
	
```sh

    #for prod purposes
    mysql -uroot < testdata/dump.sql

    #for dev purposes
    mysql -uroot < testdata/dump_test.sql

```

### List of End-Points-Url


```sh



```


### Mini-How-To on running the api binary

	[x] Prior to running the server, db must be configured first 
	
    [x] The api can accept a json format configuration
	
	[x] Fields:
	
		- http_port = port to run the http server (default: 8989)
		
		- driver    = database details for mysql  (user/pass/dbname/host/port)
		
		- showlog   = flag for dev't log on std-out
		
	[x] Sanity check
	    
		go test ./...
	
	[x] Run from the console

```sh
	./benjerry-icecream --config '{
                "http_port":"8989",
                    "driver":{
                    "user":"benjerry",
                    "pass":"rxxxx",
                    "port":"3306",
                    "name":"benjerry",
                    "host":"127.0.0.1"},
                "showlog":true}'


```

### Notes



### Reference
	

### License

[MIT](https://bayugyug.mit-license.org/)

