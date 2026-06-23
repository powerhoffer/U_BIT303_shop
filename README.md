# YUTANK Shop Management System

[Chinese version](README-zh.md)

## Environment Requirements

This project is a GoFrame monolithic backend with a Vue admin frontend. The following environment has been verified for local development:

- OS: Ubuntu 20.04 on WSL
- Go: 1.25.5
- GoFrame CLI: 2.10.0
- GoFrame framework: 2.9.5
- MySQL: 8.0.x
- Node.js: 22.22.2
- npm: 10.9.7

Default local database settings:

    Host: 127.0.0.1
    Port: 3306
    User: root
    Password: 123456
    Database: yutank_shop

Note: yutank_shop is the default local development database name.

## Start the Project

### 1. Initialize the Database

Enter the project directory:

    cd /path/to/yutank-shop

Verify the MySQL connection:

    mysql -uroot -p123456
    SELECT 1;

Import the initialization SQL:

    mysql -uroot -p123456 < hack/init.sql

### 2. Start the Backend

Start the backend service first. The default API address is:

    http://127.0.0.1:8000

Command:

    cd /path/to/yutank-shop
    /usr/local/go/bin/go run main.go

Health check:

    http://127.0.0.1:8000/health

Swagger:

    http://127.0.0.1:8000/swagger/

### 3. Start the Admin Frontend

Then start the admin frontend. The default frontend address is:

    http://localhost:8080

Command:

    cd /path/to/yutank-shop/frontend_manage
    npm install
    npm run dev

### 4. Log In to the Admin Panel

Open the admin login page in your browser:

    http://localhost:8080/#/login

Test account:

    Username: root
    Password: 123456

After a successful login, you will enter the admin dashboard.

## Available Features

- Dashboard: view the current employee, personal credits, and management module overview.
- Employee List: create employees, edit employee information, enable or disable employees, and reset employee passwords.
- My Credits: view the current employee credit balance and credit records.
- Credit Operations: select an employee, add or deduct credits, and view the employee credit records.
- Category List: view current goods category information.
- Goods Management: create goods, edit goods, view details, and update shelf status.
- Frontend APIs: goods browsing and shopping cart APIs are available for later employee-facing pages.
- User menu in the top-right corner: change the current account password or log out.

## Notes

- The admin frontend sends backend API requests through the /backend route.
- Employee-facing APIs use the /frontend route.
- If login fails, first confirm that the target environment has the root employee account and that its password is 123456.
