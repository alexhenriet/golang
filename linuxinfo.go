package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Details struct
type Details struct {
	osName     string
	osType     string
	osBased    string
	osOrigin   string
	osStatus   string
	osHomepage string
}

var distributions = make(map[string]string)

func getMapValues(strMap map[string]string) []string {
	var values []string
	for _, value := range strMap {
		values = append(values, value)
	}
	return values
}

func getMapKey(strMap map[string]string, search string) string {
	for key, value := range strMap {
		if strings.ToLower(value) == strings.ToLower(search) {
			return key
		}
	}
	return ""
}

func searchDistribution(search string) map[string]string {
	results := make(map[string]string)
	for key, name := range distributions {
		if strings.Contains(strings.ToLower(name), strings.ToLower(search)) {
			results[key] = name
		}
	}
	return results
}

func loadDistributions(url string) {

	bodyStr := downloadCache(url)
	re := regexp.MustCompile(`(?i)(?:<b>[\d]+. <a href=")(.+?)(?:">)(.+?)(?:</a>)`)
	matches := re.FindAllStringSubmatch(bodyStr, -1)
	for _, match := range matches {
		distributions[match[1]] = match[2]
	}
}

func getDetails(url string, name string, key string) Details {
	var details Details
	bodyStr := downloadCache(url + key)
	if bodyStr == "" {
		return details
	}
	details.osName = name

	re := regexp.MustCompile(`<b>Type d'OS:</b> <a href=".+?">(.+?)</a>`)
	match := re.FindStringSubmatch(bodyStr)
	if len(match) > 1 {
		details.osType = match[1]
	}

	re = regexp.MustCompile(`<b>Basée sur:</b> (.+?)<br />`)
	match = re.FindStringSubmatch(bodyStr)
	if len(match) > 1 {
		dirty := match[1]
		re2 := regexp.MustCompile(`<.+?>`)
		clean := re2.ReplaceAllString(dirty, "")
		details.osBased = clean
	}

	re = regexp.MustCompile(`<b>Origine:</b> <a href=".+?">(.+?)</a>`)
	match = re.FindStringSubmatch(bodyStr)
	if len(match) > 1 {
		details.osOrigin = match[1]
	}

	re = regexp.MustCompile(`<b>Statut:</b> <font.+?>(.+?)</font>`)
	match = re.FindStringSubmatch(bodyStr)
	if len(match) > 1 {
		details.osStatus = match[1]
	}

	re = regexp.MustCompile(`(?m)<tr class="Background">
    <th class="Info">Home Page</th>
    <td class="Info"><a href=".+?">(.+?)</a></td>
  </tr>`)
	match = re.FindStringSubmatch(bodyStr)
	if len(match) > 1 {
		details.osHomepage = match[1]
	}

	return details
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func downloadCache(url string) string {
	cacheDir := "cache/"
	_ = os.Mkdir(cacheDir, 0644)
	cacheFile := cacheDir + getMD5Hash(url) + ".cache.txt"
	data, err := ioutil.ReadFile(cacheFile)
	if err == nil {
		return string(data) // return cached data
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 5.1.1; Nexus 5 Build/LMY48B; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/43.0.2357.65 Mobile Safari/537.36")
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
