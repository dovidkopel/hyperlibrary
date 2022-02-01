export PATH=$PATH:/home/dkopel/go/src/github.com/dovidkopel/fabric-samples/bin/
export FABRIC_CFG_PATH=/home/dkopel/go/src/github.com/dovidkopel/fabric-samples/config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_MSPCONFIGPATH=/home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_LOCALMSPID=Org1MSP

cd /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network && scripts/deployCC.sh mychannel hyperledger /home/dkopel/IdeaProjects/hyper-library go 1.1 1 Init


peer chaincode list --instantiated -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel --peerAddresses localhost:7051 --tlsRootCertFiles /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt


peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n hyperledger --peerAddresses localhost:7051 --tlsRootCertFiles /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -isInit -c '{"function":"Init","Args":[]}'


peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n hyperledger --peerAddresses localhost:7051 --tlsRootCertFiles /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"ListBooks","Args":[]}'


peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n hyperledger --peerAddresses localhost:7051 --tlsRootCertFiles /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles /home/dkopel/go/src/github.com/dovidkopel/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"Invoke","Args":["create", "{\"isbn\":\"abcd9944\",\"author\":\"Foo Bar\",\"title\":\"Something 1\",\"genre\": 0}"]}'

