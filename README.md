# mintyplex-api documentaion

get onboard plexer
this handles Mintyplex's authentication, authorization, users, storage, etc.

✔️ - done
🚧 - pending/in progress
📋 - backlog

# tasks

## user management 🚧

- Design user profile functionality: This allows users to view and edit their profile information (e.g., name, bio, profile picture). ✔️
- Implement profile update functionality: This allows users to modify their profile information and store the changes securely. ✔️
- Integrate profile picture upload/management (optional): This allows users to upload and manage their profile pictures. ✔️

## product management ✔️

## user dashboard 📋

# documentation for Mintyplex

use: `https://mintyplex-api.onrender.com/` [GET] to check if the API is up. it's no big deal if it doesn't give a `200` at first.

slide to:

- [environment setup](#environment-setup)
- [api endpoints](#api-endpoints)
- [user endpoints](#user-endpoints)
- [product endpoints](#product-endpoints)

## environment setup

In the root directory of the application, create an 'app.env' file with the following;

```
    # MONGO DB SRV Record
MONGODB_SRV_RECORD="mongodb+srv://minty:plexerbaby@mpacluster0.b8ire6p.mongodb.net/?retryWrites=true&w=majority&appName=mpacluster0"

BASE_URL="https://mintyplex-api.onrender.com"

# DB Variables
MONGODB_DATABASE="minty"
USER_COLLECTION="users"
PRODUCT_COLLECTION="products"
USERNAME_COLLECTION="username"

AVATAR_BUCKET="avatars"
AVATAR_COLLECTION="avatars.files"

COVER_BUCKET="covers"
COVER_COLLECTION="covers.files"

JWT_SECRET_KEY="ThisIsMySecretKey"

CLOUDINARY_URL=cloudinary://421383923439176:HlqD4kGIszGt3wfqVe4RzSDjJkI@dymctedcp
CLOUDINARY_API_SECRET="HlqD4kGIszGt3wfqVe4RzSDjJkI"
CLOUDINARY_CLOUD_NAME="dymctedcp"
CLOUDINARY_API_KEY="421383923439176"
```

## API-Endpoints

## user endpoints

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

### get all users - `GET`
    `https://mintyplex-api.onrender.com/api/v1/user/users/`



### get a user profile - `GET`

    https://mintyplex-api.onrender.com/api/v1/user/profile/:id

this route takes the user's id as a parameter. e.g `https://mintyplex-api.onrender.com/api/v1/user/profile/660c3aafe22f82232121bbd9`

response example-

```
   {
  "error": false,
  "message": "User Profile",
  "user": {
    "wallet_address": "xion186n0xxs96rzvnrc8ld66zkywc54xvta0mc5ewx5yvx8tde4xvals8xxekr",
    "bio": "came 2 fuck shit up n l(i)(ea)ve",
    "x_link": "www.x.com/ephraimxblack",
    "avatar": "https://res.cloudinary.com/dymctedcp/image/upload/v1717613524/tgpmqpef4u0jivgzy6d3.jpg",
    "products": [
      {
        "ID": "6641236c278f995dce8acf2c",
        "UserId": "xion186n0xxs96rzvnrc8ld66zkywc54xvta0mc5ewx5yvx8tde4xvals8xxekr",
        "CoverImage": "",
        "Name": "rugged dreamer",
        "Price": 619,
        "Discount": 69,
        "Description": "the lands are blue, blood is purple, buildings are leaves, mansions are dunghills",
        "Categories": "art",
        "Quantity": 666,
        "Tags": [
          "ancient",
          "forever"
        ],
        "CreatedAt": 1715544940,
        "UpdatedAt": 1715544940
      },
      {
        "ID": "66412419278f995dce8acf2f",
        "UserId": "xion186n0xxs96rzvnrc8ld66zkywc54xvta0mc5ewx5yvx8tde4xvals8xxekr",
        "CoverImage": "",
        "Name": "Alas Alan",
        "Price": 700,
        "Discount": 9,
        "Description": "Lost brother is finally found and reconciled",
        "Categories": "photography",
        "Quantity": 50,
        "Tags": [
          "humans",
          "affection"
        ],
        "CreatedAt": 1715545113,
        "UpdatedAt": 1715545113
      },
}
```

### edit a user profile - `PUT`

    `https://mintyplex-api.onrender.com/api/v1/user/profile/:id`

this route enables users to edit already existing information
request example-
```
    {
            "email": "seanP@gmail.com",
            "bio": "i love crypto",
            "x_link": "x.com/seanP",
    }
```


## product endpoints

### add a product - `POST`

        `https://mintyplex-api.onrender.com/api/v1/product/:id`

`*the ':id' is the user's id which is their wallet address`. to parse/pass this request, recommended is Postman, use `form-data` under the `Body`, key and value example as follows

request example-
    {
  name: Rugged Ruddy
  price: 300.99
  discount: 99.0
  description: This is the tale of the burgundy.
  categories: art
  categories: forever
  quantity: 10,
  tags: crip
  tags: lrip
  //to select image, select key type as file and select image
  image: <selected file>
}

### get all products - `GET`
        `https://mintyplex-api.onrender.com/api/v1/product/`
this route all gets all existing products in the database

### get one product - `GET`
        `https://mintyplex-api.onrender.com/api/v1/product/:id`
this route gets one existing product in the database


### update product - `PUT`
        `https://mintyplex-api.onrender.com/api/v1/product/:id/:uid`
example `https://mintyplex-api.onrender.com/api/v1/product/66198ca9102c15fb33d65490/bion186n0xxs96rzvnrc8ld66zkywc54xvta0mc5ewx5yvx8tde4xvals8xxekrzs`
this route updates existing product in the database
request example-

```
    {
  "name": "Rugged Ruddy",
  "price": 300.99,
  "discount": 99.0,
  "description": "This is the tale of the burgundy.",
  "categories": ["axeless", "dull"],
  "quantity": 10,
  "tags": ["crip", "lrip"]
}
```

### reserve username - `POST`
    `https://mintyplex-api.onrender.com/api/v1/product/cover/66460cd9cfd1da94c95c5238`

this route gets the cover of every product cover. `66460cd9cfd1da94c95c5238` is the `id` of the product.

### reserve username - `POST`
    https://mintyplex-api.onrender.com/api/v1/reserve/

request example-

```
    {
        "username": "koxdy",
        "email": "epphraim@gmail.com"
    }
```