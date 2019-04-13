# xobj
xobj is a small go library which unifies the access to various data formats or markup languages in an object oriented way. Usually JSON or XML are used to simply represent some kind of serialized objects. However you need to use different apis to work with it, which is especially a problem, if you want to work with it in a generalized way. The common approach is to use a kind of marshaller/unmarshaller which works with concrete data structures. The main benefit of this is type safety, but has the following drawbacks:

* xml and json marshaller/unmarshaller use often a DOM, from which the type safe object hierarchy is inflated which doubles the amount of memory allocations
* the mapping to a fixed object structure will cause information loss, e.g. if you add information to the document in a future version, older versions will discard it in the marshalling/unmarshalling process.

## goals
* 
