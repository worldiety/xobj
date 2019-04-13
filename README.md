# xobj
xobj is a go library which unifies the access to various data formats or markup languages in an object oriented way. Usually JSON or XML are used to represent some kind of serialized objects. However you need to use different apis to work with it, which is especially a problem, if you want to handle it in a generalized way. The common approach is to use a kind of marshaller/unmarshaller which works with concrete data structures. The main benefit of this is type safety, but has the following drawbacks:

* xml and json marshaller/unmarshaller use often a DOM, from which the type safe object hierarchy is inflated which doubles the amount of memory allocations
* the mapping to a fixed object structure will cause information loss, e.g. if you add nodes to the document in a future version, older versions will discard it in the marshalling/unmarshalling process.
* you cannot treat your typesafe structures in a generic way, besides reflection

## goals
* a nearly lossless representation for json and xml documents for reading. Also a struct hierarchy is wrappable (reflection)
* a common simply query and modification api
* a quite efficient serialization format
* a simple load/store stack, including network support, especially suited for quick prototyping
* support for some primitive types, using heuristics for type conversion

## non-goals
* a natural fit for json or xml, preserving every detail and order
* a lossless representation for json and xml for a read/write cycle
* replace xquery or xpath
* replacement for json, xml, protobuf, flatbuffers, captai'n proto etc.
* streaming support
* best performance in terms of cpu usage or memory allocations
* support every data format or a formats primitive


## mapping rules by example

```xml
<person>
    <!-- a comment -->
    <name firstname="true">Frodo</name>
    <age>42</age>
    <address>
        <city>Hobbingen</city>
    </address>
    <hobbies>
        <hobby>eating</hobby>
        <hobby>sleeping</hobby>   
        <hobby>reading</hobby> 
        <else>and another loss</else>
    </hobbies>
    <missing>null</missing>
</person>
```

```json
{
  ".comment": "a comment",
  "name": "Frodo",
  "name.firstname": true,
  "age": 42,
  "address": {
    "city": "Hobbingen"
  },
  "hobbies": ["eating","sleeping","reading","and another loss"],
  "missing": null
}
```

As you can see, the mapping between json and xml takes a lot of assumptions, which cannot be applied in general
but should be a good fit in a lot of situations:

* a json attribute and a xml text node in an element are interchangeable. However converting back and forth will
cause a loss of type information for the json (e.g. boolean, number or string). The string *null* in xml will be interpreted
as `null` in json.
* xml attributes have no direct correspondence in json, so they are encoded as *node name*.*attrib name*.
* xml comments are encoded as *.comment*
* json types are guessed, from their possible natural interpretation
* All types can be duck typed, causing further data loss. Rules:
  * `null`: the json *null* value and the string "null" will evaluate to nil or null
  * `float`: a bool will evaluate to 0 or 1 (true) and integers to their nearest floating point representation.
   Anything else will evaluate to NaN.
  * `int`: a bool will evaluate to 0 or 1 (true) and floats will be truncated. Anything else will be evaluated to 0.
  * `string`: a bool will be evaluated to "true" or "false". Integers and floats as usual. *null* will become "null".
  * `bool`: the number 1 or the strings "true","1" or any "tRuE" will evaluate to true, everything else false.




