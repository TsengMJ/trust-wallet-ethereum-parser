## Methods

### WebSocket:

Route: ws://localhost:8080/ws

- GetCurrentBlock

  ```js
  Message: {
    "action": "GetCurrentBlock"
  }
  Response: {
    "data": {
        "action": "GetCurrentBlock",
        "block": Object // Same as the response from ethereum eth_getBlockByHash
    },
    "error": String
  }
  ```

- Subscribe

  ```js
  Message: {
    "action": "Subscribe",
    "address": String // hex address with 0x
  }
  Response: {
    "data": {
        "action": "Subscribe",
        "subscribed": Boolean
    },
    "error": String
  }
  ```

- UnSubscribe
  ```js
  Message: {
    "action": "UnSubscribe",
    "address": String // hex address with 0x
  }
  Response: {
    "data": {
        "action": "UnSubscribe",
        "subscribed": Boolean
    },
    "error": String
  }
  ```

### Rest Api:

- GetCurrentBlock

  ```js
  Method: Get;
  Route: 'http://localhost:8080/current-block';
  Response: {
    "data": {
        "block": Object // Same as the response from ethereum eth_getBlockByHash
    },
    "error": String
  }
  ```

- GetTransactionsByAddress

  ```js
  Method: Get;
  Route: 'http://localhost:8080/transaction/:address';
  Param: {
    "address" String
  }
  Response: {
    "data": {
        "block": Object // Same as the response from ethereum eth_getBlockByHash
    },
    "error": String
  }
  ```
