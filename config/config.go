// config exposes configuration values to be used by the server.
//configuration can be initialized from env variables or configuration file in future.
package config

import (
	"time"
)



// Delay specifies the duration
const Delay = 5 * time.Second

// Addr specifies the address at which the sever will lister
const Addr string = ":8080"
