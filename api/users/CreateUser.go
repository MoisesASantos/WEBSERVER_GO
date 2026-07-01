package users

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/MoisesASantos/WEBSERVER_GO/admin/config"
)

type returnVals struct {
    ID uuid `json:"id"`
    CreatedAt time.time `json:"created_at"`
    UpdatedAt time.time `json:"updated_at"`
    Email string `json:"email"`
}

func UsersRequestHandler(w http.ResponseWriter, r *http.Request) error {

	//Get the request body
	type requestbody struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	data := requestbody{}
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	//create the user using a function created per sqlc
	apiconfig := config.ApiConfig{}
	userResult, err := apiconfig.Db.CreateUSer(r.Context(), data.Email)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	//now create a response

	payload := returnVals{
		ID: 		userResult.ID,
		CreatedAt: 	userResult.CreatedAt,
		UpdatedAt: 	userResult.UpdatedAt,
		Email:		userResult.Email,
	}

	err = config.RespondWithJSON(w, 200, payload)
	if err != nil {
		return err
	}
	return nil
}
