package api

type contextKey int

//InventoryTransactionKey is the context key for the inventory database transaction for a request
const InventoryTransactionKey contextKey = 0

//SkywardTransactionKey is the context key for the skyward database transaction for a request
const SkywardTransactionKey contextKey = 1

//UserKey is the context key for the user for a request
const UserKey contextKey = 2
