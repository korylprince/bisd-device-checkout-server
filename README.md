# Info

This is the back end server for [bisd-device-checkout-client](https://github.com/korylprince/bisd-device-checkout-client), a small portal to assist in device checkout at Bullard ISD.

# Install

```
go get github.com/korylprince/bisd-device-checkout-server
```

Create a MySQL database with `model.sql`. (This matches [pyInventory](https://github.com/korylprince/pyInventory).)

# Configuration

    INVENTORY_LDAPSERVER="ad1.example.com"
    INVENTORY_LDAPPORT="389"
    INVENTORY_LDAPBASEDN="OU=base,DC=example,DC=com"
    INVENTORY_LDAPGROUP="Admin Group"
    INVENTORY_LDAPSECURITY="starttls"
    INVENTORY_SQLDRIVER="mysql"
    INVENTORY_INVENTORYDSN="username:password@tcp(server:3306)/database?parseTime=true"
    INVENTORY_SKYWARDDSN="DRIVER={Progress};HostName=server;DATABASENAME=database;PORTNUMBER=12501;LogonID=username;PASSWORD=password"
    INVENTORY_LISTENADDR=":8080"
    INVENTORY_PREFIX="/inventory" #URL prefix
