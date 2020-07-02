# kea-dhcp-event
Kea DHCP Event hook

This program will post back to a custom CMDB when a Dell iDRAC management device obtains a DHCP4 lease.

Requirements:
  * zorun/kea-hook-runscript compiled and added to your hooks library
  * env variables:
    - KEA_LEASE4_ADDRESS
    - KEA_QUERY4_OPTION60
    - KEA_HOOK_DEBUG
    - KEA_FAKE_ALLOCATION
    - CMDB_DISCOVER_URL
    - KEA_CMDB_TOKEN
	* post to cmdb discovery url
	* print debug messages if KEA_HOOK_DEBUG is true,enabled, etc...
  * built payload from json: 
```
{
  "host": {
    "ip_address": "11.17.120.13",
    "vclass": "iDRAC"
    }
}
```
Kea Configuration:
```
{
  "library": "/usr/lib64/kea/hooks/kea-hook-runscript.so",
  "parameters": {
    "script": "/usr/local/bin/kea-dhcp-event",
    "wait": false
  }
},
```

Usage:
```
CMDB_URL="http://127.0.0.1:3000" KEA_CMDB_TOKEN="someapitoken" KEA_LEASE4_ADDRESS="11.17.120.13" KEA_QUERY4_OPTION60="iDRAC" ./kea-dhcp-event lease4_select
```

With Debug Logging:
```
time="2020-04-24T12:45:03-07:00" level=debug msg="[./kea-dhcp-hook lease4_select]"
time="2020-04-24T12:45:03-07:00" level=debug msg=11.17.120.13
time="2020-04-24T12:45:03-07:00" level=debug msg=iDRAC
time="2020-04-24T12:45:03-07:00" level=debug msg="{\"okay\}"
time="2020-04-24T12:45:03-07:00" level=debug msg="&{{11.17.120.13 iDRAC}}"
time="2020-04-24T12:45:03-07:00" level=debug msg="&{0xc0001150b0 http://127.0.0.1:3000 someapitoken}"
```

