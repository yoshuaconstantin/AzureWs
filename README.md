
## Welcome to AzureWs !

AzureWs is a Golang-based API designed for use with the Azure Reborn Flutter API. As the project develops, new features and updates will be added to the API.
## How it Works

AzureWs provides a range of features, including a full system check to ensure proper authentication, the ability to upload images as bytes, and a secure user login flow. Requests are also secured and require an authentication token.


## Benefits

AzureWs is designed to be easy to read, with a well-planned API flow for better handling. The API is optimized for security and provides a seamless experience for both developers and end-users.

## Collaboration

If you're interested in contributing to this project, please don't hesitate to contact me. We welcome new collaborators and are always open to new ideas.

## Requirement
- Go lang
- PostgreSQL
- PgAdmin
- Any IDE

## What's inside
- Golang
- Restful API
- Docker
## Tech Stack

**Client:** Flutter, Java, *Any

**Server:** Go, Restfull


## Roadmap

- Design and plan the API: Decide on the endpoints, payloads, response formats, and any authentication/authorization requirements.

- Implement secure user authentication and authorization: Use secure and reliable methods for user authentication, such as JWT tokens, and implement proper 

- Implement input validation and error handling: Validate all input data to ensure it meets the API requirements and handle errors gracefully.

- Write clean, easy-to-read code: Follow best practices for Go code style, documentation, and organization. Consider using packages and modules for better code organization.

- Implement API endpoint handlers: Write the code to handle each endpoint's logic and make sure they follow the API's design and requirements.

- Test the API: Create tests to validate the API's functionality, performance, and security. Use automated testing tools and perform manual testing to find and fix any issues.

- Deploy the API: Choose a reliable cloud hosting platform or server to deploy the API and make it publicly available. Configure the server with proper security settings and monitoring.

- Maintain and update the API: Regularly update the API to fix bugs, add features, and improve performance. Keep track of security vulnerabilities and apply security patches as needed.

## How-To-Run

Just run like normal Go.

DB will be included later

```bash
  Go mod tidy
  Go build
  Go run main.go
```
    
## Feedback

If you have any feedback, please reach out to us at joshuaconstantine.k@gmail.com


## Authors

- [@yoshuaconstantin](https://github.com/yoshuaconstantin)


## Example API Reference [WIP]

### All Request need JWT Token as Header

#### Create New Account

```http
  POST /api/add_user
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `username` | `string` | **Required**.  |
| `password` | `string` | **Required**.  |

#### CreateAccount(username, password)

After checking if username already taken or not, this will trigger init function


#### Login

```http
  GET /api/login
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `username`      | `string` | **Required**. |
| `password`      | `string` | **Required**. |

#### Login(username, password)

This login info if succes will generate token to request everything

#### Dashboards

```http
  GET /api/home/dashboards
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `token`      | `string` | **Required**. To Aunth and get UserId |

#### GetDashboards(token)

After request succes user will get their dashboards data

#### Upload Profile Image

```http
  POST /api/home/user/profile/image
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `token`      | `string` | **Required**. To Aunth and get UserId |
| `data`      | `byte` | **Required**. Byte to convert into Image Url|

#### UploadProfileImage(token, data)

This will convert byte image to actual image and store into local dir, and store generated image url to database with it's userId

#### Update Profile Data

```http
  POST /api/home/user/profile
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `token`      | `string` | **Required**. To Aunth and get UserId|
| `data`      | `array` | **Required**. Array with needed value below|
| `data [nickname]`      | `string` | **Not Required**. To save account nickname|
| `data [age]`      | `string` | **Not Required**. To save account age |
| `data [gender]`      | `string` | **Not Required**. To save account gender|
| `data [nation]`      | `string` | **Not Required**. To save account nationality|
| `data [imageUrl]`      | `string` | **Not Required**. To save account Image Url|


#### UpdateProfileData(UserId, Nickname, Age, Gender, ImageUrl, Token)

This API is for updating account profile information

#### Get Profile Data

```http
  GET /api/home/user/profile
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `token`      | `string` | **Required**. To Aunth and get UserId|

#### GetProfileData(Token)

This API is for get all user profile data information


#### Community Chat WebSocket

```http
  Connect to Url /community_chat
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `User_id`      | `string` | **Required**. To save into tables as primary UserId|
| `nickname`      | `string` | **Required**. To know who's the sender of chat|
| `message`      | `string` | **Required**. body message|
| `nation`      | `string` | **Not Required**. To inform nationality of the sender|

#### Community_chat(json{user_id,nickname,message,nation})

Because this service using websocket, users need to connect it first then send the json type
it will also return same model but without user_id in it. it's happen realtime


## Badges


[![AGPL License](https://img.shields.io/badge/license-AGPL-blue.svg)](http://www.gnu.org/licenses/agpl-3.0)

