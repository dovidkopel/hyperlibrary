import os
import subprocess
import json


def get_org_certs():
    return {
        1: [
            'localhost:7051',
            '{}/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt'.format(fabric_samples),
            '{}/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp'.format(fabric_samples)
            ],
        2: [
            'localhost:9051',
            '{}/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt'.format(fabric_samples),
            '{}/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp'.format(fabric_samples)
            ]
    }


orgs = [1, 2]
fabric_samples = '/home/dkopel/go/src/github.com/dovidkopel/fabric-samples/'
channel = 'mychannel'
orderer = 'localhost:7050'
orderer_hostname = 'orderer.example.com'
org_certs = get_org_certs()
name = 'hyperlibrary'
out_file = '{}.tar.gz'.format(name)
version = 1.6
label = '{}_{}'.format(name, version)
sequence = 2
init_func = 'Init'
ca = '{}/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem'.format(fabric_samples)
cwd = os.getcwd()


def set_org(org):
    os.environ['CORE_PEER_LOCALMSPID'] = "Org{}MSP".format(org)
    os.environ['CORE_PEER_ADDRESS'] = get_org_certs()[org][0]
    os.environ['CORE_PEER_TLS_ROOTCERT_FILE'] = get_org_certs()[org][1]
    os.environ['CORE_PEER_MSPCONFIGPATH'] = get_org_certs()[org][2]


def package():
    os.system('peer lifecycle chaincode package {} --path {}/pkg/backend --lang golang --label {}'.format(out_file, cwd, label))
    print('Packaged')


def get_hash():
    return subprocess.getoutput('sha256sum {}'.format(out_file))


def check_installed():
    out = subprocess.getoutput('peer lifecycle chaincode queryinstalled').split('\n')
    hash = get_hash()
    for pkg in out:
        if hash in pkg:
            return True

    return False


def install():
    if check_installed():
        print('Already installed')
    else:
        out = subprocess.getoutput('peer lifecycle chaincode install {}'.format(out_file))
        if out.startswith('Error:'):
            print('Already installed!')
        else:
            print('Installed')


def approve(org):
    print('Approving for {}'.format(org))
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
    '--sequence {} '.format(sequence) + \
    '--init-required'))


def check_commit_readiness():
    out = subprocess.getoutput('peer lifecycle chaincode checkcommitreadiness ' + \
    '--channelID {} '.format(channel) + \
    '--name {} '.format(name) + \
    '--version {} '.format(version) + \
    '--sequence {} --init-required --output json'.format(sequence))
    print(out)
    if 'Error' not in out:
        approvals = json.loads(out)
        print(approvals['approvals'])
        return approvals['approvals']


def get_peers():
    certs = ''
    for org in get_org_certs().values():
        certs += '--peerAddresses {} --tlsRootCertFiles {} '.format(org[0], org[1])
    return certs


def commit():
    print('Committing!')


    out = subprocess.getoutput('peer lifecycle chaincode commit ' + \
    '-o {} '.format(orderer) + \
    '--ordererTLSHostnameOverride {} --tls '.format(orderer_hostname) + \
    '--cafile {} '.format(ca) + \
    '--channelID {} '.format(channel) + \
    '--name {} '.format(name) + \
    get_peers() + \
    '--version {} '.format(version) + \
    '--sequence {} --init-required'.format(sequence))
    print(out)


def get_installed():
    out = subprocess.getoutput('peer lifecycle chaincode queryinstalled ' + \
                               '--output json'.format(sequence))

    print(out)


def get_committed():
    out = subprocess.getoutput('peer lifecycle chaincode querycommitted ' + \
                               '--channelID {} '.format(channel) + \
                               '--name {} '.format(name) + \
                               '--output json'.format(sequence))

    print(out)
    return json.loads(out)


def invoke_init():
    out = subprocess.getoutput('peer chaincode invoke '+ \
    '-o {} '.format(orderer) + \
    '--ordererTLSHostnameOverride {}'.format(orderer_hostname) + \
    '--tls --cafile {} '.format(ca) + \
    '-C {} '.format(channel) + \
    '-n {} '.format(name) + \
    get_peers() + \
    '--isInit -c \'{"function":"{}","Args":[]}\''.format(init_func))
    print(out)


# get_committed()
# get_installed()
# exit()
package()

for org in orgs:
    set_org(org)
    install()

approved = 0

for org in orgs:
    org_name = 'Org{}MSP'.format(org)
    set_org(org)

    readiness = check_commit_readiness()
    if readiness[org_name]:
        print('Already approved!')
        approved += 1
        break
    else:
        approve(org)

if approved == len(orgs):
    commit()
    get_committed()
    invoke_init()
else:
    print('Not fully approved!')