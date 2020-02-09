# v2ray-rules-file
Convert the rules files of formats V2Ray or AutoProxy to geosite.dat files.

## Usage
```
Usage of v2ray-rules-file:
  -format string
        Format of the site files. (default "v2ray")
  -output string
        Path of the output .dat file. (default "geosite.dat")
  -sites string
        Folder storing site files. (default "sites")
```

## Formats
### V2Ray
Same as the domains in the config.json file. See https://www.v2ray.com/chapter_02/03_routing.html#ruleobject for details.
```
# Subdomain: matches example.com and it's subdomains
example.com
domain:example.com

# Plaintext: matches strings containing "example.com"
plain:example.com

# Full domain: matches "example.com" only
full:example.com

# Regular expression: matches example*.com
regex:example.*\.com
regexp:example.*\.com
```

### AutoProxy
v2ray-rules-file supports part of the AutoProxy format, which is used in GFWList and Adblock Plus.
```
! Plaintext: matches strings containing "example.com"
example.com

! Subdomain: matches example.com and it's subdomains
||example.com

! Full domain: matches "example.com" only
|example.com|

! Regular expression: matches example*.com
/example.*\.com/

! Wildcard character(*) will be converted to regular expression rule, for example
example*.com
! is the same as /example.*\.com/
|example*.com|
! is the same as /^example.*\.com$/

! Start anchor and end anchor are not supported and will be converted to plaintext rules, for example
|example.com
example.com|
! are the same as example.com
|example*.com
example*.com|
! are the same as /example.*\.com/

! Rules starting with @@ will be ignored
@@||example.com
```