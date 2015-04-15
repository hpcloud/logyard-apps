#!/bin/bash

set -e


DEVDEPLOYPATH=$(pwd)/release/
COMPONENTS="logyard_sieve applog_endpoint apptail docker_events systail"


function deploy-local-vm {
	COMPONENT=${1}
	cp $GOPATH/bin/$COMPONENT $DEVDEPLOYPATH
}

mkdir -p $DEVDEPLOYPATH

for COMPONENT in $COMPONENTS; do
	deploy-local-vm $COMPONENT
done
