#!/bin/sh

export INVENTORY_INVENTORYDSN="$INVENTORY_USER:$INVENTORY_PASS@tcp($INVENTORY_HOST:$INVENTORY_PORT)/$INVENTORY_DATABASE?parseTime=true"
export INVENTORY_SKYWARDDSN="DRIVER={Progress};HostName=$SKYWARD_HOST;PORTNUMBER=$SKYWARD_PORT;DATABASENAME=$SKYWARD_DATABASE;LogonID=$SKYWARD_USER;PASSWORD=$SKYWARD_PASS"

/$GO_PROJECT_NAME
