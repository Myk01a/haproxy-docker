package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ip, err := internalIP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<h1>Привіт! Застосунок працює на порту %s з IP-адресою %s!</h1>", os.Getenv("APPPORT"), ip)
}

// Функція для отримання внутрішньої IP-адреси
func internalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			// Виберемо першу внутрішню IP-адресу IPv4
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("не вдалося знайти внутрішню IP-адресу")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APPPORT")), nil))
}

