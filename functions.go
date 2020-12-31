package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

func main() {
	/*
		html := `<h1>Dogs</h1> <a href="medor.html">Medor</a>, <a href="rex.html">Rex</a>, <a href="rintintin.html">Rintintin</a>`
		matchRepetitivePattern(html)
		fmt.Printf("- %v -\n", getMD5Hash(html))
		fmt.Printf("- %v -\n", getRandomUserAgent())
	*/
	/*
			html := `<tr class="bordered">
		    <th class="light">Cat</th>
		    <td class="light"><a href="isidor.html">Isidor</a></td>
		  </tr>`
			matchMultiline(html)
	*/
	/*
		html := `Hello <a href="dolly.html">Dolly</a>`
		fmt.Printf("%v", removeHtmlTags(html))
	*/
}

// Regexp Samples

func matchRepetitivePattern(html string) {
	re := regexp.MustCompile(`<h1>Dogs</h1> (<a href=".+?">.+?</a>)(, <a href=".+?">.+?</a>)*`)
	match := re.FindString(html)
	fmt.Printf("%v", match)
}

func matchMultiline(html string) {
	re := regexp.MustCompile(`(?m)<tr class="bordered">
    <th class="light">Cat</th>
    <td class="light"><a href=".+?">(.+?)</a></td>
  </tr>`)
	match := re.FindStringSubmatch(html)
	fmt.Printf("%v", match)
}

func matchWithNonCapture(html string) {
	results := make(map[string]string)
	re := regexp.MustCompile(`(?i)(?:<b>[\d]+. <a href=")(.+?)(?:">)(.+?)(?:</a>)`)
	matches := re.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		results[match[1]] = match[2]
	}
}

func removeHTMLTags(html string) string {
	re2 := regexp.MustCompile(`<.+?>`)
	return re2.ReplaceAllString(html, "")
}

// Other samples

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func downloadCache(url string) string {
	cacheDir := "cache/"
	_ = os.Mkdir(cacheDir, 0700)
	cacheFile := cacheDir + getMD5Hash(url) + ".cache.txt"
	data, err := ioutil.ReadFile(cacheFile)
	if err == nil {
		return string(data) // return cached data
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", getRandomUserAgent())
	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return ""
	}
	_ = ioutil.WriteFile(cacheFile, body, 0644)
	return string(body)
}

func getRandomUserAgent() string {
	var agents = []string{
		"Mozilla/5.0 (X11; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
		"Mozilla/5.0 (X11; U; Linux amd64; rv:5.0) Gecko/20100101 Firefox/81.0 (Debian)",
	}
	rand.Seed(time.Now().UnixNano())
	return agents[rand.Intn(len(agents))]
}
