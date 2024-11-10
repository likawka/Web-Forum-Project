# Web Forum Project

## Run

### Docker

Docker is essential for this project, facilitating easy deployment and scalability. Before running the project, users need to update Docker to the latest version by following the instructions provided at the following link: [https://docs.docker.com/engine/install/](https://docs.docker.com/engine/install/)

### To run

1. Open terminal, clone this repo using command

```
git clone https://github.com/0-LY/Forum-test.git
```

2. Change current folder to folder with the project
3. To start the docker, enter:

```
bash scripts/run.sh
```

- To delete the docker and all data, enter:

```
bash scripts/stop.sh
```

OR you can run the project, execute `go run .` in your terminal.
Make sure you have all the necessary third-party packages installed.

Access the application via a web browser at `https://localhost:8080`

## Objectives

This project is geared towards creating a web forum with the following objectives:

- Facilitating communication between users.
- Associating categories with posts.
- Allowing users to like and dislike posts and comments.
- Implementing a filtering mechanism for posts.

## SQLite

To achieve the objectives, SQLite is chosen for data storage. It's a robust embedded database software, suitable for managing user data, posts, comments, etc. Refer to the provided entity relationship diagram and customize it according to project requirements. Utilize at least one SELECT, one CREATE, and one INSERT query.

## Authentication

User authentication is crucial. Users can register with unique credentials, including email, username, and password (with encrypted storage as a bonus task). Sessions are managed using cookies, allowing one active session per user with an expiration date.

## Communication

Users can create posts and comments for communication purposes. Only registered users have this privilege. When creating a post, users can associate one or more categories with it. All posts and comments are visible to all users, registered or not.

## Likes and Dislikes

Registered users can like or dislike posts and comments, with the count visible to all users.

## Filter

A filtering mechanism enhances user experience by allowing users to filter displayed posts by categories, created posts.

## Security

- **Security-First Approach**: We have rigorously focused on securing user data and interactions within the forum. Security measures have been deeply integrated into every layer of the application.
- **HTTPS Enabled**: Our forum is fully operational over HTTPS, ensuring secure communication through:

  - **Encrypted Connections**: We have generated SSL certificates for our website, establishing a verified and secure digital identity. Both self-signed certificates and those issued by Certificate Authorities (CAs) are supported.
  - **Secure Cipher Suites**: After extensive research, we have implemented the most secure and recommended cipher suites, ensuring robust encryption for all communications.
- **Rate Limiting Implemented**: To protect against brute-force attacks and ensure equitable resource usage, rate limiting has been integrated. This measure effectively prevents excessive requests from overloading our services.
- **Data Encryption**: All client passwords are securely encrypted, providing a high level of credential security.

## Authentication

Implemented integration of new authentication methods into the forum, enabling registration and login functionalities using Google and GitHub as authentication tools.
