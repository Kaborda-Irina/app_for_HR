db = db.getSiblingDB('api_db');
db.createUser(
    {
        user: 'api_user',
        pwd: 'api1234',
        roles: [{role: 'readWrite', db: 'api_db'}],
    },
);
db.createCollection('users', 'salaries');
