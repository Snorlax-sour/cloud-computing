# Cloud-Native Web Ordering System Project

This project focuses on building a robust and scalable web-ordering system by leveraging **cloud-native principles** and **containerization technologies**. Our core objective is to **integrate** various components into a cohesive, functional system that operates seamlessly.

## Envisioned Architecture

The architecture described in this section represents our ideal plan for this web ordering system. Some functionalities and components are currently still under development or in the planning phase, and have not yet been fully implemented. The envisioned architecture for this web-ordering system (referred to as "Web-ordering-system-software_engineer" throughout this document) will utilize **Docker** to containerize its fundamental components:

* **Frontend**: Responsible for the user interface and user experience, providing the interactive elements for customers to place orders.
* **Backend**: Handles all core business logic, API services, and interactions with the database, acting as the brain of the system.
* **Database**: A dedicated SQL server (such as PostgreSQL or MySQL) will be used for persistent data storage and retrieval, ensuring data integrity and availability.

**Docker Compose will be employed to orchestrate and manage these interconnected Docker containers, ensuring they can communicate and function together as a unified system.** This approach facilitates independent development, deployment, and scaling of each component.

## Key Technologies

* **Docker**: For containerizing individual application components.
* **Docker Compose**: For defining and running multi-container Docker applications, managing their lifecycle and networking.
* **SQL Database**: SQLite

## Getting Started

1.  Clone this repository: $ git clone https://github.com/Snorlax-sour/cloud-computing.git
2.  Navigate to the project directory: $ cd ./cloud-computing
3.  Build and start the containers: $ docker compose up --build
4.  Access the frontend at: http://localhost:80
5.  Backend listen: http://localhost:5000

``` bash
git clone https://github.com/Snorlax-sour/cloud-computing.git
cd ./cloud-computing
docker compose up --build
```
---



