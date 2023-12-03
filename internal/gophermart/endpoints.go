package gophermart

const (
	// POST - user registration
	registerEndpoint = "/api/user/register"

	// POST - user autentification
	loginEndpoint = "/api/user/login"

	// GET - getting all user's orders
	// POST - download user's orders
	ordersEndpoint = "/api/user/orders"

	// GET - getting user's loyality balance
	balanceEndpoint = "/api/user/balance"

	// POST - user's loyality points withdraw
	balanceWithdrawEndpoint = "/api/user/balance/withdraw"

	// GET - information about loyality withdrawals
	allWithdrawalsEndpoint = "/api/user/withdrawals"
)
