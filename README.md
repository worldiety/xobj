# xobj
xobj is a go library which unifies the access to various data formats or markup languages in an object oriented way. 
Usually JSON or XML are used to represent some kind of serialized objects. 
However you need to use different apis to work with it, which is especially a problem, 
if you want to handle it in a generalized way. 
The common approach is to use a kind of marshaller/unmarshaller which works with 
concrete data structures. The main benefit of this is type safety, but has the following drawbacks:

* xml and json marshaller/unmarshaller use often a DOM, from which the type safe object 
hierarchy is inflated which doubles the amount of memory allocations
* the mapping to a fixed object structure will cause information loss, e.g. if you add 
nodes to the document in a future version, older versions will discard it in the 
marshalling/unmarshalling process.
* you cannot treat your typesafe structures in a generic way, besides reflection

## goals
* a nearly lossless representation for json and xml documents for reading.
 Also a struct hierarchy is wrappable (reflection)
* a common simply query and modification api
* a quite efficient serialization format
* a simple load/store stack, including network support, especially suited for quick prototyping
* support for some primitive types, using heuristics for type conversion

## non-goals
* a natural fit for json or xml, preserving every detail
* a entirely lossless representation for json and xml for a read/write cycle
* replace xquery or xpath
* replacement for json, xml, protobuf, flatbuffers, captai'n proto etc.
* streaming support
* best performance in terms of cpu usage or memory allocations
* support every data format or a formats primitive


## tldr
Actually *xobj* provides an API to load JSON into a contract which is more convenient to work
with. XML is transformed using [jsonml](https://github.com/worldiety/jsonml). Also it
provides some helper methods, to make prototyping faster.
