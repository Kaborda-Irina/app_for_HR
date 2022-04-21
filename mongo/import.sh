#! /bin/bash

mongoimport -u api_user -p "api1234" --db api_db --collection users --file /docker-entrypoint-initdb.d/users.json
mongoimport -u api_user -p "api1234" --db api_db --collection salaries --file /docker-entrypoint-initdb.d/salaries.json