package main
 
import (
	"fmt"
        "io/ioutil"
        "net/http"
	"github.com/miekg/dns" 
	"regexp"
	"flag"
)

func main() {
	
	// CONFIG
	
    config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
    c := new(dns.Client)
    m := new(dns.Msg)

	// Your target

    var domain string
    flag.StringVar(&domain, "d", "s3.shopify.com", "target.com.br")
    flag.Parse()
	
	// Find CNAME
	
    m.SetQuestion(string(domain) + ".", dns.TypeCNAME)
    m.RecursionDesired = true
    r, _, _ := c.Exchange(m, config.Servers[0]+":"+config.Port)

	// Request

	url := "http://" + (r.Answer[0].(*dns.CNAME).Target)
	fmt.Printf("S3 Domain: %s \n", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Checking

	re := regexp.MustCompile("NoSuchBucket")
	vuln := re.FindString(string(content))
	if (vuln == "NoSuchBucket") {
		fmt.Println("S3 Subdomain Takeover VULNERABLE! \n")
	} else {
		fmt.Println("Not Vulnerable! \n")
	}
}
