package http

import (
	"bufio"
	"fmt"
)

func WriteSimpleResponse(w *bufio.Writer, status int, body string) {
	statusText := statusText(status)
	fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", status, statusText)
	fmt.Fprintf(w, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(w, "Connection: close\r\n")
	fmt.Fprintf(w, "Content-Type: text/plain; charset=utf-8\r\n")
	fmt.Fprint(w, "\r\n")
	fmt.Fprint(w, body)
}

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 404:
		return "Not Found"
	case 400:
		return "Bad Request"
	case 500:
		return "Internal Server Error"
	default:
		return ""
	}
}
