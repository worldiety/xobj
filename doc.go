// The package xobj provides a convenient key/value api for []interface{} and map[string]interface{} as
// used by the json package to represent JSONObjects.
//
// API discussion
//
// However the API has a dependency problem: even if everything is defined as an interface, a different implementation
// cannot satisfy our interface contract, because go cannot duck type recursive interface definitions, see also
// https://github.com/golang/go/issues/8082 for details.
//
// Alternatives are to use a contract without recursion and using helper functions instead, but that would
// be detrimental anyway. Also anonymous interfaces would not be of any help, because we require the recursive
// definition. Also a type alias would not work, because an alias cannot be recursive and the compiler even rejects it.
package xobj
