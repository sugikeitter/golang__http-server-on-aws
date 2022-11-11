package main

import (
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"time"
)

var tmpl = template.Must(template.New("index.html").
	//Parse("<html><body><div>{{.Time}}</div>{{.Name}}さん</body></html>"))
	ParseFiles("template/index.html"))

type Result struct {
	Time       string
	Name       string
	PrivateIps string
	AwsAz      string
}

func handler(w http.ResponseWriter, r *http.Request) {
	result := Result{
		Time:       currentTime(),
		Name:       r.FormValue("name"),
		PrivateIps: myPrivateIps(),
		AwsAz:      awsAz(),
	}
	tmpl.Execute(w, result)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func currentTime() string {
	t := time.Now()
	return t.Format("2006/01/02 15:04:05.000")
}

func myPrivateIps() string {
	netInterfaceAddresses, _ := net.InterfaceAddrs()

	ips := []string{}
	for _, netInterfaceAddress := range netInterfaceAddresses {
		networkIp, ok := netInterfaceAddress.(*net.IPNet)
		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
			ips = append(ips, networkIp.IP.String())
		}
	}
	return fmt.Sprintf("%s", ips)
}

// TODO IMDS v2
// TODO ECS の場合
func awsAz() string {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/placement/availability-zone")
	if err != nil {
		// log
		return "ERROR - http.get"
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "ERROR - io.ReadAll"
	}
	return string(body[:])
}
