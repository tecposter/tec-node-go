# The Go implementation of the TecPoster Server

## APIs

* post.commit
* post.fetch
* post.list
* post.search
* post.create
* post.edit
* draft.save
* draft.fetch
* draft.list
* user.reg
* user.login
* user.refresh-token

## Data Schema

* post
	* pid
	* uid
	* changed ?
	* pcid (post commit id): [pid]-[index]
* postCommit
	* pcid (post commit id)
	* pid
	* uid
	* created
	* content
* draft
	* pid
	* uid
	* changed
	* content
* [uid]/post
	* pid
	* changed
	* title
* [uid]/draft
	* pid
	* changed
	* title
* [uid]/txn
	* timestamp (int)
	* pcid
* content
	* type: markdown / text / html
	* body: ""


## Request & Response Data Structure

Request

```
{
	"cmd": "post.commit",
	"token": "[token]",
	"params: {
	}
}
```

Response

```
// default
{
	"cmd": "post.commit",
	"status": "ok",
	"data": {
	}
}

// error
{
	"cmd": "post.commit",
	"status": "error",
	"data": {
		"message": "[error message]"
	}
}
```

## User

### user.reg

request

```
{
	"cmd": "user.reg",
	"token": "",
	"params": {
		"email": "zhanjh@126.com",
		"username": "zhanjh",
		"password": "123456789"
	}
}
```

Error

* Error email format
* Username too short - minimum length is 7
* Password too short - minimum length is 7
* Email already exists
* Username already exists

Success

```
{
	"cmd": "user.reg",
	"status": "ok",
	"data": {
		"email": "zhanjh@126.com",
    "uid": "3HrupZrJFPJnhqvhDXDmEb",
    "username": "zhanjh"
	}
}
```

### user.login

request

```
{
	"cmd": "user.login",
	"token": "",
	"params": {
		"email": "zhanjh@126.com",
		"password": "xxx"
	}
}
```

Error

* Email not found
* Incorrect password

Success

```
{
	"cmd": "user.login",
	"status": "ok",
	"data": {
		"token": "***"
	}
}
```

### user.refresh-token

request

```
{
	"cmd": "user.refresh-token",
	"token": "***",
}
```

Error

* Invalid token

Success

```
{
	"cmd": "user.refresh-token",
	"status": "ok",
	"data": {
		"token": "***"
	}
}
```
