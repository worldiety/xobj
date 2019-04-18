package xobj

// design discussion: it is a bad performance and memory thing, that we cannot use map[string]interface{} and
// []interface{} directly but need to wrap that. On the other hand, we need that for go mobile. An alternative
// would be something like
type Document struct {
	Root map[string]interface{} // this member would not get an accessor in go mobile
}

// selector would be something like /myobject/persons[0]/@name
func (d Document) AsString(selector string) string {
	return ""
}

// selector would be something like /myobject (== amount of keys) or /myobject/persons (== amount entries)
func (d Document) Size(selector string) int {
	return 0
}

// /myobject/persons[0]/ name HelloName
func (d Document) Put(selector string, key string, value string) {

}
