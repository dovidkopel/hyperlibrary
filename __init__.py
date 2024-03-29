import os


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


def get_sequence(increment: bool = False):
    if os.path.exists('sequence'):
        with open('sequence', 'r') as f:
            seq = int(f.read())
            f.close()

            if increment:
                with open('sequence', 'w') as f:
                    f.write(str(seq+1))
                    f.close()
                    return seq+1
            else:
                return seq
    else:
        with open('sequence', 'w') as f:
            f.write("1")
            f.close()
            return 1


orgs = [1, 2]
fabric_samples = '/home/dkopel/go/src/github.com/dovidkopel/fabric-samples/'
channel = 'mychannel'
orderer = 'localhost:7050'
orderer_hostname = 'orderer.example.com'
org_certs = get_org_certs()
name = 'hyperlibrary'
out_file = '{}.tar.gz'.format(name)
version = 1.0
label = '{}_{}'.format(name, version)
# sequence = get_sequence()
init_func = 'Init'
ca = '{}/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem'.format(fabric_samples)
cwd = os.getcwd()

os.environ['CORE_PEER_TLS_ENABLED'] = 'true'
os.environ['FABRIC_CFG_PATH'] = '{}/config/'.format(fabric_samples)
os.environ['CORE_PEER_MSPCONFIGPATH'] = '{}/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp'.format(fabric_samples)
os.environ['CORE_PEER_TLS_ROOTCERT_FILE'] = '{}/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt'.format(fabric_samples)
os.environ['CORE_PEER_ADDRESS'] = 'localhost:7051'
os.environ['CORE_PEER_LOCALMSPID'] = 'Org1MSP'
BIN = '{}/bin/'.format(fabric_samples)

if BIN not in os.environ['PATH']:
    os.environ['PATH'] = '{}:{}'.format(os.environ['PATH'], BIN)


def set_org(org):
    os.environ['CORE_PEER_LOCALMSPID'] = "Org{}MSP".format(org)
    os.environ['CORE_PEER_ADDRESS'] = get_org_certs()[org][0]
    os.environ['CORE_PEER_TLS_ROOTCERT_FILE'] = get_org_certs()[org][1]
    os.environ['CORE_PEER_MSPCONFIGPATH'] = get_org_certs()[org][2]


def get_peers():
    certs = ''
    for org in get_org_certs().values():
        certs += '--peerAddresses {} --tlsRootCertFiles {} '.format(org[0], org[1])
    return certs