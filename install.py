import json
import subprocess

from __init__ import *
import re


def package():
    ret = os.system("cd backend && go build")

    if ret == 0:
        if os.path.exists(out_file):
            os.remove(out_file)

        os.system('peer lifecycle chaincode package {} --path {}/backend --lang golang --label {}'.format(out_file, cwd, label))
        print('Packaged: {}'.format(get_hash()))
    else:
        raise(Exception("Code build error!"))


def get_hash():
    out = subprocess.getoutput('sha256sum {}'.format(out_file))
    m = re.match('(.+)\s'+out_file, out)
    if m:
        return m.group(1).strip()


def check_installed():
    out = subprocess.getoutput('peer lifecycle chaincode queryinstalled').split('\n')
    hash = get_hash()
    for pkg in out:
        if hash in pkg:
            return True

    return False


def install():
    if check_installed():
        out = subprocess.getoutput('peer lifecycle chaincode install {}'.format(out_file))
        print('Already installed')
    else:
        out = subprocess.getoutput('peer lifecycle chaincode install {}'.format(out_file))
        if out.startswith('Error:'):
            print('Already installed!')
        else:
            print('Installed {}'.format(get_sequence()))


def approve(org):
    print('Approving for {} sequence {}'.format(org, get_sequence()))
    set_org(org)
    print(subprocess.getoutput(
    'peer lifecycle chaincode approveformyorg ' + \
    '-o {} '.format(orderer) + \
    '--ordererTLSHostnameOverride {} '.format(orderer_hostname) + \
    '--tls --cafile {} '.format(ca) + \
    '--channelID {} '.format(channel) + \
    '--name {} '.format(name) + \
    '--version {} '.format(version) + \
    '--package-id {}:{} '.format(label, get_hash()) + \
    '--sequence {} '.format(get_sequence()) + \
    '--init-required'))


def check_commit_readiness():
    print('Checking commit readiness for sequence {}'.format(get_sequence()))
    out = subprocess.getoutput('peer lifecycle chaincode checkcommitreadiness ' + \
    '--channelID {} '.format(channel) + \
    '--name {} '.format(name) + \
    '--version {} '.format(version) + \
    '--sequence {} --init-required --output json'.format(get_sequence()))
    print(out)
    if 'Error' not in out:
        approvals = json.loads(out)
        print(approvals['approvals'])
        return approvals['approvals']


def commit():
    print('Committing! sequence {}'.format(get_sequence()))
    out = subprocess.getoutput('peer lifecycle chaincode commit ' + \
    '-o {} '.format(orderer) + \
    '--ordererTLSHostnameOverride {} --tls '.format(orderer_hostname) + \
    '--cafile {} '.format(ca) + \
    '--channelID {} '.format(channel) + \
    '--name {} '.format(name) + \
    get_peers() + \
    '--version {} '.format(version) + \
    '--sequence {} --init-required'.format(get_sequence()))
    print(out)
    get_sequence(True)


def get_installed():
    out = subprocess.getoutput('peer lifecycle chaincode queryinstalled --output json')

    print(out)


def get_committed():
    out = subprocess.getoutput('peer lifecycle chaincode querycommitted ' + \
                               '--channelID {} '.format(channel) + \
                               '--name {} '.format(name) + \
                               '--output json')

    print(out)
    return json.loads(out)


def invoke_init():
    out = subprocess.getoutput('peer chaincode invoke '+ \
    '-o {} '.format(orderer) + \
    '--ordererTLSHostnameOverride {} '.format(orderer_hostname) + \
    '--tls --cafile {} '.format(ca) + \
    '-C {} '.format(channel) + \
    '-n {} '.format(name) + \
    get_peers() + \
    '--isInit -c \'{"function":"'+init_func+'","Args":[]}\'')
    print(out)


def install_chaincode():
    seq = get_sequence()
    package()

    for org in orgs:
        set_org(org)
        install()

    for org in orgs:
        org_name = 'Org{}MSP'.format(org)
        set_org(org)

        readiness = check_commit_readiness()
        if readiness[org_name]:
            print('Already approved!')
            break
        else:
            approve(org)

    # if approved == len(orgs):
    commit()
    # get_committed()

    if seq == 1:
        invoke_init()
# get_committed()
# get_installed()
# exit()

# install_chaincode()

# get_committed()
# get_installed()
# check_commit_readiness()
# package()
# invoke_init()


if __name__ =="__main__":
    install_chaincode()
