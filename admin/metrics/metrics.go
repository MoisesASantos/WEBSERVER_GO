package metrics


import (
	"fmt"
	"net/http"
	"github.com/MoisesASantos/WEBSERVER_GO/admin/config"
)

func (cfg *config.ApiConfig) MetricRequestHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	body := fmt.Sprintf(`
	<html>
	  <body>
	    <h1>Welcome, Chirpy Admin</h1>
	    <p>Chirpy has been visited %d times!</p>
	  </body>
	</html>
	`, cfg.FileserverHits.Load())

	w.Write([]byte(body))
}

