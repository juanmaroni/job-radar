package responses

import "net/http"

func OkHandler(w http.ResponseWriter, r *http.Request, content []byte) string {
	var message string = "200 OK"
	buildResponse(w, http.StatusOK, message, content)

	return message
}

func ForbiddenHandler(w http.ResponseWriter, r *http.Request) string {
	var message string = "403 Forbidden"
	buildResponse(w, http.StatusForbidden, message, nil)

	return message
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) string {
	var message string = "404 Not Found"
	buildResponse(w, http.StatusNotFound, message, nil)
	
	return message
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) string {
	var message string = "405 Method Not Allowed"
	buildResponse(w, http.StatusMethodNotAllowed, message, nil)

	return message
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) string {
	var message string = "500 Internal Server Error"
	buildResponse(w, http.StatusInternalServerError, message, nil)

	return message
}

func buildResponse(w http.ResponseWriter, statusCode int, message string, content []byte) {
	w.WriteHeader(statusCode)

	if content == nil {
		content = []byte(message)
	}

	w.Write(content)
}
