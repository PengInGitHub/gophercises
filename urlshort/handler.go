package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

//handler is the process that runs in response to a request
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//fallback: an alternative plan that may be used in an emergency

	//return http.HandlerFunc, which is a type of func that takes a ResponseWritter and a Request
	//no need to cast
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		//if there is a match with the path, it would be redirected
		if destination, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, destination, http.StatusFound)
			return //stop the rest code from running
		}
		//otherwise fallback.ServeHTTP() would be called
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	//get YAML array from []byte
	yamlArray, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, fmt.Errorf("Got error in yaml.Unmarshal: %s", err)
	}
	pathToURLs := buildMap(yamlArray)

	//return a map handler using the map
	return MapHandler(pathToURLs, fallback), nil
}

//used to parse YAML
type pathURL struct {
	Path string
	URL  string
}

//parse YAML - return YAML array
func parseYaml(data []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(data, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}

//convert YAML array into map
func buildMap(pathURLs []pathURL) map[string]string {
	pathToURLs := make(map[string]string)
	for _, pathURL := range pathURLs {
		pathToURLs[pathURL.Path] = pathURL.URL
	}
	return pathToURLs
}
