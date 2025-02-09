# DeepWildcard for step-ca

DeepWildcard is a microservice for step-ca to overcome step-ca wildcard limitations.

## Integrate to step-ca

Integration to step-ca can be done with editing `ca.json`. 
You need to add options to provisioner.
Url path `authenticate/x509` is hard-coded into DeepWildcard

```json
{
  "...-...-...": {},
  "authority": {
    "provisioners": [
      {
        "type": "ACME",
        "name": "default",
        "challenges": [
          "http-01"
        ],
        "forceCN": true,
        "options": {
          "webhooks": [
            {
              "name": "deepWildcard",
              "url": "http://127.0.0.1:10443/authenticate/x509",
              "kind": "AUTHORIZING",
              "certType": "X509",
              "disableTLSClientAuth": true
            }
          ]
        }
      }
    ]
  },
  "...-...-...-...": {}
}
```

## Config

```yaml
# Address to listen.
# Defaults to ":10443" which means listen on all ip addresses on port 10443
address: ":10443"

# Validator is a submodule that allows or denies certificates
# Decision is made by rules
# Default action if no allow rules matched is to deny
# Deny rules always take precedence over allow rules
# See Syntax in README.MD
validator:
  dns:
    # Rules to allow. Defaults to empty list
    allow: 
    - "example.com"     # EXACT: matches only `example.com` in SAN
    - "*.example.com"   # CHILD: matches single level of subdomain, no wildcards
    - "+.example.com"   # GRAND: matches 2 or more levels of subdomain, no wildcards
    # Rules to deny. Defaults to empty list. Same syntax as `allow` list
    deny: 
    - "www.example.com" # EXACT: matches only `example.com` in SAN
    - "+.example.com"   # GRAND: matches 2 or more levels of subdomain, no wildcards
```

## DNS Rules

Rules define which domain names are allowed or denied. Rules are case-insensitive.
Additional restrictions apply, see [allowed domain names](#allowed-domain-names).

There are few kinds of rules:

* __Exact__: If rule string starts with letter or digit followed by letter or digit, it matches exact domain name supplied in request.
  * _Example_: `constoso.com`, `www.contoso.com`
  * _Wrong_: `com`, `local`, `servername`, `*.local`, `www.*.local`
  * _Note_: Rule applies only to SANs exactly. See __child__ and __grand__ rule kinds to get work with subdomains.

* __Child__: If rule string starts with `*.`, this rule match single level subdomain, except wildcard.
  * _Example_: `*.constoso.com`
  * _Wrong_: `*.local`, `www.*.local`
  * _Note_: this matches exactly one subdomain level. rule `*.dom.com` will match `any.dom.com`, but WILL NOT match `dom.com`, `any.any.dom.com`. See __grand__ rule kind for sub-sub-domains.

* __Grand__: If rule string starts with `+.`, this rule match 2 or more level subdomain (sub-domains and sub-sub-domains), except wildcard.
  * _Example_: `+.constoso.com`
  * _Wrong_: `+.local`, `www.+.local`
  * _Note_: this matches 2 + level. rule `+.dom.com` will match `any.any.dom.com` and `any.any.any.dom.com` but WILL NOT match `any.dom.com`, `dom.com`. When you need both subdomains and sub-sub-domains, define two rules: __child__ and __grand__

## Allowed domain names

* Only correct fully qualified dns names and wildcards for them are allowed:
* * Domain name must not start nor end with dot
* * * _good_ : `server-name.tld`
* * * _bad_ : `.server-name.tld`
* * * _bad_ : `server-name.tld.`
* * Domain name must have at least two domain components
* * * _good_ : `server-name.tld`
* * * _bad_ : `server-name-tld`
* * Domain components must not contain characters except `A`-`Z`, `0`-`9` and `-`
* * * _good_ : `xn--80aceu4cs.tld`
* * * _bad_ : `гитхаб.tld`
* * Domain components must not start nor end with hyphen
* * * _good_ : `server-name.tld`
* * * _bad_ : `server-name-.tld`
* * * _bad_ : `-server-name.tld`
* * Domain components must not contain underscore
* * * _bad_ : `_vlmcsd._tcp.domain.tld`
* * * _bad_ : `server_name.tld`
* * Domain components must be at least 1 character long
* * * _good_ : `server-name.1.tld`
* * * _bad_ : `server-name..tld`
* * Top level domain component must not contain digits, exchept when starts with `xn--`
* * * _good_ : `server-name.xn--j1amh`
* * * _bad_ : `server-name.3ld`
* * * _bad_ : `server-name.t3d`
* * * _bad_ : `server-name.tl3`
* * Top level domain component must not be `arpa`
* * * _bad_ : `server-name.arpa`
* * Domain component must be no longer 63 characters
* * * _good_ : `server-name.tld`
* * * _bad_ : `this-label-is-way-too-long-because-it-exceeds-63-characters-in-length.tld`
