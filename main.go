package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {
	apiBasePath := "/auth"
	websiteBasePath := "/auth"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			// https://try.supertokens.com is for demo purposes. Replace this with the address of your core instance (sign up on supertokens.com), or self host a core.
			ConnectionURI: "http://localhost:3567",
			// APIKey: <API_KEY(if configured)>,
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "supertoken-core",
			APIDomain:       "http://localhost:3000",
			WebsiteDomain:   "http://localhost:8080",
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			session.Init(nil),
		},
	})

	if err != nil {
		panic(err.Error())
	}

	http.ListenAndServe(":3000", corsMiddleware(
		supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter,
			r *http.Request) {
			// TODO: Handle your APIs..
			if r.URL.Path == "/signup" && r.Method == "POST" {
				handleSignup(rw, r)
			} else if r.URL.Path == "/signin" && r.Method == "POST" {
				fmt.Println("/Signin API run")
				// handleSignin(rw, r)
			} else {
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte("Not Found"))
			}

		}))))

}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			// we add content-type + other headers used by SuperTokens
			response.Header().Set("Access-Control-Allow-Headers",
				strings.Join(append([]string{"Content-Type"},
					supertokens.GetAllCORSHeaders()...), ","))
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
	})
}

func handleSignup(rw http.ResponseWriter, r *http.Request) {
	// Decode JSON body
	var signupData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&signupData)
	if err != nil {
		fmt.Println("Error in decode ", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a user with SuperTokens
	createUserResponse, err := emailpassword.SignUp("public", signupData.Email, signupData.Password)
	if err != nil {
		fmt.Println("Error Create a user with SuperTokens: ", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the result
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(createUserResponse)
}

// func handleSignin(rw http.ResponseWriter, r *http.Request) {
// 	var userData map[string]string
// 	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
// 		fmt.Println("Error in decode request data, ", err)
// 		http.Error(rw, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Creating JSON body to send to SuperTokens for signin
// 	signinBody, _ := json.Marshal(map[string]interface{}{
// 		"email":    userData["email"],
// 		"password": userData["password"],
// 	})

// 	// Make request to SuperTokens
// 	resp, err := http.Post("http://localhost:3567/recipe/emailpassword/signin", "application/json", bytes.NewBuffer(signinBody))
// 	if err != nil || resp.StatusCode != http.StatusOK {
// 		fmt.Println("Error in ")
// 		http.Error(rw, "Failed to sign in user", http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprintf(rw, "Signin successful")
// }
