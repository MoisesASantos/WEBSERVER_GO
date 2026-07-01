-- name: GetOneChirp :one
SELECT * FROM chirp WHERE id = $1; 
