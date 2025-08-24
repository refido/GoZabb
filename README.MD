# GoZabb

GoZabb is a wrapper for the Zabbix API written in Go. It provides a simple way to interact with Zabbix servers, allowing you to perform CRUD operations on hosts and other entities.

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/refido/GoZabb.git
    ```

2. Add values in your environment variables (.env)

3. To download dependencies, use the following command:

    ```bash
    go mod tidy
    ```

4. To run, use:

    ```bash
    go run main.go
    ```

    or

    ```bash
    go run .
    ```

## To-Do-List

- [x] Sanity check
- [ ] Better HTTP Router
- [ ] Middleware & Logging
- [ ] Struct Models
- [ ] Token Reuse / Expiry Handling
- [ ] Add Update Host (host.update) endpoint to complete CRUD.
- [ ] Try another resource (e.g. item.get, trigger.get).

## References

- [Zabbix API Documentation](https://www.zabbix.com/documentation/current/en/manual/api/reference)
