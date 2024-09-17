# Wizzl - Backend

A high-performance backend for the Wizzl app, built using [Go](https://github.com/golang/go) and [Fiber](https://github.com/gofiber/fiber). Wizzl is a simple messaging platform available as both a website and app.

## Project Description

The core logic of the app resides in the `app` folder. This includes HTTP handlers, request validation structs, services, and WebSocket event handlers. You'll find utility functions and globally accessible structs in the `pkg` folder. The application's routes are organized under the `routes` folder. Additionally, we store email templates in the `templates` folder.

We also provide a `Makefile` that simplifies handler creation. Using the command `make handler name=custom_handler`, you can generate a new handler file under `app/handlers` with a pre-defined struct and method.

### Usage

Wizzl is live at [wizzl.app](https://wizzl.app), and anyone can use it. If you'd like to run your own instance of Wizzl, follow these steps:

Ensure you have `docker` and `docker-compose` installed. Then, simply run:

```bash
docker compose up
```

This will start the backend service. By default, Wizzl runs on `127.0.0.1:3000`, but you can modify the `docker-compose.yml` file to change the port.

### Database

We use a **MySQL MariaDB** database alongside [GORM](https://gorm.io) as our **ORM**, which is a perfect fit for the app. **Redis** is also integrated to improve performance, cache chat data, and group users, enabling real-time permission-based message delivery.

All database models are located in the `database/models` folder. To add new models, define them in this folder, and manage relationships in the `relations.go` file for easier maintenance.

### Security

Our backend employs `bcrypt` for password encryption and `AES` for message encryption, ensuring user privacy. The encryption keys are **securely** stored.

## Support us

We truly appreciate any support, whether it’s through contributions, feedback, or spreading the word about **Wizzl**. Every bit helps us improve the platform and continue developing new features. 

If you’d like to support us financially, you can do so by donating through our [Ko-fi](https://ko-fi.com/bndrmrtn) page. Your donations help cover development costs, server expenses, and allow us to focus more on building a better experience for everyone.