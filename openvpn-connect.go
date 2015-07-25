package main

import (
    "io"
    "os"
    "fmt"
    "net"
    "strings"
    "flag"
    "text/template"
    "github.com/abh/geoip"
    "strconv"
)


type Client struct {
    Name string
    Vpn_ip string
    Real_ip string
    Country string
    Real_port string
    Connected string
    Upload string
    Download string
}

type Data struct {
    Clients []Client
    Update string
}

const page = `<!DOCTYPE html>
        <html lang="en">
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
    serverPtr := flag.String("server", "127.0.0.1", "IP Server VPN")
    portPtr := flag.String("port", "5555", "Port server VPN")
    outputPtr := flag.String("file", "./vpn_page.html", "Output file name")
    flag.Parse()
    host := fmt.Sprint(*serverPtr, ":", *portPtr)
    tmp := make(map[string]Client)
    Session := make([]Client,0)
    conn, err := net.Dial("tcp", host)
    checkError(err)
    defer conn.Close()
    writer(conn, "state\n")
    _ = reader(conn)
    //fmt.Println(state)
    writer(conn, "status 3\n")

    status3 := strings.Split(reader(conn), "\r")
    //fmt.Println(status3)
    update := Parser3(status3, tmp)
    Session = map_to_slice(tmp)
    //fmt.Println(Session, update)
    fd, err := os.Create(*outputPtr)
    t, _ := template.New("vpn").Parse(page)
    data := Data{ Clients:Session, Update: update}
    err = t.Execute(fd, data)
    checkError(err)
}

func Parser3(lines []string, session map[string]Client) string{
    file := "/usr/share/GeoIP/GeoIP.dat"
    gi, err := geoip.Open(file)
    if err != nil {
        fmt.Printf("Could not open GeoIP database please install in /usr/share/GeoIP/\n")
    }

    var update string
    for i := range lines {
        var c Client
        tmp := strings.Split(lines[i], "\t")
        if i == 1 {
            update = tmp[1]
        }
        if i > 2 {

            c.Name = tmp[1]
            c.Vpn_ip = tmp[3]
            tmp1 := strings.Split(tmp[2], ":")
            c.Real_ip = tmp1[0]
            c.Real_port = tmp1[1]
            upload, _ := strconv.ParseInt(tmp[5], 10, 32)
            download, _ := strconv.ParseInt(tmp[4], 10, 32)
            c.Upload =  fmt.Sprint(upload/1000, " Kb")
            c.Download = fmt.Sprint(download/1000, " Kb")
            if gi != nil {
                country, _ := gi.GetCountry(c.Real_ip)
                if len(country) < 2 {
                    c.Country = "Lan"
                } else {
                    c.Country = country
                }
            }
            if len(tmp) > 6{
                c.Connected = tmp[6]
            }
            session[tmp[1]] = c
        }
    }
    return update
}



func map_to_slice(in map[string]Client) []Client {
    var out []Client
    for _, client := range in {

        out = append(out, client)

    }
    return out
}

func reader(r io.Reader) string {
    buf := make([]byte, 1024)
    var output string
    for {
        n , err := r.Read(buf[:])
        if err != nil {
            return ""
        }
        if strings.HasSuffix(string(buf[0:n]), "END\r\n") {
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
