package normalize_email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var DomainTagsMap = map[string]string{
	// Google only has two Gmail domains: https://en.wikipedia.org/wiki/List_of_Google_domains
	"gmail.com":      "+",
	"googlemail.com": "+",
	"google.com":     "+", // corporate email addresses; TODO presumably country domains also receive corporate email?
	// Microsoft
	"outlook.com": "+",
	"hotmail.com": "+",
	"live.com":    "+",
	// Fastmail - https://www.fastmail.com/help/receive/addressing.html TODO: whatever@username.fastmail.com -> username@fastmail.com
	"fastmail.com": "+",
	"fastmail.fm":  "+",
	// Yahoo Mail Plus accounts, per https://en.wikipedia.org/wiki/Yahoo!_Mail#Email_domains, use hyphens - http://www.cnet.com/forums/discussions/did-yahoo-break-disposable-email-addresses-mail-plus-395088/
	"yahoo.com.ar":         "-",
	"yahoo.com.au":         "-",
	"yahoo.at":             "-",
	"yahoo.be/fr":          "-",
	"yahoo.be/nl":          "-",
	"yahoo.com.br":         "-",
	"ca.yahoo.com":         "-",
	"qc.yahoo.com":         "-",
	"yahoo.com.co":         "-",
	"yahoo.com.hr":         "-",
	"yahoo.cz":             "-",
	"yahoo.dk":             "-",
	"yahoo.fi":             "-",
	"yahoo.fr":             "-",
	"yahoo.de":             "-",
	"yahoo.gr":             "-",
	"yahoo.com.hk":         "-",
	"yahoo.hu":             "-",
	"yahoo.co.in/yahoo.in": "-",
	"yahoo.co.id":          "-",
	"yahoo.ie":             "-",
	"yahoo.co.il":          "-",
	"yahoo.it":             "-",
	"yahoo.co.jp":          "-",
	"yahoo.com.my":         "-",
	"yahoo.com.mx":         "-",
	"yahoo.ae":             "-",
	"yahoo.nl":             "-",
	"yahoo.co.nz":          "-",
	"yahoo.no":             "-",
	"yahoo.com.ph":         "-",
	"yahoo.pl":             "-",
	"yahoo.pt":             "-",
	"yahoo.ro":             "-",
	"yahoo.ru":             "-",
	"yahoo.com.sg":         "-",
	"yahoo.co.za":          "-",
	"yahoo.es":             "-",
	"yahoo.se":             "-",
	"yahoo.ch/fr":          "-",
	"yahoo.ch/de":          "-",
	"yahoo.com.tw":         "-",
	"yahoo.co.th":          "-",
	"yahoo.com.tr":         "-",
	"yahoo.co.uk":          "-",
	"yahoo.com":            "-",
	"yahoo.com.vn":         "-",
}

type Option struct {
	ForceRemoveDots bool
	ForceRemoveTags bool
	DetectProvider  bool
}

var httpClient *http.Client = &http.Client{
	Timeout: time.Duration(5 * time.Second),
}

func Normalize(email string, optionsPtr *Option, callback chan string) string {
	email = strings.ToLower(strings.Trim(email, " "))

	emailParts := strings.Split(email, "@")
	user := emailParts[0]
	domain := emailParts[1]

	var options Option
	if optionsPtr != nil {
		options = *optionsPtr
	}

	if options.ForceRemoveTags {
		tags := []string{"-", "+", "="}
		for _, tag := range tags {
			idx := strings.Index(user, tag)
			if idx != -1 {
				user = strings.Replace(user, user[idx:], "", -1)
			}
		}
	} else {
		separator, ok := DomainTagsMap[domain]
		if ok {
			user = strings.Split(user, separator)[0]
		}
	}

	if domain == "googlemail.com" {
		domain = "gmail.com"
	}

	var isGoogleMail = domain == "gmail.com" || domain == "google.com"

	if options.ForceRemoveDots || isGoogleMail {
		user = strings.Replace(user, ".", "", -1)
	}

	if options.DetectProvider {
		go DetectProvider(user, domain, callback)
	}

	return user + "@" + domain
}

func DetectProvider(user, domain string, callback chan string) {
	url := fmt.Sprintf("%s?url=%s", "http://enclout.com/api/v1/dns/show.json", domain)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Println("err [normalize-email] fetching response from enclout", err)
		callback <- user + "@" + domain
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err [normalize-email] reading response", err)
		callback <- user + "@" + domain
		return
	}

	var data struct {
		DNSEntries []struct {
			RData string `json:"RData"`
			Type  string `json:"Type"`
		} `json:"dns_entries"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("err [normalize-email] unmarshaling response", err)
		callback <- user + "@" + domain
		return
	}

	for _, dnsEntry := range data.DNSEntries {
		if dnsEntry.Type == "MX" {
			r, _ := regexp.Compile("aspmx.*google.*\\.com\\.?$")
			if r.MatchString(dnsEntry.RData) {
				user = strings.Split(user, "+")[0]
				user = strings.Replace(user, ".", "", -1)
			}

			r, _ = regexp.Compile("\\.messagingengine\\.com\\.?$")
			if r.MatchString(dnsEntry.RData) {
				// dots are significant -
				// https://www.fastmail.com/help/account/changeusername.html
				user = strings.Split(user, "+")[0]
			}
		}
	}

	callback <- user + "@" + domain
	return
}
