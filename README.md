# mintyplex-api documentaion
get onboard plexer
this handles Mintyplex's authentication, authorization, users, storage, etc.

- api: user facing part of the application

# documentation for Mintyplex
use: https://mintyplex-api.onrender.com/ [GET] to check if the API is up. it's no big deal if it doesn't give a `200` at first.

slide to:

- [Environment Setup](#Environment-Setup)
- [API Endpoints](#API-Endpoints)

## environment setup

In the root directory of the application, create an 'app.env' file with the following;

    # MONGO DB SRV Record
MONGODB_SRV_RECORD="mongodb+srv://minty:plexerbaby@mpacluster0.b8ire6p.mongodb.net/?retryWrites=true&w=majority&appName=mpacluster0"

    # DB Variables
MONGODB_DATABASE="minty"
USER_COLLECTION="users"

AVATAR_BUCKET="avatars"
AVATAR_COLLECTION="avatars.files"


## API-Endpoints

### add user profile - `POST`
    https://mintyplex-api.onrender.com/api/v1/user/profile/
    
request example-
```
    {
    "wallet_address": "xion186n0xxs96rzvnrc8ld66zkywc54xvta0mc5ewx5yvx8tde4xvals8xxekr",
    "email": "john.doe@example.com",
    "x_link": "x.com/elstsp",
    "bio": "i love crypto"
    }
```


### get a user profile - `GET`
    https://mintyplex-api.onrender.com/api/v1/user/profile/:id
this route takes the user's id as a parameter. e.g `https://mintyplex-api.onrender.com/api/v1/user/profile/660c3aafe22f82232121bbd9`

response example-

    {
        "error": false,
        "message": "User Profile",
        "user": {
                "id": "660c3aafe22f82232121bbd9",
                "wallet_address": "xion186n0xxs96rzvnrc8ld66zkywc54xvta0mc5ewx5yvx8tde4xvals8xxekr",
                "email": "seanP@gmail.com",
                "avatar": "/api/v1/user/avatar/",
                "bio": "i love crypto",
                "x_link": "x.com/seanP",
                "created_at": 1712077487,
                "updated_at": 1712077487
        }
    }

### edit a user profile - `PUT`
    https://mintyplex-api.onrender.com/api/v1/user/profile/:id
this route enables users to edit already existing information
request example- 

    {
            "email": "seanP@gmail.com",
            "bio": "i love crypto",
            "x_link": "x.com/seanP",
    }