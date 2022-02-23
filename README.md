# Hyperlibrary

This is project for playing with Hyperledger Fabric. It deals with a fictitious library with various functions.

The project has two components. The chaincode and client code. The system is based on using the [Fabric Samples](https://github.com/hyperledger/fabric-samples) repository.

I've created several helper scripts for developing. Use the `restart.py` script to bootstrap the system.
The location of the Fabric Samples repo should be updated in the `__init__.py` file. Update the `fabric_samples` variable.

To push updates of the chaincode you can use the `install.py` command.

Use the `run.sh` command to run the client application. If you wish it is setup to emulate multiple clients concurrently.
You can pass in the argument which client to run.

Books are identified by its unique ISBN. In order to populate physical books into the library you need to `purchase` them.
