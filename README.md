# rdf2graph

Sample code to play with turtle files and create a graph based representation

## Example

```shell
curl -s https://schema.org/version/latest/schemaorg-current-http.ttl | ./rdf2graph http://schema.org/PostalAddress
<http://schema.org/PostalAddress>
        <http://www.w3.org/1999/02/22-rdf-syntax-ns#type>: [<http://www.w3.org/2000/01/rdf-schema#Class>]
        <http://www.w3.org/2000/01/rdf-schema#label>: ["PostalAddress"^^]
        <http://www.w3.org/2000/01/rdf-schema#comment>: ["The mailing address."^^]
        <http://www.w3.org/2000/01/rdf-schema#subClassOf>: [<http://schema.org/ContactPoint>]
                -> <http://schema.org/domainIncludes> -> <http://schema.org/addressLocality>
                -> <http://schema.org/domainIncludes> -> <http://schema.org/addressCountry>
                -> <http://schema.org/rangeIncludes> -> <http://schema.org/deliveryAddress>
                -> <http://schema.org/rangeIncludes> -> <http://schema.org/billingAddress>
                -> <http://schema.org/rangeIncludes> -> <http://schema.org/servicePostalAddress>
                -> <http://schema.org/domainIncludes> -> <http://schema.org/postOfficeBoxNumber>
                -> <http://schema.org/rangeIncludes> -> <http://schema.org/originAddress>
                -> <http://schema.org/rangeIncludes> -> <http://schema.org/itemLocation>
                -> <http://schema.org/domainIncludes> -> <http://schema.org/streetAddress>
                -> <http://schema.org/rangeIncludes> -> <http://schema.org/location>
                -> <http://schema.org/domainIncludes> -> <http://schema.org/postalCode>
                -> <http://schema.org/rangeIncludes> -> <http://schema.org/gameLocation>
                -> <http://schema.org/domainIncludes> -> <http://schema.org/addressRegion>
                -> <http://schema.org/rangeIncludes> -> <http://schema.org/address>
```
