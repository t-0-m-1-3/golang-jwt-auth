package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/auth0-community/auth0"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	jose "gopkg.in/square/go-jose.v2"
)

func main() {
	// initialize the gorilla/mux router
	r := mux.NewRouter()

	// ON the default page serve the static index page

	r.Handle("/", http.FileServer(http.Dir("./views/")))

	// setup the server to serve static assets as well
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// API is going to consist of three routes
	// /status which will call to make sure the API is up and running
	// /products which will retrieve a list of products that the user can leave feedback on
	// /products/{slug}/feedback will caputre user feedback on products

	r.Handle("/status", StatusHandler).Methods("GET")
	r.Handle("/products", jwtMiddleware.Handler(ProductsHandler)).Methods("GET")
	r.Handle("/products/{slug}/feedback", jwtMiddleware.Handler(AddFeedbackHandler)).Methods("POST")

	// run the application on port 3000
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))

	r.Handle("/get-token", GetTokenHandler).Methods("GET")
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := []byte("vLKf4UMiv31uDHhvSWbrWnQqrRVyxQlP")
		secretProvider := auth0.NewKeyProvider(secret)
		audience := []string{"ofp"}

		configuration := auth0.NewConfiguration(secretProvider, audience, "https://ofp.auth0.com/", jose.HS256)
		validator := auth0.NewValidator(configuration)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// the NotImplemented handler, will handler API end points hit, it will return the msg
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

// create a new product type with a struct containing the relative information
type Product struct {
	Id          int
	Name        string
	Slug        string
	Description string
}

// create a catalog of products

var products = []Product{
	Product{Id: 1, Name: "Scatterplot", Slug: "scatter-plot", Description: "basic usage of scatterplots"},
	Product{Id: 2, Name: "BoxPlot", Slug: "box-plot", Description: "using box plots for distributions"},
	Product{Id: 3, Name: "Time Series Analysis", Slug: "time-series-analysis", Description: "charting series of data along an axis"},
	Product{Id: 3, Name: "Regression Analysis", Slug: "regression-analysis", Description: "plotting regressions using static and series data"},
	Product{Id: 4, Name: "Cluster Analysis", Slug: "cluster-analysis", Description: "plotting clusters of data to analyze different centroids"},
	Product{Id: 5, Name: "Decision Trees", Slug: "decision-trees", Description: "creating a decision tree to analyze probability paths"},
	Product{Id: 6, Name: "Matrix Calculations", Slug: "matrix-calculations", Description: "using packages to peform matrix calculations"},
}

// the status handler will be invoked on the status route
// it will return that the route it up and running
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
})

// the products handler will be called with a GET request on /products
// this will return a list or products available for review

var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// convert the slice to json
	payload, _ := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

// the feedback handler will add either positive or negative feedback to the products
// we would normally save this to a database but for now we fake it
var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var product Product
	vars := mux.Vars(r)
	slug := vars["slug"]

	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}

	w.Header().Set("Content-Type", "applciation/json")
	if product.Slug != "" {
		payload, _ := json.Marshal(product)
		w.Write([]byte(payload))
	} else {
		w.Write([]byte("Product Not Found"))
	}
})

// setup a global string for our secret

var mySigningKey = []byte("secret")

// handlers
var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": true,
		"name":  "John Doe",
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	// create a map to store the claims
	// claims := token.Claims(jwt.MapClaims)

	// Set token claims

	// claims["admin"] = true
	// claims["name"] = "Tom"
	// claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// sign the token with our secret
	tokenString, _ := token.SignedString(mySigningKey)

	// Finally write the token to the browser window
	w.Write([]byte(tokenString))
})

// jwt handler
var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})
