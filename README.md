# iplocate.go

A small tool to gather and list location of attempted intruders caught by Fail2Ban auth protection.

This app will scan the auth.log and fail2ban.log located in `/var/log/` (linux/raspbian) at a given interval. 

Location and interval are set in the `./settings/config.json` file.

App will either list 10 most recent attempts and a sorted list of countries of origin - or you can list all attempts. Select by entering 1 or 2 followed by enter. 0 followed by enter to exit.

App must be run using SUDO or as root due to the files it need access to demanding it.

Written in Go 1.11.2 (linux/ubuntu) - intended target environment, any Debian based Linux distro really - but I suspect most Linux distros will work fine. Originally runs on my little Raspberry Pi 3B - but will be equally happy on a Cloud Virtual Machine, like a t2.micro from AWS for instance.

## Build

Only imports from standard library in this app. No need to get any 3rd party libraries or frameworks.

First make sure you have Go installed and configured correctly, then enter `go build iplocate.go` whilst in your local directory of this repo. Finally run using `sudo ./iplocate` and follow onscreen instructions.

Before use, ensure Fail2Ban is installed and configured correctly and running smoothly. Ensure that both relevant log files exist at given location. If in another location, first update this apps config file before running the app.

**Contact:**

location   | name/handle
-----------|---------
github:    | rDybing
Linked In: | Roy Dybing
MeWe:      | Roy Dybing

---

## Releases

- Version format: [major release].[new feature(s)].[bugfix patch-version]
- Date format: yyyy-mm-dd

#### v.0.3.0: 2019-07-04
- Added browse 10 entries at a time when listing all, and made it look pretty.

#### v.0.2.1: 2019-07-03
- Fixed where it would not parse fail2ban.log if formatting deviated even slightly from mine.
- Fixed so that the two Go-routines did not max out two cores... *cough*

#### v.0.2.0: 2019-07-03
- Added tally of countries.
- Added 10 most recent continuos update at set interval.
- Added List all as a separate option.
- Added app control using numbers + enter.

#### v.0.1.1: 2019-07-02
- Sorts on time descending (most recent first).
- Prettied the output somewhat.

#### v.0.1.0: 2019-07-01
- Only gets IPs from the fail2ban.log file, ignores auth.log for now. 
- No choice in how to display result, will squirt result out in a very haphazard fashion.
- No leave to run in background triggering on time interval. Will have to manually run the app to update.

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

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions: The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

---

ʕ◔ϖ◔ʔ