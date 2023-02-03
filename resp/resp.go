package resp

//resp is protocol to communicate via network
//https://redis.io/docs/reference/protocol-spec/
const (
	//resp-type header
	SIMPLE  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'

	//suffix or prefix
	cr   = '\r'
	lf   = '\n'
	crlf = "\r\n"

	// ok   = "ok"
	// null = -1
)
