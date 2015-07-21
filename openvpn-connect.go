package main

import (
    "io"
    "os"
    "fmt"
    "net"
    "strings"
    "flag"
    "text/template"
)


type Client struct {
    Name string
    Vpn_ip string
    Real_ip string
    Real_port string
    Connected string
    Upload string
    Download string
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
                <h1>Update</h1>
              </div>
              <div class="row">
                <div class="col-md-12">
                  <table class="table">
                    <thead>
                      <tr>
                        <th>Client name</th>
                        <th>VPN IP</th>
                        <th>Real IP</th>
                        <th>Real port</th>
                        <th>Upload</th>
                        <th>Download</th>
                        <th>Connected</th>
                      </tr>
                    </thead>
                    <tbody>
                    {{range .}}
                        <tr>
                            <td>{{.Name}}</td>
                            <td>{{.Vpn_ip}}</td>
                            <td>{{.Real_ip}}</td>
                            <td>{{.Real_port}}</td>
                            <td>{{.Upload}}</td>
                            <td>{{.Download}}</td>
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
    //skip_namePtr := flag.String("skip", "", "Skip dirs")
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
    writer(conn, "status\n")
    status := strings.Split(reader(conn), "\r\n")
    //fmt.Println(status)
    _ = parser(status, tmp)
    Session = map_to_slice(tmp)
    //fmt.Println(Session, update)
    t, _ := template.New("vpn").Parse(page)
    err = t.Execute(os.Stdout, Session)

    if err != nil {
        panic(err)
    }
}




func parser(lines []string, session map[string]Client) string{
    var update string
    virtual := false
    name := false
    for i := range lines {
        tmp := strings.Split(lines[i], ",")
        if tmp[0] == "Updated" {
            update = tmp[1]
        }

        if tmp[0] == "ROUTING TABLE" {
            continue
        }

        if tmp[0] == "Virtual Address" {
            virtual = true
            name = false
            continue
        }

        if tmp[0] == "Common Name" {
            virtual = false
            name = true
            continue
        }

        if virtual == true && name == false {
            //fmt.Println(tmp[1], len(tmp[1]))

                c := session[tmp[1]]
                c.Name = tmp[1]
                c.Vpn_ip = tmp[0]
                tmp1 := strings.Split(tmp[2], ":")
                c.Real_ip = tmp1[0]
                c.Real_port = tmp1[1]
                session[tmp[1]] = c

        }
        if virtual == false && name == true {
            //fmt.Println(tmp[0], len(tmp[0]))

                c := session[tmp[0]]
                c.Connected = tmp[4]
                c.Upload = tmp[3]
                c.Download = tmp[2]
                session[tmp[0]] = c

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
