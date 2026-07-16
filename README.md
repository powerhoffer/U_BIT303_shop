# YUTANK Shop System

[Chinese version](README-zh.md)

YUTANK is a Credits-based product redemption system for company employees. The project includes an admin panel, an employee storefront, and a GoFrame backend service. It covers employee and Credits management, goods and stock management, order processing, and administrator access control.

## Project Components

| Component | Technology | Default URL | Purpose |
| --- | --- | --- | --- |
| Backend | GoFrame 2 + MySQL | `http://127.0.0.1:8000` | Provides authentication, storefront, and management APIs |
| Admin Panel | Vue 2 + Element UI | `http://localhost:8080` | Manages employees, Credits, goods, stock, orders, and access control |
| Employee Storefront | Vue 2 + Element UI | `http://localhost:9528` | Browses products, redeems products with Credits, and manages personal orders |

The backend uses the following standard response format:

```json
{
  "code": 0,
  "message": "Success",
  "data": {}
}
```

Authenticated requests include the token in this request header:

```text
Authorization: Bearer <token>
```

## Current Features

### Admin Panel

- **Dashboard**: View the current administrator and an overview of the system modules.
- **Employees**: View details, create, edit, enable or disable, reset passwords, and soft-delete employees.
- **Credit Operations**: Add or deduct employee Credits and view employee Credit records.
- **Categories**: View goods categories.
- **Goods**: Filter, create, edit, and view goods; upload goods images; and control on-shelf or off-shelf status.
- **Stock**: Increase or decrease goods stock and filter stock records by goods, date, and change type.
- **Orders**: Filter orders by order number, employee, status, and date; view item details; and complete or cancel pending orders.
- **Access Control**: Manage administrators, roles, and API permissions, including status control, role assignment, and permission assignment.
- **Account Security**: Change the current administrator password and log out.

### Employee Storefront

- Register, log in, and log out of an employee account.
- Browse products by category or keyword and view product details.
- Add products to the cart, change quantities, or remove products.
- Submit redemption orders using Credits.
- View order lists and details, and cancel eligible orders.
- View the personal Credits balance and records.
- Change the current employee account password.

## Technology Stack

- Go `1.25.5`
- GoFrame CLI `2.10.0`
- GoFrame `2.9.5`
- MySQL `8.0.x`
- Vue `2.6.10`
- Element UI `2.13.2`
- Node.js `22.22.2`
- npm `10.9.7`

These versions have been verified in the current local development environment. Run all commands in WSL Ubuntu 20.04 where possible.

## Project Structure

```text
BIT303_shop/
├── api/                         # Backend API request and response contracts
│   ├── backend/                 # Admin panel and employee account API contracts
│   ├── frontend/                # Employee storefront API contracts
│   └── base/                    # Shared API contracts
├── frontend_manage/             # Admin panel Vue project
├── frontend_web/                # Employee storefront Vue project
├── hack/init.sql                # Database schema and initial data
├── internal/                    # GoFrame controllers, logic, services, DAOs, and models
├── manifest/config/config.yaml  # Backend runtime configuration
├── resource/public/upload/      # Goods images uploaded from the admin panel
└── main.go                      # Backend entry point
```

## Prerequisites

Install the required tools and confirm that they are available:

```bash
go version
mysql --version
node --version
npm --version
```

If MySQL is not running, start it in Ubuntu:

```bash
sudo service mysql start
```

## Database Configuration

The default configuration file is `manifest/config/config.yaml`:

```yaml
database:
  default:
    type: mysql
    host: 127.0.0.1
    port: "3306"
    user: root
    pass: "123456"
    name: bit303_shop
```

**Important:** Update the database address, username, and password for the target environment. The default configuration is intended for local development only and must not be used directly in production.

## Database Initialization

Enter the project root and import the initialization SQL:

```bash
cd /path/to/BIT303_shop
mysql -uroot -p123456 < hack/init.sql
```

`hack/init.sql` creates the `bit303_shop` database, business tables, base categories, permission data, the default administrator, and the demo employee account. The script uses idempotent operations and can be run again after schema updates.

## Start the Project

The backend, admin panel, and employee storefront run as three separate services. Use a separate terminal for each service.

### 1. Start the Backend

```bash
cd /path/to/BIT303_shop
go mod download
go run main.go
```

After startup, the following URLs are available:

- Health check: `http://127.0.0.1:8000/health`
- Swagger: `http://127.0.0.1:8000/swagger/`
- OpenAPI JSON: `http://127.0.0.1:8000/api.json`

### 2. Start the Admin Panel

```bash
cd /path/to/BIT303_shop/frontend_manage
npm install
npm run dev
```

Login URL: `http://localhost:8080/#/login`

### 3. Start the Employee Storefront

```bash
cd /path/to/BIT303_shop/frontend_web
npm install
npm run dev
```

Storefront URL: `http://localhost:9528`

Both frontend applications send requests to `http://127.0.0.1:8000` by default. Set `VUE_APP_BASE_API` when the backend is deployed at a different address.

## Test Accounts

| Application | Identity | Username | Password | Description |
| --- | --- | --- | --- | --- |
| Admin Panel | Super Administrator | `root` | `123456` | Display name: `System Administrator`; has all management permissions |
| Employee Storefront | Demo Employee | `root` | `123456` | Can be used to test the employee storefront; new employees can also register through the storefront |

Administrators and employees are separate identity types. Even when the demo usernames are the same, their tokens cannot be used interchangeably. Change all default passwords immediately after delivery or deployment.

## Admin Panel Guide

### Employees and Credits

1. Log in to the admin panel with the super administrator account.
2. Create or maintain employee accounts under **Employees**.
3. Select an employee under **Credit Operations** and grant the employee initial Credits.
4. View the employee's Credit changes on the same page.

### Goods and Stock

1. Confirm that the target category exists and is enabled under **Categories**.
2. Open **Goods** and create a new goods record.
3. Select a category and enter the English goods name, Credits Price, stock, and description.
4. Upload a goods image with the upload control, or enter an image URL manually.
5. Save the goods and confirm that its status is on shelf so that it appears in the employee storefront.
6. Manage subsequent stock adjustments and stock records under **Stock**.

### Orders

1. Filter orders under **Orders** by order number, employee ID, status, or date.
2. Open the order details to view the goods snapshots and total Credits.
3. Use **Complete** or **Cancel** for pending orders.
4. Completed or cancelled orders cannot be processed again.

### Administrator Permissions

1. Maintain API permissions under **Access Control / Permissions**.
2. Create roles and assign permissions under **Roles**.
3. Create administrators and assign roles under **Admins**.
4. Super administrators are not restricted by ordinary roles. The backend permission middleware performs the final authorization check for ordinary administrators.

## Employee Storefront Guide

1. Open the storefront to browse on-shelf products without logging in.
2. Log in or register an employee account to access Cart, My Orders, My Credits, and Account Security.
3. Add products to the cart and confirm their quantities.
4. When an order is submitted, the system validates the Credits balance and goods stock.
5. A successful order deducts Credits and stock. Cancelling an eligible order restores them according to the business rules.
6. View order status and details under **My Orders**, and view the balance and records under **My Credits**.

## API Boundaries

- `/backend/admin/...`: Administrator login, logout, profile, and password APIs; uses an administrator token.
- `/backend/employee/...`: Employee registration, login, profile, and password APIs; management APIs require an administrator token and the corresponding permission.
- `/backend/points/...`: Employee personal Credits and administrator Credit operation APIs.
- `/backend/goods/...`, `/backend/stock/...`, and `/backend/order/...`: Admin goods, stock, and order APIs.
- `/backend/role/...` and `/backend/permission/...`: Administrator RBAC APIs.
- `/frontend/category/...` and `/frontend/goods/...`: Public storefront browsing APIs.
- `/frontend/cart/...` and `/frontend/order/...`: Cart and order APIs that require an employee token.

Refer to Swagger for complete paths, parameters, and response fields.

## Build and Verification

### Backend Tests

```bash
cd /path/to/BIT303_shop
go test ./...
```

### Admin Panel Production Build

```bash
cd /path/to/BIT303_shop/frontend_manage
NODE_OPTIONS=--openssl-legacy-provider npm run build:prod
```

The build output is written to `frontend_manage/dist`.

### Employee Storefront Production Build

```bash
cd /path/to/BIT303_shop/frontend_web
npm run build
```

The build output is written to `frontend_web/dist`. If Node.js reports an OpenSSL compatibility error, run:

```bash
NODE_OPTIONS=--openssl-legacy-provider npm run build
```

## Data and Delivery Notes

- Business records, including goods, orders, employees, Credits, roles, and permissions, are stored in MySQL and are not synchronized automatically through Git.
- Images uploaded from the admin panel are stored under `resource/public/upload`. Upload records and goods image URLs are stored in the database.
- Back up both the MySQL database and `resource/public/upload` when migrating the complete environment.
- Importing only `hack/init.sql` into a new environment creates the base accounts, categories, and permissions, but it does not include formal goods that administrators later created locally.
- Create formal goods through the admin panel in the target environment, or migrate them through a verified database backup.
- `Credits` is the standard term displayed in the UI. Some legacy database and JSON fields still use `points` in their names; this does not affect the UI wording.

## Troubleshooting

### Login Reports an Invalid Username or Password

- Confirm that both the backend and MySQL are running.
- Confirm that the latest `hack/init.sql` has been imported.
- Use an administrator account in the admin panel and an employee account in the employee storefront.
- If the account is disabled or soft-deleted, an administrator must restore or recreate it.

### An API Returns 401 or 403

- `401` usually means that the token is missing, expired, or belongs to the wrong identity type.
- `403` usually means that an ordinary administrator does not have the required API permission.
- Log out and log in again. If the API still returns `403`, ask a super administrator to review the role and permission assignments.

### No Products Appear in the Storefront

- The initialization SQL does not contain goods created later by administrators.
- Create goods in the admin panel, upload their images, set their stock, and confirm that they are on shelf.

### Goods Images Do Not Load

- Confirm that the backend is running because it serves `/upload/...` images.
- Confirm that the matching files exist under `resource/public/upload` and that the runtime user has read permission.
- Confirm that the image URLs in the database match the URLs returned by the upload API.

### OpenSSL Error During Frontend Startup or Build

Vue CLI 4 may require a compatibility option with newer Node.js versions:

```bash
export NODE_OPTIONS=--openssl-legacy-provider
npm run dev
```

### Port Already in Use

Check the default ports:

```bash
lsof -i :8000
lsof -i :8080
lsof -i :9528
```

Stop the process using the port or configure a different port for the affected service. If the backend address changes, also update `VUE_APP_BASE_API` for both frontend applications.
