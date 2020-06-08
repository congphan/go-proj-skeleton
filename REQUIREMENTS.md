## Requirements

- Implement a Rest API with CRUD functionality. 
- Database: MySQL or PostgreSQL.
- Unit test whatever you can.

**Nice to have**:
- Automate setup by using Docker.
- Secure the API with your choice of authentication mechanism.

**Notes**:  
**Please refrain from releasing your code to public so that everyone has equal opportunities.**

## Data models

### User

- has accounts

### Account

- has at least the following fields:
    - `name`: name of account
    - `bank`: name of bank (3 possible values: `VCB`, `ACB`, `VIB`)
- has transactions

### Transaction

- has at least the following fields:
    - `amount`:  amount of money
    - `transaction_type`: type of transaction (2 possible values: `withdraw`, `deposit`)

## Detail of endpoints

#### 1. Get transactions of an user

- URL path: `/api/users/<user_id>/transactions`
- HTTP method: `GET`
- Request:
    - Parameters:
        |Name|Required|Data type|Description|
        | --- | --- | --- | --- |
        |`user_id`|Yes|Integer|User's ID|
        |`account_id`|No|Integer|Account's ID|
    - Note: When `account_id` is not specified, return all transactions of the user.
    - Please have validations for required fields

- Response:
    - Content type: `application/json` 
    - HTTP status: `200 OK`
    - Body: Array of user's transactions, each of which has the following fields:

        |Name|Data type|Description|
        | --- | --- | --- |
        | `id` |Integer| Transaction's ID |
        | `account_id` |Integer| Account's id |
        | `amount` |Decimal| Amount of money |
        | `bank` |String| Bank's name |
        | `transaction_type` |String| Type of transaction |
        | `created_at` |String| Created date of transaction |

- Example:  GET `/api/users/1/transactions?account_id=1`
  - Response:
    ```json
    [{
      "id": 1,
      "account_id": 1,
      "amount": 100000.00,
      "bank": "VCB",
      "transaction_type": "deposit",
      "created_at": "2020-02-10 20:00:00 +0700"
    }, { ... }]
    ```

#### 2. Create a transaction for an user
- URL path: `/api/users/<user_id>/transactions`
- HTTP method: `POST`
- Request:
    - Parameters:

        |Name|Required|Data type|Description|
        | --- | --- | --- | --- |
        |`user_id`|Yes|Integer|User's ID|

    - Body:

        |Name|Required|Data type|Description|
        | --- | --- | --- | --- |
        |`account_id`|Yes|Integer|Account's ID|
        | `amount`|Yes|Decimal| Amount of money |
        | `transaction_type`|Yes |String| Type of transaction |
    - Please have validations for required fields

- Response:
    - Content type: `application/json` 
    - HTTP status: `201 Created`
    - Body: Details of the created transaction with the following fields:

        |Name|Data type|Description|
        | --- | --- | --- |
        | `id` |Integer| Transaction's ID |
        | `account_id` |Integer| Account's id |
        | `amount` |Decimal| Amount of transaction |
        | `bank` |String| Bank's name |
        | `transaction_type` |String| Type of transaction |
        | `created_at` |String| Created date of transaction |

- Example: POST `/api/users/1/transactions`
  - Request body:
    ```json
    {
      "account_id": 2,
      "amount": 100000.00,
      "transaction_type": "deposit"
    }
     ```  
  - Response
    ```json
    {
      "id": 10,
      "account_id": 2,
      "amount": 100000.00,
      "bank": "VCB",
      "transaction_type": "deposit",
      "created_at": "2020-02-10 20:10:00 +0700"
    }
    ```
#### 3. Think about PUT and DELETE by yourself

## What we are looking for

Pay close attention to the structure and the dependency between packages of your application. This determines how testable your code is. It also allows you to create an API structure that is highly composable, modular and easy to read.

On the whole, we care about attention to detail, the idiomatic use of Go, and an understanding of what an API's role is.

Happy Coding!
