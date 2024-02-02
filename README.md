# Redis-Go
### Toy "Redis" written in Go with in functionality for basic Redis commands, key-value storage, RDB file parsing, and expiry.
Supports the following Redis CLI commands: PING, ECHO, GET, SET, KEYS *, CONFIG GET, and takes in dirname/file flags.

### Prerequisites

You will need redis-cli to utilize this as it will expect RESP format conversion (see https://redis.io/docs/reference/protocol-spec/). Install using Brew:
  ```sh
  brew install redis
  ```

### Usage
1. Clone the repo
   ```sh
   git clone https://github.com/connorpalermo/codecrafters-redis-go.git
   ```
2. Run the spawn_redis_server.sh script
   ```sh
   ./spawn_redis_server.sh
   ```
3. Try some example redis-cli commands!
   ```sh
   % redis-cli ping
   PONG
   % redis-cli echo yoohoo
   yoohoo
   % redis-cli set Hello World!
   OK
   % redis-cli get Hello
   World!
   ```
4. You can also parse RDB Files:
   ```sh
   ./spawn_redis_server.sh --dir /tmp/rdbfiles2764249175 --dbfilename pineapple.rdb # example file with multiple keys
   % redis-cli GET orange
   apple
   % redis-cli GET raspberry
   blueberry # keys from file
   ```
5. Can handle expiry as well (or parse Expiry from RDB File)
   ```sh
   % redis-cli set elephants world px 100
   OK
   % redis-cli get elephants
   world
   # Sleeping for 101ms
   % redis-cli get elephants
   nil
   ```
_Will likely clean this code up at some point, initially created to learn more about Go._