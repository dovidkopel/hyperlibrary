from __init__ import *
import subprocess
import json
import re
import time


def invoke_init(func, args):
    cmd = 'peer chaincode invoke '+ \
    '-o {} '.format(orderer) + \
    '--ordererTLSHostnameOverride {} '.format(orderer_hostname) + \
    '--tls --cafile {} '.format(ca) + \
    '-C {} '.format(channel) + \
    '-n {} '.format(name) + \
    get_peers() + \
    ' --waitForEvent -c \'{"function":"'+func+'","Args":'+json.dumps(args)+'}\''
    # print(cmd)
    out = subprocess.getoutput(cmd)
    print(out)
    m = re.match('.+status:200 payload:"(.+)"', out)
    if m:
        payload = json.loads(re.sub('\\\\"', '"', m.group(1)))
        return payload
    else:
        print(out)


# invoke_init('ListBooks', [])
print(invoke_init('PurchaseBook', ['abcd45454', '1', '10.50']))
# time.sleep(5)
print(invoke_init('ListBooks', []))