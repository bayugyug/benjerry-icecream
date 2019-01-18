## benjerry-icecream

* [x] Sample golang rest-api ( ben&jerry icecream )


### Pre-Requisite
	
	- Please run this in your command line to ensure packages are in-place.
	  (normally these will be handled when compiling the api binary)
	
		go get -u -v github.com/go-chi/chi
		go get -u -v github.com/go-chi/chi/middleware
		go get -u -v github.com/go-chi/cors
		go get -u -v github.com/go-chi/jwtauth
		go get -u -v github.com/go-chi/render
		go get -u -v github.com/dgrijalva/jwt-go
		go get -u -v github.com/go-sql-driver/mysql


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
	
	- Refer the testdata/*.sql
	
```sh

    #prod purposes
    mysql -uroot < testdata/db_prod.sql
    mysql -uroot < testdata/dump_prod.sql

    #dev purposes
    mysql -uroot < testdata/db_dev.sql
    mysql -uroot < testdata/dump_dev.sql

```

### List of End-Points-Url


```bash

		#User Create
		curl -v -X POST 'http://127.0.0.1:8989/v1/api/user'  -d '{
				  "user":"ben@jerry.com",
				  "pass":"8888"
				  }'

				  {"Code":200,"Status":"Create successful","Otp":"06370","OtpExpiry":"2019-01-18 10:47:36"}

		#User Create (Invalid parameters)
		curl -v -X POST 'http://127.0.0.1:8989/v1/api/user'  -d '{
				  "user":"ben@jerry.com",
				  "pass":"888"
				  }'

				  {"Code":206,"Status":"User/Pass must at least 4 characters"}

		#User OTP
		curl -v -X POST 'http://127.0.0.1:8989/v1/api/otp'     -d '{"user":"ben@jerry.com","otp":"06370"}'
				
				  {"Code":200,"Status":"Otp successful"}

		#User OTP (Invalid parameters)
		curl -v -X POST 'http://127.0.0.1:8989/v1/api/otp'     -d '{"user":"ben@jerry.com","otp":"x09733"}'
				
				  {"Code":403,"Status":"Otp mismatch or invalid"}

						  

		#User Login
		curl -v -X POST 'http://127.0.0.1:8989/v1/api/login'  -d '{
				  "user":"ben@jerry.com",
				  "pass":"8888"
				  }'

				  {"Code":200,"Status":"Login Successfull","Token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTAzODUxNTAsInVzZXIiOiJiZW5AamVycnkuY29tIiwidXVpZCI6IjBkOGI2NDZmYzExN2QxNTM1NWMzZDM0MmVjZWE2MzdiNyJ9.ByEHyvvP7b_TaAmDKlrqlH7hWE3tEQe8dE3SNrBH0xw"}


		#User Login (Invalid parameters)
		curl -v -X POST 'http://127.0.0.1:8989/v1/api/login'  -d '{
				  "user":"ben@jerry.comx",
				  "pass":"8888"
				  }'

				  {"Code":404,"Status":"Record not found"}
					  
				  
	
		#User Delete (Invalid parameters)
		curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X DELETE 'http://127.0.0.1:8989/v1/api/user/ben@jerry.comx' 

			      {"Code":403,"Status":"Invalid token"}

		#User Delete
		curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X DELETE 'http://127.0.0.1:8989/v1/api/user/ben@jerry.com' 
			  
			     {"Code":200,"Status":"Delete successful"}

		
		#Icecream Create
		curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X POST 'http://127.0.0.1:8989/v1/api/icecream' -d '{"name": "01-Vanilla Toffee Bar Crunch",
				"image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
				"image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
				"description": "Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
				"story": "Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars we’re using as part of our commitment to source Fairtrade Certified and non-GMO ingredients. We love it and know you will too!",
				"sourcing_values": [
					"Fairtrade",
					"Responsibly Sourced Packaging",
					"Caring Dairy"
				],
				"ingredients": [
					"vegetable oil (canola",
					"safflower",
					"and/or sunflower oil)",
					"guar gum",
					"carrageenan"
				],
				"allergy_info": "may contain wheat, peanuts and other tree nuts",
				"dietary_certifications": "Kosher"}'


				{"Code":200,"Status":"Create successful","ProductID":"154"}

		#Icecream Delete
		curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X DELETE 'http://127.0.0.1:8989/v1/api/icecream/154' 
				
				 {"Code":200,"Status":"Delete successful"}

		#Icecream Update		 
		curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/icecream/154' -d '{"name": "01-Vanilla Toffee Bar Crunch",
				"image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
				"image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
				"description": "updated Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
				"story": "updated Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars we’re using as part of our commitment to source Fairtrade Certified and non-GMO ingredients. We love it and know you will too!",
					"sourcing_values": [
					"updated Fairtrade",
					"xResponsibly Sourced Packaging",
					"yCaring Dairy"
				],
				"ingredients": [
					"yvegetable oil (canola",
					"xsafflower",
					"and/or sunflower oil)",
					"guar gum",
					"zcarrageenan"
				],
				"allergy_info": "--updated may contain wheat, peanuts and other tree nuts",
				"dietary_certifications": "--updated Kosher"}'
				 
				{"Code":200,"Status":"Update successful"}
	
		#Icecream Get
		curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X GET 'http://127.0.0.1:8989/v1/api/icecream/154'

				{
				  "Code": 200,
				  "Status": "Record found",
				  "Result": {
					"productId": "154",
					"name": "01-Vanilla Toffee Bar Crunch",
					"story": "updated Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars we’re using as part of our commitment to source Fairtrade Certified and non-GMO ingredients. We love it and know you will too!",
					"description": "updated Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
					"image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
					"image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
					"dietary_certifications": "--updated Kosher",
					"allergy_info": "--updated may contain wheat, peanuts and other tree nuts",
					"sourcing_values": [
					  "updated Fairtrade",
					  "xResponsibly Sourced Packaging",
					  "yCaring Dairy"
					],
					"ingredients": [
					  "and/or sunflower oil)",
					  "guar gum",
					  "xsafflower",
					  "yvegetable oil (canola",
					  "zcarrageenan"
					]
				  }
				}		
		
```


### Mini-How-To on running the api binary

	[x] Prior to running the server, db must be configured first 
	
    [x] The api can accept a json format configuration
	
	[x] Fields:
	
		- http_port = port to run the http server (default: 8989)
		
		- driver    = database details for mysql  (user/pass/dbname/host/port)
		
		- dump_file = json file for prep-data     (loaded to db during server start-up)

		- showlog   = flag for dev't log on std-out
		
	[x] Sanity check
	    
		go test ./...
	
	[x] Run from the console

```sh
	./benjerry-icecream --config '{
                "http_port":"8989",
		"dump_file":"./testdata/icecream.json",
                    "driver":{
                    "user":"benjerry",
                    "pass":"rxxxx",
                    "port":"3306",
                    "name":"benjerry",
                    "host":"127.0.0.1"},
                "showlog":true}'


```

### Notes

	[x] During http-server start-up, it will load the test-data from the dump_file parameter.
	

### Reference
	

### License

[MIT](https://bayugyug.mit-license.org/)

