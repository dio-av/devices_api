package rest

// var (
// 	// state carries an internal token during the oauth2 workflow
// 	// we just need a non empty initial value
// 	state = "foobar" // Don't make this a global in production.

// 	// the credentials for this API (adapt values when registering API)
// 	clientID     = "" // <= enter registered API client ID here
// 	clientSecret = "" // <= enter registered API client secret here

// 	//  unused in this example: the signer of the delivered token
// 	issuer = "https://accounts.google.com"

// 	// the Google login URL
// 	authURL = "https://accounts.google.com/o/oauth2/v2/auth"

// 	// the Google OAuth2 resource provider which delivers access tokens
// 	tokenURL    = "https://www.googleapis.com/oauth2/v4/token"
// 	userInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"

// 	// our endpoint to be called back by the redirected client
// 	callbackURL = "http://127.0.0.1:12345/api/auth/callback"

// 	// the description of the OAuth2 flow
// 	endpoint = oauth2.Endpoint{
// 		AuthURL:  authURL,
// 		TokenURL: tokenURL,
// 	}

// 	config = oauth2.Config{
// 		ClientID:     clientID,
// 		ClientSecret: clientSecret,
// 		Endpoint:     endpoint,
// 		RedirectURL:  callbackURL,
// 	}
// )

// func configureAPI(api *operations.OauthSampleAPI) http.Handler {
// 	// configure the api here
// 	api.ServeError = errors.ServeError

// 	// Set your custom logger if needed. Default one is log.Printf
// 	// Expected interface func(string, ...interface{})
// 	//
// 	// Example:
// 	api.Logger = log.Printf

// 	api.JSONConsumer = runtime.JSONConsumer()

// 	api.JSONProducer = runtime.JSONProducer()

// 	api.OauthSecurityAuth = func(token string, scopes []string) (*models.Principal, error) {
// 		// This handler is called by the runtime whenever a route needs authentication
// 		// against the 'OAuthSecurity' scheme.
// 		// It is passed a token extracted from the Authentication Bearer header, and
// 		// the list of scopes mentioned by the spec for this route.

// 		// NOTE: in this simple implementation, we do not check scopes against
// 		// the signed claims in the JWT token.
// 		// So whatever the required scope (passed a parameter by the runtime),
// 		// this will succeed provided we get a valid token.

// 		// authenticated validates a JWT token at userInfoURL
// 		ok, err := authenticated(token)
// 		if err != nil {
// 			return nil, errors.New(401, "error authenticate")
// 		}
// 		if !ok {
// 			return nil, errors.New(401, "invalid token")
// 		}

// 		// returns the authenticated principal (here just filled in with its token)
// 		prin := models.Principal(token)
// 		return &prin, nil
// 	}

// 	api.GetAuthCallbackHandler = operations.GetAuthCallbackHandlerFunc(func(params operations.GetAuthCallbackParams) middleware.Responder {
// 		// implements the callback operation
// 		token, err := callback(params.HTTPRequest)
// 		if err != nil {
// 			return middleware.NotImplemented("operation .GetAuthCallback error")
// 		}
// 		log.Println("Token", token)
// 		return operations.NewGetAuthCallbackDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(token)})
// 	})

// 	api.GetLoginHandler = operations.GetLoginHandlerFunc(func(params operations.GetLoginParams) middleware.Responder {
// 		// implements the login operation
// 		login(params.HTTPRequest)
// 		return middleware.NotImplemented("operation .GetLogin has not yet been implemented")
// 	})

// 	api.CustomersCreateHandler = customers.CreateHandlerFunc(func(params customers.CreateParams, principal *models.Principal) middleware.Responder {
// 		// other API endpoint ...
// 		log.Println("hit customer API")
// 		return middleware.NotImplemented("operation customers.Create has not yet been implemented")
// 	})

// 	api.CustomersGetIDHandler = customers.GetIDHandlerFunc(func(params customers.GetIDParams, principal *models.Principal) middleware.Responder {
// 		// other API endpoint ...
// 		log.Println("hit customer API")
// 		return middleware.NotImplemented("operation customers.GetID has not yet been implemented")
// 	})

// 	api.ServerShutdown = func() {}

// 	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
// }

// func login(r *http.Request) string {
// 	// implements the login with a redirection and an access token
// 	var accessToken string
// 	wG := r.Context().Value(ctxResponseWriter).(http.ResponseWriter)
// 	http.Redirect(wG, r, config.AuthCodeURL(state), http.StatusFound)
// 	return accessToken
// }

// func callback(r *http.Request) (string, error) {
// 	// we expect the redirected client to call us back
// 	// with 2 query params: state and code.
// 	// We use directly the Request params here, since we did not
// 	// bother to document these parameters in the spec.

// 	if r.URL.Query().Get("state") != state {
// 		log.Println("state did not match")
// 		return "", fmt.Errorf("state did not match")
// 	}

// 	myClient := &http.Client{}

// 	parentContext := context.Background()
// 	ctx := oidc.ClientContext(parentContext, myClient)

// 	authCode := r.URL.Query().Get("code")
// 	log.Printf("Authorization code: %v\n", authCode)

// 	// Exchange converts an authorization code into a token.
// 	// Under the hood, the oauth2 client POST a request to do so
// 	// at tokenURL, then redirects...
// 	oauth2Token, err := config.Exchange(ctx, authCode)
// 	if err != nil {
// 		log.Println("failed to exchange token", err.Error())
// 		return "", fmt.Errorf("failed to exchange token")
// 	}

// 	// the authorization server's returned token
// 	log.Println("Raw token data:", oauth2Token)
// 	return oauth2Token.AccessToken, nil
// }

// func authenticated(token string) (bool, error) {
// 	// validates the token by sending a request at userInfoURL
// 	bearToken := "Bearer " + token
// 	req, err := http.NewRequest("GET", userInfoURL, nil)
// 	if err != nil {
// 		return false, fmt.Errorf("http request: %v", err)
// 	}

// 	req.Header.Add("Authorization", bearToken)

// 	cli := &http.Client{}
// 	resp, err := cli.Do(req)
// 	if err != nil {
// 		return false, fmt.Errorf("http request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	_, err = io.ReadAll(resp.Body)
// 	if err != nil {
// 		return false, fmt.Errorf("fail to get response: %v", err)
// 	}
// 	if resp.StatusCode != 200 {
// 		return false, nil
// 	}
// 	return true, nil
// }
