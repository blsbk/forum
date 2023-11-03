## Forum
* `bbilisbe`
* `aaidana`
* `masalen` 

### Description:
This project involves creating a web forum with features like, post categorization, liking and disliking posts and comments, and post filtering. The data is managed using SQLite, a popular embedded database, and the project encourages optimizing performance with an entity relationship diagram. User authentication allows registration, login sessions with cookie management, and optional UUID usage. Users can create posts and comments, associate categories with posts, and engage with likes and dislikes. A filter system allows users to sort posts by categories.


### Usage
Clone the repository:
```
git clone git@git.01.alem.school:bbilisbe/forum.git
```

#### Run with docker
Go to the downloaded repository: 
```
cd forum
```
Build a program:
```
make build
```
Run a program:
```
make run
```
Open the link in browser
```
https://127.0.0.1:7070
```
Stop a program:
```
make stop
```
#### Run without docker
Go to the downloaded repository: 
```
cd forum
```
Run a program:
```
go run ./cmd/
```
Open the link in browser
```
https://127.0.0.1:7070
```
 