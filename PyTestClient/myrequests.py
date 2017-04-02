import requests
import json
import sys

"""
This program will test basic functionality of the url shortener application.
Exit code will indicate status of the test.
Functionality being tested:
    1. Request homepage;
    2. Create short url;
    3. Verify short url.
"""

# colors used for command line output
class bcolors:
    OKGREEN = '\033[92m'
    FAIL = '\033[91m'
    BOLD = '\033[1m'
    ENDC = '\033[0m'

# Where to find the application
base_url = 'http://localhost:5100/'

# Url to be shortened
short_url = 'carpathians'
long_url = 'http://www.geo.arizona.edu/geo5xx/geo527.bck/Carpathians/index.html'

# json payload for creation of short url
data = {'shorturl': short_url, 'longurl': long_url }


## test homepage
print bcolors.BOLD + 'Request homepage:' + bcolors.ENDC
home = requests.get(base_url)

if home.status_code == 200:

    print bcolors.OKGREEN + 'Success!' + bcolors.ENDC

else:

    print bcolors.FAIL
    print 'Error!'
    print 'Status: ' + str(home.status_code)
    print 'Response: ' + home.text
    print bcolors.ENDC
    sys.exit(1)


## test creating new shortened url.
print bcolors.BOLD + 'Create new short url:' + bcolors.ENDC
create = requests.post( base_url + 'Create', data=json.dumps(data) )

if create.status_code == 200:

    print bcolors.OKGREEN + 'Success!' + bcolors.ENDC
    print 'short url: ' + short_url + ' - long url: ' + long_url

else:

    print bcolors.FAIL
    print 'Error!'
    print 'Status: ' + str(create.status_code)
    print 'Response: ' + create.text
    print bcolors.ENDC
    sys.exit(1)


## verify newly created short url
print bcolors.BOLD + 'Verify created short url:' + bcolors.ENDC
verify = requests.get( base_url + short_url, allow_redirects=False)

if verify.status_code in [301, 302] and verify.headers['Location'] == long_url:

    print bcolors.OKGREEN + 'Success!' + bcolors.ENDC
    print short_url + ' redirects to ' + long_url
    sys.exit(0)

else:

    print bcolors.FAIL
    print 'Error!'
    print 'Status: ' + str(verify.status_code)
    print 'Response: ' + verify.text
    print 'Headers: '
    print verify.headers
    print bcolors.ENDC
    sys.exit(1)
