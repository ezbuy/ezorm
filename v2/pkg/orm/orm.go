package orm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ezOrmObjs       = make(map[string]func() EzOrmObj)
	ezOrmObjsByID   = make(map[string]func(id string) (result EzOrmObj, err error))
	ezOrmObjsRemove = make(map[string]func(id string) (err error))
	Indexers        = make(map[string]func())
	postSetupHooks  = make(map[string]func())
)

func RegisterEzOrmObj(namespace, classname string, constructor func() EzOrmObj) {
	ezOrmObjs[namespace+"."+classname] = constructor
}

func RegisterEzOrmObjByID(namespace, classname string, f func(id string) (result EzOrmObj, err error)) {
	ezOrmObjsByID[namespace+"."+classname] = f
}

func RegisterEzOrmObjRemove(namespace, classname string, f func(id string) (err error)) {
	ezOrmObjsRemove[namespace+"."+classname] = f
}

func RegisterIndexer(namespace, classname string, indexer func()) {
	Indexers[namespace+"."+classname] = indexer
}

func RegisterSetupPostHooks(name string, className string, fn func()) {
	postSetupHooks[name+"."+className] = fn
}

func GetPostHooks(name string, className string) (func(), bool) {
	if post, ok := postSetupHooks[name+"."+className]; ok {
		return post, true
	}
	return nil, false
}

func NewEzOrmObjObj(namespace, classname string) EzOrmObj {
	constructor, ok := ezOrmObjs[namespace+"."+classname]
	if !ok {
		return nil
	}

	return constructor()
}

func NewEzOrmObjByID(namespace, classname, id string) (result EzOrmObj, err error) {
	f, ok := ezOrmObjsByID[namespace+"."+classname]
	if !ok {
		return nil, nil
	}

	return f(id)
}

func RemoveEzOrmObj(namespace, classname, id string) (err error) {
	f, ok := ezOrmObjsRemove[namespace+"."+classname]
	if !ok {
		return errors.New(namespace + "." + classname + " remove func not found")
	}

	return f(id)
}

type EzOrmObj interface {
	Id() string
	GetClassName() string
	GetNameSpace() string
}

type SearchObj interface {
	IsSearchEnabled() bool
	GetSearchTip() string
}

var (
	DateTimeLayout = "2006-01-02 15:04"
	DateLayout     = "2006-01-02"
	TimeLayout     = "15:04"
)

func I64DateTime(c int64) string {
	if c == 0 {
		return ""
	}
	return time.Unix(c, 0).Format(DateTimeLayout)
}

func I64Date(c int64) string {
	if c == 0 {
		return ""
	}
	return time.Unix(c, 0).Format(DateLayout)
}

func I64Time(c int64) string {
	if c == 0 {
		return ""
	}
	return time.Unix(c, 0).Format(TimeLayout)
}

func I32Time(c int32) string {
	return I64Time(int64(c))
}

func XGetQueryString(word string, fields []string) map[string]interface{} {
	queryString := make(map[string]interface{})
	queryString["default_operator"] = "AND"
	queryString["fields"] = fields
	queryString["query"] = word

	return queryString
}

func XGetQuery(key string, data map[string]interface{}) map[string]interface{} {
	query := make(map[string]interface{})
	query[key] = data
	args := make(map[string]interface{})
	args["query"] = query
	return args
}

func parseTime(layout, str string) time.Time {
	now := time.Now()
	t, _ := time.ParseInLocation(layout, str, now.Location())
	return t
}

func XGetSearchObj(word string, fields []string, params map[string]string, termKeys map[string]bool, dateKeys map[string]bool) map[string]interface{} {
	terms := make(map[string]string)
	ranges := make(map[string]map[string]int64)

	for k, v := range params {
		if v == "" {
			continue
		}
		if _, ok := termKeys[k]; ok {
			terms[k] = v
			continue
		}

		if isStart, ok := dateKeys[k]; ok {
			intVal := parseTime(DateLayout, v)
			if isStart {
				fieldName := k[0 : len(k)-5]
				if dateVal, ok := ranges[fieldName]; ok {
					dateVal["gte"] = intVal.Unix()
					ranges[fieldName] = dateVal
				} else {
					ranges[fieldName] = map[string]int64{
						"gte": intVal.Unix(),
						"lt":  intVal.AddDate(0, 0, 1).Unix(),
					}
				}
			} else {
				fieldName := k[0 : len(k)-3]
				if dateVal, ok := ranges[fieldName]; ok {
					dateVal["lt"] = intVal.AddDate(0, 0, 1).Unix()
				} else {
					ranges[fieldName] = map[string]int64{
						"gte": intVal.AddDate(0, 0, -1).Unix() + 1,
						"lt":  intVal.AddDate(0, 0, 1).Unix(),
					}
				}
			}
		}
	}

	if len(terms) == 0 && len(ranges) == 0 {
		query := XGetQuery("query_string", XGetQueryString(word, fields))
		PrintToJson(query)
		return query
	}

	filtered := make(map[string]interface{})
	if word != "" {
		filtered["query"] = map[string]interface{}{
			"query_string": XGetQueryString(word, fields),
		}
	}

	filter := make(map[string]interface{})
	var must []interface{}

	for k, v := range terms {
		must = append(must, map[string]interface{}{
			"term": map[string]string{
				k: v,
			},
		})
	}

	for k, v := range ranges {
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				k: v,
			},
		})
	}
	filter["bool"] = map[string]interface{}{
		"must": must,
	}
	filtered["filter"] = filter
	query := XGetQuery("filtered", filtered)
	PrintToJson(query)
	return query
}

func XGetMoreSearchObj(word string, fields []string, params map[string]interface{}, termKeys map[string]bool, dateKeys map[string]bool) map[string]interface{} {
	terms := make(map[string]interface{})
	ranges := make(map[string]map[string]int64)

	for k, v := range params {
		switch v := v.(type) {
		case string:
			if _, ok := termKeys[k]; ok {
				terms[k] = v
				continue
			}

			if isStart, ok := dateKeys[k]; ok {
				intVal := parseTime(DateLayout, v)
				if isStart {
					fieldName := k[0 : len(k)-5]
					if dateVal, ok := ranges[fieldName]; ok {
						dateVal["gte"] = intVal.Unix()
						ranges[fieldName] = dateVal
					} else {
						ranges[fieldName] = map[string]int64{
							"gte": intVal.Unix(),
							"lt":  intVal.AddDate(0, 0, 1).Unix(),
						}
					}
				} else {
					fieldName := k[0 : len(k)-3]
					if dateVal, ok := ranges[fieldName]; ok {
						dateVal["lt"] = intVal.AddDate(0, 0, 1).Unix()
					} else {
						ranges[fieldName] = map[string]int64{
							"gte": intVal.AddDate(0, 0, -1).Unix() + 1,
							"lt":  intVal.AddDate(0, 0, 1).Unix(),
						}
					}
				}
			}
		case []string:
			if len(v) == 0 {
				continue
			}
			if _, ok := termKeys[k]; ok {
				terms[k] = v
				continue
			}
		}
	}

	if len(terms) == 0 && len(ranges) == 0 {
		return XGetQuery("query_string", XGetQueryString(word, fields))
	}

	filtered := make(map[string]interface{})
	if word != "" {
		filtered["query"] = map[string]interface{}{
			"query_string": XGetQueryString(word, fields),
		}
	}

	filter := make(map[string]interface{})
	var must []interface{}
	var should []interface{}

	for k, v := range terms {
		switch v := v.(type) {
		case string:
			must = append(must, map[string]interface{}{
				"term": map[string]string{
					k: v,
				},
			})
		case []string:
			for _, val := range v {
				should = append(should, map[string]interface{}{
					"term": map[string]string{
						k: val,
					},
				})
			}
		}
	}

	for k, v := range ranges {
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				k: v,
			},
		})
	}
	filter["bool"] = map[string]interface{}{
		"must":   must,
		"should": should,
	}
	filtered["filter"] = filter

	return XGetQuery("filtered", filtered)
}

func XSortFieldsFilter(sortFields []string) (rtn []string) {
	rtn = make([]string, 0, len(sortFields))
	for _, s := range sortFields {
		if s != "" {
			rtn = append(rtn, s)
		}
	}
	return
}

func UniqURLParams(url_ string) string {
	parsedURL, _ := url.Parse(url_)
	values := parsedURL.Query()
	for k, v := range values {
		if len(v) > 1 {
			values[k] = []string{v[0]}
		}
	}
	parsedURL.RawQuery = values.Encode()
	r := parsedURL.String()
	return r
}

func ToJsonString(obj interface{}) string {
	bs, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

func PrintToJson(obj interface{}) {
	fmt.Println(ToJsonString(obj))
}

func GetIDFromSingleResult(single *mongo.SingleResult) (string, error) {
	var result map[string]interface{}
	err := single.Decode(&result)
	if err != nil {
		return "", fmt.Errorf("decode error: %w", err)
	}
	if id, ok := result["_id"].(string); ok {
		return id, nil
	}
	return "", errors.New("id not found")
}
