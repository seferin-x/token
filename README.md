# token
A (really) simple package for service to service authorization.

Basically a wrapper for github.com/o1egl/paseto with a method for storing arbitrary service keys in token payload (further authentication keys, rbac, user details, etc).

## install

go get github.com/seferin-x/token

## use

Create a new tokenmaker on each service app with a 32 character random symmetric key, usually stored in a secret or env variable. 

`t, err := NewTokenMaker("12345678901234567890123456789012")`

Create your service token with any keys you would like to pass:

`
    payload := make(map[string]interface{})
    payload["service_name"] = "my service"

    t.CreateToken(payload)
`

Send your token with requests in the "authorization" header:
 
`request.Header.Set("authorization", t.Value)`

Use the token maker on receiving service to verify token in middleware for example:

`
	authHeader := ctx.GetHeader(authorizationHeaderKey)

	if len(authHeader) == 0 {
		err := errors.New("no authorization header")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponseJson(err))
	}

	payload, err := t.VerifyToken(authHeader)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponseJson(err))
	}
`

Or use the gin middleware provided in `github.com/seferin-x/token/middleware`:

`router.GET("/api/alive").Use(middleware.GinAuthMiddleware(*server.token))`