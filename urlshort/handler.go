package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

//handler is the process that runs in response to a request
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	//fallback: an alternative plan that may be used in an emergency
	return func(w http.ResponseWriter, r *http.Request) {
		//if there is a match with the path, it would be redirected
		//otherwise fallback.ServeHTTP() would be called
		path := r.URL.Path
		if destination, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, destination, http.StatusFound)
			return //stop the rest code from running
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, fmt.Errorf("Got error in yaml.Unmarshal: %s", err)
	}
	pathToUrls := buildMap(pathUrls)

	//return a map handler using the map
	return MapHandler(pathToUrls, fallback), nil
}

type pathUrl struct {
	Path string
	Url  string
}

//parse YAML - yaml.Unmarshal
func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

//convert YAML array into map
func buildMap(pathUrls []pathUrl) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pathUrl := range pathUrls {
		pathToUrls[pathUrl.Path] = pathUrl.Url
	}
	return pathToUrls
}
