#!/usr/bin/python

import os
import glob
import base64
import dns.update
import dns.query
import dns.tsigkeyring
from sys import exit
from pprint import pprint

KEYRING=".keyring"

if not os.path.isdir(KEYRING):
    print("Not a directory: {}".format(KEYRING))
    exit(1)

keys = {}
for filepath in glob.iglob(os.path.join(KEYRING, '*.private')):
    if not os.path.islink(filepath):
        continue
    kname = os.path.basename(filepath).split('.')[0]
    with open(filepath, 'r') as f:
        for line in f.readlines():
            if line.startswith('PrivateKey: '):
                key = line.split(':')[1][1:]
                # yes... :-)
                key = base64.b64decode(key)
                key = base64.b64encode(key)
                keys[kname] = key.decode('ascii')

print("Found {} private keys:".format(len(keys.keys())))
pprint(keys)

print("\nBuilding keyring...")
keyring = dns.tsigkeyring.from_text(keys)

update = dns.update.Update('example.org', keyring=keyring, keyname='ed25519')
update.add('hello', 30, 'TXT', '"hello world!"')

print("Sending query...")
response = dns.query.tcp(update, '127.0.0.1', port=53000)

print("Response received:\n")
print(response)
