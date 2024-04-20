## Simple Text share API

- Able to add Text upto 10MB
- Able to get random ID associated with it
- Able to retrive the original Text


- User Login
- Able to save text
- Check how many times it get retrived
- Restrict to particular user only
- expire the text


## Things that need to implement

- GORM
- REDIS
- RATE LIMITER ALGORITHM
- AUTHENTICATION, 3rd Party as well



## Routes needed

POST /api/public/add
BODY JSON
    {"text": "", encrypted : bool}
GET /api/public/get?id={id}&metaInfo=true


