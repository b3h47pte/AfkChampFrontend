/*
 * 'json' provides basic functionality in reading/writing json data from/to a request.
 */
package utility

import(
  "net/http"
  "encoding/json"
)

// ReadJsonFromRequestBodyMap takes in a request that has JSON data in its body and returns a map if it parses it properly.
func ReadJsonFromRequestBodyMap(r *http.Request) (map[string]interface{}, error) {
  decoder := json.NewDecoder(r.Body)
  var postData map[string]interface{}
  err := decoder.Decode(&postData)
  if err != nil {
    return nil, err
  }
  return postData, nil
}

// ReadJsonFromRequestBodyStruct reads the JSON request into the given data structure.
func ReadJsonFromRequestBodyStruct(r *http.Request, s interface{}) error {
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(s)
  return err
}

// WriteJsonToResponse writes the given structure into the response
func WriteJsonToResponse(w http.ResponseWriter, data interface{}) error {
  encoder := json.NewEncoder(w)
  return encoder.Encode(data)
}