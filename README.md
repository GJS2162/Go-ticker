# Concurrent Data Routing to Redis using Go Tickers

## Project Overview

This project demonstrates the use of Go's concurrency features to route data to a Redis database using tickers. Two goroutines (G1 and G2) continuously send evenly spaced ticks one second apart. Whenever G1 has passed 3 ticks, it writes data to the Redis database. Whenever G2 has passed 7 ticks, it alerts G1 to send data to the Redis database. This project utilizes Go channels for synchronization and communication between the goroutines.

## Features

- **Concurrency**: Utilize Go's goroutines and channels for concurrent execution and synchronization.
- **Redis Integration**: Interact with a Redis database to store tick data.
- **Time-Based Operations**: Implement tickers to trigger periodic actions.
- **Docker**: Use Docker to run Redis locally for a consistent development environment.

## Requirements

- Go 1.16 or later
- Docker

## Setup

1. **Clone the Repository**:
    ```sh
    git clone https://github.com/GJS2162/Go-ticker.git
    cd go-ticker-redis
    ```

2. **Run Redis Locally via Docker**:
    ```sh
    docker run --name redis -p 6379:6379 -d redis
    ```

3. **Install Go Redis Package**:
    ```sh
    go get github.com/go-redis/redis/v8
    ```

## Running the Application

1. **Compile and Run**:
    ```sh
    go run main.go
    ```

2. **Observe the Output**:
    The application will print data written to Redis in the console, indicating successful writes triggered by the tickers.

## Checking Data in Redis

1. **Access the Redis CLI**:
    ```sh
    docker exec -it redis redis-cli
    ```

2. **List All Keys**:
    ```sh
    KEYS *
    ```

3. **Get Value of a Specific Key**:
    ```sh
    GET <key>
    ```

    Replace `<key>` with the actual key returned by the `KEYS *` command. The key will contain information such as the source (G1 or G2), tick count, and timestamp.

4. **Example**:
    ```sh
    GET G1-2024-06-28T12:34:56Z
    ```

    This command will return the data written by G1 at the specified timestamp.

## Conclusion

This project showcases how to effectively use Go's concurrency features and integrate with a Redis database, providing a foundation for building more complex and scalable applications. By using Docker to manage Redis, it ensures a consistent and reliable development environment.
