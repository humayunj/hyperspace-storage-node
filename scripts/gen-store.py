import sys

from eth_account import Account
import json
private_key = input("Enter private key: ")
password = input("Enter passphrase: ")
if len(private_key) == 0 or len(password) == 0:
    print("exiting")
    sys.exit()
acc = Account.encrypt(private_key,password)

with open("keystore",'w') as f:
    f.write(json.dumps(acc))
print("generated keystore as \"keystore\"")