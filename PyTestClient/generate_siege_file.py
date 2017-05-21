#!/usr/bin/env python
'''
Generate test file for a Siege benchmark:
    1. Send new short urls to application;
    2. Write file with list of urls.

Generate file with:
python generate_siege_file.py

Benchmark with:
siege -c 500 -r 100 -f benchmark_url_list.txt
'''
import requests
import json
import random
import string

# Where to find application
app_url = 'http://localhost:5100/'
# How many shortened urls to insert
insert_keys = 1000
# How many shortened urls to search
search_keys = 5000

# Generate random string with given length
def random_string(length):
    return ''.join(random.choice(string.ascii_letters) for x in range(length))

'''
Create random shortened urls.
'''

created_keys = []

for i in range(0, insert_keys):

    shorturl = random_string(6)
    longurl = 'http://localhost/content/' + random_string(22)

    data = {'shorturl': shorturl, 'longurl': longurl}
    r = requests.post(app_url + 'Create',data=json.dumps(data))

    created_keys.append(shorturl)


'''
Create file with list of shortened urls to benchmark against.
Follows Pareto principle, 20%% of urls will receive 80%% of the requests.
'''

benchmark_url_list = ""
file = open('benchmark_url_list.txt', 'w')

for i in range(0, search_keys):

    r = random.random()
    if r > 0.4: # 60%
        benchmark_url_list += app_url + created_keys[random.randint(0, insert_keys*0.01)] + '\n'
    elif r > 0.2: # 20%
        benchmark_url_list += app_url + created_keys[random.randint((insert_keys*0.01)+1, insert_keys*0.2)] + '\n'
    else: # 20%
        benchmark_url_list += app_url + created_keys[random.randint((insert_keys*0.2)+1, insert_keys-1)] + '\n'
#    else:
#        benchmark_url_list += app_url + 'non-existant-url' + str(random.randint(1, 1000)) + '\n'

file.write(benchmark_url_list)
file.close()
