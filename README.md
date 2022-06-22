# library-management
This is the codebase for online library management in GoLang. It contains three microservices for users, books and book-issue management. This repo is equipped with EFK logging.

The microservices can be run using the individual commands given below with the service description or they can be run collectively with the Docker.

### user-service:
This microservice accounts for user signup and login along with the authentication related tasks.

### book-service:
This microservice accounts for the books and authors related details management. It provides a structured book searching functionality with several filters and parameters.

### management-service:
This microservice is accountable for book-issue and availability management and services and maintains the record for books issue and returns.

