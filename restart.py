from __init__ import *
import re
from install import *


network_down = '{}/fabcar/networkDown.sh'.format(fabric_samples)
network_up = '{}/fabcar/startFabric.sh'.format(fabric_samples)
key_path = '{}/test-network/organizations/peerOrganizations/org1.example.com/msp/keystore'.format(fabric_samples)
replace_pattern = '\\$\\{FABRIC_SDK_GO_PROJECT_PATH\\}/test-network/organizations/peerOrganizations/org1.example.com/msp/keystore/.+'
replace_str = '${FABRIC_SDK_GO_PROJECT_PATH}/test-network/organizations/peerOrganizations/org1.example.com/msp/keystore/'


def stop():
    os.system("cd {}/fabcar && ./networkDown.sh".format(fabric_samples))


def start():
    os.system("cd {}/fabcar && ./startFabric.sh".format(fabric_samples))


def update_config():
    for f in os.listdir(key_path):
        key_file = f

    if key_file:
        print('New key file is "{}"'.format(key_file))
        with open('client/config.yaml', 'r') as f:
            config = f.read()
            updated_config = re.sub(replace_pattern, replace_str+key_file, config)
            f.close()
            with open('client/config.yaml', 'w') as f:
                f.write(updated_config)
                f.close()


def update_sequence():
    global sequence
    if os.path.exists("sequence"):
        os.remove("sequence")

    with open('sequence', 'w') as f:
        f.write("1")
        f.close()
        sequence = get_sequence(False)


if os.path.exists('client/wallet/*'):
    os.remove('client/wallet/*')

stop()
start()
update_config()
update_sequence()
install_chaincode()
