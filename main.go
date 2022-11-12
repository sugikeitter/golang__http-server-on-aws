package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

var tmpl = template.Must(template.New("index.html").
	//Parse("<html><body><div>{{.Time}}</div>{{.Name}}さん</body></html>"))
	ParseFiles("template/index.html"))

var privateIps = myPrivateIps()
var awsAz = awsAzFromMetadata()
var counter = 0

type Result struct {
	Time       string
	Counter    int
	Name       string
	PrivateIps string
	AwsAz      string
	H3Color    string
}

type EcsTaskMeta struct {
	AvailabilityZone string `json:"AvailabilityZone"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	result := Result{
		Time:       currentTime(),
		Counter:    count(),
		Name:       r.FormValue("name"),
		PrivateIps: privateIps,
		AwsAz:      awsAz,
		H3Color:    "33, 119, 218",
	}
	tmpl.Execute(w, result)
}

func handleIcon(w http.ResponseWriter, r *http.Request) {}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", handleIcon)
	fmt.Println(currentTime(), " start!!")
	http.ListenAndServe(":8080", nil)
}

func currentTime() string {
	t := time.Now()
	return t.Format("2006/01/02 15:04:05.000")
}

func myPrivateIps() string {
	netInterfaceAddresses, _ := net.InterfaceAddrs()

	var ips []string
	for _, netInterfaceAddress := range netInterfaceAddresses {
		networkIp, ok := netInterfaceAddress.(*net.IPNet)
		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
			ips = append(ips, networkIp.IP.String())
		}
	}
	return fmt.Sprintf("%s", ips)
}

// TODO IMDS v2
func awsAzFromMetadata() string {
	// ECS のメタデータから取得できれば利用する
	// TODO `az != ""` で条件あってる？
	if az := awsAzFromEcsMeta(); az != "" {
		return az
	}
	// TODO EKS の場合

	// EC2 インスタンス
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := client.Get("http://169.254.169.254/latest/meta-data/placement/availability-zone")
	if err != nil {
		// log
		fmt.Println(currentTime(), " ERROR - http.get from IMDS")
		return "ERROR - http.get"
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(currentTime(), " ERROR - io.ReadAll from IMDS")
		return "ERROR - io.ReadAll"
	}
	return string(body[:])
}

func awsAzFromEcsMeta() string {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := client.Get(os.Getenv("ECS_CONTAINER_METADATA_URI_V4") + "/task")
	if err != nil {
		// log
		fmt.Println(currentTime(), " ERROR - http.get from ECS_CONTAINER_METADATA_URI_V4")
		return ""
	}
	defer resp.Body.Close()
	var ecsTaskMeta EcsTaskMeta
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&ecsTaskMeta); err != nil {
		// log
		fmt.Println(currentTime(), " ERROR - Decode ECS_CONTAINER_METADATA_URI_V4")
		return ""
	}
	return ecsTaskMeta.AvailabilityZone
}

func count() int {
	counter++
	fmt.Println(currentTime(), " Count: ", counter)
	return counter
}
