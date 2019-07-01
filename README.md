# iplocate.go

A small tool to gather and list location of attempted intruders caught by Fail2Ban auth protection.

This app will scan the auth.log and fail2ban.log located in `/var/log/` (linux/raspbian) at a given interval. 

Location and interval are set in the `./settings/config.json` file.

List can be displayed alphabetically or by country with the most failed login attempts.

App must be run using SUDO or as root due to the files it need access to demanding it.

Written in Go 1.11.2 (linux/ubuntu) - intended target environment, any debian based Linux distro really. Originally runs on my little Raspberry Pi 3B - but will be equally happy on a Cloud Virtual Machine, like a t2.micro from AWS for instance.

No 3rd party libraries used. To build a static binary enter `go build iplocate.go` and then run using `sudo ./iplocate`

Before use, ensure Fail2Ban is installed and configured correctly and running smoothly. Ensure that both relevant log files exist at given location. If in another location, first update this apps config file before running the app.

**To contact author:**

location   | name/handle
-----------|---------
github:    | rDybing
Linked In: | Roy Dybing
MeWe:      | Roy Dybing

---

## Releases

#### v.1.0.0: Soon™
- First release


---

## Fail2Ban Details

Fail2Ban is an intrusion prevention software framework that protects computer servers from brute-force attacks. 

More info here: https://www.fail2ban.org/wiki/index.php/Main_Page

---

## IP Location API Details

Location data is gathered from https://www.iplocation.net/

API examples and documentation at https://www.ip2location.com/web-service/ip2location

This app use their WS3 package, with requests to the API looking like this: 

`https://api.ip2location.com/v2/?ip=123.123.123.123&key=demo&package=WS3`

This app use the free demo key - which limits ip-lookups to 20 per day (24h period). For information on paid subscription, please visit their site.

#### Body JSON returns:

```javascript
{
	"country_code":		"",
	"country_name":		"",
	"region_name":		"",
	"city_name":		"",
	"credits_consumed":	2
}
```
---

## License: MIT

**Copyright © 2019 Roy Dybing** 

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

---

ʕ◔ϖ◔ʔ