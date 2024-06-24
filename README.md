# SkillQ

Simple skill querying website built with [React](https://react.dev/) & [Go](https://go.dev/).

## Setup

To setup the application, there are a couple of things that will be required:

1. [Bun](https://bun.sh/). This is used by the [client](./client/) to manage dependencies.
2. [Go](https://go.dev/). This is the programming language used for the backend.
3. [Docker](https://www.docker.com/). This is for running supported services required by the backend, such as the database, etc.

Once you have the above setup, proceed to setup and download the required dependencies:

1. Within the [client](./client/) directory, run:

    ```shell
    bun install
    ```

2. Within the [server](./server/) directory, run:

    ```shell
    go mod
    ```

## Running the application

1. First, run docker with:

    ```shell
    docker compose up
    ```

2. Secondly, in a separate terminal session, run the backend application:

    ```shell
    cd server/app/cmd
    go run main.go
    ```

3. Third, in a separate terminal session, run the frontend application:

    ```shell
    cd client
    bun run dev
    ```

That's it, now navigate to the browser on <http://localhost:5173/> and play around with the application :).
