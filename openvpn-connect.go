package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/abh/geoip"
)

type Client struct {
	Name      string
	Vpn_ip    string
	Real_ip   string
	Country   string
	Real_port string
	Connected string
	Upload    string
	Download  string
}

type Data struct {
	Clients []Client
	Update  string
}

const page = `<!DOCTYPE html>
        <html lang="it">
          <head>
            <meta charset="utf-8">
            <meta http-equiv="X-UA-Compatible" content="IE=edge">
            <meta name="viewport" content="width=device-width, initial-scale=1">
            <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
            <meta name="description" content="">
            <meta name="author" content="">
            <title>VPN Monitor</title>
            <!-- Latest compiled and minified CSS -->
            <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
            <!-- Optional theme -->
            <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap-theme.min.css">
          </head>
          <body role="document">
            <div class="container theme-showcase" role="main">

              <!-- Main jumbotron for a primary marketing message or call to action -->
                <div class="jumbotron text-center">
                  <h1>GoVPN Monitor</h1>
                </div>
              <div class="page-header">
                <h3>Update: {{ .Update }}</h3>
              </div>
              <div class="row">
                <div class="col-md-12">
                  <table class="table-striped table">
                    <thead>
                      <tr>
                        <th>Hostname</th>
                        <th>VPN IP Address</th>
                        <th>Real IP Address</th>
                        <th>Real port</th>
                        <th>Country</th>
                        <th>Upload</th>
                        <th>Download</th>
                        <th>Connected Since</th>
                      </tr>
                    </thead>
                    <tbody>
                    {{range .Clients}}
                        <tr>
                            <td><strong>{{.Name}}</strong></td>
                            <td>{{.Vpn_ip}}</td>
                            <td>{{.Real_ip}}</td>
                            <td>{{.Real_port}}</td>
                            <td>{{.Country}}</td>
                            <td style="color: red;">{{.Upload}}</td>
                            <td style="color: green;">{{.Download}}</td>
                            <td>{{.Connected}}</td>
                        </tr>
                    {{end}}
                    </tbody>
                  </table>
                </div>
              </div>
            </div> <!-- /container -->
            <!-- Bootstrap core JavaScript
            ================================================== -->
            <!-- Placed at the end of the document so the pages load faster -->
            <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
            <!-- Latest compiled and minified JavaScript -->
            <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js"></script>
          </body>
        </html> `

func main() {
<<<<<<< HEAD
	//Vars
	tmp := make(map[string]Client)
	Session := make([]Client, 0)

	//Command line params
	serverPtr := flag.String("server", "127.0.0.1", "IP Server VPN")
	portPtr := flag.String("port", "5555", "Port server VPN")
	outputPtr := flag.String("file", "./vpn_page.html", "Output file name")
	flag.Parse()

	//Server connection
	host := fmt.Sprint(*serverPtr, ":", *portPtr)
	conn, err := net.Dial("tcp", host)
	checkError(err)
	defer conn.Close()

	//Clear first data
	buf := make([]byte, 256)
	_, err = conn.Read(buf[:])
	if err != nil {
		fmt.Println("Init error!")
	}

	//Get data
	writer(conn, "status 3\n")
	status3 := strings.Split(reader(conn), "\r\n")
	//fmt.Println("status: ", status3)
	writer(conn, "exit\n")

	//Parse data and make html
	update := Parser(status3, tmp)
=======
	serverPtr := flag.String("server", "192.168.0.106", "IP Server VPN")
	portPtr := flag.String("port", "5555", "Port server VPN")
	outputPtr := flag.String("file", "./vpn_page.html", "Output file name")
	flag.Parse()
	host := fmt.Sprint(*serverPtr, ":", *portPtr)
	tmp := make(map[string]Client)
	Session := make([]Client, 0)
	conn, err := net.Dial("tcp", host)
	checkError(err)
	defer conn.Close()
	writer(conn, "state\n")
	_ = reader(conn)
	writer(conn, "status 3\n")
	status3 := strings.Split(reader(conn), "\r\n")
	//fmt.Println("status3: ", status3)
	writer(conn, "exit\n")
	update := Parser3(status3, tmp)
>>>>>>> af6e5f1124a91023500ef623f73c6da469c51d12
	Session = map_to_slice(tmp)
	fd, err := os.Create(*outputPtr)
	t, _ := template.New("vpn").Parse(page)
	data := Data{Clients: Session, Update: update}
	err = t.Execute(fd, data)
	checkError(err)
}

<<<<<<< HEAD
func Parser(lines []string, session map[string]Client) string {
	//Vars
	var c Client
	var update string
	file := "/usr/share/GeoIP/GeoIP.dat"

	//Init geoip
=======
func Parser3(lines []string, session map[string]Client) string {
	file := "/usr/share/GeoIP/GeoIP.dat"
>>>>>>> af6e5f1124a91023500ef623f73c6da469c51d12
	gi, err := geoip.Open(file)
	if err != nil {
		fmt.Printf("Could not open GeoIP database please install in /usr/share/GeoIP/\n")
	}

<<<<<<< HEAD
	//Parse Client
	for i := range lines {
		if getLineTitle(lines[i]) == "TIME" {
			update = getLineData(lines[i])[0]
		}
		if getLineTitle(lines[i]) == "CLIENT_LIST" {
			data := getLineData(lines[i])
			c.Name = data[0]
			c.Vpn_ip = data[2]
			tmp := strings.Split(data[1], ":")
			c.Real_ip = tmp[0]
			c.Real_port = tmp[1]
			upload, _ := strconv.ParseInt(data[5], 10, 32)
			download, _ := strconv.ParseInt(data[4], 10, 32)
=======
	var update string
	for i := range lines {
		var c Client
		tmp := strings.Split(lines[i], "\t")
		if i == 1 {
			update = tmp[1]
		}
		if i > 2 && len(tmp) == 9 {
			c.Name = tmp[1]
			c.Vpn_ip = tmp[3]
			tmp1 := strings.Split(tmp[2], ":")
			c.Real_ip = tmp1[0]
			c.Real_port = tmp1[1]
			upload, _ := strconv.ParseInt(tmp[5], 10, 32)
			download, _ := strconv.ParseInt(tmp[4], 10, 32)
>>>>>>> af6e5f1124a91023500ef623f73c6da469c51d12
			c.Upload = fmt.Sprint(upload/1000, " Kb")
			c.Download = fmt.Sprint(download/1000, " Kb")
			if gi != nil {
				country, _ := gi.GetCountry(c.Real_ip)
				if len(country) < 2 {
					c.Country = "Lan"
				} else {
					c.Country = country
				}
			}
<<<<<<< HEAD
			if len(data) > 6 {
				c.Connected = data[6]
			}
			session[data[0]] = c
			//fmt.Println(c)

		}
	}
	return update

}

func getLineTitle(line string) string {
	return strings.Split(line, "\t")[0]
}

func getLineData(line string) []string {
	return strings.Split(line, "\t")[1:]
}

=======
			if len(tmp) > 6 {
				c.Connected = tmp[6]
			}
			session[tmp[1]] = c
		}
	}
	return update
}

>>>>>>> af6e5f1124a91023500ef623f73c6da469c51d12
func map_to_slice(in map[string]Client) []Client {
	var out []Client
	for _, client := range in {
		out = append(out, client)
	}
	return out
}

func reader(r io.Reader) string {
<<<<<<< HEAD
	buf := make([]byte, 256)
=======
	buf := make([]byte, 1024)
>>>>>>> af6e5f1124a91023500ef623f73c6da469c51d12
	var output string
	for {
		n, err := r.Read(buf[:])
		if err != nil {
<<<<<<< HEAD
			return "--"
		}
		if strings.HasSuffix(string(buf[0:n]), "\n") {
=======
			return ""
		}
		if strings.HasSuffix(string(buf[0:n]), "END\r\n") {
>>>>>>> af6e5f1124a91023500ef623f73c6da469c51d12
			break
		}
		output = fmt.Sprint(output, string(buf[0:n]))
	}
	return output
}

func writer(conn io.Writer, s string) {
	_, err := conn.Write([]byte(s))
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}
